package jwt

import (
	"time"
)

type JWTServiceProvider interface {
	GenerateAuthToken(additionalClaims map[string]string, duration time.Duration) (authToken string, err error)
	ParseToken(token string) (claims map[string]interface{}, err error)
	IsAuthTokenValid(authToken string) (valid bool)
}
