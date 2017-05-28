package api

import (
	"net/http"

	"bitbucket.org/rawfish-dev/wedding-rsvp/server/domain"
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/services/jwt"
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/services/session"

	"github.com/gin-gonic/gin"
)

const (
	authHeaderKey = "X-Auth-Header"
)

// SessionMiddleware rejects requests without the correct auth header value and packs it into the context if present
func SessionMiddleware(authService jwt.JWTServiceProvider, sessionService session.SessionServiceProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the auth header is present
		authToken := c.Request.Header.Get(authHeaderKey)
		if authToken == "" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		exists, err := sessionService.IsSessionValid(authToken)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if !exists {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(domain.ContextAuthToken, authToken)

		c.Next()
	}
}
