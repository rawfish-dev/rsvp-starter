package rsvp_test

import (
	"fmt"
	"strings"

	"github.com/rawfish-dev/rsvp-starter/server/domain"
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"
	"github.com/rawfish-dev/rsvp-starter/server/mock"
	serviceErrors "github.com/rawfish-dev/rsvp-starter/server/services/errors"
	"github.com/rawfish-dev/rsvp-starter/server/services/postgres"
	. "github.com/rawfish-dev/rsvp-starter/server/services/rsvp"

	"github.com/Sirupsen/logrus"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
)

var _ = Describe("RSVP", func() {

	var ctrl *gomock.Controller
	var mockRSVPStorage *mock_interfaces.MockRSVPStorage
	var testRSVPService interfaces.RSVPServiceProvider

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		ctxlogger := logrus.New()
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", ctxlogger)

		mockRSVPStorage = mock_interfaces.NewMockRSVPStorage(ctrl)
		testRSVPService = NewService(ctx, mockRSVPStorage)
	})

	Context("creation", func() {

		var req *domain.RSVPCreateRequest

		BeforeEach(func() {
			req = &domain.RSVPCreateRequest{
				BaseRSVP: domain.BaseRSVP{
					FullName:          "mitten lin",
					Attending:         true,
					GuestCount:        1,
					SpecialDiet:       true,
					Remarks:           "some remarks",
					MobilePhoneNumber: "91234123",
				},
				InvitationPrivateID: "some-private-id",
			}
		})

		It("should create an rsvp given valid values", func() {
			mockRSVPStorage.EXPECT().InsertRSVP(req).Return(&domain.RSVP{
				BaseRSVP:            req.BaseRSVP,
				ID:                  1,
				InvitationPrivateID: req.InvitationPrivateID,
				Completed:           true,
			}, nil)

			newRSVP, err := testRSVPService.CreateRSVP(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(newRSVP).ToNot(BeNil())
			Expect(newRSVP.ID).ToNot(BeZero())
			Expect(newRSVP.InvitationPrivateID).To(Equal("some-private-id"))
			Expect(newRSVP.FullName).To(Equal("mitten lin"))
			Expect(newRSVP.Attending).To(BeTrue())
			Expect(newRSVP.GuestCount).To(Equal(1))
			Expect(newRSVP.SpecialDiet).To(BeTrue())
			Expect(newRSVP.Remarks).To(Equal("some remarks"))
			Expect(newRSVP.MobilePhoneNumber).To(Equal("91234123"))
			Expect(newRSVP.Completed).To(BeTrue())
		})

		It("should not allow rsvps with duplicate private ids", func() {
			mockRSVPStorage.EXPECT().InsertRSVP(req).Return(
				nil, postgres.NewPostgresRSVPPrivateIDUniqueConstraintError())

			duplicateRSVP, err := testRSVPService.CreateRSVP(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal("rsvp already exists for invitation"))
			Expect(duplicateRSVP).To(BeNil())
		})

		It("should return an error if full name is too short", func() {
			req.FullName = "a"

			// Validation should catch it before any attempt to storage is made
			mockRSVPStorage.EXPECT().InsertRSVP(gomock.Any()).Times(0)

			newRSVP, err := testRSVPService.CreateRSVP(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp full name must be between %v to %v characters", GreetingMinLength, GreetingMaxLength)))
			Expect(newRSVP).To(BeNil())
		})

		It("should return an error if full name is too long", func() {
			req.FullName = strings.Repeat("a", GreetingMaxLength+1)

			// Validation should catch it before any attempt to storage is made
			mockRSVPStorage.EXPECT().InsertRSVP(gomock.Any()).Times(0)

			newRSVP, err := testRSVPService.CreateRSVP(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp full name must be between %v to %v characters", GreetingMinLength, GreetingMaxLength)))
			Expect(newRSVP).To(BeNil())
		})

		It("should return an error if attending is true guest count is 0", func() {
			req.Attending = true
			req.GuestCount = 0

			// Validation should catch it before any attempt to storage is made
			mockRSVPStorage.EXPECT().InsertRSVP(gomock.Any()).Times(0)

			newRSVP, err := testRSVPService.CreateRSVP(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp guest count must be between %v to %v", MaximumGuestCountMin, MaximumGuestCountMax)))
			Expect(newRSVP).To(BeNil())
		})

		It("should return an error if guest count is above the maximum allowed", func() {
			req.Attending = true
			req.GuestCount = MaximumGuestCountMax + 1

			// Validation should catch it before any attempt to storage is made
			mockRSVPStorage.EXPECT().InsertRSVP(gomock.Any()).Times(0)

			newRSVP, err := testRSVPService.CreateRSVP(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp guest count must be between %v to %v", MaximumGuestCountMin, MaximumGuestCountMax)))
			Expect(newRSVP).To(BeNil())
		})

		It("should not return an error if remarks are empty", func() {
			req.Remarks = ""

			mockRSVPStorage.EXPECT().InsertRSVP(req).Return(&domain.RSVP{
				BaseRSVP:            req.BaseRSVP,
				ID:                  1,
				InvitationPrivateID: req.InvitationPrivateID,
				Completed:           true,
			}, nil)

			newRSVP, err := testRSVPService.CreateRSVP(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(newRSVP).ToNot(BeNil())
			Expect(newRSVP.ID).ToNot(BeZero())
			Expect(newRSVP.Remarks).To(BeEmpty())
		})

		It("should return an error if remarks are too long", func() {
			req.Remarks = strings.Repeat("a", NoteMaxLength+1)

			// Validation should catch it before any attempt to storage is made
			mockRSVPStorage.EXPECT().InsertRSVP(gomock.Any()).Times(0)

			newRSVP, err := testRSVPService.CreateRSVP(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp remarks must be less than %v characters", NoteMaxLength)))
			Expect(newRSVP).To(BeNil())
		})

		It("should return an error if the mobile phone number is too short", func() {
			req.MobilePhoneNumber = "1234567" // Min is 8

			// Validation should catch it before any attempt to storage is made
			mockRSVPStorage.EXPECT().InsertRSVP(gomock.Any()).Times(0)

			newRSVP, err := testRSVPService.CreateRSVP(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp mobile phone number must be between %v to %v in length and contain only numbers", MobilePhoneNumberMinLength, MobilePhoneNumberMaxLength)))
			Expect(newRSVP).To(BeNil())
		})

		It("should return an error if the mobile phone number is too long", func() {
			req.MobilePhoneNumber = "+65 1234567890 121023" // Max is 20

			// Validation should catch it before any attempt to storage is made
			mockRSVPStorage.EXPECT().InsertRSVP(gomock.Any()).Times(0)

			newRSVP, err := testRSVPService.CreateRSVP(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp mobile phone number must be between %v to %v in length and contain only numbers", MobilePhoneNumberMinLength, MobilePhoneNumberMaxLength)))
			Expect(newRSVP).To(BeNil())
		})

		// Temporarily allowing invalid numbers to reduce front end validation
		// It("should return an error if mobile phone number is invalid", func() {
		// 	req.MobilePhoneNumber = "9824abcd@" // Only accepts numbers

		// 	newRSVP, err := testRSVPService.CreateRSVP(req)
		// 	Expect(err).To(HaveOccurred())
		// 	Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
		// 	Expect(err.Error()).To(Equal(
		// 		fmt.Sprintf("rsvp mobile phone number must be between %v to %v in length and contain only numbers", MobilePhoneNumberMinLength, MobilePhoneNumberMaxLength)))
		// 	Expect(newRSVP).To(BeNil())
		// })
	})

	// Context("retrieval", func() {

	// 	It("should return all rsvps sorted by updated at asc", func() {
	// 		newRSVP, err := testRSVPService.CreateRSVP(req)
	// 		Expect(err).ToNot(HaveOccurred())
	// 		Expect(newRSVP).ToNot(BeNil())
	// 		Expect(newRSVP.ID).ToNot(BeZero())

	// 		newInvitation2 := testRSVPService.CreateTestInvitation()
	// 		Expect(err).ToNot(HaveOccurred())
	// 		Expect(newInvitation2).ToNot(BeNil())

	// 		req.InvitationPrivateID = newInvitation2.PrivateID

	// 		newRSVP2, err := testRSVPService.CreateRSVP(req)
	// 		Expect(err).ToNot(HaveOccurred())
	// 		Expect(newRSVP2).ToNot(BeNil())
	// 		Expect(newRSVP2.ID).ToNot(BeZero())

	// 		newInvitation3 := testRSVPService.CreateTestInvitation()
	// 		Expect(err).ToNot(HaveOccurred())
	// 		Expect(newInvitation3).ToNot(BeNil())

	// 		req.InvitationPrivateID = newInvitation3.PrivateID

	// 		newRSVP3, err := testRSVPService.CreateRSVP(req)
	// 		Expect(err).ToNot(HaveOccurred())
	// 		Expect(newRSVP3).ToNot(BeNil())
	// 		Expect(newRSVP3.ID).ToNot(BeZero())

	// 		allRSVPs, err := testRSVPService.ListRSVPs()
	// 		Expect(err).ToNot(HaveOccurred())
	// 		Expect(allRSVPs).To(HaveLen(3))
	// 		Expect(allRSVPs[0].ID).To(Equal(newRSVP3.ID))
	// 		Expect(allRSVPs[1].ID).To(Equal(newRSVP2.ID))
	// 		Expect(allRSVPs[2].ID).To(Equal(newRSVP.ID))
	// 	})

	// 	It("should return an empty slice if no rsvps exist", func() {
	// 		allRSVPs, err := testRSVPService.ListRSVPs()
	// 		Expect(err).ToNot(HaveOccurred())
	// 		Expect(allRSVPs).To(BeEmpty())
	// 	})
	// })

	Context("updating", func() {

		var baseRSVP domain.BaseRSVP
		var updateReq *domain.RSVPUpdateRequest

		BeforeEach(func() {
			baseRSVP = domain.BaseRSVP{
				FullName:          "some updated full name",
				Attending:         false,
				GuestCount:        2,
				SpecialDiet:       false,
				Remarks:           "some updated remarks",
				MobilePhoneNumber: "91234125",
			}
			updateReq = &domain.RSVPUpdateRequest{
				BaseRSVP:            baseRSVP,
				ID:                  1,
				InvitationPrivateID: "some-private-id",
			}
		})

		It("should update a rsvp given valid values", func() {
			rsvp := &domain.RSVP{
				BaseRSVP:            baseRSVP,
				ID:                  1,
				InvitationPrivateID: "some-private-id",
				Completed:           true,
			}

			modifiedRSVP := *rsvp

			modifiedRSVP.FullName = updateReq.FullName
			modifiedRSVP.Attending = updateReq.Attending
			modifiedRSVP.GuestCount = updateReq.GuestCount
			modifiedRSVP.SpecialDiet = updateReq.SpecialDiet
			modifiedRSVP.Remarks = updateReq.Remarks
			modifiedRSVP.MobilePhoneNumber = updateReq.MobilePhoneNumber

			gomock.InOrder(
				mockRSVPStorage.EXPECT().FindRSVPByID(int64(1)).Return(
					rsvp, nil),
				mockRSVPStorage.EXPECT().UpdateRSVP(&modifiedRSVP).Return(
					&modifiedRSVP, nil),
			)

			updatedRSVP, err := testRSVPService.UpdateRSVP(updateReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(updatedRSVP).ToNot(BeNil())
			Expect(updatedRSVP.ID).To(Equal(int64(1)))
			Expect(updatedRSVP.InvitationPrivateID).To(Equal("some-private-id"))
			Expect(updatedRSVP.FullName).To(Equal("some updated full name"))
			Expect(updatedRSVP.Attending).To(BeFalse())
			Expect(updatedRSVP.GuestCount).To(Equal(2))
			Expect(updatedRSVP.SpecialDiet).To(BeFalse())
			Expect(updatedRSVP.Remarks).To(Equal("some updated remarks"))
			Expect(updatedRSVP.MobilePhoneNumber).To(Equal("91234125"))
			Expect(updatedRSVP.Completed).To(BeTrue())
		})

		It("should return an error if the rsvp cannot be found", func() {
			mockRSVPStorage.EXPECT().FindRSVPByID(int64(123123123)).Return(
				nil, postgres.NewPostgresRecordNotFoundError())

			updateReq.ID = 123123123

			updatedRSVP, err := testRSVPService.UpdateRSVP(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(RSVPNotFoundError{}))
			Expect(updatedRSVP).To(BeNil())
		})

		It("should return an error if full name is too short", func() {
			// Validation should catch it before any attempt to storage is made
			mockRSVPStorage.EXPECT().InsertRSVP(gomock.Any()).Times(0)

			updateReq.FullName = "a"

			updatedRSVP, err := testRSVPService.UpdateRSVP(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp full name must be between %v to %v characters", GreetingMinLength, GreetingMaxLength)))
			Expect(updatedRSVP).To(BeNil())
		})

		It("should return an error if full name is too long", func() {
			// Validation should catch it before any attempt to storage is made
			mockRSVPStorage.EXPECT().InsertRSVP(gomock.Any()).Times(0)

			updateReq.FullName = strings.Repeat("a", GreetingMaxLength+1)

			updatedRSVP, err := testRSVPService.UpdateRSVP(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp full name must be between %v to %v characters", GreetingMinLength, GreetingMaxLength)))
			Expect(updatedRSVP).To(BeNil())
		})

		It("should return an error if attending is true guest count is 0", func() {
			// Validation should catch it before any attempt to storage is made
			mockRSVPStorage.EXPECT().InsertRSVP(gomock.Any()).Times(0)

			updateReq.Attending = true
			updateReq.GuestCount = 0

			updatedRSVP, err := testRSVPService.UpdateRSVP(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp guest count must be between %v to %v", MaximumGuestCountMin, MaximumGuestCountMax)))
			Expect(updatedRSVP).To(BeNil())
		})

		It("should return an error if guest count is above the maximum allowed", func() {
			// Validation should catch it before any attempt to storage is made
			mockRSVPStorage.EXPECT().InsertRSVP(gomock.Any()).Times(0)

			updateReq.Attending = true
			updateReq.GuestCount = MaximumGuestCountMax + 1

			updatedRSVP, err := testRSVPService.UpdateRSVP(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp guest count must be between %v to %v", MaximumGuestCountMin, MaximumGuestCountMax)))
			Expect(updatedRSVP).To(BeNil())
		})

		It("should not return an error if remarks are empty", func() {
			rsvp := &domain.RSVP{
				BaseRSVP:            baseRSVP,
				ID:                  1,
				InvitationPrivateID: "some-private-id",
				Completed:           true,
			}
			rsvp.Remarks = ""

			updateReq.Remarks = ""

			gomock.InOrder(
				mockRSVPStorage.EXPECT().FindRSVPByID(int64(1)).Return(
					rsvp, nil),
				mockRSVPStorage.EXPECT().UpdateRSVP(rsvp).Return(
					rsvp, nil),
			)

			updatedRSVP, err := testRSVPService.UpdateRSVP(updateReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(updatedRSVP).ToNot(BeNil())
			Expect(updatedRSVP.ID).ToNot(BeZero())
			Expect(updatedRSVP.Remarks).To(BeEmpty())
		})

		It("should return an error if remarks are too long", func() {
			// Validation should catch it before any attempt to storage is made
			mockRSVPStorage.EXPECT().InsertRSVP(gomock.Any()).Times(0)

			updateReq.Remarks = strings.Repeat("a", NoteMaxLength+1)

			updatedRSVP, err := testRSVPService.UpdateRSVP(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp remarks must be less than %v characters", NoteMaxLength)))
			Expect(updatedRSVP).To(BeNil())
		})

		It("should return an error if the mobile phone number is too short", func() {
			// Validation should catch it before any attempt to storage is made
			mockRSVPStorage.EXPECT().InsertRSVP(gomock.Any()).Times(0)

			updateReq.MobilePhoneNumber = "1234567" // Min is 8

			updatedRSVP, err := testRSVPService.UpdateRSVP(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp mobile phone number must be between %v to %v in length and contain only numbers", MobilePhoneNumberMinLength, MobilePhoneNumberMaxLength)))
			Expect(updatedRSVP).To(BeNil())
		})

		It("should return an error if the mobile phone number is too long", func() {
			// Validation should catch it before any attempt to storage is made
			mockRSVPStorage.EXPECT().InsertRSVP(gomock.Any()).Times(0)

			updateReq.MobilePhoneNumber = "+65 1234567890 121023" // Max is 20

			updatedRSVP, err := testRSVPService.UpdateRSVP(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp mobile phone number must be between %v to %v in length and contain only numbers", MobilePhoneNumberMinLength, MobilePhoneNumberMaxLength)))
			Expect(updatedRSVP).To(BeNil())
		})

		// Temporarily allowing invalid numbers to reduce front end validation
		// It("should return an error if mobile phone number is invalid", func() {
		// 	updateReq.MobilePhoneNumber = "9824abcd@" // Only accepts numbers

		// 	updatedRSVP, err := testRSVPService.UpdateRSVP(updateReq)
		// 	Expect(err).To(HaveOccurred())
		// 	Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
		// 	Expect(err.Error()).To(Equal(
		// 		fmt.Sprintf("rsvp mobile phone number must be between %v to %v in length and contain only numbers", MobilePhoneNumberMinLength, MobilePhoneNumberMaxLength)))
		// 	Expect(updatedRSVP).To(BeNil())
		// })
	})

	Context("deletion", func() {

		It("should delete a rsvp", func() {
			rsvp := &domain.RSVP{
				BaseRSVP: domain.BaseRSVP{
					FullName:          "some updated full name",
					Attending:         false,
					GuestCount:        2,
					SpecialDiet:       false,
					Remarks:           "some updated remarks",
					MobilePhoneNumber: "91234125",
				},
				ID:                  1,
				InvitationPrivateID: "some-private-id",
				Completed:           true,
			}

			gomock.InOrder(
				mockRSVPStorage.EXPECT().FindRSVPByID(int64(1)).Return(rsvp, nil),
				mockRSVPStorage.EXPECT().DeleteRSVP(rsvp).Return(nil),
			)

			err := testRSVPService.DeleteRSVPByID(1)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return an error if the rsvp cannot be found", func() {
			mockRSVPStorage.EXPECT().FindRSVPByID(int64(123123123)).Return(
				nil, postgres.NewPostgresRecordNotFoundError())

			err := testRSVPService.DeleteRSVPByID(123123123)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(RSVPNotFoundError{}))
		})
	})
})
