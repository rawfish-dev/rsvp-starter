package guest

import (
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/interfaces"
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/services/base"
)

type service struct {
	baseService  *base.Service
	guestStorage interfaces.GuestStorage
}

func NewService(baseService *base.Service, guestStorage interfaces.GuestStorage) GuestServiceProvider {
	return &service{baseService, guestStorage}
}
