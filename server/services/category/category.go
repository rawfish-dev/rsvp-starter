package category

import (
	"fmt"

	"github.com/rawfish-dev/rsvp-starter/server/domain"
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"
	"github.com/rawfish-dev/rsvp-starter/server/services"
	"github.com/rawfish-dev/rsvp-starter/server/services/base"
	serviceErrors "github.com/rawfish-dev/rsvp-starter/server/services/errors"
	"github.com/rawfish-dev/rsvp-starter/server/services/postgres"
	"github.com/rawfish-dev/rsvp-starter/server/utils"
)

const (
	TagMinLength = 1
	TagMaxLength = 100
)

var _ services.CategoryServiceProvider = new(service)

type service struct {
	baseService     *base.Service
	categoryStorage interfaces.CategoryStorage
}

func NewService(baseService *base.Service, categoryStorage interfaces.CategoryStorage) *service {
	return &service{
		baseService:     baseService,
		categoryStorage: categoryStorage,
	}
}

func (s *service) CreateCategory(req *domain.CategoryCreateRequest) (*domain.Category, error) {
	errorMessages := validateCategoryCreateRequest(req)
	if len(errorMessages) > 0 {
		return nil, serviceErrors.NewValidationError(errorMessages)
	}

	newCategory, err := s.categoryStorage.InsertCategory(req)
	if err != nil {
		switch err.(type) {
		case postgres.PostgresCategoryTagUniqueConstraintError:
			errorMessage := []string{"category tag already exists"}
			return nil, serviceErrors.NewValidationError(errorMessage)
		}

		return nil, serviceErrors.NewGeneralServiceError()
	}

	return newCategory, nil
}

func (s *service) ListCategories() ([]domain.Category, error) {
	categories, err := s.categoryStorage.FindAllCategories()
	if err != nil {
		s.baseService.Error("guest service - unable to list all categories")
		return nil, serviceErrors.NewGeneralServiceError()
	}

	return categories, nil
}

func (s *service) UpdateCategory(req *domain.CategoryUpdateRequest) (*domain.Category, error) {
	errorMessages := validateCategoryUpdateRequest(req)
	if len(errorMessages) > 0 {
		return nil, serviceErrors.NewValidationError(errorMessages)
	}

	category, err := s.categoryStorage.FindCategoryByID(req.ID)
	if err != nil {
		switch err.(type) {
		case postgres.PostgresRecordNotFoundError:
			return nil, NewCategoryNotFoundError()
		}

		return nil, serviceErrors.NewGeneralServiceError()
	}

	category.Tag = req.Tag

	updatedCategory, err := s.categoryStorage.UpdateCategory(category)
	if err != nil {
		switch err.(type) {
		case postgres.PostgresCategoryTagUniqueConstraintError:
			errorMessage := []string{"category tag already exists"}
			return nil, serviceErrors.NewValidationError(errorMessage)
		}

		return nil, serviceErrors.NewGeneralServiceError()
	}

	return updatedCategory, nil
}

func (s *service) DeleteCategory(categoryID int64) error {
	err := s.categoryStorage.DeleteCategoryByID(categoryID)
	if err != nil {
		switch err.(type) {
		case postgres.PostgresRecordNotFoundError:
			return NewCategoryNotFoundError()
		}

		return serviceErrors.NewGeneralServiceError()
	}

	return nil
}

func validateCategoryCreateRequest(req *domain.CategoryCreateRequest) (errorMessages []string) {
	if !utils.IsWithin(len(req.Tag), TagMinLength, TagMaxLength) {
		errorMessages = append(errorMessages, fmt.Sprintf("category tag must be between %v to %v characters", TagMinLength, TagMaxLength))
	}

	return errorMessages
}

func validateCategoryUpdateRequest(req *domain.CategoryUpdateRequest) (errorMessages []string) {
	if req.ID <= 0 {
		errorMessages = append(errorMessages, "category id is invalid")
	}

	if !utils.IsWithin(len(req.Tag), TagMinLength, TagMaxLength) {
		errorMessages = append(errorMessages, fmt.Sprintf("category tag must be between %v to %v characters", TagMinLength, TagMaxLength))
	}

	return errorMessages
}