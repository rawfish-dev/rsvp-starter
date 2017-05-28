package guest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bitbucket.org/rawfish-dev/wedding-rsvp/server/testhelpers"

	"testing"
)

func TestGuest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Guest Suite")
}

var _ = BeforeEach(func() {
	testhelpers.TruncateTestPostgresDB()
})
