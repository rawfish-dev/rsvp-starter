package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rawfish-dev/rsvp-starter/server/domain"
	serviceErrors "github.com/rawfish-dev/rsvp-starter/server/services/errors"
	"github.com/rawfish-dev/rsvp-starter/server/services/invitation"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

func createInvitation(api *API) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctxlogger := logrus.New()
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", ctxlogger)

		invitationService := api.InvitationServiceFactory(ctx)

		var invitationCreateRequest domain.InvitationCreateRequest
		err := c.BindJSON(&invitationCreateRequest)
		if err != nil {
			ctxlogger.Errorf("invitation api - unable to create new invitation while unwrapping request due to %v", err)
			c.JSON(domain.NewInvalidJSONBodyError())
			return
		}

		newInvitation, err := invitationService.CreateInvitation(&invitationCreateRequest)
		if err != nil {
			switch err.(type) {
			case serviceErrors.ValidationError:
				ctxlogger.Errorf("invitation api - unable to create new invitation due to validation error %v", err)
				c.JSON(domain.NewCustomBadRequestError(err.Error()))
				return
			}

			ctxlogger.Errorf("invitation api - unable to create new invitation due to %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, newInvitation)
		return
	}
}

func listInvitations(api *API) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctxlogger := logrus.New()
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", ctxlogger)

		rsvpService := api.RSVPServiceFactory(ctx)
		invitationService := api.InvitationServiceFactory(ctx)

		allRSVPs, err := rsvpService.ListRSVPs()
		if err != nil {
			ctxlogger.Errorf("invitation api - unable to retrieve all rsvps due to %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		allInvitations, err := invitationService.ListInvitations(allRSVPs)
		if err != nil {
			ctxlogger.Errorf("invitation api - unable to list all invitations due to %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, allInvitations)
		return
	}
}

func updateInvitation(api *API) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctxlogger := logrus.New()
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", ctxlogger)

		invitationService := api.InvitationServiceFactory(ctx)

		var invitationUpdateRequest domain.InvitationUpdateRequest
		err := c.BindJSON(&invitationUpdateRequest)
		if err != nil {
			ctxlogger.Errorf("invitation api - unable to update invitation while unwrapping request due to %v", err)
			c.JSON(domain.NewInvalidJSONBodyError())
			return
		}

		if c.Param("id") != fmt.Sprintf("%v", invitationUpdateRequest.ID) {
			ctxlogger.Warnf("invitation api - unable to update invitation as params id %v don't match request id %v", c.Param("id"), invitationUpdateRequest.ID)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		updatedInvitation, err := invitationService.UpdateInvitation(&invitationUpdateRequest)
		if err != nil {
			switch err.(type) {
			case serviceErrors.ValidationError:
				ctxlogger.Errorf("invitation api - unable to update invitation due to validation error %v", err)
				c.JSON(domain.NewCustomBadRequestError(err.Error()))
				return
			}

			ctxlogger.Errorf("invitation api - unable to update invitation due to %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, updatedInvitation)
		return
	}
}

func deleteInvitation(api *API) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctxlogger := logrus.New()
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", ctxlogger)

		invitationService := api.InvitationServiceFactory(ctx)

		invitationIDStr := c.Param("id")
		invitationID, err := strconv.ParseInt(invitationIDStr, 10, 64)
		if err != nil {
			ctxlogger.Warnf("invitation api - unable to delete invitation as params id %v could not be converted due to %v", c.Param("id"), err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		err = invitationService.DeleteInvitationByID(invitationID)
		if err != nil {
			switch err.(type) {
			case invitation.InvitationNotFoundError:
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			ctxlogger.Errorf("invitation api - unable to delete invitation due to %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		return
	}
}
