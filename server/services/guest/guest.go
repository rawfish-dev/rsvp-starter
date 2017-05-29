package guest

import (
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"
	"github.com/rawfish-dev/rsvp-starter/server/services/base"
)

type service struct {
	baseService  *base.Service
	guestStorage interfaces.GuestStorage
}

func NewService(baseService *base.Service, guestStorage interfaces.GuestStorage) GuestServiceProvider {
	return &service{baseService, guestStorage}
}
