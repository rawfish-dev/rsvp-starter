package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rawfish-dev/rsvp-starter/server/config"
	"github.com/rawfish-dev/rsvp-starter/server/domain"
	"github.com/rawfish-dev/rsvp-starter/server/services/base"
	"github.com/rawfish-dev/rsvp-starter/server/services/category"
	serviceErrors "github.com/rawfish-dev/rsvp-starter/server/services/errors"
	"github.com/rawfish-dev/rsvp-starter/server/services/postgres"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

func createCategory(c *gin.Context) {
	loadedConfig := config.LoadConfig()

	baseService := base.NewService(logrus.New())
	postgresService := postgres.NewService(baseService, loadedConfig.Postgres)
	categoryService := category.NewService(baseService, postgresService)

	var categoryCreateRequest domain.CategoryCreateRequest
	err := c.BindJSON(&categoryCreateRequest)
	if err != nil {
		baseService.Errorf("category api - unable to create new category while unwrapping request due to %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	newCategory, err := categoryService.CreateCategory(&categoryCreateRequest)
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
	loadedConfig := config.LoadConfig()

	baseService := base.NewService(logrus.New())
	postgresService := postgres.NewService(baseService, loadedConfig.Postgres)
	categoryService := category.NewService(baseService, postgresService)

	allCategories, err := categoryService.ListCategories()
	if err != nil {
		baseService.Errorf("category api - unable to list all categories due to %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, allCategories)
	return
}

func updateCategory(c *gin.Context) {
	loadedConfig := config.LoadConfig()

	baseService := base.NewService(logrus.New())
	postgresService := postgres.NewService(baseService, loadedConfig.Postgres)
	categoryService := category.NewService(baseService, postgresService)

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

	updatedCategory, err := categoryService.UpdateCategory(&categoryUpdateRequest)
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
	loadedConfig := config.LoadConfig()

	baseService := base.NewService(logrus.New())
	postgresService := postgres.NewService(baseService, loadedConfig.Postgres)
	categoryService := category.NewService(baseService, postgresService)

	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		baseService.Warnf("category api - unable to delete category as params id %v could not be converted due to %v", c.Param("id"), err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = categoryService.DeleteCategoryByID(categoryID)
	if err != nil {
		switch err.(type) {
		case category.CategoryNotFoundError:
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		baseService.Errorf("category api - unable to delete category due to %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	return
}
