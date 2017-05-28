package testhelpers

import (
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/config"
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/services/jwt"
)

func NewTestJWTService() jwt.JWTServiceProvider {

	testConfig := config.TestConfig().JWT

	testBaseService := NewTestBaseService()

	return jwt.NewService(testBaseService, testConfig)
}
