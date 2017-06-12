package category_test

import (
	"fmt"
	"strings"

	"github.com/rawfish-dev/rsvp-starter/server/domain"
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"
	"github.com/rawfish-dev/rsvp-starter/server/mock"
	. "github.com/rawfish-dev/rsvp-starter/server/services/category"
	serviceErrors "github.com/rawfish-dev/rsvp-starter/server/services/errors"
	"github.com/rawfish-dev/rsvp-starter/server/services/postgres"

	"github.com/Sirupsen/logrus"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
)

var _ = Describe("Category", func() {

	var ctrl *gomock.Controller
	var mockCategoryStorage *mock_interfaces.MockCategoryStorage
	var testCategoryService interfaces.CategoryServiceProvider

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		ctxlogger := logrus.New()
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", ctxlogger)

		mockCategoryStorage = mock_interfaces.NewMockCategoryStorage(ctrl)
		testCategoryService = NewService(ctx, mockCategoryStorage)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("creation", func() {

		var req *domain.CategoryCreateRequest

		BeforeEach(func() {
			req = &domain.CategoryCreateRequest{
				Tag: "some tag",
			}
		})

		It("should create a category given valid values", func() {
			mockCategoryStorage.EXPECT().InsertCategory(req).Return(&domain.Category{
				ID:  1,
				Tag: "some tag",
			}, nil)

			newCategory, err := testCategoryService.CreateCategory(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(newCategory).ToNot(BeNil())
			Expect(newCategory.ID).To(Equal(int64(1)))
			Expect(newCategory.Tag).To(Equal("some tag"))
			Expect(newCategory.Total).To(Equal(0))
		})

		It("should not allow categories with duplicate tags", func() {
			mockCategoryStorage.EXPECT().InsertCategory(req).
				Return(nil, postgres.NewPostgresCategoryTagUniqueConstraintError())

			newCategory, err := testCategoryService.CreateCategory(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal("category tag already exists"))
			Expect(newCategory).To(BeNil())
		})

		It("should return an error if the tag is too short", func() {
			// Validation should catch it before any attempt to storage is made
			mockCategoryStorage.EXPECT().InsertCategory(req).Times(0)

			req.Tag = ""

			newCategory, err := testCategoryService.CreateCategory(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("category tag must be between %v to %v characters", TagMinLength, TagMaxLength)))
			Expect(newCategory).To(BeNil())
		})

		It("should return an error if the tag is too long", func() {
			// Validation should catch it before any attempt to storage is made
			mockCategoryStorage.EXPECT().InsertCategory(req).Times(0)

			req.Tag = strings.Repeat("a", TagMaxLength+1)

			newCategory, err := testCategoryService.CreateCategory(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("category tag must be between %v to %v characters", TagMinLength, TagMaxLength)))
			Expect(newCategory).To(BeNil())
		})
	})

	Context("retrieval", func() {

		It("should return all categories sorted alphabetically by tags", func() {
			mockCategoryStorage.EXPECT().ListCategories().Return(
				[]domain.Category{
					{
						Tag: "some tag 3",
					},
					{
						Tag: "some tag 2",
					},
					{
						Tag: "some tag",
					},
				}, nil)

			allCategories, err := testCategoryService.ListCategories()
			Expect(err).ToNot(HaveOccurred())
			Expect(allCategories).To(HaveLen(3))
			Expect(allCategories[0].Tag).To(Equal("some tag 3"))
			Expect(allCategories[1].Tag).To(Equal("some tag 2"))
			Expect(allCategories[2].Tag).To(Equal("some tag"))
		})

		It("should return an empty slice if no categories exist", func() {
			mockCategoryStorage.EXPECT().ListCategories().Return(
				[]domain.Category{}, nil)

			allCategories, err := testCategoryService.ListCategories()
			Expect(err).ToNot(HaveOccurred())
			Expect(allCategories).To(BeEmpty())
		})
	})

	Context("updating", func() {

		var updateReq *domain.CategoryUpdateRequest

		BeforeEach(func() {
			updateReq = &domain.CategoryUpdateRequest{
				ID:  1,
				Tag: "some updated tag",
			}
		})

		It("should update a tag given valid values", func() {
			category := &domain.Category{
				ID:    1,
				Tag:   "some tag",
				Total: 0,
			}

			gomock.InOrder(
				mockCategoryStorage.EXPECT().FindCategoryByID(int64(1)).Return(
					category, nil),
				mockCategoryStorage.EXPECT().UpdateCategory(category).Return(
					&domain.Category{
						ID:    1,
						Tag:   "some updated tag",
						Total: 0,
					}, nil),
			)

			updatedCategory, err := testCategoryService.UpdateCategory(updateReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(updatedCategory).ToNot(BeNil())
			Expect(updatedCategory.ID).To(Equal(int64(1)))
			Expect(updatedCategory.Tag).To(Equal("some updated tag"))
			Expect(updatedCategory.Total).To(Equal(0))
		})

		It("should return an error if the category id cannot be found", func() {
			mockCategoryStorage.EXPECT().FindCategoryByID(int64(123123123123)).Return(
				nil, postgres.NewPostgresRecordNotFoundError())

			updateReq.ID = 123123123123

			updatedCategory, err := testCategoryService.UpdateCategory(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(CategoryNotFoundError{}))
			Expect(updatedCategory).To(BeNil())
		})

		It("should not allow updating categories with duplicate tags", func() {
			category := &domain.Category{
				ID:    1,
				Tag:   "some tag",
				Total: 0,
			}

			gomock.InOrder(
				mockCategoryStorage.EXPECT().FindCategoryByID(int64(1)).Return(
					category, nil),
				mockCategoryStorage.EXPECT().UpdateCategory(category).Return(
					nil, postgres.NewPostgresCategoryTagUniqueConstraintError()),
			)

			updatedCategory, err := testCategoryService.UpdateCategory(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal("category tag already exists"))
			Expect(updatedCategory).To(BeNil())
		})

		It("should return an error while updating if the tag is too short", func() {
			// Validation should catch it before any attempt to storage is made
			mockCategoryStorage.EXPECT().FindCategoryByID(int64(1)).Times(0)

			updateReq.Tag = ""

			newCategory, err := testCategoryService.UpdateCategory(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("category tag must be between %v to %v characters", TagMinLength, TagMaxLength)))
			Expect(newCategory).To(BeNil())
		})

		It("should return an error while updating if the tag is too long", func() {
			// Validation should catch it before any attempt to storage is made
			mockCategoryStorage.EXPECT().FindCategoryByID(int64(1)).Times(0)

			updateReq.Tag = strings.Repeat("a", TagMaxLength+1)

			newCategory, err := testCategoryService.UpdateCategory(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("category tag must be between %v to %v characters", TagMinLength, TagMaxLength)))
			Expect(newCategory).To(BeNil())
		})

		XIt("should not allow the Total field to be updated", func() {

		})
	})

	Context("deletion", func() {

		It("should allow deleting of a category that has no invitations linked to it", func() {
			category := &domain.Category{
				ID:    1,
				Tag:   "some tag",
				Total: 0,
			}

			gomock.InOrder(
				mockCategoryStorage.EXPECT().FindCategoryByID(int64(1)).Return(
					category, nil),
				mockCategoryStorage.EXPECT().DeleteCategory(category).Return(nil),
			)

			err := testCategoryService.DeleteCategoryByID(1)
			Expect(err).ToNot(HaveOccurred())
		})

		XIt("should not allow deleting of a category that already has invitations linked to it", func() {

		})

		It("should return an error if the category id cannot be found", func() {
			mockCategoryStorage.EXPECT().FindCategoryByID(int64(123123123)).Return(
				nil, postgres.NewPostgresRecordNotFoundError())

			err := testCategoryService.DeleteCategoryByID(123123123)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(CategoryNotFoundError{}))
		})
	})
})
