package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rawfish-dev/rsvp-starter/server/config"
	"github.com/rawfish-dev/rsvp-starter/server/domain"
	"github.com/rawfish-dev/rsvp-starter/server/services/base"
	serviceErrors "github.com/rawfish-dev/rsvp-starter/server/services/errors"
	"github.com/rawfish-dev/rsvp-starter/server/services/guest"
	"github.com/rawfish-dev/rsvp-starter/server/services/postgres"
	"github.com/rawfish-dev/rsvp-starter/server/services/twilio"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func sendInvitationSMS(c *gin.Context) {
	loadedConfig := config.LoadConfig()

	baseService := base.NewService(logrus.New())
	postgresService := postgres.NewService(baseService, loadedConfig.Postgres)
	guestService := guest.NewService(baseService, postgresService)
	twilioService := twilio.NewService(baseService)

	var invitationSMSRequest domain.InvitationSMSRequest
	err := c.BindJSON(&invitationSMSRequest)
	if err != nil {
		baseService.Errorf("invitation api - unable to send invitation SMS while unwrapping request due to %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	invitation, err := guestService.RetrieveInvitationByPrivateID(invitationSMSRequest.PrivateID)
	if err != nil {
		switch err.(type) {
		case guest.InvitationNotFoundError:
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		baseService.Errorf("invitation api - unable to send invitation SMS due to %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if invitation.MobilePhoneNumber == "" {
		baseService.Error("invitation api - unable to send invitation SMS as there was no mobile phone number attached to the invitation")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// TODO:: Inject this somehow, kind of ugly to hardcode domain here
	fullRSVPLink := "https://jennykevinweddingbells.com/rsvp/" + invitation.PrivateID

	// Construct invitation message
	message := fmt.Sprintf(`
	Dear %v, we are tying the knot and would love to have you to attend our wedding lunch on Saturday, 18th Feb 2017!
	Please find the RSVP form and further details at %v
	- Jenny & Kevin  
	`, invitation.Greeting, fullRSVPLink)

	success, err := twilioService.SendSMS(invitation.MobilePhoneNumber, message)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !success {
		baseService.Error("invitation api - unable to successfully send invitation SMS even though there were no errors")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	return
}

func createInvitation(c *gin.Context) {
	loadedConfig := config.LoadConfig()

	baseService := base.NewService(logrus.New())
	postgresService := postgres.NewService(baseService, loadedConfig.Postgres)
	guestService := guest.NewService(baseService, postgresService)

	var invitationCreateRequest domain.InvitationCreateRequest
	err := c.BindJSON(&invitationCreateRequest)
	if err != nil {
		baseService.Errorf("invitation api - unable to create new invitation while unwrapping request due to %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	newInvitation, err := guestService.CreateInvitation(&invitationCreateRequest)
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
	guestService := guest.NewService(baseService, postgresService)

	allInvitations, err := guestService.ListInvitations()
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
	guestService := guest.NewService(baseService, postgresService)

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

	updatedInvitation, err := guestService.UpdateInvitation(&invitationUpdateRequest)
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
	guestService := guest.NewService(baseService, postgresService)

	invitationIDStr := c.Param("id")
	invitationID, err := strconv.ParseInt(invitationIDStr, 10, 64)
	if err != nil {
		baseService.Warnf("invitation api - unable to delete invitation as params id %v could not be converted due to %v", c.Param("id"), err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = guestService.DeleteInvitation(invitationID)
	if err != nil {
		switch err.(type) {
		case guest.InvitationNotFoundError:
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		baseService.Errorf("invitation api - unable to delete invitation due to %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	return
}
