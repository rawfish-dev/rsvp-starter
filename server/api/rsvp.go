package api

import (
	"fmt"
	"net/http"
	"strconv"

	"bitbucket.org/rawfish-dev/wedding-rsvp/server/config"
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/domain"
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/services/base"
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/services/guest"
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/services/postgres"
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/services/security"
	serviceErrors "github.com/rawfish-dev/react-redux-basics/server/services/errors"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

// The handler for creating RSVPs not requiring authentication
func guestCreateRSVP(c *gin.Context) {
	loadedConfig := config.Load()

	baseService := base.NewService(logrus.New())
	securityService := security.NewService(baseService)
	postgresService := postgres.NewService(baseService, loadedConfig.Postgres)
	guestService := guest.NewService(baseService, postgresService)

	var rsvpCreateRequest domain.RSVPCreateRequest
	err := c.BindJSON(&rsvpCreateRequest)
	if err != nil {
		baseService.Errorf("rsvp api - unable to create new guest rsvp while unwrapping request due to %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	baseService.Infof("rsvp api - receiving create guest rsvp %+v", rsvpCreateRequest)

	///////////////////
	// Private RSVPs //
	///////////////////

	// For private invitations, check if the user has already RSVP-ed
	if rsvpCreateRequest.InvitationPrivateID != "" {

		privateRSVP, err := guestService.RetrievePrivateRSVP(rsvpCreateRequest.InvitationPrivateID)
		if err != nil {
			switch err.(type) {
			case guest.RSVPNotFoundError:
				// Do nothing, means guest has not RSVP-ed yet
			default:
				// All other errors should be considered failures
				baseService.Errorf("rsvp api - unable to retrieve private rsvp due to %v", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			// Proceed to RSVP
			newRSVP, err := guestService.CreateRSVP(&rsvpCreateRequest)
			if err != nil {
				switch err.(type) {
				case serviceErrors.ValidationError:
					baseService.Errorf("rsvp api - unable to create new rsvp due to validation error %v", err)
					c.JSON(domain.NewCustomBadRequestError(err.Error()))
					return
				}

				baseService.Errorf("rsvp api - unable to create new rsvp due to %v", err)
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
	newRSVP, err := guestService.CreateRSVP(&rsvpCreateRequest)
	if err != nil {
		switch err.(type) {
		case serviceErrors.ValidationError:
			baseService.Errorf("rsvp api - unable to create new rsvp due to validation error %v", err)
			c.JSON(domain.NewCustomBadRequestError(err.Error()))
			return
		}

		baseService.Errorf("rsvp api - unable to create new rsvp due to %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, newRSVP)
	return
}

// The handler for fetching their own RSVPs not requiring authentication
func guestGetRSVP(c *gin.Context) {
	loadedConfig := config.Load()

	baseService := base.NewService(logrus.New())
	postgresService := postgres.NewService(baseService, loadedConfig.Postgres)
	guestService := guest.NewService(baseService, postgresService)

	// Only private invitations can be fetched
	invitationPrivateID := c.Param("id")
	if invitationPrivateID == "" {
		baseService.Warn("rsvp api - unable to retrieve private rsvp with a blank invitation private id")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// If a RSVP record can be found, the guest has already RSVP-ed
	privateRSVP, err := guestService.RetrievePrivateRSVP(invitationPrivateID)
	if err != nil {
		switch err.(type) {
		case guest.RSVPNotFoundError:

			// In the event the RSVP cannot be found, check if the invitation exists
			invitation, err := guestService.RetrieveInvitationByPrivateID(invitationPrivateID)
			if err != nil {
				switch err.(type) {
				case guest.InvitationNotFoundError:
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
					FullName:          invitation.Greeting,
					Attending:         true,
					GuestCount:        invitation.MaximumGuestCount,
					SpecialDiet:       false,
					Remarks:           "",
					MobilePhoneNumber: invitation.MobilePhoneNumber,
				},
				InvitationPrivateID: invitation.PrivateID,
				Completed:           false,
				UpdatedAt:           invitation.UpdatedAt,
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

func createRSVP(c *gin.Context) {
	loadedConfig := config.Load()

	baseService := base.NewService(logrus.New())
	postgresService := postgres.NewService(baseService, loadedConfig.Postgres)
	guestService := guest.NewService(baseService, postgresService)

	var rsvpCreateRequest domain.RSVPCreateRequest
	err := c.BindJSON(&rsvpCreateRequest)
	if err != nil {
		baseService.Errorf("rsvp api - unable to create new rsvp while unwrapping request due to %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	newRSVP, err := guestService.CreateRSVP(&rsvpCreateRequest)
	if err != nil {
		switch err.(type) {
		case serviceErrors.ValidationError:
			baseService.Errorf("rsvp api - unable to create new rsvp due to validation error %v", err)
			c.JSON(domain.NewCustomBadRequestError(err.Error()))
			return
		}

		baseService.Errorf("rsvp api - unable to create new rsvp due to %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, newRSVP)
	return
}

func listRSVPs(c *gin.Context) {
	loadedConfig := config.Load()

	baseService := base.NewService(logrus.New())
	postgresService := postgres.NewService(baseService, loadedConfig.Postgres)
	guestService := guest.NewService(baseService, postgresService)

	allRSVPs, err := guestService.ListRSVPs()
	if err != nil {
		baseService.Errorf("rsvp api - unable to list all rsvps due to %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, allRSVPs)
	return
}

func updateRSVP(c *gin.Context) {
	loadedConfig := config.Load()

	baseService := base.NewService(logrus.New())
	postgresService := postgres.NewService(baseService, loadedConfig.Postgres)
	guestService := guest.NewService(baseService, postgresService)

	var rsvpUpdateRequest domain.RSVPUpdateRequest
	err := c.BindJSON(&rsvpUpdateRequest)
	if err != nil {
		baseService.Errorf("rsvp api - unable to update rsvp while unwrapping request due to %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if c.Param("id") != fmt.Sprintf("%v", rsvpUpdateRequest.ID) {
		baseService.Warnf("rsvp api - unable to update rsvp as params id %v don't match request id %v", c.Param("id"), rsvpUpdateRequest.ID)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	updatedRSVP, err := guestService.UpdateRSVP(&rsvpUpdateRequest)
	if err != nil {
		switch err.(type) {
		case serviceErrors.ValidationError:
			baseService.Errorf("rsvp api - unable to update rsvp due to validation error %v", err)
			c.JSON(domain.NewCustomBadRequestError(err.Error()))
			return
		}

		baseService.Errorf("rsvp api - unable to update rsvp due to %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, updatedRSVP)
	return
}

func deleteRSVP(c *gin.Context) {
	loadedConfig := config.Load()

	baseService := base.NewService(logrus.New())
	postgresService := postgres.NewService(baseService, loadedConfig.Postgres)
	guestService := guest.NewService(baseService, postgresService)

	rsvpIDStr := c.Param("id")
	rsvpID, err := strconv.ParseInt(rsvpIDStr, 10, 64)
	if err != nil {
		baseService.Warnf("rsvp api - unable to delete rsvp as params id %v could not be converted due to %v", c.Param("id"), err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = guestService.DeleteRSVP(rsvpID)
	if err != nil {
		switch err.(type) {
		case guest.RSVPNotFoundError:
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		baseService.Errorf("rsvp api - unable to delete rsvp due to %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	return
}
