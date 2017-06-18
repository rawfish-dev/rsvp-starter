package api

import (
	"github.com/rawfish-dev/rsvp-starter/server/config"
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"
	"github.com/rawfish-dev/rsvp-starter/server/services/cache"
	"github.com/rawfish-dev/rsvp-starter/server/services/category"
	"github.com/rawfish-dev/rsvp-starter/server/services/invitation"
	"github.com/rawfish-dev/rsvp-starter/server/services/jwt"
	"github.com/rawfish-dev/rsvp-starter/server/services/postgres"
	"github.com/rawfish-dev/rsvp-starter/server/services/rsvp"
	"github.com/rawfish-dev/rsvp-starter/server/services/security"
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

func NewAPI(config config.Config) *API {
	// Setup storage factories
	categoryStorageFactory := func(ctx context.Context) interfaces.CategoryStorage {
		return postgres.NewService(ctx, config.Postgres)
	}
	invitationStorageFactory := func(ctx context.Context) interfaces.InvitationStorage {
		return postgres.NewService(ctx, config.Postgres)
	}
	rsvpStorageFactory := func(ctx context.Context) interfaces.RSVPStorage {
		return postgres.NewService(ctx, config.Postgres)
	}

	// Setup service factories
	jwtServiceFactory := func(ctx context.Context) interfaces.JWTServiceProvider {
		return jwt.NewService(ctx, config.JWT)
	}
	cacheServiceFactory := func(ctx context.Context) interfaces.CacheServiceProvider {
		return cache.NewService(ctx)
	}
	sessionServiceFactory := func(ctx context.Context) interfaces.SessionServiceProvider {
		return session.NewService(ctx, config.Session, jwtServiceFactory(ctx), cacheServiceFactory(ctx))
	}
	securityServiceFactory := func(ctx context.Context) interfaces.SecurityServiceProvider {
		return security.NewService(ctx)
	}
	categoryServiceFactory := func(ctx context.Context) interfaces.CategoryServiceProvider {
		return category.NewService(ctx, categoryStorageFactory(ctx))
	}
	invitationServiceFactory := func(ctx context.Context) interfaces.InvitationServiceProvider {
		return invitation.NewService(ctx, invitationStorageFactory(ctx))
	}
	rsvpServiceFactory := func(ctx context.Context) interfaces.RSVPServiceProvider {
		return rsvp.NewService(ctx, rsvpStorageFactory(ctx))
	}

	return &API{
		Router:                   gin.New(),
		HTTPPort:                 config.HTTPPort,
		JWTServiceFactory:        jwtServiceFactory,
		CacheServiceFactory:      cacheServiceFactory,
		SessionServiceFactory:    sessionServiceFactory,
		SecurityServiceFactory:   securityServiceFactory,
		CategoryServiceFactory:   categoryServiceFactory,
		InvitationServiceFactory: invitationServiceFactory,
		RSVPServiceFactory:       rsvpServiceFactory,
		CategoryStorageFactory:   categoryStorageFactory,
		InvitationStorageFactory: invitationStorageFactory,
		RSVPStorageFactory:       rsvpStorageFactory,
	}
}
