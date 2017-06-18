package api

import (
	"net/http"
	"strings"

	"github.com/rawfish-dev/rsvp-starter/server/domain"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

func createSession(api *API) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctxlogger := logrus.New()
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", ctxlogger)

		securityService := api.SecurityServiceFactory(ctx)
		sessionService := api.SessionServiceFactory(ctx)

		var sessionCreateRequest domain.SessionCreateRequest
		err := c.BindJSON(&sessionCreateRequest)
		if err != nil {
			ctxlogger.Errorf("session api - unable to create new session while unwrapping request due to %v", err)
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
			ctxlogger.Warn("session api - unable to create new session due to unrecognised credentials")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authToken, err := sessionService.CreateWithExpiry(sessionCreateRequest.Username)
		if err != nil {
			ctxlogger.Errorf("session api - unable to create new session due to %v", err)
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
}

func destroySession(api *API) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctxlogger := logrus.New()
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", ctxlogger)

		sessionService := api.SessionServiceFactory(ctx)

		authToken, exists := c.Get(domain.ContextAuthToken)
		if !exists {
			ctxlogger.Error("session api - context does not contain the auth token")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if authToken == nil {
			ctxlogger.Error("session api - context has a blank auth token")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		err := sessionService.Destroy(authToken.(string))
		if err != nil {
			ctxlogger.Errorf("session api - unable to destroy session due to %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Should return 200 OK by default
		return
	}
}
