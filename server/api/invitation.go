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
)

func createInvitation(c *gin.Context) {
	loadedConfig := config.LoadConfig()

	baseService := base.NewService(logrus.New())
	postgresService := postgres.NewService(baseService, loadedConfig.Postgres)
	invitationService := invitation.NewService(baseService, postgresService)

	var invitationCreateRequest domain.InvitationCreateRequest
	err := c.BindJSON(&invitationCreateRequest)
	if err != nil {
		baseService.Errorf("invitation api - unable to create new invitation while unwrapping request due to %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	newInvitation, err := invitationService.CreateInvitation(&invitationCreateRequest)
	if err != nil {
		switch err.(type) {
		case serviceErrors.ValidationError:
			baseService.Errorf("invitation api - unable to create new invitation due to validation error %v", err)
			c.JSON(domain.NewCustomBadRequestError(err.Error()))
			return
		}

		baseService.Errorf("invitation api - unable to create new invitation due to %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, newInvitation)
	return
}

func listInvitations(c *gin.Context) {
	loadedConfig := config.LoadConfig()

	baseService := base.NewService(logrus.New())
	postgresService := postgres.NewService(baseService, loadedConfig.Postgres)
	rsvpService := rsvp.NewService(baseService, postgresService)
	invitationService := invitation.NewService(baseService, postgresService)

	allRSVPs, err := rsvpService.ListRSVPs()
	if err != nil {
		baseService.Errorf("invitation api - unable to retrieve all rsvps due to %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	allInvitations, err := invitationService.ListInvitations(allRSVPs)
	if err != nil {
		baseService.Errorf("invitation api - unable to list all invitations due to %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, allInvitations)
	return
}

func updateInvitation(c *gin.Context) {
	loadedConfig := config.LoadConfig()

	baseService := base.NewService(logrus.New())
	postgresService := postgres.NewService(baseService, loadedConfig.Postgres)
	invitationService := invitation.NewService(baseService, postgresService)

	var invitationUpdateRequest domain.InvitationUpdateRequest
	err := c.BindJSON(&invitationUpdateRequest)
	if err != nil {
		baseService.Errorf("invitation api - unable to update invitation while unwrapping request due to %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if c.Param("id") != fmt.Sprintf("%v", invitationUpdateRequest.ID) {
		baseService.Warnf("invitation api - unable to update invitation as params id %v don't match request id %v", c.Param("id"), invitationUpdateRequest.ID)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	updatedInvitation, err := invitationService.UpdateInvitation(&invitationUpdateRequest)
	if err != nil {
		switch err.(type) {
		case serviceErrors.ValidationError:
			baseService.Errorf("invitation api - unable to update invitation due to validation error %v", err)
			c.JSON(domain.NewCustomBadRequestError(err.Error()))
			return
		}

		baseService.Errorf("invitation api - unable to update invitation due to %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, updatedInvitation)
	return
}

func deleteInvitation(c *gin.Context) {
	loadedConfig := config.LoadConfig()

	baseService := base.NewService(logrus.New())
	postgresService := postgres.NewService(baseService, loadedConfig.Postgres)
	invitationService := invitation.NewService(baseService, postgresService)

	invitationIDStr := c.Param("id")
	invitationID, err := strconv.ParseInt(invitationIDStr, 10, 64)
	if err != nil {
		baseService.Warnf("invitation api - unable to delete invitation as params id %v could not be converted due to %v", c.Param("id"), err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = invitationService.DeleteInvitation(invitationID)
	if err != nil {
		switch err.(type) {
		case invitation.InvitationNotFoundError:
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		baseService.Errorf("invitation api - unable to delete invitation due to %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	return
}
