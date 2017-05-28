package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}
