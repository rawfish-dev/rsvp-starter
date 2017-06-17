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
	. "github.com/rawfish-dev/rsvp-starter/server/services/rsvp"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
)

var _ = Describe("RSVP", func() {

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

		var createRSVPreq domain.RSVPCreateRequest
		var rsvp domain.RSVP

		BeforeEach(func() {
			createRSVPreq = domain.RSVPCreateRequest{
				BaseRSVP: domain.BaseRSVP{
					FullName:          "Mitten Lin",
					Attending:         true,
					GuestCount:        2,
					SpecialDiet:       true,
					Remarks:           "some remarks",
					MobilePhoneNumber: "91234123",
				},
				InvitationPrivateID: "some-private-id",
				ReCAPTCHAToken:      "some-recaptcha-token",
			}

			rsvp = domain.RSVP{
				BaseRSVP:            createRSVPreq.BaseRSVP,
				ID:                  1,
				InvitationPrivateID: "some-private-id",
				Completed:           true,
				UpdatedAt:           "2017-12-13",
			}
		})

		It("should return 200 OK and create a rsvp given valid values", func() {
			testAPI.RSVPServiceFactory = func(ctx context.Context) interfaces.RSVPServiceProvider {
				mockRSVPService := mock_interfaces.NewMockRSVPServiceProvider(ctrl)
				mockRSVPService.EXPECT().CreateRSVP(&createRSVPreq).
					Return(&rsvp, nil)

				return mockRSVPService
			}

			reqBytes, err := json.Marshal(createRSVPreq)
			Expect(err).ToNot(HaveOccurred())

			responseBytes := HitEndpoint(testAPI, "POST", "/api/rsvps", bytes.NewBuffer(reqBytes), http.StatusOK)

			var newRSVP domain.RSVP
			err = json.Unmarshal(responseBytes, &newRSVP)
			Expect(err).ToNot(HaveOccurred())

			Expect(newRSVP).To(Equal(rsvp))
		})

		It("should return 400 Bad Request when a validation error occurs", func() {
			testAPI.RSVPServiceFactory = func(ctx context.Context) interfaces.RSVPServiceProvider {
				mockRSVPService := mock_interfaces.NewMockRSVPServiceProvider(ctrl)
				mockRSVPService.EXPECT().CreateRSVP(&createRSVPreq).
					Return(nil, serviceErrors.NewValidationError([]string{"some validation error"}))

				return mockRSVPService
			}

			reqBytes, err := json.Marshal(createRSVPreq)
			Expect(err).ToNot(HaveOccurred())

			responseBytes := HitEndpoint(testAPI, "POST", "/api/rsvps", bytes.NewBuffer(reqBytes), http.StatusBadRequest)

			var badRequestError domain.CustomBadRequestError
			err = json.Unmarshal(responseBytes, &badRequestError)
			Expect(err).ToNot(HaveOccurred())
			Expect(badRequestError.Error).To(Equal("some validation error"))
		})

		It("should return 500 Internal Server Error when an unknown service error occurs", func() {
			testAPI.RSVPServiceFactory = func(ctx context.Context) interfaces.RSVPServiceProvider {
				mockRSVPService := mock_interfaces.NewMockRSVPServiceProvider(ctrl)
				mockRSVPService.EXPECT().CreateRSVP(&createRSVPreq).
					Return(nil, serviceErrors.NewGeneralServiceError())

				return mockRSVPService
			}

			reqBytes, err := json.Marshal(createRSVPreq)
			Expect(err).ToNot(HaveOccurred())

			HitEndpoint(testAPI, "POST", "/api/rsvps", bytes.NewBuffer(reqBytes), http.StatusInternalServerError)
		})
	})

	Context("retrieval", func() {

		It("should return 200 OK and the list of rsvps", func() {
			rsvps := []domain.RSVP{
				{
					BaseRSVP: domain.BaseRSVP{
						FullName:          "Mitten Lin",
						Attending:         true,
						GuestCount:        1,
						SpecialDiet:       true,
						Remarks:           "some remarks",
						MobilePhoneNumber: "91234121",
					},
					ID:                  1,
					InvitationPrivateID: "some-private-id",
					Completed:           true,
					UpdatedAt:           "2017-12-11",
				},
				{
					BaseRSVP: domain.BaseRSVP{
						FullName:          "Mitten Lin",
						Attending:         true,
						GuestCount:        2,
						SpecialDiet:       true,
						Remarks:           "some remarks",
						MobilePhoneNumber: "91234122",
					},
					ID:                  2,
					InvitationPrivateID: "some-private-id",
					Completed:           true,
					UpdatedAt:           "2017-12-12",
				},
				{
					BaseRSVP: domain.BaseRSVP{
						FullName:          "Mitten Lin",
						Attending:         true,
						GuestCount:        3,
						SpecialDiet:       true,
						Remarks:           "some remarks",
						MobilePhoneNumber: "91234123",
					},
					ID:                  3,
					InvitationPrivateID: "some-private-id",
					Completed:           true,
					UpdatedAt:           "2017-12-13",
				},
			}

			testAPI.RSVPServiceFactory = func(ctx context.Context) interfaces.RSVPServiceProvider {
				mockRSVPService := mock_interfaces.NewMockRSVPServiceProvider(ctrl)
				mockRSVPService.EXPECT().ListRSVPs().
					Return(rsvps, nil)

				return mockRSVPService
			}

			responseBytes := HitEndpoint(testAPI, "GET", "/api/rsvps", nil, http.StatusOK)

			var rsvpList []domain.RSVP
			err := json.Unmarshal(responseBytes, &rsvpList)
			Expect(err).ToNot(HaveOccurred())

			Expect(rsvpList).To(Equal(rsvps))
		})

		It("should return 500 Internal Server Error when an unknown service error occurs", func() {
			testAPI.RSVPServiceFactory = func(ctx context.Context) interfaces.RSVPServiceProvider {
				mockRSVPService := mock_interfaces.NewMockRSVPServiceProvider(ctrl)
				mockRSVPService.EXPECT().ListRSVPs().
					Return(nil, serviceErrors.NewGeneralServiceError())

				return mockRSVPService
			}

			HitEndpoint(testAPI, "GET", "/api/rsvps", nil, http.StatusInternalServerError)
		})
	})

	Context("updating", func() {

		var updateRSVPReq domain.RSVPUpdateRequest

		BeforeEach(func() {
			updateRSVPReq = domain.RSVPUpdateRequest{
				BaseRSVP: domain.BaseRSVP{
					FullName:          "Mitten Linz",
					Attending:         false,
					GuestCount:        3,
					SpecialDiet:       false,
					Remarks:           "some updated remarks",
					MobilePhoneNumber: "91234123",
				},
				ID:                  1,
				InvitationPrivateID: "some-private-id-3",
			}
		})

		It("should return 200 OK and update a rsvp given valid values", func() {
			updatedRSVP := &domain.RSVP{
				BaseRSVP:            updateRSVPReq.BaseRSVP,
				ID:                  updateRSVPReq.ID,
				InvitationPrivateID: updateRSVPReq.InvitationPrivateID,
				Completed:           true,
				UpdatedAt:           "2017-12-13",
			}

			testAPI.RSVPServiceFactory = func(ctx context.Context) interfaces.RSVPServiceProvider {
				mockRSVPService := mock_interfaces.NewMockRSVPServiceProvider(ctrl)
				mockRSVPService.EXPECT().UpdateRSVP(&updateRSVPReq).
					Return(updatedRSVP, nil)

				return mockRSVPService
			}

			reqBytes, err := json.Marshal(updateRSVPReq)
			Expect(err).ToNot(HaveOccurred())

			responseBytes := HitEndpoint(testAPI, "PUT", "/api/rsvps/1", bytes.NewBuffer(reqBytes), http.StatusOK)

			var rsvp domain.RSVP
			err = json.Unmarshal(responseBytes, &rsvp)
			Expect(err).ToNot(HaveOccurred())

			Expect(rsvp).To(Equal(*updatedRSVP))
		})

		It("should return 400 Bad Request if the id in the URL does not match the update req ID", func() {
			testAPI.RSVPServiceFactory = func(ctx context.Context) interfaces.RSVPServiceProvider {
				mockRSVPService := mock_interfaces.NewMockRSVPServiceProvider(ctrl)
				mockRSVPService.EXPECT().UpdateRSVP(&updateRSVPReq).Times(0)

				return mockRSVPService
			}

			reqBytes, err := json.Marshal(updateRSVPReq)
			Expect(err).ToNot(HaveOccurred())

			HitEndpoint(testAPI, "PUT", "/api/rsvps/2", bytes.NewBuffer(reqBytes), http.StatusBadRequest)
		})

		It("should return 400 Bad Request when a validation error occurs", func() {
			testAPI.RSVPServiceFactory = func(ctx context.Context) interfaces.RSVPServiceProvider {
				mockRSVPService := mock_interfaces.NewMockRSVPServiceProvider(ctrl)
				mockRSVPService.EXPECT().UpdateRSVP(&updateRSVPReq).
					Return(nil, serviceErrors.NewValidationError([]string{"some validation error"}))

				return mockRSVPService
			}

			reqBytes, err := json.Marshal(updateRSVPReq)
			Expect(err).ToNot(HaveOccurred())

			responseBytes := HitEndpoint(testAPI, "PUT", "/api/rsvps/1", bytes.NewBuffer(reqBytes), http.StatusBadRequest)

			var badRequestError domain.CustomBadRequestError
			err = json.Unmarshal(responseBytes, &badRequestError)
			Expect(err).ToNot(HaveOccurred())
			Expect(badRequestError.Error).To(Equal("some validation error"))
		})

		It("should return 500 Internal Server Error when an unknown service error occurs", func() {
			testAPI.RSVPServiceFactory = func(ctx context.Context) interfaces.RSVPServiceProvider {
				mockRSVPService := mock_interfaces.NewMockRSVPServiceProvider(ctrl)
				mockRSVPService.EXPECT().UpdateRSVP(&updateRSVPReq).
					Return(nil, serviceErrors.NewGeneralServiceError())

				return mockRSVPService
			}

			reqBytes, err := json.Marshal(updateRSVPReq)
			Expect(err).ToNot(HaveOccurred())

			HitEndpoint(testAPI, "PUT", "/api/rsvps/1", bytes.NewBuffer(reqBytes), http.StatusInternalServerError)
		})
	})

	Context("deletion", func() {

		It("should return 200 OK and delete a rsvp given a valid id", func() {
			testAPI.RSVPServiceFactory = func(ctx context.Context) interfaces.RSVPServiceProvider {
				mockRSVPService := mock_interfaces.NewMockRSVPServiceProvider(ctrl)
				mockRSVPService.EXPECT().DeleteRSVPByID(int64(1)).Return(nil)

				return mockRSVPService
			}

			HitEndpoint(testAPI, "DELETE", "/api/rsvps/1", nil, http.StatusOK)
		})

		It("should return 400 Bad Request if the id is not valid", func() {
			testAPI.RSVPServiceFactory = func(ctx context.Context) interfaces.RSVPServiceProvider {
				mockRSVPService := mock_interfaces.NewMockRSVPServiceProvider(ctrl)
				mockRSVPService.EXPECT().DeleteRSVPByID("abc").Times(0)

				return mockRSVPService
			}

			HitEndpoint(testAPI, "DELETE", "/api/rsvps/abc", nil, http.StatusBadRequest)
		})

		It("should return 404 Not Found if the id cannot be found", func() {
			testAPI.RSVPServiceFactory = func(ctx context.Context) interfaces.RSVPServiceProvider {
				mockRSVPService := mock_interfaces.NewMockRSVPServiceProvider(ctrl)
				mockRSVPService.EXPECT().DeleteRSVPByID(int64(1)).
					Return(NewRSVPNotFoundError())

				return mockRSVPService
			}

			HitEndpoint(testAPI, "DELETE", "/api/rsvps/1", nil, http.StatusNotFound)
		})

		It("should return 500 Internal Server Error when an unknown service error occurs", func() {
			testAPI.RSVPServiceFactory = func(ctx context.Context) interfaces.RSVPServiceProvider {
				mockRSVPService := mock_interfaces.NewMockRSVPServiceProvider(ctrl)
				mockRSVPService.EXPECT().DeleteRSVPByID(int64(1)).
					Return(serviceErrors.NewGeneralServiceError())

				return mockRSVPService
			}

			HitEndpoint(testAPI, "DELETE", "/api/rsvps/1", nil, http.StatusInternalServerError)
		})
	})
})
