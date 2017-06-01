package jwt

import (
	"fmt"
	"time"

	"github.com/rawfish-dev/rsvp-starter/server/config"
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"
	"github.com/rawfish-dev/rsvp-starter/server/services/base"
	serviceErrors "github.com/rawfish-dev/rsvp-starter/server/services/errors"

	gjwt "github.com/dgrijalva/jwt-go"
)

var _ interfaces.JWTServiceProvider = new(service)

type service struct {
	baseService *base.Service
	jwtConfig   config.JWTConfig
}

func NewService(baseService *base.Service, jwtConfig config.JWTConfig) *service {
	return &service{baseService, jwtConfig}
}

func (s *service) GenerateAuthToken(additionalClaims map[string]string, duration time.Duration) (authToken string, err error) {
	if duration.Seconds() <= 0 {
		s.baseService.Error("jwt service - unable to generate JWT as expiry was less than 0")
		return "", fmt.Errorf("expiry time must be more than 0, was %v", duration.Seconds())
	}

	currentTime := time.Now()
	expiryTime := currentTime.Add(duration)

	baseJWT := gjwt.New(gjwt.GetSigningMethod("HS256"))

	// Write additional claims first in case base claims are present
	for claimKey, claimValue := range additionalClaims {
		baseJWT.Claims[claimKey] = claimValue
	}

	// Set base claims
	baseJWT.Claims["iss"] = s.jwtConfig.TokenIssuer
	baseJWT.Claims["iat"] = currentTime
	baseJWT.Claims["exp"] = expiryTime.Unix()

	// Sign and get the complete encoded token as a string
	return baseJWT.SignedString([]byte(s.jwtConfig.HMACSecret))
}

func (s *service) ParseToken(token string) (claims map[string]interface{}, err error) {
	baseJWT, err := s.parseJWTString(token)

	switch err.(type) {
	case nil:
		break

	case *gjwt.ValidationError:
		validationErr := err.(*gjwt.ValidationError)

		switch validationErr.Errors {
		case gjwt.ValidationErrorMalformed:
			s.baseService.Warnf("jwt service - jwt is malformed - %v", validationErr.Errors)

		case gjwt.ValidationErrorUnverifiable:
			s.baseService.Warnf("jwt service - jwt is unverifiable - %v", validationErr.Errors)

		case gjwt.ValidationErrorSignatureInvalid:
			s.baseService.Warnf("jwt service - jwt signature is invalid - %v", validationErr.Errors)

		case gjwt.ValidationErrorExpired:
			s.baseService.Warnf("jwt service - jwt has expired - %v", validationErr.Errors)

		default:
			s.baseService.Warnf("jwt service - jwt could not be parsed - %v", validationErr.Errors)
		}

		return nil, NewJWTInvalidError()

	default:
		s.baseService.Errorf("jwt service - jwt could not be parsed - %v", err.Error())
		return nil, serviceErrors.NewGeneralServiceError()
	}

	// JWT had no errors, proceed to check issuer
	jwtClaimedIssuer := baseJWT.Claims["iss"].(string)
	if jwtClaimedIssuer != s.jwtConfig.TokenIssuer {
		s.baseService.Errorf("jwt service - jwt was valid but issued by %v instead of %v", jwtClaimedIssuer, s.jwtConfig.TokenIssuer)
		return nil, NewJWTInvalidError()
	}

	return baseJWT.Claims, nil
}

func (s *service) IsAuthTokenValid(token string) (valid bool) {
	baseJWT, err := s.parseJWTString(token)
	if err != nil {
		s.baseService.Warnf("jwt service - jwt could not be parsed due to %v", err)
		return false
	}

	return baseJWT.Valid
}

func (s *service) parseJWTString(token string) (baseJWT *gjwt.Token, err error) {
	return gjwt.Parse(token, func(baseJWT *gjwt.Token) (interface{}, error) {
		return []byte(s.jwtConfig.HMACSecret), nil
	})
}
