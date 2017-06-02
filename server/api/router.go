package api

import (
	"fmt"
	"net/http"

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

		apiNameSpace.POST("/p_rsvps", guestCreateRSVP)
		apiNameSpace.GET("/p_rsvps/:id", guestGetRSVP)
	}

	ctx := context.Background()
	apiNameSpace.Use(SessionMiddleware(a.SessionServiceFactory(ctx)))

	// Auth required
	{
		apiNameSpace.DELETE("/sessions", destroySession(a))

		apiNameSpace.POST("/categories", createCategory)
		apiNameSpace.GET("/categories", listCategories)
		apiNameSpace.PUT("/categories/:id", updateCategory)
		apiNameSpace.DELETE("/categories/:id", deleteCategory)

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
	a.Router.Run(fmt.Sprintf(":%v", a.HTTPPort))
}
