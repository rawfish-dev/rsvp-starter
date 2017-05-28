package testhelpers

import (
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/services/cache"
)

func NewTestCacheService() cache.CacheServiceProvider {

	testBaseService := NewTestBaseService()

	return cache.NewService(testBaseService)
}
