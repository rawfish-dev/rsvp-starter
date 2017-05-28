package base

import (
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/interfaces"
)

type Service struct {
	interfaces.Logger
}

func NewService(logger interfaces.Logger) *Service {
	return &Service{logger}
}
