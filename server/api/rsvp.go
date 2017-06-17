package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rawfish-dev/rsvp-starter/server/config"
	"github.com/rawfish-dev/rsvp-starter/server/domain"
	"github.com/rawfish-dev/rsvp-starter/server/services/base"
	serviceErrors "github.com/rawfish-dev/rsvp-starter/server/services/errors"
	"github.com/rawfish-dev/rsvp-starter/server/services/invitation"
	"github.com/rawfish-dev/rsvp-starter/server/services/postgres"
	"github.com/rawfish-dev/rsvp-starter/server/services/rsvp"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
)

// The handler for creating RSVPs not requiring authentication
func guestCreateRSVP(api *API) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctxlogger := logrus.New()
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", ctxlogger)

		securityService := api.SecurityServiceFactory(ctx)
		rsvpService := api.RSVPServiceFactory(ctx)

		var rsvpCreateRequest domain.RSVPCreateRequest
		err := c.BindJSON(&rsvpCreateRequest)
		if err != nil {
			ctxlogger.Errorf("rsvp api - unable to create new guest rsvp while unwrapping request due to %v", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		ctxlogger.Infof("rsvp api - receiving create guest rsvp %+v", rsvpCreateRequest)

		///////////////////
		// Private RSVPs //
		///////////////////

		// For private invitations, check if the user has already RSVP-ed
		if rsvpCreateRequest.InvitationPrivateID != "" {

			privateRSVP, err := rsvpService.RetrievePrivateRSVP(rsvpCreateRequest.InvitationPrivateID)
			if err != nil {
				switch err.(type) {
				case rsvp.RSVPNotFoundError:
					// Do nothing, means guest has not RSVP-ed yet
				default:
					// All other errors should be considered failures
					ctxlogger.Errorf("rsvp api - unable to retrieve private rsvp due to %v", err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}

				// Proceed to RSVP
				newRSVP, err := rsvpService.CreateRSVP(&rsvpCreateRequest)
				if err != nil {
					switch err.(type) {
					case serviceErrors.ValidationError:
						ctxlogger.Errorf("rsvp api - unable to create new rsvp due to validation error %v", err)
						c.JSON(domain.NewCustomBadRequestError(err.Error()))
						return
					}

					ctxlogger.Errorf("rsvp api - unable to create new rsvp due to %v", err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}

				c.JSON(http.StatusOK, newRSVP)
				return
			}

			// If a RSVP record can be found, the guest has already RSVP-ed
			c.JSON(http.StatusOK, privateRSVP)
			return
		}

		//////////////////
		// Public RSVPs //
		//////////////////

		if !securityService.VerifyReCAPTCHA(rsvpCreateRequest.ReCAPTCHAToken) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		rsvpCreateRequest.GuestCount = 1                              // Set default guest count to 1
		rsvpCreateRequest.InvitationPrivateID = uuid.NewV4().String() // Give public rsvps a private id so they can reference it later

		// Proceed to RSVP
		newRSVP, err := rsvpService.CreateRSVP(&rsvpCreateRequest)
		if err != nil {
			switch err.(type) {
			case serviceErrors.ValidationError:
				ctxlogger.Errorf("rsvp api - unable to create new rsvp due to validation error %v", err)
				c.JSON(domain.NewCustomBadRequestError(err.Error()))
				return
			}

			ctxlogger.Errorf("rsvp api - unable to create new rsvp due to %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, newRSVP)
		return
	}
}

// The handler for fetching their own RSVPs not requiring authentication
func guestGetRSVP(api *API) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctxlogger := logrus.New()
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", ctxlogger)

		loadedConfig := config.LoadConfig()

		baseService := base.NewService(logrus.New())
		postgresService := postgres.NewService(ctx, loadedConfig.Postgres)
		invitationService := invitation.NewService(baseService, postgresService)
		rsvpService := rsvp.NewService(baseService, postgresService)

		// Only private invitations can be fetched
		invitationPrivateID := c.Param("id")
		if invitationPrivateID == "" {
			baseService.Warn("rsvp api - unable to retrieve private rsvp with a blank invitation private id")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// If a RSVP record can be found, the guest has already RSVP-ed
		privateRSVP, err := rsvpService.RetrievePrivateRSVP(invitationPrivateID)
		if err != nil {
			switch err.(type) {
			case rsvp.RSVPNotFoundError:

				// In the event the RSVP cannot be found, check if the invitation exists
				retrievedInvitation, err := invitationService.RetrieveInvitationByPrivateID(invitationPrivateID)
				if err != nil {
					switch err.(type) {
					case invitation.InvitationNotFoundError:
						c.AbortWithStatus(http.StatusNotFound)
						return
					}

					baseService.Errorf("rsvp api - unable to retrieve private rsvp due to %v", err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}

				// Invitation exists but the guest has not yet RSVP-ed
				privateRSVP = &domain.RSVP{
					BaseRSVP: domain.BaseRSVP{
						FullName:          retrievedInvitation.Greeting,
						Attending:         true,
						GuestCount:        retrievedInvitation.MaximumGuestCount,
						SpecialDiet:       false,
						Remarks:           "",
						MobilePhoneNumber: retrievedInvitation.MobilePhoneNumber,
					},
					InvitationPrivateID: retrievedInvitation.PrivateID,
					Completed:           false,
					UpdatedAt:           retrievedInvitation.UpdatedAt,
				}

				c.JSON(http.StatusOK, privateRSVP)
				return
			}

			baseService.Errorf("rsvp api - unable to retrieve private rsvp due to %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Guest has already completed the RSVP
		privateRSVP.Completed = true

		c.JSON(http.StatusOK, privateRSVP)
		return
	}
}

func createRSVP(api *API) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctxlogger := logrus.New()
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", ctxlogger)

		rsvpService := api.RSVPServiceFactory(ctx)

		var rsvpCreateRequest domain.RSVPCreateRequest
		err := c.BindJSON(&rsvpCreateRequest)
		if err != nil {
			ctxlogger.Errorf("rsvp api - unable to create new rsvp while unwrapping request due to %v", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		newRSVP, err := rsvpService.CreateRSVP(&rsvpCreateRequest)
		if err != nil {
			switch err.(type) {
			case serviceErrors.ValidationError:
				ctxlogger.Errorf("rsvp api - unable to create new rsvp due to validation error %v", err)
				c.JSON(domain.NewCustomBadRequestError(err.Error()))
				return
			}

			ctxlogger.Errorf("rsvp api - unable to create new rsvp due to %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, newRSVP)
		return
	}
}

func listRSVPs(api *API) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctxlogger := logrus.New()
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", ctxlogger)

		rsvpService := api.RSVPServiceFactory(ctx)

		allRSVPs, err := rsvpService.ListRSVPs()
		if err != nil {
			ctxlogger.Errorf("rsvp api - unable to list all rsvps due to %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, allRSVPs)
		return
	}
}

func updateRSVP(api *API) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctxlogger := logrus.New()
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", ctxlogger)

		rsvpService := api.RSVPServiceFactory(ctx)

		var rsvpUpdateRequest domain.RSVPUpdateRequest
		err := c.BindJSON(&rsvpUpdateRequest)
		if err != nil {
			ctxlogger.Errorf("rsvp api - unable to update rsvp while unwrapping request due to %v", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if c.Param("id") != fmt.Sprintf("%v", rsvpUpdateRequest.ID) {
			ctxlogger.Warnf("rsvp api - unable to update rsvp as params id %v don't match request id %v", c.Param("id"), rsvpUpdateRequest.ID)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		updatedRSVP, err := rsvpService.UpdateRSVP(&rsvpUpdateRequest)
		if err != nil {
			switch err.(type) {
			case serviceErrors.ValidationError:
				ctxlogger.Errorf("rsvp api - unable to update rsvp due to validation error %v", err)
				c.JSON(domain.NewCustomBadRequestError(err.Error()))
				return
			}

			ctxlogger.Errorf("rsvp api - unable to update rsvp due to %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, updatedRSVP)
		return
	}
}

func deleteRSVP(api *API) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctxlogger := logrus.New()
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", ctxlogger)

		rsvpService := api.RSVPServiceFactory(ctx)

		rsvpIDStr := c.Param("id")
		rsvpID, err := strconv.ParseInt(rsvpIDStr, 10, 64)
		if err != nil {
			ctxlogger.Warnf("rsvp api - unable to delete rsvp as params id %v could not be converted due to %v", c.Param("id"), err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		err = rsvpService.DeleteRSVPByID(rsvpID)
		if err != nil {
			switch err.(type) {
			case rsvp.RSVPNotFoundError:
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			ctxlogger.Errorf("rsvp api - unable to delete rsvp due to %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		return
	}
}
