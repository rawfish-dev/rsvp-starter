package testhelpers

import (
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/logger"
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/services/base"

	"github.com/satori/go.uuid"
)

func NewTestBaseService() *base.Service {
	testContextID := "test-" + uuid.NewV4().String()
	testLogger := logger.NewLoggerWithContext(testContextID)

	return base.NewService(testLogger)
}
