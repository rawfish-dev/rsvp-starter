package testhelpers

import (
	"fmt"

	"github.com/rawfish-dev/rsvp-starter/server/domain"
	"github.com/rawfish-dev/rsvp-starter/server/services"
	"github.com/rawfish-dev/rsvp-starter/server/services/category"

	"github.com/satori/go.uuid"
)

var _ services.CategoryServiceProvider = new(TestCategoryService)

type TestCategoryService struct {
	services.CategoryServiceProvider
}

func NewTestCategoryService() *TestCategoryService {
	testBaseService := NewTestBaseService()
	testPostgresService := NewTestPostgresService()

	testCategoryService := &TestCategoryService{
		CategoryServiceProvider: category.NewService(testBaseService, testPostgresService),
	}

	return testCategoryService
}

func (t *TestCategoryService) CreateTestCategory() *domain.Category {
	req := &domain.CategoryCreateRequest{
		Tag: "some tag " + uuid.NewV4().String(),
	}

	category, err := t.CreateCategory(req)
	if err != nil {
		panic(fmt.Sprintf("test category service - failed to create test category %v", err))
	}

	return category
}
