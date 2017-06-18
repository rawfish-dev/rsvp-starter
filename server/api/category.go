package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/rawfish-dev/rsvp-starter/server/domain"
	"github.com/rawfish-dev/rsvp-starter/server/services/category"
	serviceErrors "github.com/rawfish-dev/rsvp-starter/server/services/errors"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

func createCategory(api *API) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctxlogger := logrus.New()
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", ctxlogger)

		categoryService := api.CategoryServiceFactory(ctx)

		var categoryCreateRequest domain.CategoryCreateRequest
		err := c.BindJSON(&categoryCreateRequest)
		if err != nil {
			ctxlogger.Errorf("category api - unable to create new category while unwrapping request due to %v", err)
			c.JSON(domain.NewInvalidJSONBodyError())
			return
		}

		newCategory, err := categoryService.CreateCategory(&categoryCreateRequest)
		if err != nil {
			switch err.(type) {
			case serviceErrors.ValidationError:
				ctxlogger.Errorf("category api - unable to create new category due to validation error %v", err)
				c.JSON(domain.NewCustomBadRequestError(err.Error()))
				return
			}

			ctxlogger.Errorf("category api - unable to create new category due to %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, newCategory)
		return
	}
}

func listCategories(api *API) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctxlogger := logrus.New()
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", ctxlogger)

		categoryService := api.CategoryServiceFactory(ctx)

		allCategories, err := categoryService.ListCategories()
		if err != nil {
			ctxlogger.Errorf("category api - unable to list all categories due to %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, allCategories)
		return
	}
}

func updateCategory(api *API) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctxlogger := logrus.New()
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", ctxlogger)

		categoryService := api.CategoryServiceFactory(ctx)

		var categoryUpdateRequest domain.CategoryUpdateRequest
		err := c.BindJSON(&categoryUpdateRequest)
		if err != nil {
			ctxlogger.Errorf("category api - unable to update category while unwrapping request due to %v", err)
			c.JSON(domain.NewInvalidJSONBodyError())
			return
		}

		if c.Param("id") != fmt.Sprintf("%v", categoryUpdateRequest.ID) {
			ctxlogger.Warnf("category api - unable to update category as params id %v don't match request id %v", c.Param("id"), categoryUpdateRequest.ID)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		updatedCategory, err := categoryService.UpdateCategory(&categoryUpdateRequest)
		if err != nil {
			switch err.(type) {
			case serviceErrors.ValidationError:
				ctxlogger.Errorf("category api - unable to update category due to validation error %v", err)
				c.JSON(domain.NewCustomBadRequestError(err.Error()))
				return
			}

			ctxlogger.Errorf("category api - unable to update category due to %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, updatedCategory)
		return
	}
}

func deleteCategory(api *API) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctxlogger := logrus.New()
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", ctxlogger)

		categoryService := api.CategoryServiceFactory(ctx)

		categoryIDStr := c.Param("id")
		categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
		if err != nil {
			ctxlogger.Warnf("category api - unable to delete category as params id %v could not be converted due to %v", c.Param("id"), err)
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

			ctxlogger.Errorf("category api - unable to delete category due to %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		return
	}
}
