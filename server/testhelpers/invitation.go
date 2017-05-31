package testhelpers

import (
	"fmt"

	"github.com/rawfish-dev/rsvp-starter/server/domain"
	"github.com/rawfish-dev/rsvp-starter/server/services"
	"github.com/rawfish-dev/rsvp-starter/server/services/invitation"
	"github.com/rawfish-dev/rsvp-starter/server/utils"

	"github.com/satori/go.uuid"
)

var _ services.InvitationServiceProvider = new(TestInvitationService)

type TestInvitationService struct {
	services.InvitationServiceProvider
	services.CategoryServiceProvider
}

func NewTestInvitationService() *TestInvitationService {
	testBaseService := NewTestBaseService()
	testPostgresService := NewTestPostgresService()
	testCategoryService := NewTestCategoryService()

	testInvitationService := &TestInvitationService{
		InvitationServiceProvider: invitation.NewService(testBaseService, testPostgresService),
		CategoryServiceProvider:   testCategoryService,
	}

	return testInvitationService
}

func (t *TestInvitationService) CreateTestInvitation() *domain.Invitation {
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
		panic(fmt.Sprintf("test invitation service - failed to create test invitation %v", err))
	}

	return invitation
}
