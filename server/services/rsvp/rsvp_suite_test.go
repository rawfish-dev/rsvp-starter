package rsvp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRsvp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rsvp Suite")
}
