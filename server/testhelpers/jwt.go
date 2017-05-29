package testhelpers

import (
	"github.com/rawfish-dev/rsvp-starter/server/config"
	"github.com/rawfish-dev/rsvp-starter/server/services/jwt"
)

func NewTestJWTService() jwt.JWTServiceProvider {

	testConfig := config.LoadConfig().JWT

	testBaseService := NewTestBaseService()

	return jwt.NewService(testBaseService, testConfig)
}
