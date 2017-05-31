package invitation_test

import (
	"fmt"
	"strings"

	"github.com/rawfish-dev/rsvp-starter/server/domain"
	serviceErrors "github.com/rawfish-dev/rsvp-starter/server/services/errors"
	. "github.com/rawfish-dev/rsvp-starter/server/services/invitation"
	"github.com/rawfish-dev/rsvp-starter/server/testhelpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Invitation", func() {

	var testGuestService *testhelpers.TestGuestService
	var testCategory *domain.Category
	var req *domain.InvitationCreateRequest

	BeforeEach(func() {
		testGuestService = testhelpers.NewTestGuestService()
		testCategory = testGuestService.CreateTestCategory()
		req = &domain.InvitationCreateRequest{
			BaseInvitation: domain.BaseInvitation{
				CategoryID:        testCategory.ID,
				Greeting:          "ah ma and ah gong",
				MaximumGuestCount: 2,
				Notes:             "some notes",
				MobilePhoneNumber: "91231234",
			},
		}
	})

	Context("creation", func() {

		It("should create an invitation given valid values", func() {
			newInvitation, err := testGuestService.CreateInvitation(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(newInvitation).ToNot(BeNil())
			Expect(newInvitation.ID).ToNot(BeZero())
			Expect(newInvitation.CategoryID).To(Equal(testCategory.ID))
			Expect(newInvitation.PrivateID).ToNot(BeEmpty())
			Expect(newInvitation.Greeting).To(Equal("ah ma and ah gong"))
			Expect(newInvitation.MaximumGuestCount).To(Equal(2))
			Expect(newInvitation.Status).To(BeEquivalentTo(domain.NotSent))
			Expect(newInvitation.Notes).To(Equal("some notes"))
			Expect(newInvitation.MobilePhoneNumber).To(Equal("91231234"))
		})

		It("should populate an invitation's mobile number with the default extension if blank", func() {
			req.MobilePhoneNumber = ""

			newInvitation, err := testGuestService.CreateInvitation(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(newInvitation).ToNot(BeNil())
			Expect(newInvitation.ID).ToNot(BeZero())
			Expect(newInvitation.MobilePhoneNumber).To(Equal("+65"))
		})

		It("should not allow invitations with duplicate greetings", func() {
			newInvitation, err := testGuestService.CreateInvitation(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(newInvitation).ToNot(BeNil())

			// Change all necessary fields except for greeting
			req = &domain.InvitationCreateRequest{
				BaseInvitation: domain.BaseInvitation{
					CategoryID:        testCategory.ID,
					Greeting:          "ah ma and ah gong",
					MaximumGuestCount: 2,
					Notes:             "some notes",
					MobilePhoneNumber: "91231235",
				},
			}

			duplicateInvitation, err := testGuestService.CreateInvitation(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal("greeting already exists"))
			Expect(duplicateInvitation).To(BeNil())
		})

		// Temporarily relaxed requirements for mobile phone number
		// It("should not allow invitations with duplicate mobile phone numbers", func() {
		// 	newInvitation, err := testGuestService.CreateInvitation(req)
		// 	Expect(err).ToNot(HaveOccurred())
		// 	Expect(newInvitation).ToNot(BeNil())

		// 	// Change all necessary fields except for mobile phone number
		// 	req = &domain.InvitationCreateRequest{
		// 		BaseInvitation: domain.BaseInvitation{
		// 			CategoryID:        testCategory.ID,
		// 			Greeting:          "ah ma and ah gong 2",
		// 			MaximumGuestCount: 2,
		// 			Notes:             "some notes",
		// 			MobilePhoneNumber: "91231234",
		// 		},
		// 	}

		// 	duplicateInvitation, err := testGuestService.CreateInvitation(req)
		// 	Expect(err).To(HaveOccurred())
		// 	Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
		// 	Expect(err.Error()).To(Equal("mobile phone number already exists"))
		// 	Expect(duplicateInvitation).To(BeNil())
		// })

		It("should return an error if greeting is too short", func() {
			req.Greeting = "a"

			newInvitation, err := testGuestService.CreateInvitation(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation greeting must be between %v to %v characters", GreetingMinLength, GreetingMaxLength)))
			Expect(newInvitation).To(BeNil())
		})

		It("should return an error if greeting is too long", func() {
			req.Greeting = strings.Repeat("a", GreetingMaxLength+1)

			newInvitation, err := testGuestService.CreateInvitation(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation greeting must be between %v to %v characters", GreetingMinLength, GreetingMaxLength)))
			Expect(newInvitation).To(BeNil())
		})

		It("should return an error if maximum guest count is too small", func() {
			req.MaximumGuestCount = MaximumGuestCountMin - 1

			newInvitation, err := testGuestService.CreateInvitation(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation maximum guest count must be between %v to %v", MaximumGuestCountMin, MaximumGuestCountMax)))
			Expect(newInvitation).To(BeNil())
		})

		It("should return an error if maximum guest count is too large", func() {
			req.MaximumGuestCount = MaximumGuestCountMax + 1

			newInvitation, err := testGuestService.CreateInvitation(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation maximum guest count must be between %v to %v", MaximumGuestCountMin, MaximumGuestCountMax)))
			Expect(newInvitation).To(BeNil())
		})

		It("should return an error if note is too long", func() {
			req.Notes = strings.Repeat("a", NoteMaxLength+1)

			newInvitation, err := testGuestService.CreateInvitation(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation note must be less than %v characters", NoteMaxLength)))
			Expect(newInvitation).To(BeNil())
		})

		It("should return an error if mobile phone number is too long", func() {
			req.MobilePhoneNumber = "+65 1234567890 121023" // Max is 20

			newInvitation, err := testGuestService.CreateInvitation(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation mobile phone number must be less than %v in length", MobilePhoneNumberMaxLength)))
			Expect(newInvitation).To(BeNil())
		})

		// Temporarily allowing invalid numbers to reduce front end validation
		// It("should return an error if mobile phone number is invalid", func() {
		// 	req.MobilePhoneNumber = "9824abcd@" // Only accepts numbers

		// 	newInvitation, err := testGuestService.CreateInvitation(req)
		// 	Expect(err).To(HaveOccurred())
		// 	Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
		// 	Expect(err.Error()).To(Equal(
		// 		fmt.Sprintf("invitation mobile phone number must be between %v to %v in length and contain only numbers", MobilePhoneNumberMinLength, MobilePhoneNumberMaxLength)))
		// 	Expect(newInvitation).To(BeNil())
		// })
	})

	Context("retrieval", func() {

		// TODO:: Improve, this doesn't completely test the updated at since created at and id are all in their original states
		It("should return all invitations sorted by updated at asc", func() {
			newInvitation, err := testGuestService.CreateInvitation(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(newInvitation).ToNot(BeNil())

			req.Greeting += " a"
			req.MobilePhoneNumber += "1"
			newInvitation2, err := testGuestService.CreateInvitation(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(newInvitation2).ToNot(BeNil())

			req.Greeting += " b"
			req.MobilePhoneNumber += "2"
			newInvitation3, err := testGuestService.CreateInvitation(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(newInvitation3).ToNot(BeNil())

			allInvitations, err := testGuestService.ListInvitations()
			Expect(err).ToNot(HaveOccurred())
			Expect(allInvitations).To(HaveLen(3))
			Expect(allInvitations[0].ID).To(Equal(newInvitation3.ID))
			Expect(allInvitations[1].ID).To(Equal(newInvitation2.ID))
			Expect(allInvitations[2].ID).To(Equal(newInvitation.ID))
		})

		It("should return an empty slice if no categories exist", func() {
			allInvitations, err := testGuestService.ListInvitations()
			Expect(err).ToNot(HaveOccurred())
			Expect(allInvitations).To(BeEmpty())
		})

		It("should return the status 'RA' if the guests have RSVP-ed as not attending", func() {
			newRSVP := testGuestService.CreateTestRSVP(false)
			Expect(newRSVP).ToNot(BeNil())
			Expect(newRSVP.ID).ToNot(BeZero())

			allInvitations, err := testGuestService.ListInvitations()
			Expect(err).ToNot(HaveOccurred())
			Expect(allInvitations).To(HaveLen(1))

			onlyInvitation := allInvitations[0]
			Expect(onlyInvitation.Status).To(BeEquivalentTo("RN"))
		})

		It("should return the status 'RA' if the guests have RSVP-ed as attending", func() {
			newRSVP := testGuestService.CreateTestRSVP(true)
			Expect(newRSVP).ToNot(BeNil())
			Expect(newRSVP.ID).ToNot(BeZero())

			allInvitations, err := testGuestService.ListInvitations()
			Expect(err).ToNot(HaveOccurred())
			Expect(allInvitations).To(HaveLen(1))

			onlyInvitation := allInvitations[0]
			Expect(onlyInvitation.Status).To(BeEquivalentTo("RA"))
		})
	})

	Context("updating", func() {

		var newCategory *domain.Category
		var newInvitation *domain.Invitation
		var updateReq *domain.InvitationUpdateRequest

		BeforeEach(func() {
			newCategory = testGuestService.CreateTestCategory()

			newInvitation, _ = testGuestService.CreateInvitation(req)
			Expect(newInvitation.ID).ToNot(BeZero())

			updateReq = &domain.InvitationUpdateRequest{
				BaseInvitation: domain.BaseInvitation{
					CategoryID:        newCategory.ID,
					Greeting:          "ah ma and ah gong updated",
					MaximumGuestCount: 3,
					Notes:             "some updated notes",
					MobilePhoneNumber: "91231236",
				},
				ID:     newInvitation.ID,
				Status: domain.Sent,
			}
		})

		It("should update a invitation given valid values", func() {
			updatedInvitation, err := testGuestService.UpdateInvitation(updateReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(updatedInvitation).ToNot(BeNil())
			Expect(updatedInvitation.ID).To(Equal(newInvitation.ID))
			Expect(updatedInvitation.PrivateID).To(Equal(newInvitation.PrivateID))
			Expect(updatedInvitation.CategoryID).To(Equal(newCategory.ID))
			Expect(updatedInvitation.Greeting).To(Equal("ah ma and ah gong updated"))
			Expect(updatedInvitation.MaximumGuestCount).To(Equal(3))
			Expect(updatedInvitation.Notes).To(Equal("some updated notes"))
			Expect(updatedInvitation.MobilePhoneNumber).To(Equal("91231236"))
			Expect(updatedInvitation.Status).To(BeEquivalentTo(domain.Sent))
		})

		It("should return an error if the invitation cannot be found", func() {
			updateReq.ID = 123123123

			updatedInvitation, err := testGuestService.UpdateInvitation(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(InvitationNotFoundError{}))
			Expect(updatedInvitation).To(BeNil())
		})

		It("should not allow invitations with duplicate greetings", func() {
			// Change all necessary fields except for greeting
			req = &domain.InvitationCreateRequest{
				BaseInvitation: domain.BaseInvitation{
					CategoryID:        testCategory.ID,
					Greeting:          newInvitation.Greeting,
					MaximumGuestCount: 2,
					Notes:             "some notes",
					MobilePhoneNumber: "91231235",
				},
			}

			duplicateInvitation, err := testGuestService.CreateInvitation(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal("greeting already exists"))
			Expect(duplicateInvitation).To(BeNil())
		})

		// Temporarily relaxed requirements for mobile phone number
		// It("should not allow invitations with duplicate mobile phone numbers", func() {
		// 	// Change all necessary fields except for mobile phone number
		// 	req = &domain.InvitationCreateRequest{
		// 		BaseInvitation: domain.BaseInvitation{
		// 			CategoryID:        testCategory.ID,
		// 			Greeting:          "ah ma and ah gong 2",
		// 			MaximumGuestCount: 2,
		// 			Notes:             "some notes",
		// 			MobilePhoneNumber: newInvitation.MobilePhoneNumber,
		// 		},
		// 	}

		// 	duplicateInvitation, err := testGuestService.CreateInvitation(req)
		// 	Expect(err).To(HaveOccurred())
		// 	Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
		// 	Expect(err.Error()).To(Equal("mobile phone number already exists"))
		// 	Expect(duplicateInvitation).To(BeNil())
		// })

		It("should return an error if greeting is too short", func() {
			updateReq.Greeting = "a"

			updatedInvitation, err := testGuestService.UpdateInvitation(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation greeting must be between %v to %v characters", GreetingMinLength, GreetingMaxLength)))
			Expect(updatedInvitation).To(BeNil())
		})

		It("should return an error if greeting is too long", func() {
			updateReq.Greeting = strings.Repeat("a", GreetingMaxLength+1)

			updatedInvitation, err := testGuestService.UpdateInvitation(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation greeting must be between %v to %v characters", GreetingMinLength, GreetingMaxLength)))
			Expect(updatedInvitation).To(BeNil())
		})

		It("should return an error if maximum guest count is too small", func() {
			updateReq.MaximumGuestCount = MaximumGuestCountMin - 1

			updatedInvitation, err := testGuestService.UpdateInvitation(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation maximum guest count must be between %v to %v", MaximumGuestCountMin, MaximumGuestCountMax)))
			Expect(updatedInvitation).To(BeNil())
		})

		It("should return an error if maximum guest count is too large", func() {
			updateReq.MaximumGuestCount = MaximumGuestCountMax + 1

			updatedInvitation, err := testGuestService.UpdateInvitation(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation maximum guest count must be between %v to %v", MaximumGuestCountMin, MaximumGuestCountMax)))
			Expect(updatedInvitation).To(BeNil())
		})

		It("should return an error if note is too long", func() {
			updateReq.Notes = strings.Repeat("a", NoteMaxLength+1)

			updatedInvitation, err := testGuestService.UpdateInvitation(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation note must be less than %v characters", NoteMaxLength)))
			Expect(updatedInvitation).To(BeNil())
		})

		It("should return an error if mobile phone number is too long", func() {
			updateReq.MobilePhoneNumber = "+65 1234567890 121023" // Max is 20

			updatedInvitation, err := testGuestService.UpdateInvitation(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation mobile phone number must be less than %v in length", MobilePhoneNumberMaxLength)))
			Expect(updatedInvitation).To(BeNil())
		})

		// Temporarily allowing invalid numbers to reduce front end validation
		// It("should return an error if mobile phone number is invalid", func() {
		// 	updateReq.MobilePhoneNumber = "9824abcd@" // Only accepts numbers

		// 	updatedInvitation, err := testGuestService.UpdateInvitation(updateReq)
		// 	Expect(err).To(HaveOccurred())
		// 	Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
		// 	Expect(err.Error()).To(Equal(
		// 		fmt.Sprintf("invitation mobile phone number must be between %v to %v in length and contain only numbers", MobilePhoneNumberMinLength, MobilePhoneNumberMaxLength)))
		// 	Expect(updatedInvitation).To(BeNil())
		// })

		It("should return an error if status is invalid", func() {
			updateReq.Status = domain.RSVPStatus("INVALID")

			updatedInvitation, err := testGuestService.UpdateInvitation(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal("status is invalid"))
			Expect(updatedInvitation).To(BeNil())
		})
	})

	Context("deletion", func() {

		It("should delete an invitation", func() {
			newInvitation, err := testGuestService.CreateInvitation(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(newInvitation).ToNot(BeNil())

			// Ensure that the invitation exists
			allInvitations, err := testGuestService.ListInvitations()
			Expect(err).ToNot(HaveOccurred())
			Expect(allInvitations).To(HaveLen(1))

			err = testGuestService.DeleteInvitation(newInvitation.ID)
			Expect(err).ToNot(HaveOccurred())

			// Ensure that the invitation no longer exists
			allInvitations, err = testGuestService.ListInvitations()
			Expect(err).ToNot(HaveOccurred())
			Expect(allInvitations).To(HaveLen(0))
		})

		It("should return an error if the invitation cannot be found", func() {
			err := testGuestService.DeleteInvitation(123123123)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(InvitationNotFoundError{}))
		})
	})
})
