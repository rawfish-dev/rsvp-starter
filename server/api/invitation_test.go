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
	serviceErrors "github.com/rawfish-dev/rsvp-starter/server/services/errors"
	. "github.com/rawfish-dev/rsvp-starter/server/services/invitation"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
)

var _ = Describe("Invitation", func() {

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

		var createInvitationReq domain.InvitationCreateRequest
		var invitation domain.Invitation

		BeforeEach(func() {
			createInvitationReq = domain.InvitationCreateRequest{
				BaseInvitation: domain.BaseInvitation{
					CategoryID:        1,
					Greeting:          "Mitten",
					MaximumGuestCount: 2,
					Notes:             "some notes",
					MobilePhoneNumber: "91234123",
				},
			}

			invitation = domain.Invitation{
				BaseInvitation: createInvitationReq.BaseInvitation,
				ID:             1,
				PrivateID:      "some-private-id",
				Status:         domain.NotSent,
				UpdatedAt:      "2017-12-13",
			}
		})

		It("should return 200 OK and create an invitation given valid values", func() {
			testAPI.InvitationServiceFactory = func(ctx context.Context) interfaces.InvitationServiceProvider {
				mockInvitationService := mock_interfaces.NewMockInvitationServiceProvider(ctrl)
				mockInvitationService.EXPECT().CreateInvitation(&createInvitationReq).
					Return(&invitation, nil)

				return mockInvitationService
			}

			reqBytes, err := json.Marshal(createInvitationReq)
			Expect(err).ToNot(HaveOccurred())

			responseBytes := HitEndpoint(testAPI, "POST", "/api/invitations", bytes.NewBuffer(reqBytes), http.StatusOK)

			var newInvitation domain.Invitation
			err = json.Unmarshal(responseBytes, &newInvitation)
			Expect(err).ToNot(HaveOccurred())

			Expect(newInvitation).To(Equal(invitation))
		})

		It("should return 400 Bad Request when invalid JSON is passed", func() {
			testAPI.InvitationServiceFactory = func(ctx context.Context) interfaces.InvitationServiceProvider {
				mockInvitationService := mock_interfaces.NewMockInvitationServiceProvider(ctrl)
				mockInvitationService.EXPECT().CreateInvitation(&createInvitationReq).Times(0)

				return mockInvitationService
			}

			reqBytes, err := json.Marshal(`{`)
			Expect(err).ToNot(HaveOccurred())

			responseBytes := HitEndpoint(testAPI, "POST", "/api/invitations", bytes.NewBuffer(reqBytes), http.StatusBadRequest)

			var badRequestError domain.CustomBadRequestError
			err = json.Unmarshal(responseBytes, &badRequestError)
			Expect(err).ToNot(HaveOccurred())
			Expect(badRequestError.Error).To(Equal("JSON request was invalid"))
		})

		It("should return 400 Bad Request when a validation error occurs", func() {
			testAPI.InvitationServiceFactory = func(ctx context.Context) interfaces.InvitationServiceProvider {
				mockInvitationService := mock_interfaces.NewMockInvitationServiceProvider(ctrl)
				mockInvitationService.EXPECT().CreateInvitation(&createInvitationReq).
					Return(nil, serviceErrors.NewValidationError([]string{"some validation error"}))

				return mockInvitationService
			}

			reqBytes, err := json.Marshal(createInvitationReq)
			Expect(err).ToNot(HaveOccurred())

			responseBytes := HitEndpoint(testAPI, "POST", "/api/invitations", bytes.NewBuffer(reqBytes), http.StatusBadRequest)

			var badRequestError domain.CustomBadRequestError
			err = json.Unmarshal(responseBytes, &badRequestError)
			Expect(err).ToNot(HaveOccurred())
			Expect(badRequestError.Error).To(Equal("some validation error"))
		})

		It("should return 500 Internal Server Error when an unknown service error occurs", func() {
			testAPI.InvitationServiceFactory = func(ctx context.Context) interfaces.InvitationServiceProvider {
				mockInvitationService := mock_interfaces.NewMockInvitationServiceProvider(ctrl)
				mockInvitationService.EXPECT().CreateInvitation(&createInvitationReq).
					Return(nil, serviceErrors.NewGeneralServiceError())

				return mockInvitationService
			}

			reqBytes, err := json.Marshal(createInvitationReq)
			Expect(err).ToNot(HaveOccurred())

			HitEndpoint(testAPI, "POST", "/api/invitations", bytes.NewBuffer(reqBytes), http.StatusInternalServerError)
		})
	})

	Context("retrieval", func() {

		It("should return 200 OK and the list of invitations", func() {
			invitations := []domain.Invitation{
				{
					BaseInvitation: domain.BaseInvitation{
						CategoryID:        1,
						Greeting:          "Mitten 1",
						MaximumGuestCount: 1,
						Notes:             "some notes 1",
						MobilePhoneNumber: "91234123 1",
					},
					ID:        1,
					PrivateID: "some-private-id-1",
					Status:    domain.NotSent,
					UpdatedAt: "2017-12-11",
				},
				{
					BaseInvitation: domain.BaseInvitation{
						CategoryID:        2,
						Greeting:          "Mitten 2",
						MaximumGuestCount: 2,
						Notes:             "some notes 2",
						MobilePhoneNumber: "91234123 2",
					},
					ID:        2,
					PrivateID: "some-private-id-2",
					Status:    domain.NotSent,
					UpdatedAt: "2017-12-12",
				},
				{
					BaseInvitation: domain.BaseInvitation{
						CategoryID:        3,
						Greeting:          "Mitten 3",
						MaximumGuestCount: 3,
						Notes:             "some notes 3",
						MobilePhoneNumber: "91234123 3",
					},
					ID:        3,
					PrivateID: "some-private-id-3",
					Status:    domain.NotSent,
					UpdatedAt: "2017-12-13",
				},
			}

			testAPI.RSVPServiceFactory = func(ctx context.Context) interfaces.RSVPServiceProvider {
				mockRSVPService := mock_interfaces.NewMockRSVPServiceProvider(ctrl)
				mockRSVPService.EXPECT().ListRSVPs().
					Return(nil, nil)

				return mockRSVPService
			}

			testAPI.InvitationServiceFactory = func(ctx context.Context) interfaces.InvitationServiceProvider {
				mockInvitationService := mock_interfaces.NewMockInvitationServiceProvider(ctrl)
				mockInvitationService.EXPECT().ListInvitations(nil).
					Return(invitations, nil)

				return mockInvitationService
			}

			responseBytes := HitEndpoint(testAPI, "GET", "/api/invitations", nil, http.StatusOK)

			var invitationList []domain.Invitation
			err := json.Unmarshal(responseBytes, &invitationList)
			Expect(err).ToNot(HaveOccurred())

			Expect(invitationList).To(Equal(invitations))
		})

		It("should return 500 Internal Server Error when an unknown rsvp service error occurs", func() {
			testAPI.RSVPServiceFactory = func(ctx context.Context) interfaces.RSVPServiceProvider {
				mockRSVPService := mock_interfaces.NewMockRSVPServiceProvider(ctrl)
				mockRSVPService.EXPECT().ListRSVPs().
					Return(nil, serviceErrors.NewGeneralServiceError())

				return mockRSVPService
			}

			testAPI.InvitationServiceFactory = func(ctx context.Context) interfaces.InvitationServiceProvider {
				mockInvitationService := mock_interfaces.NewMockInvitationServiceProvider(ctrl)
				mockInvitationService.EXPECT().ListInvitations(nil).Times(0)

				return mockInvitationService
			}

			HitEndpoint(testAPI, "GET", "/api/invitations", nil, http.StatusInternalServerError)
		})

		It("should return 500 Internal Server Error when an unknown invitation service error occurs", func() {
			testAPI.RSVPServiceFactory = func(ctx context.Context) interfaces.RSVPServiceProvider {
				mockRSVPService := mock_interfaces.NewMockRSVPServiceProvider(ctrl)
				mockRSVPService.EXPECT().ListRSVPs().
					Return(nil, nil)

				return mockRSVPService
			}

			testAPI.InvitationServiceFactory = func(ctx context.Context) interfaces.InvitationServiceProvider {
				mockInvitationService := mock_interfaces.NewMockInvitationServiceProvider(ctrl)
				mockInvitationService.EXPECT().ListInvitations(nil).
					Return(nil, serviceErrors.NewGeneralServiceError())

				return mockInvitationService
			}

			HitEndpoint(testAPI, "GET", "/api/invitations", nil, http.StatusInternalServerError)
		})
	})

	Context("updating", func() {

		var updateInvitationReq domain.InvitationUpdateRequest

		BeforeEach(func() {
			updateInvitationReq = domain.InvitationUpdateRequest{
				BaseInvitation: domain.BaseInvitation{
					CategoryID:        1,
					Greeting:          "Mitten 1",
					MaximumGuestCount: 1,
					Notes:             "some notes 1",
					MobilePhoneNumber: "91234123 1",
				},
				ID:     1,
				Status: domain.Sent,
			}
		})

		It("should return 200 OK and update an invitation given valid values", func() {
			updatedInvitation := &domain.Invitation{
				BaseInvitation: updateInvitationReq.BaseInvitation,
				ID:             updateInvitationReq.ID,
				PrivateID:      "some-private-id",
				Status:         updateInvitationReq.Status,
			}

			testAPI.InvitationServiceFactory = func(ctx context.Context) interfaces.InvitationServiceProvider {
				mockInvitationService := mock_interfaces.NewMockInvitationServiceProvider(ctrl)
				mockInvitationService.EXPECT().UpdateInvitation(&updateInvitationReq).
					Return(updatedInvitation, nil)

				return mockInvitationService
			}

			reqBytes, err := json.Marshal(updateInvitationReq)
			Expect(err).ToNot(HaveOccurred())

			responseBytes := HitEndpoint(testAPI, "PUT", "/api/invitations/1", bytes.NewBuffer(reqBytes), http.StatusOK)

			var invitation domain.Invitation
			err = json.Unmarshal(responseBytes, &invitation)
			Expect(err).ToNot(HaveOccurred())

			Expect(invitation).To(Equal(*updatedInvitation))
		})

		It("should return 400 Bad Request if the id in the URL does not match the update req ID", func() {
			testAPI.InvitationServiceFactory = func(ctx context.Context) interfaces.InvitationServiceProvider {
				mockInvitationService := mock_interfaces.NewMockInvitationServiceProvider(ctrl)
				mockInvitationService.EXPECT().UpdateInvitation(&updateInvitationReq).Times(0)

				return mockInvitationService
			}

			reqBytes, err := json.Marshal(updateInvitationReq)
			Expect(err).ToNot(HaveOccurred())

			HitEndpoint(testAPI, "PUT", "/api/invitations/2", bytes.NewBuffer(reqBytes), http.StatusBadRequest)
		})

		It("should return 400 Bad Request when invalid JSON is passed", func() {
			testAPI.InvitationServiceFactory = func(ctx context.Context) interfaces.InvitationServiceProvider {
				mockInvitationService := mock_interfaces.NewMockInvitationServiceProvider(ctrl)
				mockInvitationService.EXPECT().UpdateInvitation(&updateInvitationReq).Times(0)

				return mockInvitationService
			}

			reqBytes, err := json.Marshal(`{`)
			Expect(err).ToNot(HaveOccurred())

			responseBytes := HitEndpoint(testAPI, "PUT", "/api/invitations/1", bytes.NewBuffer(reqBytes), http.StatusBadRequest)

			var badRequestError domain.CustomBadRequestError
			err = json.Unmarshal(responseBytes, &badRequestError)
			Expect(err).ToNot(HaveOccurred())
			Expect(badRequestError.Error).To(Equal("JSON request was invalid"))
		})

		It("should return 400 Bad Request when a validation error occurs", func() {
			testAPI.InvitationServiceFactory = func(ctx context.Context) interfaces.InvitationServiceProvider {
				mockInvitationService := mock_interfaces.NewMockInvitationServiceProvider(ctrl)
				mockInvitationService.EXPECT().UpdateInvitation(&updateInvitationReq).
					Return(nil, serviceErrors.NewValidationError([]string{"some validation error"}))

				return mockInvitationService
			}

			reqBytes, err := json.Marshal(updateInvitationReq)
			Expect(err).ToNot(HaveOccurred())

			responseBytes := HitEndpoint(testAPI, "PUT", "/api/invitations/1", bytes.NewBuffer(reqBytes), http.StatusBadRequest)

			var badRequestError domain.CustomBadRequestError
			err = json.Unmarshal(responseBytes, &badRequestError)
			Expect(err).ToNot(HaveOccurred())
			Expect(badRequestError.Error).To(Equal("some validation error"))
		})

		It("should return 500 Internal Server Error when an unknown service error occurs", func() {
			testAPI.InvitationServiceFactory = func(ctx context.Context) interfaces.InvitationServiceProvider {
				mockInvitationService := mock_interfaces.NewMockInvitationServiceProvider(ctrl)
				mockInvitationService.EXPECT().UpdateInvitation(&updateInvitationReq).
					Return(nil, serviceErrors.NewGeneralServiceError())

				return mockInvitationService
			}

			reqBytes, err := json.Marshal(updateInvitationReq)
			Expect(err).ToNot(HaveOccurred())

			HitEndpoint(testAPI, "PUT", "/api/invitations/1", bytes.NewBuffer(reqBytes), http.StatusInternalServerError)
		})
	})

	Context("deletion", func() {

		It("should return 200 OK and delete a category given a valid id", func() {
			testAPI.InvitationServiceFactory = func(ctx context.Context) interfaces.InvitationServiceProvider {
				mockInvitationService := mock_interfaces.NewMockInvitationServiceProvider(ctrl)
				mockInvitationService.EXPECT().DeleteInvitationByID(int64(1)).Return(nil)

				return mockInvitationService
			}

			HitEndpoint(testAPI, "DELETE", "/api/invitations/1", nil, http.StatusOK)
		})

		It("should return 400 Bad Request if the id is not valid", func() {
			testAPI.InvitationServiceFactory = func(ctx context.Context) interfaces.InvitationServiceProvider {
				mockInvitationService := mock_interfaces.NewMockInvitationServiceProvider(ctrl)
				mockInvitationService.EXPECT().DeleteInvitationByID("abc").Times(0)

				return mockInvitationService
			}

			HitEndpoint(testAPI, "DELETE", "/api/invitations/abc", nil, http.StatusBadRequest)
		})

		It("should return 404 Not Found if the id cannot be found", func() {
			testAPI.InvitationServiceFactory = func(ctx context.Context) interfaces.InvitationServiceProvider {
				mockInvitationService := mock_interfaces.NewMockInvitationServiceProvider(ctrl)
				mockInvitationService.EXPECT().DeleteInvitationByID(int64(1)).
					Return(NewInvitationNotFoundError())

				return mockInvitationService
			}

			HitEndpoint(testAPI, "DELETE", "/api/invitations/1", nil, http.StatusNotFound)
		})

		It("should return 500 Internal Server Error when an unknown service error occurs", func() {
			testAPI.InvitationServiceFactory = func(ctx context.Context) interfaces.InvitationServiceProvider {
				mockInvitationService := mock_interfaces.NewMockInvitationServiceProvider(ctrl)
				mockInvitationService.EXPECT().DeleteInvitationByID(int64(1)).
					Return(serviceErrors.NewGeneralServiceError())

				return mockInvitationService
			}

			HitEndpoint(testAPI, "DELETE", "/api/invitations/1", nil, http.StatusInternalServerError)
		})
	})
})
