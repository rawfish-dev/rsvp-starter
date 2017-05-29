package api

import (
	"sync"

	"github.com/rawfish-dev/rsvp-starter/server/config"

	"github.com/gin-gonic/gin"
)

type API struct {
	Router   *gin.Engine
	HttpPort int
}

var singletonAPI *API
var once sync.Once

func NewAPI(apiConfig config.Config) *API {

	once.Do(func() {
		singletonAPI = &API{
			Router:   gin.New(),
			HttpPort: apiConfig.HTTPPort,
		}
	})

	return singletonAPI
}
