package session

import (
	"github.com/rawfish-dev/rsvp-starter/server/config"
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"
	serviceErrors "github.com/rawfish-dev/rsvp-starter/server/services/errors"
	"github.com/rawfish-dev/rsvp-starter/server/services/jwt"

	"golang.org/x/net/context"
)

var _ interfaces.SessionServiceProvider = new(service)

type service struct {
	ctx           context.Context
	sessionConfig config.SessionConfig
	jwtService    interfaces.JWTServiceProvider
	cacheService  interfaces.CacheServiceProvider
}

func NewService(ctx context.Context,
	sessionConfig config.SessionConfig,
	jwtService interfaces.JWTServiceProvider,
	cacheService interfaces.CacheServiceProvider) *service {
	return &service{ctx, sessionConfig, jwtService, cacheService}
}

func (s *service) CreateWithExpiry(username string) (authToken string, err error) {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	// Store username as an additional claim
	additionalClaims := make(map[string]string)
	additionalClaims["username"] = username

	authToken, err = s.jwtService.GenerateAuthToken(additionalClaims, s.sessionConfig.Duration)
	if err != nil {
		return "", err
	}

	expiryInSeconds := int(s.sessionConfig.Duration.Seconds())
	err = s.cacheService.SetWithExpiry(username, authToken, expiryInSeconds)
	if err != nil {
		ctxLogger.Errorf("session service - unable to set auth token with expiry in cache due to %v", err)
		return "", err
	}

	return authToken, nil
}

func (s *service) IsSessionValid(authToken string) (valid bool, err error) {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	claims, err := s.jwtService.ParseToken(authToken)
	if err != nil {
		switch err.(type) {
		case jwt.JWTInvalidError:
			ctxLogger.Warn("session service - auth token was invalid")
			return false, nil
		}

		ctxLogger.Errorf("session service - unable to parse auth token due to %v", err)
		return false, serviceErrors.NewGeneralServiceError()
	}

	username, ok := claims["username"]
	if !ok {
		ctxLogger.Error("session service - could not find username claim in auth token")
		return false, nil
	}

	exists, err := s.cacheService.Exists(username.(string))
	if err != nil {
		ctxLogger.Errorf("session service - unable to check if session is active for %v due to %v", username, err)
		return false, serviceErrors.NewGeneralServiceError()
	}
	if !exists {
		return false, nil
	}

	// TODO:: Consider logins from multiple browsers causing multiple auth tokens

	return true, nil
}

func (s *service) Destroy(authToken string) (err error) {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	claims, err := s.jwtService.ParseToken(authToken)
	if err != nil {
		ctxLogger.Errorf("session service - unable to parse auth token due to %v", err)
		return
	}

	username, ok := claims["username"]
	if !ok {
		ctxLogger.Error("session service - could not find username claim in auth token")
		return
	}

	s.cacheService.Delete(username.(string))

	return nil
}
