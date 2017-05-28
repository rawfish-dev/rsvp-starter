package session

import (
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/config"
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/services/base"
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/services/cache"
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/services/jwt"
	serviceErrors "github.com/rawfish-dev/react-redux-basics/server/services/errors"
)

type service struct {
	baseService   *base.Service
	sessionConfig config.SessionConfig
	jwtService    jwt.JWTServiceProvider
	cacheService  cache.CacheServiceProvider
}

func NewService(baseService *base.Service,
	sessionConfig config.SessionConfig,
	jwtService jwt.JWTServiceProvider,
	cacheService cache.CacheServiceProvider) SessionServiceProvider {
	return &service{baseService, sessionConfig, jwtService, cacheService}
}

func (s *service) CreateWithExpiry(username string) (authToken string, err error) {
	// Store username as an additional claim
	additionalClaims := make(map[string]string)
	additionalClaims["username"] = username

	authToken, err = s.jwtService.GenerateAuthToken(additionalClaims, s.sessionConfig.DurationSeconds)
	if err != nil {
		return "", err
	}

	err = s.cacheService.SetWithExpiry(username, authToken, s.sessionConfig.DurationSeconds)
	if err != nil {
		s.baseService.Errorf("session service - unable to set auth token with expiry in cache due to %v", err)
		return "", err
	}

	return authToken, nil
}

func (s *service) IsSessionValid(authToken string) (valid bool, err error) {
	claims, err := s.jwtService.ParseToken(authToken)
	if err != nil {
		switch err.(type) {
		case jwt.JWTInvalidError:
			s.baseService.Warn("session service - auth token was invalid")
			return false, nil
		}

		s.baseService.Errorf("session service - unable to parse auth token due to %v", err)
		return false, serviceErrors.NewGeneralServiceError()
	}

	username, ok := claims["username"]
	if !ok {
		s.baseService.Error("session service - could not find username claim in auth token")
		return false, nil
	}

	exists, err := s.cacheService.Exists(username.(string))
	if err != nil {
		s.baseService.Errorf("session service - unable to check if session is active for %v due to %v", username, err)
		return false, serviceErrors.NewGeneralServiceError()
	}
	if !exists {
		return false, nil
	}

	// TODO:: Consider logins from multiple browsers causing multiple auth tokens

	return true, nil
}

func (s *service) Destroy(authToken string) (err error) {
	claims, err := s.jwtService.ParseToken(authToken)
	if err != nil {
		s.baseService.Errorf("session service - unable to parse auth token due to %v", err)
		return
	}

	username, ok := claims["username"]
	if !ok {
		s.baseService.Error("session service - could not find username claim in auth token")
		return
	}

	s.cacheService.Delete(username.(string))

	return nil
}
