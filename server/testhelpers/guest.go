package testhelpers

import (
	"fmt"

	"github.com/rawfish-dev/rsvp-starter/server/domain"
	"github.com/rawfish-dev/rsvp-starter/server/services/guest"
	"github.com/rawfish-dev/rsvp-starter/server/utils"

	"github.com/satori/go.uuid"
)

var _ guest.GuestServiceProvider = new(TestGuestService)

type TestGuestService struct {
	guest.GuestServiceProvider
}

func NewTestGuestService() *TestGuestService {
	testBaseService := NewTestBaseService()
	testPostgresService := NewTestPostgresService()

	testGuestService := &TestGuestService{
		GuestServiceProvider: guest.NewService(testBaseService, testPostgresService),
	}

	return testGuestService
}

func (t *TestGuestService) CreateTestCategory() *domain.Category {
	req := &domain.CategoryCreateRequest{
		Tag: "some tag " + uuid.NewV4().String(),
	}

	category, err := t.CreateCategory(req)
	if err != nil {
		panic(fmt.Sprintf("test guest service - failed to create test category %v", err))
	}

	return category
}

func (t *TestGuestService) CreateTestInvitation() *domain.Invitation {
	testCategory := t.CreateTestCategory()

	req := &domain.InvitationCreateRequest{
		BaseInvitation: domain.BaseInvitation{
			CategoryID:        testCategory.ID,
			Greeting:          "ah ma and ah gong " + uuid.NewV4().String(),
			MaximumGuestCount: 2,
			Notes:             "some notes",
			MobilePhoneNumber: utils.GenerateRandomPhoneNumber(),
		},
	}

	invitation, err := t.CreateInvitation(req)
	if err != nil {
		panic(fmt.Sprintf("test guest service - failed to create test invitation %v", err))
	}

	return invitation
}

func (t *TestGuestService) CreateTestRSVP(attending bool) *domain.RSVP {
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
		panic(fmt.Sprintf("test guest service - failed to create test rsvp %v", err))
	}

	return rsvp
}
