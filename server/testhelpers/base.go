package testhelpers

import (
	"github.com/rawfish-dev/rsvp-starter/server/logger"
	"github.com/rawfish-dev/rsvp-starter/server/services/base"

	"github.com/satori/go.uuid"
)

func NewTestBaseService() *base.Service {
	testContextID := "test-" + uuid.NewV4().String()
	testLogger := logger.NewLoggerWithContext(testContextID)

	return base.NewService(testLogger)
}
