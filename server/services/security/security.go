package security

import (
	"strings"

	"github.com/rawfish-dev/rsvp-starter/server/interfaces"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

var _ interfaces.SecurityServiceProvider

type service struct {
	ctx context.Context
}

var (
	// TODO:: Store in DB/ENV so they can be changed easily
	credentials = map[string]string{
		"kevin": "$2a$10$Aq2PBEw7JeYFEvesyZBDV.Q5tyFoyLObTaMa03VHgxIfwzRRrZH0W",
		"jenny": "$2a$10$vOZqrM2Kv5TI0hIdfpWA4.dqIGYV5G7nOl7sybAWPr8gupdiwRsja",
	}
)

func NewService(ctx context.Context) *service {
	return &service{ctx}
}

func (s *service) ValidateCredentials(username, password string) (valid bool) {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	downcasedUsername := strings.ToLower(username)

	encryptedPassword, ok := credentials[downcasedUsername]
	if !ok {
		ctxLogger.Warnf("security service - unable to find username %v", downcasedUsername)
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
	if err != nil {
		ctxLogger.Warnf("security service - password did not match encrypted password %v", err)
		return false
	}

	return true
}
