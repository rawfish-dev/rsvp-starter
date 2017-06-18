package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/rawfish-dev/rsvp-starter/server/api"
	"github.com/rawfish-dev/rsvp-starter/server/config"
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"
	"github.com/rawfish-dev/rsvp-starter/server/mock"
	serviceErrors "github.com/rawfish-dev/rsvp-starter/server/services/errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
)

var _ = Describe("Middleware", func() {

	var ctrl *gomock.Controller
	var testAPI *api.API

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("should return 500 Internal Server Error when an unknown session service error occurs", func() {
		testConfig := config.LoadConfig()
		testAPI = api.NewAPI(testConfig)

		testAPI.SessionServiceFactory = func(ctx context.Context) interfaces.SessionServiceProvider {
			mockSessionService := mock_interfaces.NewMockSessionServiceProvider(ctrl)
			mockSessionService.EXPECT().IsSessionValid("").
				Return(false, serviceErrors.NewGeneralServiceError())

			return mockSessionService
		}

		testAPI.InitRoutes()

		// Just use any protected route with any request body
		reqBytes, err := json.Marshal(`{}`)
		Expect(err).ToNot(HaveOccurred())

		HitEndpoint(testAPI, "POST", "/api/invitations", bytes.NewBuffer(reqBytes), http.StatusInternalServerError)
	})

	It("should return 401 Unauthorized when the session is invalid", func() {
		testConfig := config.LoadConfig()
		testAPI = api.NewAPI(testConfig)

		testAPI.SessionServiceFactory = func(ctx context.Context) interfaces.SessionServiceProvider {
			mockSessionService := mock_interfaces.NewMockSessionServiceProvider(ctrl)
			mockSessionService.EXPECT().IsSessionValid("").
				Return(false, nil)

			return mockSessionService
		}

		testAPI.InitRoutes()

		// Just use any protected route with any request body
		reqBytes, err := json.Marshal(`{}`)
		Expect(err).ToNot(HaveOccurred())

		HitEndpoint(testAPI, "POST", "/api/invitations", bytes.NewBuffer(reqBytes), http.StatusUnauthorized)
	})
})
