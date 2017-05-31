package rsvp_test

import (
	"fmt"
	"strings"

	"github.com/rawfish-dev/rsvp-starter/server/domain"
	serviceErrors "github.com/rawfish-dev/rsvp-starter/server/services/errors"
	. "github.com/rawfish-dev/rsvp-starter/server/services/rsvp"
	"github.com/rawfish-dev/rsvp-starter/server/testhelpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rsvp", func() {

	var testGuestService *testhelpers.TestGuestService
	var testInvitation *domain.Invitation
	var req *domain.RSVPCreateRequest

	BeforeEach(func() {
		testGuestService = testhelpers.NewTestGuestService()
		testInvitation = testGuestService.CreateTestInvitation()
		req = &domain.RSVPCreateRequest{
			BaseRSVP: domain.BaseRSVP{
				FullName:          testInvitation.Greeting,
				Attending:         true,
				GuestCount:        1,
				SpecialDiet:       true,
				Remarks:           "some remarks",
				MobilePhoneNumber: "91234123",
			},
			InvitationPrivateID: testInvitation.PrivateID,
		}
	})

	Context("creation", func() {

		It("should create an rsvp given valid values", func() {
			newRSVP, err := testGuestService.CreateRSVP(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(newRSVP).ToNot(BeNil())
			Expect(newRSVP.ID).ToNot(BeZero())
			Expect(newRSVP.InvitationPrivateID).To(Equal(testInvitation.PrivateID))
			Expect(newRSVP.FullName).To(Equal(testInvitation.Greeting))
			Expect(newRSVP.Attending).To(BeTrue())
			Expect(newRSVP.GuestCount).To(Equal(1))
			Expect(newRSVP.SpecialDiet).To(BeTrue())
			Expect(newRSVP.Remarks).To(Equal("some remarks"))
			Expect(newRSVP.MobilePhoneNumber).To(Equal("91234123"))
			Expect(newRSVP.Completed).To(BeTrue())
		})

		It("should not allow rsvps with duplicate private ids", func() {
			newRSVP, err := testGuestService.CreateRSVP(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(newRSVP).ToNot(BeNil())

			// Change all necessary fields except for private id
			req = &domain.RSVPCreateRequest{
				BaseRSVP: domain.BaseRSVP{
					FullName:          "mitten and friends",
					Attending:         false,
					GuestCount:        2,
					SpecialDiet:       false,
					Remarks:           "some remarks again",
					MobilePhoneNumber: "91234124",
				},
				InvitationPrivateID: testInvitation.PrivateID,
			}

			duplicateRSVP, err := testGuestService.CreateRSVP(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal("rsvp already exists for invitation"))
			Expect(duplicateRSVP).To(BeNil())
		})

		It("should return an error if full name is too short", func() {
			req.FullName = "a"

			newRSVP, err := testGuestService.CreateRSVP(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp full name must be between %v to %v characters", GreetingMinLength, GreetingMaxLength)))
			Expect(newRSVP).To(BeNil())
		})

		It("should return an error if full name is too long", func() {
			req.FullName = strings.Repeat("a", GreetingMaxLength+1)

			newRSVP, err := testGuestService.CreateRSVP(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp full name must be between %v to %v characters", GreetingMinLength, GreetingMaxLength)))
			Expect(newRSVP).To(BeNil())
		})

		It("should return an error if attending is true guest count is 0", func() {
			req.Attending = true
			req.GuestCount = 0

			newRSVP, err := testGuestService.CreateRSVP(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp guest count must be between %v to %v", MaximumGuestCountMin, MaximumGuestCountMax)))
			Expect(newRSVP).To(BeNil())
		})

		It("should return an error if guest count is above the maximum allowed", func() {
			req.Attending = true
			req.GuestCount = MaximumGuestCountMax + 1

			newRSVP, err := testGuestService.CreateRSVP(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp guest count must be between %v to %v", MaximumGuestCountMin, MaximumGuestCountMax)))
			Expect(newRSVP).To(BeNil())
		})

		It("should return not return an error if remarks are empty", func() {
			req.Remarks = ""

			newRSVP, err := testGuestService.CreateRSVP(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(newRSVP).ToNot(BeNil())
			Expect(newRSVP.ID).ToNot(BeZero())
			Expect(newRSVP.Remarks).To(BeEmpty())
		})

		It("should return an error if remarks are too long", func() {
			req.Remarks = strings.Repeat("a", NoteMaxLength+1)

			newRSVP, err := testGuestService.CreateRSVP(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp remarks must be less than %v characters", NoteMaxLength)))
			Expect(newRSVP).To(BeNil())
		})

		It("should return an error if the mobile phone number is too short", func() {
			req.MobilePhoneNumber = "1234567" // Min is 8

			newRSVP, err := testGuestService.CreateRSVP(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp mobile phone number must be between %v to %v in length and contain only numbers", MobilePhoneNumberMinLength, MobilePhoneNumberMaxLength)))
			Expect(newRSVP).To(BeNil())
		})

		It("should return an error if the mobile phone number is too long", func() {
			req.MobilePhoneNumber = "+65 1234567890 121023" // Max is 20

			newRSVP, err := testGuestService.CreateRSVP(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp mobile phone number must be between %v to %v in length and contain only numbers", MobilePhoneNumberMinLength, MobilePhoneNumberMaxLength)))
			Expect(newRSVP).To(BeNil())
		})

		// Temporarily allowing invalid numbers to reduce front end validation
		// It("should return an error if mobile phone number is invalid", func() {
		// 	req.MobilePhoneNumber = "9824abcd@" // Only accepts numbers

		// 	newRSVP, err := testGuestService.CreateRSVP(req)
		// 	Expect(err).To(HaveOccurred())
		// 	Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
		// 	Expect(err.Error()).To(Equal(
		// 		fmt.Sprintf("rsvp mobile phone number must be between %v to %v in length and contain only numbers", MobilePhoneNumberMinLength, MobilePhoneNumberMaxLength)))
		// 	Expect(newRSVP).To(BeNil())
		// })
	})

	Context("retrieval", func() {

		It("should return all rsvps sorted by updated at asc", func() {
			newRSVP, err := testGuestService.CreateRSVP(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(newRSVP).ToNot(BeNil())
			Expect(newRSVP.ID).ToNot(BeZero())

			newInvitation2 := testGuestService.CreateTestInvitation()
			Expect(err).ToNot(HaveOccurred())
			Expect(newInvitation2).ToNot(BeNil())

			req.InvitationPrivateID = newInvitation2.PrivateID

			newRSVP2, err := testGuestService.CreateRSVP(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(newRSVP2).ToNot(BeNil())
			Expect(newRSVP2.ID).ToNot(BeZero())

			newInvitation3 := testGuestService.CreateTestInvitation()
			Expect(err).ToNot(HaveOccurred())
			Expect(newInvitation3).ToNot(BeNil())

			req.InvitationPrivateID = newInvitation3.PrivateID

			newRSVP3, err := testGuestService.CreateRSVP(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(newRSVP3).ToNot(BeNil())
			Expect(newRSVP3.ID).ToNot(BeZero())

			allRSVPs, err := testGuestService.ListRSVPs()
			Expect(err).ToNot(HaveOccurred())
			Expect(allRSVPs).To(HaveLen(3))
			Expect(allRSVPs[0].ID).To(Equal(newRSVP3.ID))
			Expect(allRSVPs[1].ID).To(Equal(newRSVP2.ID))
			Expect(allRSVPs[2].ID).To(Equal(newRSVP.ID))
		})

		It("should return an empty slice if no rsvps exist", func() {
			allRSVPs, err := testGuestService.ListRSVPs()
			Expect(err).ToNot(HaveOccurred())
			Expect(allRSVPs).To(BeEmpty())
		})
	})

	Context("updating", func() {

		var newRSVP *domain.RSVP
		var updateReq *domain.RSVPUpdateRequest

		BeforeEach(func() {
			newRSVP, _ = testGuestService.CreateRSVP(req)
			Expect(newRSVP.ID).ToNot(BeZero())

			updateReq = &domain.RSVPUpdateRequest{
				BaseRSVP: domain.BaseRSVP{
					FullName:          "some updated full name",
					Attending:         false,
					GuestCount:        2,
					SpecialDiet:       false,
					Remarks:           "some updated remarks",
					MobilePhoneNumber: "91234125",
				},
				ID: newRSVP.ID,
			}
		})

		It("should update a rsvp given valid values", func() {
			updatedRSVP, err := testGuestService.UpdateRSVP(updateReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(updatedRSVP).ToNot(BeNil())
			Expect(updatedRSVP.ID).To(Equal(newRSVP.ID))
			Expect(updatedRSVP.InvitationPrivateID).To(Equal(newRSVP.InvitationPrivateID))
			Expect(updatedRSVP.FullName).To(Equal("some updated full name"))
			Expect(updatedRSVP.Attending).To(BeFalse())
			Expect(updatedRSVP.GuestCount).To(Equal(2))
			Expect(updatedRSVP.SpecialDiet).To(BeFalse())
			Expect(updatedRSVP.Remarks).To(Equal("some updated remarks"))
			Expect(updatedRSVP.MobilePhoneNumber).To(Equal("91234125"))
			Expect(updatedRSVP.Completed).To(BeTrue())
		})

		It("should return an error if the rsvp cannot be found", func() {
			updateReq.ID = 123123123

			updatedRSVP, err := testGuestService.UpdateRSVP(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(RSVPNotFoundError{}))
			Expect(updatedRSVP).To(BeNil())
		})

		It("should return an error if full name is too short", func() {
			updateReq.FullName = "a"

			updatedRSVP, err := testGuestService.UpdateRSVP(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp full name must be between %v to %v characters", GreetingMinLength, GreetingMaxLength)))
			Expect(updatedRSVP).To(BeNil())
		})

		It("should return an error if full name is too long", func() {
			updateReq.FullName = strings.Repeat("a", GreetingMaxLength+1)

			updatedRSVP, err := testGuestService.UpdateRSVP(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp full name must be between %v to %v characters", GreetingMinLength, GreetingMaxLength)))
			Expect(updatedRSVP).To(BeNil())
		})

		It("should return an error if attending is true guest count is 0", func() {
			updateReq.Attending = true
			updateReq.GuestCount = 0

			updatedRSVP, err := testGuestService.UpdateRSVP(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp guest count must be between %v to %v", MaximumGuestCountMin, MaximumGuestCountMax)))
			Expect(updatedRSVP).To(BeNil())
		})

		It("should return an error if guest count is above the maximum allowed", func() {
			updateReq.Attending = true
			updateReq.GuestCount = MaximumGuestCountMax + 1

			updatedRSVP, err := testGuestService.UpdateRSVP(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp guest count must be between %v to %v", MaximumGuestCountMin, MaximumGuestCountMax)))
			Expect(updatedRSVP).To(BeNil())
		})

		It("should return not return an error if remarks are empty", func() {
			updateReq.Remarks = ""

			updatedRSVP, err := testGuestService.UpdateRSVP(updateReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(updatedRSVP).ToNot(BeNil())
			Expect(updatedRSVP.ID).ToNot(BeZero())
			Expect(updatedRSVP.Remarks).To(BeEmpty())
		})

		It("should return an error if remarks are too long", func() {
			updateReq.Remarks = strings.Repeat("a", NoteMaxLength+1)

			updatedRSVP, err := testGuestService.UpdateRSVP(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp remarks must be less than %v characters", NoteMaxLength)))
			Expect(updatedRSVP).To(BeNil())
		})

		It("should return an error if the mobile phone number is too short", func() {
			updateReq.MobilePhoneNumber = "1234567" // Min is 8

			updatedRSVP, err := testGuestService.UpdateRSVP(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp mobile phone number must be between %v to %v in length and contain only numbers", MobilePhoneNumberMinLength, MobilePhoneNumberMaxLength)))
			Expect(updatedRSVP).To(BeNil())
		})

		It("should return an error if the mobile phone number is too long", func() {
			updateReq.MobilePhoneNumber = "+65 1234567890 121023" // Max is 20

			updatedRSVP, err := testGuestService.UpdateRSVP(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("rsvp mobile phone number must be between %v to %v in length and contain only numbers", MobilePhoneNumberMinLength, MobilePhoneNumberMaxLength)))
			Expect(updatedRSVP).To(BeNil())
		})

		// Temporarily allowing invalid numbers to reduce front end validation
		// It("should return an error if mobile phone number is invalid", func() {
		// 	updateReq.MobilePhoneNumber = "9824abcd@" // Only accepts numbers

		// 	updatedRSVP, err := testGuestService.UpdateRSVP(updateReq)
		// 	Expect(err).To(HaveOccurred())
		// 	Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
		// 	Expect(err.Error()).To(Equal(
		// 		fmt.Sprintf("rsvp mobile phone number must be between %v to %v in length and contain only numbers", MobilePhoneNumberMinLength, MobilePhoneNumberMaxLength)))
		// 	Expect(updatedRSVP).To(BeNil())
		// })
	})

	Context("deletion", func() {

		It("should delete a rsvp", func() {
			newRSVP := testGuestService.CreateTestRSVP(true)
			Expect(newRSVP.ID).ToNot(BeZero())

			// Ensure that the rsvp exists
			allRSVPs, err := testGuestService.ListRSVPs()
			Expect(err).ToNot(HaveOccurred())
			Expect(allRSVPs).To(HaveLen(1))

			err = testGuestService.DeleteRSVP(newRSVP.ID)
			Expect(err).ToNot(HaveOccurred())

			// Ensure that the rsvp no longer exists
			allRSVPs, err = testGuestService.ListRSVPs()
			Expect(err).ToNot(HaveOccurred())
			Expect(allRSVPs).To(HaveLen(0))
		})

		It("should return an error if the rsvp cannot be found", func() {
			err := testGuestService.DeleteRSVP(123123123)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(RSVPNotFoundError{}))
		})
	})
})
