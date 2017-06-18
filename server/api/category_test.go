package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/rawfish-dev/rsvp-starter/server/api"
	"github.com/rawfish-dev/rsvp-starter/server/config"
	"github.com/rawfish-dev/rsvp-starter/server/domain"
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"
	"github.com/rawfish-dev/rsvp-starter/server/mock"
	. "github.com/rawfish-dev/rsvp-starter/server/services/category"
	serviceErrors "github.com/rawfish-dev/rsvp-starter/server/services/errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
)

var _ = Describe("Category", func() {

	var ctrl *gomock.Controller
	var testAPI *api.API

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		testConfig := config.LoadConfig()
		testAPI = api.NewAPI(testConfig)

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
			testAPI.CategoryServiceFactory = func(ctx context.Context) interfaces.CategoryServiceProvider {
				mockCategoryService := mock_interfaces.NewMockCategoryServiceProvider(ctrl)
				mockCategoryService.EXPECT().CreateCategory(&createCategoryReq).
					Return(&category, nil)

				return mockCategoryService
			}

			reqBytes, err := json.Marshal(createCategoryReq)
			Expect(err).ToNot(HaveOccurred())

			responseBytes := HitEndpoint(testAPI, "POST", "/api/categories", bytes.NewBuffer(reqBytes), http.StatusOK)

			var newCategory domain.Category
			err = json.Unmarshal(responseBytes, &newCategory)
			Expect(err).ToNot(HaveOccurred())

			Expect(newCategory).To(Equal(category))
		})

		It("should return 400 Bad Request when invalid JSON is passed", func() {
			testAPI.CategoryServiceFactory = func(ctx context.Context) interfaces.CategoryServiceProvider {
				mockCategoryService := mock_interfaces.NewMockCategoryServiceProvider(ctrl)
				mockCategoryService.EXPECT().CreateCategory(&createCategoryReq).Times(0)

				return mockCategoryService
			}

			reqBytes, err := json.Marshal(`{`)
			Expect(err).ToNot(HaveOccurred())

			responseBytes := HitEndpoint(testAPI, "POST", "/api/categories", bytes.NewBuffer(reqBytes), http.StatusBadRequest)

			var badRequestError domain.CustomBadRequestError
			err = json.Unmarshal(responseBytes, &badRequestError)
			Expect(err).ToNot(HaveOccurred())
			Expect(badRequestError.Error).To(Equal("JSON request was invalid"))
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

			responseBytes := HitEndpoint(testAPI, "POST", "/api/categories", bytes.NewBuffer(reqBytes), http.StatusBadRequest)

			var badRequestError domain.CustomBadRequestError
			err = json.Unmarshal(responseBytes, &badRequestError)
			Expect(err).ToNot(HaveOccurred())
			Expect(badRequestError.Error).To(Equal("some validation error"))
		})

		It("should return 500 Internal Server Error when an unknown service error occurs", func() {
			testAPI.CategoryServiceFactory = func(ctx context.Context) interfaces.CategoryServiceProvider {
				mockCategoryService := mock_interfaces.NewMockCategoryServiceProvider(ctrl)
				mockCategoryService.EXPECT().CreateCategory(&createCategoryReq).
					Return(nil, serviceErrors.NewGeneralServiceError())

				return mockCategoryService
			}

			reqBytes, err := json.Marshal(createCategoryReq)
			Expect(err).ToNot(HaveOccurred())

			HitEndpoint(testAPI, "POST", "/api/categories", bytes.NewBuffer(reqBytes), http.StatusInternalServerError)
		})
	})

	Context("retrieval", func() {

		It("should return 200 OK and the list of categories", func() {
			categories := []domain.Category{
				{
					ID:  1,
					Tag: "some-tag",
				},
				{
					ID:  2,
					Tag: "some-tag-2",
				},
				{
					ID:  3,
					Tag: "some-tag-3",
				},
			}

			testAPI.CategoryServiceFactory = func(ctx context.Context) interfaces.CategoryServiceProvider {
				mockCategoryService := mock_interfaces.NewMockCategoryServiceProvider(ctrl)
				mockCategoryService.EXPECT().ListCategories().
					Return(categories, nil)

				return mockCategoryService
			}

			responseBytes := HitEndpoint(testAPI, "GET", "/api/categories", nil, http.StatusOK)

			var categoryList []domain.Category
			err := json.Unmarshal(responseBytes, &categoryList)
			Expect(err).ToNot(HaveOccurred())

			Expect(categoryList).To(Equal(categories))
		})

		It("should return 500 Internal Server Error when an unknown service error occurs", func() {
			testAPI.CategoryServiceFactory = func(ctx context.Context) interfaces.CategoryServiceProvider {
				mockCategoryService := mock_interfaces.NewMockCategoryServiceProvider(ctrl)
				mockCategoryService.EXPECT().ListCategories().
					Return(nil, serviceErrors.NewGeneralServiceError())

				return mockCategoryService
			}

			HitEndpoint(testAPI, "GET", "/api/categories", nil, http.StatusInternalServerError)
		})
	})

	Context("updating", func() {

		var updateCategoryReq domain.CategoryUpdateRequest

		BeforeEach(func() {
			updateCategoryReq = domain.CategoryUpdateRequest{
				ID:  1,
				Tag: "some updated tag",
			}
		})

		It("should return 200 OK and update a category given valid values", func() {
			updatedCategory := &domain.Category{
				ID:  updateCategoryReq.ID,
				Tag: updateCategoryReq.Tag,
			}

			testAPI.CategoryServiceFactory = func(ctx context.Context) interfaces.CategoryServiceProvider {
				mockCategoryService := mock_interfaces.NewMockCategoryServiceProvider(ctrl)
				mockCategoryService.EXPECT().UpdateCategory(&updateCategoryReq).
					Return(updatedCategory, nil)

				return mockCategoryService
			}

			reqBytes, err := json.Marshal(updateCategoryReq)
			Expect(err).ToNot(HaveOccurred())

			responseBytes := HitEndpoint(testAPI, "PUT", "/api/categories/1", bytes.NewBuffer(reqBytes), http.StatusOK)

			var category domain.Category
			err = json.Unmarshal(responseBytes, &category)
			Expect(err).ToNot(HaveOccurred())

			Expect(category).To(Equal(*updatedCategory))
		})

		It("should return 400 Bad Request if the id in the URL does not match the update req ID", func() {
			testAPI.CategoryServiceFactory = func(ctx context.Context) interfaces.CategoryServiceProvider {
				mockCategoryService := mock_interfaces.NewMockCategoryServiceProvider(ctrl)
				mockCategoryService.EXPECT().UpdateCategory(&updateCategoryReq).Times(0)

				return mockCategoryService
			}

			reqBytes, err := json.Marshal(updateCategoryReq)
			Expect(err).ToNot(HaveOccurred())

			HitEndpoint(testAPI, "PUT", "/api/categories/2", bytes.NewBuffer(reqBytes), http.StatusBadRequest)
		})

		It("should return 400 Bad Request when invalid JSON is passed", func() {
			testAPI.CategoryServiceFactory = func(ctx context.Context) interfaces.CategoryServiceProvider {
				mockCategoryService := mock_interfaces.NewMockCategoryServiceProvider(ctrl)
				mockCategoryService.EXPECT().UpdateCategory(&updateCategoryReq).Times(0)

				return mockCategoryService
			}

			reqBytes, err := json.Marshal(`{`)
			Expect(err).ToNot(HaveOccurred())

			responseBytes := HitEndpoint(testAPI, "PUT", "/api/categories/1", bytes.NewBuffer(reqBytes), http.StatusBadRequest)

			var badRequestError domain.CustomBadRequestError
			err = json.Unmarshal(responseBytes, &badRequestError)
			Expect(err).ToNot(HaveOccurred())
			Expect(badRequestError.Error).To(Equal("JSON request was invalid"))
		})

		It("should return 400 Bad Request when a validation error occurs", func() {
			testAPI.CategoryServiceFactory = func(ctx context.Context) interfaces.CategoryServiceProvider {
				mockCategoryService := mock_interfaces.NewMockCategoryServiceProvider(ctrl)
				mockCategoryService.EXPECT().UpdateCategory(&updateCategoryReq).
					Return(nil, serviceErrors.NewValidationError([]string{"some validation error"}))

				return mockCategoryService
			}

			reqBytes, err := json.Marshal(updateCategoryReq)
			Expect(err).ToNot(HaveOccurred())

			responseBytes := HitEndpoint(testAPI, "PUT", "/api/categories/1", bytes.NewBuffer(reqBytes), http.StatusBadRequest)

			var badRequestError domain.CustomBadRequestError
			err = json.Unmarshal(responseBytes, &badRequestError)
			Expect(err).ToNot(HaveOccurred())
			Expect(badRequestError.Error).To(Equal("some validation error"))
		})

		It("should return 500 Internal Server Error when an unknown service error occurs", func() {
			testAPI.CategoryServiceFactory = func(ctx context.Context) interfaces.CategoryServiceProvider {
				mockCategoryService := mock_interfaces.NewMockCategoryServiceProvider(ctrl)
				mockCategoryService.EXPECT().UpdateCategory(&updateCategoryReq).
					Return(nil, serviceErrors.NewGeneralServiceError())

				return mockCategoryService
			}

			reqBytes, err := json.Marshal(updateCategoryReq)
			Expect(err).ToNot(HaveOccurred())

			HitEndpoint(testAPI, "PUT", "/api/categories/1", bytes.NewBuffer(reqBytes), http.StatusInternalServerError)
		})
	})

	Context("deletion", func() {

		It("should return 200 OK and delete a category given a valid id", func() {
			testAPI.CategoryServiceFactory = func(ctx context.Context) interfaces.CategoryServiceProvider {
				mockCategoryService := mock_interfaces.NewMockCategoryServiceProvider(ctrl)
				mockCategoryService.EXPECT().DeleteCategoryByID(int64(1)).Return(nil)

				return mockCategoryService
			}

			HitEndpoint(testAPI, "DELETE", "/api/categories/1", nil, http.StatusOK)
		})

		It("should return 400 Bad Request if the id is not valid", func() {
			testAPI.CategoryServiceFactory = func(ctx context.Context) interfaces.CategoryServiceProvider {
				mockCategoryService := mock_interfaces.NewMockCategoryServiceProvider(ctrl)
				mockCategoryService.EXPECT().DeleteCategoryByID("abc").Times(0)

				return mockCategoryService
			}

			HitEndpoint(testAPI, "DELETE", "/api/categories/abc", nil, http.StatusBadRequest)
		})

		It("should return 404 Not Found if the id cannot be found", func() {
			testAPI.CategoryServiceFactory = func(ctx context.Context) interfaces.CategoryServiceProvider {
				mockCategoryService := mock_interfaces.NewMockCategoryServiceProvider(ctrl)
				mockCategoryService.EXPECT().DeleteCategoryByID(int64(1)).Return(NewCategoryNotFoundError())

				return mockCategoryService
			}

			HitEndpoint(testAPI, "DELETE", "/api/categories/1", nil, http.StatusNotFound)
		})

		It("should return 500 Internal Server Error when an unknown service error occurs", func() {
			testAPI.CategoryServiceFactory = func(ctx context.Context) interfaces.CategoryServiceProvider {
				mockCategoryService := mock_interfaces.NewMockCategoryServiceProvider(ctrl)
				mockCategoryService.EXPECT().DeleteCategoryByID(int64(1)).Return(serviceErrors.NewGeneralServiceError())

				return mockCategoryService
			}

			HitEndpoint(testAPI, "DELETE", "/api/categories/1", nil, http.StatusInternalServerError)
		})
	})
})
