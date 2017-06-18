package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rawfish-dev/rsvp-starter/server/domain"
	serviceErrors "github.com/rawfish-dev/rsvp-starter/server/services/errors"
	"github.com/rawfish-dev/rsvp-starter/server/services/rsvp"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

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
