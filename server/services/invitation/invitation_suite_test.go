package invitation_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestInvitation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Invitation Suite")
}
