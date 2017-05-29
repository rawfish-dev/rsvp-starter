package base

import (
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"
)

type Service struct {
	interfaces.Logger
}

func NewService(logger interfaces.Logger) *Service {
	return &Service{logger}
}
