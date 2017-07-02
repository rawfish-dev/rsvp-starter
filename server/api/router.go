package api

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
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
		apiNameSpace.POST("/sessions", createSession(a))

		apiNameSpace.GET("/rsvps/:id", getRSVP(a))
	}

	// Initialise logger for the session service
	ctxlogger := logrus.New()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "logger", ctxlogger)
	apiNameSpace.Use(SessionMiddleware(a.SessionServiceFactory(ctx)))

	// Auth required
	{
		apiNameSpace.DELETE("/sessions", destroySession(a))

		apiNameSpace.POST("/categories", createCategory(a))
		apiNameSpace.GET("/categories", listCategories(a))
		apiNameSpace.PUT("/categories/:id", updateCategory(a))
		apiNameSpace.DELETE("/categories/:id", deleteCategory(a))

		apiNameSpace.POST("/invitations", createInvitation(a))
		apiNameSpace.GET("/invitations", listInvitations(a))
		apiNameSpace.PUT("/invitations/:id", updateInvitation(a))
		apiNameSpace.DELETE("/invitations/:id", deleteInvitation(a))

		apiNameSpace.POST("/rsvps", createRSVP(a))
		apiNameSpace.GET("/rsvps", listRSVPs(a))
		apiNameSpace.PUT("/rsvps/:id", updateRSVP(a))
		apiNameSpace.DELETE("/rsvps/:id", deleteRSVP(a))
	}
}

func (a *API) Run() {
	a.InitRoutes()

	// Begin blocking to listen for incoming requests
	a.Router.Run(fmt.Sprintf(":%v", a.HTTPPort))
}
