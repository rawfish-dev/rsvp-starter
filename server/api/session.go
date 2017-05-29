package api

import (
	"net/http"
	"strings"

	"github.com/rawfish-dev/rsvp-starter/server/config"
	"github.com/rawfish-dev/rsvp-starter/server/domain"
	"github.com/rawfish-dev/rsvp-starter/server/services/base"
	"github.com/rawfish-dev/rsvp-starter/server/services/cache"
	"github.com/rawfish-dev/rsvp-starter/server/services/jwt"
	"github.com/rawfish-dev/rsvp-starter/server/services/security"
	"github.com/rawfish-dev/rsvp-starter/server/services/session"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func createSession(c *gin.Context) {
	loadedConfig := config.LoadConfig()

	baseService := base.NewService(logrus.New())
	jwtService := jwt.NewService(baseService, loadedConfig.JWT)
	cacheService := cache.NewService(baseService)
	securityService := security.NewService(baseService)
	sessionService := session.NewService(baseService, loadedConfig.Session, jwtService, cacheService)

	var sessionCreateRequest domain.SessionCreateRequest
	err := c.BindJSON(&sessionCreateRequest)
	if err != nil {
		baseService.Errorf("session api - unable to create new session while unwrapping request due to %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	sessionCreateRequest.Username = strings.ToLower(sessionCreateRequest.Username)

	if !securityService.VerifyReCAPTCHA(sessionCreateRequest.ReCAPTCHAToken) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	valid := securityService.ValidateCredentials(sessionCreateRequest.Username, sessionCreateRequest.Password)
	if !valid {
		baseService.Warn("session api - unable to create new session due to unrecognised credentials")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	authToken, err := sessionService.CreateWithExpiry(sessionCreateRequest.Username)
	if err != nil {
		baseService.Errorf("session api - unable to create new session due to %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	sessionCreateResponse := &domain.SessionCreateResponse{
		Username:  strings.Title(sessionCreateRequest.Username),
		AuthToken: authToken,
	}

	c.JSON(http.StatusOK, sessionCreateResponse)
	return
}

func destroySession(c *gin.Context) {
	loadedConfig := config.LoadConfig()

	baseService := base.NewService(logrus.New())
	jwtService := jwt.NewService(baseService, loadedConfig.JWT)
	cacheService := cache.NewService(baseService)

	sessionService := session.NewService(baseService, loadedConfig.Session, jwtService, cacheService)

	authToken, exists := c.Get(domain.ContextAuthToken)
	if !exists {
		baseService.Error("session api - context does not contain the auth token")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if authToken == nil {
		baseService.Error("session api - context has a blank auth token")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err := sessionService.Destroy(authToken.(string))
	if err != nil {
		baseService.Errorf("session api - unable to destroy session due to %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Should return 200 OK by default
	return
}
