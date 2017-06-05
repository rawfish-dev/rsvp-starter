package api

import (
	"net/http"

	"github.com/rawfish-dev/rsvp-starter/server/domain"
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"

	"github.com/gin-gonic/gin"
)

const (
	authHeaderKey = "X-Auth-Header"
)

// SessionMiddleware rejects requests without the correct auth header value and packs it into the context if present
func SessionMiddleware(sessionService interfaces.SessionServiceProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.Request.Header.Get(authHeaderKey)

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
