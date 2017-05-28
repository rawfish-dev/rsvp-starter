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
	serviceErrors "github.com/rawfish-dev/react-redux-basics/server/services/errors"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func createCategory(c *gin.Context) {
	loadedConfig := config.Load()

	baseService := base.NewService(logrus.New())
	postgresService := postgres.NewService(baseService, loadedConfig.Postgres)
	guestService := guest.NewService(baseService, postgresService)

	var categoryCreateRequest domain.CategoryCreateRequest
	err := c.BindJSON(&categoryCreateRequest)
	if err != nil {
		baseService.Errorf("category api - unable to create new category while unwrapping request due to %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	newCategory, err := guestService.CreateCategory(&categoryCreateRequest)
	if err != nil {
		switch err.(type) {
		case serviceErrors.ValidationError:
			baseService.Errorf("category api - unable to create new category due to validation error %v", err)
			c.JSON(domain.NewCustomBadRequestError(err.Error()))
			return
		}

		baseService.Errorf("category api - unable to create new category due to %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, newCategory)
	return
}

func listCategories(c *gin.Context) {
	loadedConfig := config.Load()

	baseService := base.NewService(logrus.New())
	postgresService := postgres.NewService(baseService, loadedConfig.Postgres)
	guestService := guest.NewService(baseService, postgresService)

	allCategories, err := guestService.ListCategories()
	if err != nil {
		baseService.Errorf("category api - unable to list all categories due to %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, allCategories)
	return
}

func updateCategory(c *gin.Context) {
	loadedConfig := config.Load()

	baseService := base.NewService(logrus.New())
	postgresService := postgres.NewService(baseService, loadedConfig.Postgres)
	guestService := guest.NewService(baseService, postgresService)

	var categoryUpdateRequest domain.CategoryUpdateRequest
	err := c.BindJSON(&categoryUpdateRequest)
	if err != nil {
		baseService.Errorf("category api - unable to update category while unwrapping request due to %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if c.Param("id") != fmt.Sprintf("%v", categoryUpdateRequest.ID) {
		baseService.Warnf("category api - unable to update category as params id %v don't match request id %v", c.Param("id"), categoryUpdateRequest.ID)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	updatedCategory, err := guestService.UpdateCategory(&categoryUpdateRequest)
	if err != nil {
		switch err.(type) {
		case serviceErrors.ValidationError:
			baseService.Errorf("category api - unable to update category due to validation error %v", err)
			c.JSON(domain.NewCustomBadRequestError(err.Error()))
			return
		}

		baseService.Errorf("category api - unable to update category due to %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, updatedCategory)
	return
}

func deleteCategory(c *gin.Context) {
	loadedConfig := config.Load()

	baseService := base.NewService(logrus.New())
	postgresService := postgres.NewService(baseService, loadedConfig.Postgres)
	guestService := guest.NewService(baseService, postgresService)

	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		baseService.Warnf("category api - unable to delete category as params id %v could not be converted due to %v", c.Param("id"), err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = guestService.DeleteCategory(categoryID)
	if err != nil {
		switch err.(type) {
		case guest.CategoryNotFoundError:
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		baseService.Errorf("category api - unable to delete category due to %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	return
}
