package testhelpers

import (
	"fmt"

	"github.com/rawfish-dev/rsvp-starter/server/domain"
	"github.com/rawfish-dev/rsvp-starter/server/services"
	"github.com/rawfish-dev/rsvp-starter/server/services/rsvp"
)

var _ services.RSVPServiceProvider = new(TestRSVPService)

type TestRSVPService struct {
	services.RSVPServiceProvider
}

func NewTestRSVPService() *TestRSVPService {
	testBaseService := NewTestBaseService()
	testPostgresService := NewTestPostgresService()

	testRSVPService := &TestRSVPService{
		RSVPServiceProvider: rsvp.NewService(testBaseService, testPostgresService),
	}

	return testRSVPService
}

func (t *TestRSVPService) CreateTestRSVP(attending bool) *domain.RSVP {
	testInvitation := t.CreateTestInvitation()

	req := &domain.RSVPCreateRequest{
		BaseRSVP: domain.BaseRSVP{
			FullName:          testInvitation.Greeting,
			Attending:         attending,
			GuestCount:        1,
			SpecialDiet:       true,
			Remarks:           "some remarks",
			MobilePhoneNumber: "91234123",
		},
		InvitationPrivateID: testInvitation.PrivateID,
	}

	rsvp, err := t.CreateRSVP(req)
	if err != nil {
		panic(fmt.Sprintf("test rsvp service - failed to create test rsvp %v", err))
	}

	return rsvp
}
