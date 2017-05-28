package cache

import (
	"fmt"
	"sync"
	"time"

	"bitbucket.org/rawfish-dev/wedding-rsvp/server/services/base"
)

type service struct {
	baseService *base.Service
	storage     map[string]valueWrapper
	mutex       *sync.Mutex
}

type valueWrapper struct {
	value     string
	timestamp int64
}

var cacheService *service
var once sync.Once

func NewService(baseService *base.Service) CacheServiceProvider {
	once.Do(func() {
		cacheService = &service{
			baseService: baseService,
			storage:     make(map[string]valueWrapper),
			mutex:       &sync.Mutex{},
		}
	})

	return cacheService
}

func (s *service) Get(key string) (value string, err error) {
	wrappedValue, ok := s.storage[key]
	if !ok {
		return "", nil
	}
	return wrappedValue.value, nil
}

func (s *service) SetWithExpiry(key string, value string, expiryInSeconds int) (err error) {
	if key == "" {
		// TODO:: Make into a service error
		return fmt.Errorf("cache keys cannot be blank")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	timestamp := time.Now().UnixNano()

	s.storage[key] = valueWrapper{
		value:     value,
		timestamp: timestamp,
	}

	go func(storageKey string, timeMarker int64) {
		timer := time.NewTimer(time.Second * time.Duration(expiryInSeconds))
		<-timer.C

		s.mutex.Lock()
		defer s.mutex.Unlock()

		wrappedValue, ok := s.storage[key]
		if !ok {
			// Value already does not exist
			return
		}

		if wrappedValue.timestamp != timeMarker {
			// Value has already been set again
			return
		}

		s.baseService.Infof("cache service - expiring key %v", storageKey)

		delete(s.storage, storageKey)
	}(key, timestamp)

	return nil
}

func (s *service) Delete(key string) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.storage, key)
	return nil
}

func (s *service) Exists(key string) (exists bool, err error) {
	_, ok := s.storage[key]
	return ok, nil
}

func (s *service) Flush() (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.storage = make(map[string]valueWrapper)
	return nil
}
