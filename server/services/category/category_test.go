package category_test

import (
	"fmt"
	"strings"

	"github.com/rawfish-dev/rsvp-starter/server/domain"
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"
	"github.com/rawfish-dev/rsvp-starter/server/mock"
	"github.com/rawfish-dev/rsvp-starter/server/services/base"
	. "github.com/rawfish-dev/rsvp-starter/server/services/category"
	serviceErrors "github.com/rawfish-dev/rsvp-starter/server/services/errors"

	"github.com/Sirupsen/logrus"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Category", func() {

	var ctrl *gomock.Controller
	var mockCategoryStorage *mock_interfaces.MockCategoryStorage
	var testCategoryService interfaces.CategoryServiceProvider
	var req *domain.CategoryCreateRequest

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		testBaseService := base.NewService(logrus.New())
		mockCategoryStorage = mock_interfaces.NewMockCategoryStorage(ctrl)
		testCategoryService = NewService(testBaseService, mockCategoryStorage)
		req = &domain.CategoryCreateRequest{
			Tag: "some tag",
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("creation", func() {

		It("should create a category given valid values", func() {
			mockCategoryStorage.EXPECT().InsertCategory(req).Return(&domain.Category{
				ID:  1,
				Tag: "some tag",
			}, nil)

			newCategory, err := testCategoryService.CreateCategory(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(newCategory).ToNot(BeNil())
			Expect(newCategory.ID).To(BeNumerically("==", 1))
			Expect(newCategory.Tag).To(Equal("some tag"))
			Expect(newCategory.Total).To(Equal(0))
		})

		It("should not allow categories with duplicate tags", func() {
			newCategory, err := testCategoryService.CreateCategory(req)
			Expect(err).ToNot(HaveOccurred())

			newCategory, err = testCategoryService.CreateCategory(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal("category tag already exists"))
			Expect(newCategory).To(BeNil())
		})

		It("should return an error if the tag is too short", func() {
			req.Tag = ""

			newCategory, err := testCategoryService.CreateCategory(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("category tag must be between %v to %v characters", TagMinLength, TagMaxLength)))
			Expect(newCategory).To(BeNil())
		})

		It("should return an error if the tag is too long", func() {
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
			_, err := testCategoryService.CreateCategory(req)
			Expect(err).ToNot(HaveOccurred())

			req.Tag = "some tag 3"
			_, err = testCategoryService.CreateCategory(req)
			Expect(err).ToNot(HaveOccurred())

			req.Tag = "some tag 2"
			_, err = testCategoryService.CreateCategory(req)
			Expect(err).ToNot(HaveOccurred())

			allCategories, err := testCategoryService.ListCategories()
			Expect(err).ToNot(HaveOccurred())
			Expect(allCategories).To(HaveLen(3))
			Expect(allCategories[0].Tag).To(Equal("some tag 3"))
			Expect(allCategories[0].Total).To(Equal(0))
			Expect(allCategories[1].Tag).To(Equal("some tag 2"))
			Expect(allCategories[1].Total).To(Equal(0))
			Expect(allCategories[2].Tag).To(Equal("some tag"))
			Expect(allCategories[2].Total).To(Equal(0))
		})

		It("should return an empty slice if no categories exist", func() {
			allCategories, err := testCategoryService.ListCategories()
			Expect(err).ToNot(HaveOccurred())
			Expect(allCategories).To(BeEmpty())
		})
	})

	Context("updating", func() {

		var createdCategory *domain.Category
		var updateReq *domain.CategoryUpdateRequest

		BeforeEach(func() {
			var err error
			createdCategory, err = testCategoryService.CreateCategory(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(createdCategory).ToNot(BeNil())
			Expect(createdCategory.ID).ToNot(BeZero())

			updateReq = &domain.CategoryUpdateRequest{
				ID:  createdCategory.ID,
				Tag: "some updated tag",
			}
		})

		It("should update a tag given valid values", func() {
			updatedCategory, err := testCategoryService.UpdateCategory(updateReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(updatedCategory).ToNot(BeNil())
			Expect(updatedCategory.ID).To(Equal(createdCategory.ID))
			Expect(updatedCategory.Tag).To(Equal("some updated tag"))
			Expect(updatedCategory.Total).To(Equal(0))
		})

		It("should return an error if the category id cannot be found", func() {
			updateReq.ID = 123123123123

			updatedCategory, err := testCategoryService.UpdateCategory(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(CategoryNotFoundError{}))
			Expect(updatedCategory).To(BeNil())
		})

		It("should not allow updating categories with duplicate tags", func() {
			newCategory, err := testCategoryService.CreateCategory(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal("category tag already exists"))
			Expect(newCategory).To(BeNil())
		})

		It("should return an error while updating if the tag is too short", func() {
			updateReq.Tag = ""

			newCategory, err := testCategoryService.UpdateCategory(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("category tag must be between %v to %v characters", TagMinLength, TagMaxLength)))
			Expect(newCategory).To(BeNil())
		})

		It("should return an error while updating if the tag is too long", func() {
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

		var createdCategory *domain.Category

		BeforeEach(func() {
			var err error
			createdCategory, err = testCategoryService.CreateCategory(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(createdCategory).ToNot(BeNil())
			Expect(createdCategory.ID).ToNot(BeZero())
		})

		It("should allow deleting of a category that has no invitations linked to it", func() {
			// Ensure that the category exists
			allCategories, err := testCategoryService.ListCategories()
			Expect(err).ToNot(HaveOccurred())
			Expect(allCategories).To(HaveLen(1))
			Expect(allCategories[0].ID).To(Equal(createdCategory.ID))

			err = testCategoryService.DeleteCategory(createdCategory.ID)
			Expect(err).ToNot(HaveOccurred())

			// Ensure that the category no longer exists
			allCategories, err = testCategoryService.ListCategories()
			Expect(err).ToNot(HaveOccurred())
			Expect(allCategories).To(HaveLen(0))
		})

		XIt("should not allow deleting of a category that already has invitations linked to it", func() {

		})

		It("should return an error if the category id cannot be found", func() {
			err := testCategoryService.DeleteCategory(123123123)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(CategoryNotFoundError{}))
		})
	})
})
