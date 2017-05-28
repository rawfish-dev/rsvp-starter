package api

import (
	"fmt"
	"net/http"

	"bitbucket.org/rawfish-dev/wedding-rsvp/server/config"
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/services/base"
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/services/cache"
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/services/jwt"
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/services/session"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func (a *API) InitRoutes() {

	// Define middleware
	a.Router.Use(gin.Recovery())
	a.Router.Use(gin.Logger())

	// Service static files
	a.Router.Static("/static", "./static")
	a.Router.LoadHTMLFiles("index.html")

	// Catch all unmatched routes here
	a.Router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	apiNameSpace := a.Router.Group("api")

	// No auth required
	{
		apiNameSpace.GET("/healthcheck", healthcheck)
		apiNameSpace.POST("/sessions", createSession)

		apiNameSpace.POST("/p_rsvps", guestCreateRSVP)
		apiNameSpace.GET("/p_rsvps/:id", guestGetRSVP)
	}

	loadedConfig := config.Load()

	baseService := base.NewService(logrus.New())
	jwtService := jwt.NewService(baseService, loadedConfig.JWT)
	cacheService := cache.NewService(baseService)
	sessionService := session.NewService(baseService, loadedConfig.Session, jwtService, cacheService)

	apiNameSpace.Use(SessionMiddleware(jwtService, sessionService))

	// Auth required
	{
		apiNameSpace.DELETE("/sessions", destroySession)

		apiNameSpace.POST("/categories", createCategory)
		apiNameSpace.GET("/categories", listCategories)
		apiNameSpace.PUT("/categories/:id", updateCategory)
		apiNameSpace.DELETE("/categories/:id", deleteCategory)

		apiNameSpace.POST("/send_invitation", sendInvitationSMS)
		apiNameSpace.POST("/invitations", createInvitation)
		apiNameSpace.GET("/invitations", listInvitations)
		apiNameSpace.PUT("/invitations/:id", updateInvitation)
		apiNameSpace.DELETE("/invitations/:id", deleteInvitation)

		apiNameSpace.POST("/rsvps", createRSVP)
		apiNameSpace.GET("/rsvps", listRSVPs)
		apiNameSpace.PUT("/rsvps/:id", updateRSVP)
		apiNameSpace.DELETE("/rsvps/:id", deleteRSVP)
	}
}

func (a *API) Run() {
	a.InitRoutes()

	// Begin blocking to listen for incoming requests
	a.Router.Run(fmt.Sprintf(":%v", a.HttpPort))
}
