package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"

	. "github.com/rawfish-dev/rsvp-starter/server/api"
	"github.com/rawfish-dev/rsvp-starter/server/config"
	"github.com/rawfish-dev/rsvp-starter/server/domain"
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"
	"github.com/rawfish-dev/rsvp-starter/server/mock"
	serviceErrors "github.com/rawfish-dev/rsvp-starter/server/services/errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
)

var _ = Describe("Category", func() {

	var ctrl *gomock.Controller

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		testConfig := config.LoadConfig()
		testAPI = NewAPI(testConfig)

		testAPI.SessionServiceFactory = func(ctx context.Context) interfaces.SessionServiceProvider {
			mockSessionService := mock_interfaces.NewMockSessionServiceProvider(ctrl)
			mockSessionService.EXPECT().IsSessionValid("").Return(true, nil)

			return mockSessionService
		}

		testAPI.InitRoutes()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("creation", func() {

		var createCategoryReq domain.CategoryCreateRequest
		var category domain.Category

		BeforeEach(func() {
			createCategoryReq = domain.CategoryCreateRequest{
				Tag: "some-tag",
			}
			category = domain.Category{
				ID:  1,
				Tag: createCategoryReq.Tag,
			}
		})

		It("should return 200 OK and create a category given valid values", func() {
			testAPI.CategoryStorageFactory = func(ctx context.Context) interfaces.CategoryStorage {
				mockCategoryStorage := mock_interfaces.NewMockCategoryStorage(ctrl)
				mockCategoryStorage.EXPECT().InsertCategory(&createCategoryReq).Return(&category, nil)

				return mockCategoryStorage
			}

			reqBytes, err := json.Marshal(createCategoryReq)
			Expect(err).ToNot(HaveOccurred())

			responseBytes := HitEndpoint("POST", "/api/categories", bytes.NewBuffer(reqBytes), http.StatusOK)

			var newCategory domain.Category
			err = json.Unmarshal(responseBytes, &newCategory)
			Expect(err).ToNot(HaveOccurred())

			Expect(newCategory.ID).ToNot(Equal(int64(0)))
		})

		It("should return 400 Bad Request when a validation error occurs", func() {
			testAPI.CategoryServiceFactory = func(ctx context.Context) interfaces.CategoryServiceProvider {
				mockCategoryService := mock_interfaces.NewMockCategoryServiceProvider(ctrl)
				mockCategoryService.EXPECT().CreateCategory(&createCategoryReq).
					Return(nil, serviceErrors.NewValidationError([]string{"some validation error"}))

				return mockCategoryService
			}

			reqBytes, err := json.Marshal(createCategoryReq)
			Expect(err).ToNot(HaveOccurred())

			responseBytes := HitEndpoint("POST", "/api/categories", bytes.NewBuffer(reqBytes), http.StatusBadRequest)

			var badRequestError domain.CustomBadRequestError
			err = json.Unmarshal(responseBytes, &badRequestError)
			Expect(err).ToNot(HaveOccurred())
			Expect(badRequestError.Error).To(Equal("some validation error"))
		})

		XIt("should return 500 Internal Server Error when an unknown service error occurs", func() {

		})
	})
})
