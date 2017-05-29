package guest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/rawfish-dev/rsvp-starter/server/testhelpers"

	"testing"
)

func TestGuest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Guest Suite")
}

var _ = BeforeEach(func() {
	testhelpers.TruncateTestPostgresDB()
})
