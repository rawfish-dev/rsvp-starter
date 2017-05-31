package invitation_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/rawfish-dev/rsvp-starter/server/testhelpers"

	"testing"
)

func TestInvitation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Invitation Suite")
}

var _ = BeforeEach(func() {
	testhelpers.TruncateTestPostgresDB()
})
