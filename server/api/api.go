package api

import (
	"sync"

	"github.com/rawfish-dev/rsvp-starter/server/config"
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"
	"github.com/rawfish-dev/rsvp-starter/server/services/cache"
	"github.com/rawfish-dev/rsvp-starter/server/services/jwt"
	"github.com/rawfish-dev/rsvp-starter/server/services/postgres"
	"github.com/rawfish-dev/rsvp-starter/server/services/session"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

type API struct {
	Router   *gin.Engine
	HTTPPort int

	// Service Factories
	JWTServiceFactory        func(context.Context) interfaces.JWTServiceProvider
	CacheServiceFactory      func(context.Context) interfaces.CacheServiceProvider
	SessionServiceFactory    func(context.Context) interfaces.SessionServiceProvider
	SecurityServiceFactory   func(context.Context) interfaces.SecurityServiceProvider
	CategoryServiceFactory   func(context.Context) interfaces.CategoryServiceProvider
	InvitationServiceFactory func(context.Context) interfaces.InvitationServiceProvider
	RSVPServiceFactory       func(context.Context) interfaces.RSVPServiceProvider
	CategoryStorageFactory   func(context.Context) interfaces.CategoryStorage
	InvitationStorageFactory func(context.Context) interfaces.InvitationStorage
	RSVPStorageFactory       func(context.Context) interfaces.RSVPStorage
}

var singletonAPI *API
var once sync.Once

func NewAPI(config config.Config) *API {
	once.Do(func() {
		// Setup factories
		jwtServiceFactory := func(ctx context.Context) interfaces.JWTServiceProvider {
			return jwt.NewService(ctx, config.JWT)
		}
		cacheServiceFactory := func(ctx context.Context) interfaces.CacheServiceProvider {
			return cache.NewService(ctx)
		}
		sessionServiceFactory := func(ctx context.Context) interfaces.SessionServiceProvider {
			return session.NewService(ctx, config.Session, jwtServiceFactory(ctx), cacheServiceFactory(ctx))
		}
		categoryStorageFactory := func(ctx context.Context) interfaces.CategoryStorage {
			return postgres.NewService(ctx, config.Postgres)
		}

		singletonAPI = &API{
			Router:                 gin.New(),
			HTTPPort:               config.HTTPPort,
			JWTServiceFactory:      jwtServiceFactory,
			CacheServiceFactory:    cacheServiceFactory,
			SessionServiceFactory:  sessionServiceFactory,
			CategoryStorageFactory: categoryStorageFactory,
		}
	})

	return singletonAPI
}
