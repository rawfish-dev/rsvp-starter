package invitation_test

import (
	"fmt"
	"strings"

	"github.com/rawfish-dev/rsvp-starter/server/domain"
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"
	"github.com/rawfish-dev/rsvp-starter/server/mock"
	serviceErrors "github.com/rawfish-dev/rsvp-starter/server/services/errors"
	. "github.com/rawfish-dev/rsvp-starter/server/services/invitation"
	"github.com/rawfish-dev/rsvp-starter/server/services/postgres"

	"github.com/Sirupsen/logrus"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
)

var _ = Describe("Invitation", func() {

	var ctrl *gomock.Controller
	var mockInvitationStorage *mock_interfaces.MockInvitationStorage
	var testInvitationService interfaces.InvitationServiceProvider

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		ctxlogger := logrus.New()
		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", ctxlogger)

		mockInvitationStorage = mock_interfaces.NewMockInvitationStorage(ctrl)
		testInvitationService = NewService(ctx, mockInvitationStorage)
	})

	Context("creation", func() {

		var baseInvitation domain.BaseInvitation
		var req *domain.InvitationCreateRequest

		BeforeEach(func() {
			baseInvitation = domain.BaseInvitation{
				CategoryID:        1,
				Greeting:          "ah ma and ah gong",
				MaximumGuestCount: 2,
				Notes:             "some notes",
				MobilePhoneNumber: "91231234",
			}

			req = &domain.InvitationCreateRequest{
				BaseInvitation: domain.BaseInvitation{
					CategoryID:        1,
					Greeting:          "ah ma and ah gong",
					MaximumGuestCount: 2,
					Notes:             "some notes",
					MobilePhoneNumber: "91231234",
				},
			}
		})

		It("should create an invitation given valid values", func() {
			mockInvitationStorage.EXPECT().InsertInvitation(&domain.InvitationCreateRequest{
				BaseInvitation: baseInvitation,
			}).Return(
				&domain.Invitation{
					BaseInvitation: baseInvitation,
					ID:             1,
					PrivateID:      "some-private-id",
					Status:         domain.NotSent,
				}, nil)

			newInvitation, err := testInvitationService.CreateInvitation(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(newInvitation).ToNot(BeNil())
			Expect(newInvitation.ID).ToNot(BeZero())
			Expect(newInvitation.CategoryID).To(Equal(int64(1)))
			Expect(newInvitation.PrivateID).ToNot(BeEmpty())
			Expect(newInvitation.Greeting).To(Equal("ah ma and ah gong"))
			Expect(newInvitation.MaximumGuestCount).To(Equal(2))
			Expect(newInvitation.Status).To(BeEquivalentTo(domain.NotSent))
			Expect(newInvitation.Notes).To(Equal("some notes"))
			Expect(newInvitation.MobilePhoneNumber).To(Equal("91231234"))
		})

		It("should populate an invitation's mobile number with the default extension if blank", func() {
			req.MobilePhoneNumber = ""

			baseInvitation.MobilePhoneNumber = "+65" // Default phone extension
			mockInvitationStorage.EXPECT().InsertInvitation(&domain.InvitationCreateRequest{
				BaseInvitation: baseInvitation,
			})

			testInvitationService.CreateInvitation(req)
		})

		It("should not allow invitations with duplicate greetings", func() {
			mockInvitationStorage.EXPECT().InsertInvitation(&domain.InvitationCreateRequest{
				BaseInvitation: baseInvitation,
			}).Return(nil, postgres.NewPostgresInvitationGreetingUniqueConstraintError())

			duplicateInvitation, err := testInvitationService.CreateInvitation(req)
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

			// Validation should catch it before any attempt to storage is made
			mockInvitationStorage.EXPECT().InsertInvitation(gomock.Any()).Times(0)

			newInvitation, err := testInvitationService.CreateInvitation(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation greeting must be between %v to %v characters", GreetingMinLength, GreetingMaxLength)))
			Expect(newInvitation).To(BeNil())
		})

		It("should return an error if greeting is too long", func() {
			req.Greeting = strings.Repeat("a", GreetingMaxLength+1)

			// Validation should catch it before any attempt to storage is made
			mockInvitationStorage.EXPECT().InsertInvitation(gomock.Any()).Times(0)

			newInvitation, err := testInvitationService.CreateInvitation(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation greeting must be between %v to %v characters", GreetingMinLength, GreetingMaxLength)))
			Expect(newInvitation).To(BeNil())
		})

		It("should return an error if maximum guest count is too small", func() {
			req.MaximumGuestCount = MaximumGuestCountMin - 1

			// Validation should catch it before any attempt to storage is made
			mockInvitationStorage.EXPECT().InsertInvitation(gomock.Any()).Times(0)

			newInvitation, err := testInvitationService.CreateInvitation(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation maximum guest count must be between %v to %v", MaximumGuestCountMin, MaximumGuestCountMax)))
			Expect(newInvitation).To(BeNil())
		})

		It("should return an error if maximum guest count is too large", func() {
			req.MaximumGuestCount = MaximumGuestCountMax + 1

			// Validation should catch it before any attempt to storage is made
			mockInvitationStorage.EXPECT().InsertInvitation(gomock.Any()).Times(0)

			newInvitation, err := testInvitationService.CreateInvitation(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation maximum guest count must be between %v to %v", MaximumGuestCountMin, MaximumGuestCountMax)))
			Expect(newInvitation).To(BeNil())
		})

		It("should return an error if note is too long", func() {
			req.Notes = strings.Repeat("a", NoteMaxLength+1)

			// Validation should catch it before any attempt to storage is made
			mockInvitationStorage.EXPECT().InsertInvitation(gomock.Any()).Times(0)

			newInvitation, err := testInvitationService.CreateInvitation(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation note must be less than %v characters", NoteMaxLength)))
			Expect(newInvitation).To(BeNil())
		})

		It("should return an error if mobile phone number is too long", func() {
			req.MobilePhoneNumber = "+65 1234567890 121023" // Max is 20

			// Validation should catch it before any attempt to storage is made
			mockInvitationStorage.EXPECT().InsertInvitation(gomock.Any()).Times(0)

			newInvitation, err := testInvitationService.CreateInvitation(req)
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

	// Context("retrieval", func() {

	// 	// TODO:: Improve, not a very useful test as ordering is in the postgres layer
	// 	It("should return all invitations sorted by updated at asc", func() {
	// 		mockInvitationStorage.EXPECT().ListInvitations().Return(
	// 			[]domain.Invitation{
	// 				{
	// 					ID: 3,
	// 				},
	// 				{
	// 					ID: 2,
	// 				},
	// 				{
	// 					ID: 1,
	// 				},
	// 			},
	// 		)

	// 		allInvitations, err := testInvitationService.ListInvitations()
	// 		Expect(err).ToNot(HaveOccurred())
	// 		Expect(allInvitations).To(HaveLen(3))
	// 		Expect(allInvitations[0].ID).To(Equal(newInvitation3.ID))
	// 		Expect(allInvitations[1].ID).To(Equal(newInvitation2.ID))
	// 		Expect(allInvitations[2].ID).To(Equal(newInvitation.ID))
	// 	})

	// 	It("should return an empty slice if no categories exist", func() {
	// 		allInvitations, err := testGuestService.ListInvitations()
	// 		Expect(err).ToNot(HaveOccurred())
	// 		Expect(allInvitations).To(BeEmpty())
	// 	})

	// 	It("should return the status 'RA' if the guests have RSVP-ed as not attending", func() {
	// 		newRSVP := testGuestService.CreateTestRSVP(false)
	// 		Expect(newRSVP).ToNot(BeNil())
	// 		Expect(newRSVP.ID).ToNot(BeZero())

	// 		allInvitations, err := testGuestService.ListInvitations()
	// 		Expect(err).ToNot(HaveOccurred())
	// 		Expect(allInvitations).To(HaveLen(1))

	// 		onlyInvitation := allInvitations[0]
	// 		Expect(onlyInvitation.Status).To(BeEquivalentTo("RN"))
	// 	})

	// 	It("should return the status 'RA' if the guests have RSVP-ed as attending", func() {
	// 		newRSVP := testGuestService.CreateTestRSVP(true)
	// 		Expect(newRSVP).ToNot(BeNil())
	// 		Expect(newRSVP.ID).ToNot(BeZero())

	// 		allInvitations, err := testGuestService.ListInvitations()
	// 		Expect(err).ToNot(HaveOccurred())
	// 		Expect(allInvitations).To(HaveLen(1))

	// 		onlyInvitation := allInvitations[0]
	// 		Expect(onlyInvitation.Status).To(BeEquivalentTo("RA"))
	// 	})
	// })

	Context("updating", func() {

		var baseInvitation domain.BaseInvitation
		var updateReq *domain.InvitationUpdateRequest

		BeforeEach(func() {
			baseInvitation = domain.BaseInvitation{
				CategoryID:        1,
				Greeting:          "ah ma and ah gong",
				MaximumGuestCount: 2,
				Notes:             "some notes",
				MobilePhoneNumber: "91231234",
			}

			updateReq = &domain.InvitationUpdateRequest{
				BaseInvitation: domain.BaseInvitation{
					CategoryID:        2,
					Greeting:          "ah ma and ah gong updated",
					MaximumGuestCount: 3,
					Notes:             "some updated notes",
					MobilePhoneNumber: "91231236",
				},
				ID:     1,
				Status: domain.Sent,
			}
		})

		It("should update a invitation given valid values", func() {
			invitation := &domain.Invitation{
				BaseInvitation: baseInvitation,
				ID:             1,
				PrivateID:      "some-private-id",
				Status:         domain.NotSent,
			}

			modifiedInvitation := *invitation

			modifiedInvitation.CategoryID = updateReq.CategoryID
			modifiedInvitation.Greeting = updateReq.Greeting
			modifiedInvitation.MaximumGuestCount = updateReq.MaximumGuestCount
			modifiedInvitation.Notes = updateReq.Notes
			modifiedInvitation.MobilePhoneNumber = updateReq.MobilePhoneNumber
			modifiedInvitation.Status = domain.Sent

			gomock.InOrder(
				mockInvitationStorage.EXPECT().FindInvitationByID(int64(1)).Return(
					invitation, nil),
				mockInvitationStorage.EXPECT().UpdateInvitation(&modifiedInvitation).Return(
					&modifiedInvitation, nil),
			)

			updatedInvitation, err := testInvitationService.UpdateInvitation(updateReq)
			Expect(err).ToNot(HaveOccurred())
			Expect(updatedInvitation).ToNot(BeNil())
			Expect(updatedInvitation.ID).To(Equal(int64(1)))
			Expect(updatedInvitation.PrivateID).To(Equal("some-private-id"))
			Expect(updatedInvitation.CategoryID).To(Equal(int64(2)))
			Expect(updatedInvitation.Greeting).To(Equal("ah ma and ah gong updated"))
			Expect(updatedInvitation.MaximumGuestCount).To(Equal(3))
			Expect(updatedInvitation.Notes).To(Equal("some updated notes"))
			Expect(updatedInvitation.MobilePhoneNumber).To(Equal("91231236"))
			Expect(updatedInvitation.Status).To(BeEquivalentTo(domain.Sent))
		})

		It("should return an error if the invitation cannot be found", func() {
			gomock.InOrder(
				mockInvitationStorage.EXPECT().FindInvitationByID(int64(123123123)).Return(
					nil, postgres.NewPostgresRecordNotFoundError()),
				mockInvitationStorage.EXPECT().UpdateInvitation(gomock.Any()).Times(0),
			)

			updateReq.ID = 123123123

			updatedInvitation, err := testInvitationService.UpdateInvitation(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(InvitationNotFoundError{}))
			Expect(updatedInvitation).To(BeNil())
		})

		It("should not allow invitations with duplicate greetings", func() {
			invitation := &domain.Invitation{
				BaseInvitation: baseInvitation,
				ID:             1,
				PrivateID:      "some-private-id",
				Status:         domain.NotSent,
			}

			gomock.InOrder(
				mockInvitationStorage.EXPECT().FindInvitationByID(int64(1)).Return(invitation, nil),
				mockInvitationStorage.EXPECT().UpdateInvitation(invitation).Return(
					nil, postgres.NewPostgresInvitationGreetingUniqueConstraintError()),
			)

			duplicateInvitation, err := testInvitationService.UpdateInvitation(updateReq)
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
			// Validation should catch it before any attempt to storage is made
			mockInvitationStorage.EXPECT().FindInvitationByID(int64(1)).Times(0)

			updateReq.Greeting = "a"

			updatedInvitation, err := testInvitationService.UpdateInvitation(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation greeting must be between %v to %v characters", GreetingMinLength, GreetingMaxLength)))
			Expect(updatedInvitation).To(BeNil())
		})

		It("should return an error if greeting is too long", func() {
			// Validation should catch it before any attempt to storage is made
			mockInvitationStorage.EXPECT().FindInvitationByID(int64(1)).Times(0)

			updateReq.Greeting = strings.Repeat("a", GreetingMaxLength+1)

			updatedInvitation, err := testInvitationService.UpdateInvitation(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation greeting must be between %v to %v characters", GreetingMinLength, GreetingMaxLength)))
			Expect(updatedInvitation).To(BeNil())
		})

		It("should return an error if maximum guest count is too small", func() {
			// Validation should catch it before any attempt to storage is made
			mockInvitationStorage.EXPECT().FindInvitationByID(int64(1)).Times(0)

			updateReq.MaximumGuestCount = MaximumGuestCountMin - 1

			updatedInvitation, err := testInvitationService.UpdateInvitation(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation maximum guest count must be between %v to %v", MaximumGuestCountMin, MaximumGuestCountMax)))
			Expect(updatedInvitation).To(BeNil())
		})

		It("should return an error if maximum guest count is too large", func() {
			// Validation should catch it before any attempt to storage is made
			mockInvitationStorage.EXPECT().FindInvitationByID(int64(1)).Times(0)

			updateReq.MaximumGuestCount = MaximumGuestCountMax + 1

			updatedInvitation, err := testInvitationService.UpdateInvitation(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation maximum guest count must be between %v to %v", MaximumGuestCountMin, MaximumGuestCountMax)))
			Expect(updatedInvitation).To(BeNil())
		})

		It("should return an error if note is too long", func() {
			// Validation should catch it before any attempt to storage is made
			mockInvitationStorage.EXPECT().FindInvitationByID(int64(1)).Times(0)

			updateReq.Notes = strings.Repeat("a", NoteMaxLength+1)

			updatedInvitation, err := testInvitationService.UpdateInvitation(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal(
				fmt.Sprintf("invitation note must be less than %v characters", NoteMaxLength)))
			Expect(updatedInvitation).To(BeNil())
		})

		It("should return an error if mobile phone number is too long", func() {
			// Validation should catch it before any attempt to storage is made
			mockInvitationStorage.EXPECT().FindInvitationByID(int64(1)).Times(0)

			updateReq.MobilePhoneNumber = "+65 1234567890 121023" // Max is 20

			updatedInvitation, err := testInvitationService.UpdateInvitation(updateReq)
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
			// Validation should catch it before any attempt to storage is made
			mockInvitationStorage.EXPECT().FindInvitationByID(int64(1)).Times(0)

			updateReq.Status = domain.RSVPStatus("INVALID")

			updatedInvitation, err := testInvitationService.UpdateInvitation(updateReq)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(serviceErrors.ValidationError{}))
			Expect(err.Error()).To(Equal("status is invalid"))
			Expect(updatedInvitation).To(BeNil())
		})
	})

	Context("deletion", func() {

		It("should delete an invitation", func() {
			invitation := &domain.Invitation{
				BaseInvitation: domain.BaseInvitation{
					CategoryID:        1,
					Greeting:          "ah ma and ah gong",
					MaximumGuestCount: 2,
					Notes:             "some notes",
					MobilePhoneNumber: "91231234",
				},
				ID:        1,
				PrivateID: "some-private-id",
				Status:    domain.NotSent,
			}

			gomock.InOrder(
				mockInvitationStorage.EXPECT().FindInvitationByID(int64(1)).Return(invitation, nil),
				mockInvitationStorage.EXPECT().DeleteInvitation(invitation).Return(nil),
			)

			err := testInvitationService.DeleteInvitationByID(1)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return an error if the invitation cannot be found", func() {
			mockInvitationStorage.EXPECT().FindInvitationByID(int64(123123123)).Return(
				nil, postgres.NewPostgresRecordNotFoundError())

			err := testInvitationService.DeleteInvitationByID(123123123)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(InvitationNotFoundError{}))
		})
	})
})
