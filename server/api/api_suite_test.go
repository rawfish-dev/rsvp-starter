package api_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/rawfish-dev/rsvp-starter/server/api"
	"github.com/rawfish-dev/rsvp-starter/server/config"
	"github.com/rawfish-dev/rsvp-starter/server/testhelpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Suite")
}

var testAPI *API

var _ = BeforeEach(func() {
	testhelpers.TruncateTestPostgresDB()
})

var _ = BeforeSuite(func() {
	testConfig := config.LoadConfig()
	testAPI = NewAPI(testConfig)
	testAPI.InitRoutes()
})

func HitEndpoint(method, url string, reqBody io.Reader, expectedStatus int) (responseBody []byte) {
	request, err := http.NewRequest(method, url, reqBody)
	Ω(err).ToNot(HaveOccurred())

	response := httptest.NewRecorder()

	testAPI.Router.ServeHTTP(response, request)

	responseBody, err = ioutil.ReadAll(response.Body)
	Ω(err).ToNot(HaveOccurred())

	// When unexpected code, log for debugging
	if response.Code != expectedStatus {
		Fail(fmt.Sprintf("Unexpected status %d (expected %d) :: %s :: %s", response.Code, expectedStatus, string(responseBody), url))
	}

	return
}
