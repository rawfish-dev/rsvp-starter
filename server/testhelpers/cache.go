package testhelpers

import (
	"github.com/rawfish-dev/rsvp-starter/server/services/cache"
)

func NewTestCacheService() cache.CacheServiceProvider {

	testBaseService := NewTestBaseService()

	return cache.NewService(testBaseService)
}
