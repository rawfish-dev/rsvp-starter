package security

import (
	"strings"

	"github.com/rawfish-dev/rsvp-starter/server/services/base"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	baseService *base.Service
}

var (
	// TODO:: Store in DB/ENV so they can be changed easily
	credentials = map[string]string{
		"kevin": "$2a$10$Aq2PBEw7JeYFEvesyZBDV.Q5tyFoyLObTaMa03VHgxIfwzRRrZH0W",
		"jenny": "$2a$10$vOZqrM2Kv5TI0hIdfpWA4.dqIGYV5G7nOl7sybAWPr8gupdiwRsja",
	}
)

func NewService(baseService *base.Service) SecurityServiceProvider {
	return &service{baseService}
}

func (s *service) ValidateCredentials(username, password string) (valid bool) {
	downcasedUsername := strings.ToLower(username)

	encryptedPassword, ok := credentials[downcasedUsername]
	if !ok {
		s.baseService.Warnf("security service - unable to find username %v", downcasedUsername)
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
	if err != nil {
		s.baseService.Warnf("security service - password did not match encrypted password %v", err)
		return false
	}

	return true
}
