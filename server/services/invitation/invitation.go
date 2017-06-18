package invitation

import (
	"fmt"

	"github.com/rawfish-dev/rsvp-starter/server/domain"
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"
	serviceErrors "github.com/rawfish-dev/rsvp-starter/server/services/errors"
	"github.com/rawfish-dev/rsvp-starter/server/services/postgres"
	"github.com/rawfish-dev/rsvp-starter/server/utils"

	"golang.org/x/net/context"
)

const (
	GreetingMinLength          = 2
	GreetingMaxLength          = 100
	MaximumGuestCountMin       = 1
	MaximumGuestCountMax       = 10
	NoteMaxLength              = 500
	MobilePhoneNumberMinLength = 8
	MobilePhoneNumberMaxLength = 20
	defaultPhoneExtension      = "+65"
)

var _ interfaces.InvitationServiceProvider = new(service)

type service struct {
	ctx               context.Context
	invitationStorage interfaces.InvitationStorage
}

func NewService(ctx context.Context, invitationStorage interfaces.InvitationStorage) *service {
	return &service{ctx, invitationStorage}
}

func (s *service) CreateInvitation(req *domain.InvitationCreateRequest) (*domain.Invitation, error) {
	errorMessages := validateInvitationCreateRequest(req)
	if len(errorMessages) > 0 {
		return nil, serviceErrors.NewValidationError(errorMessages)
	}

	// Populate default extension to mobile phone number
	if req.MobilePhoneNumber == "" {
		req.MobilePhoneNumber = defaultPhoneExtension
	}

	newInvitation, err := s.invitationStorage.InsertInvitation(req)
	if err != nil {
		errorMessage := []string{err.Error()}

		switch err.(type) {
		case postgres.PostgresInvitationGreetingUniqueConstraintError, postgres.PostgresInvitationMobilePhoneNumberUniqueConstraintError:
			return nil, serviceErrors.NewValidationError(errorMessage)
		}

		return nil, serviceErrors.NewGeneralServiceError()
	}

	return newInvitation, nil
}

func (s *service) ListInvitations(rsvps []domain.RSVP) ([]domain.Invitation, error) {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	invitations, err := s.invitationStorage.ListInvitations()
	if err != nil {
		ctxLogger.Error("invitation service - unable to list all invitations")
		return nil, serviceErrors.NewGeneralServiceError()
	}

	mappedRSVPs := make(map[string]domain.RSVP)

	// Track which invitation private ids have already been RSVP-ed
	for idx := range rsvps {
		mappedRSVPs[rsvps[idx].InvitationPrivateID] = rsvps[idx]
	}

	for idx := range invitations {
		if rsvp, ok := mappedRSVPs[invitations[idx].PrivateID]; ok {
			// Map invitation status according to whether the guests are attending
			if rsvp.Attending {
				invitations[idx].Status = domain.RepliedAttending
			} else {
				invitations[idx].Status = domain.RepliedNotAttending
			}
		}
	}

	return invitations, nil
}

func (s *service) UpdateInvitation(req *domain.InvitationUpdateRequest) (*domain.Invitation, error) {
	errorMessages := validateInvitationUpdateRequest(req)
	if len(errorMessages) > 0 {
		return nil, serviceErrors.NewValidationError(errorMessages)
	}

	invitation, err := s.invitationStorage.FindInvitationByID(req.ID)
	if err != nil {
		switch err.(type) {
		case postgres.PostgresRecordNotFoundError:
			return nil, NewInvitationNotFoundError()
		}

		return nil, serviceErrors.NewGeneralServiceError()
	}

	invitation.CategoryID = req.CategoryID
	invitation.Greeting = req.Greeting
	invitation.MaximumGuestCount = req.MaximumGuestCount
	invitation.Notes = req.Notes
	invitation.MobilePhoneNumber = req.MobilePhoneNumber
	invitation.Status = req.Status

	updatedInvitation, err := s.invitationStorage.UpdateInvitation(invitation)
	if err != nil {
		errorMessage := []string{err.Error()}

		switch err.(type) {
		case postgres.PostgresInvitationGreetingUniqueConstraintError,
			postgres.PostgresInvitationMobilePhoneNumberUniqueConstraintError:
			return nil, serviceErrors.NewValidationError(errorMessage)
		}

		return nil, serviceErrors.NewGeneralServiceError()
	}

	return updatedInvitation, nil
}

func (s *service) DeleteInvitationByID(invitationID int64) error {
	invitation, err := s.invitationStorage.FindInvitationByID(invitationID)
	if err != nil {
		switch err.(type) {
		case postgres.PostgresRecordNotFoundError:
			return NewInvitationNotFoundError()
		}

		return serviceErrors.NewGeneralServiceError()
	}

	err = s.invitationStorage.DeleteInvitation(invitation)
	if err != nil {
		switch err.(type) {
		case postgres.PostgresRecordNotFoundError:
			return NewInvitationNotFoundError()
		}

		return serviceErrors.NewGeneralServiceError()
	}

	return nil
}

func (s *service) RetrieveInvitationByPrivateID(privateID string) (*domain.Invitation, error) {
	invitation, err := s.invitationStorage.FindInvitationByPrivateID(privateID)
	if err != nil {
		switch err.(type) {
		case postgres.PostgresRecordNotFoundError:
			return nil, NewInvitationNotFoundError()
		}

		return nil, serviceErrors.NewGeneralServiceError()
	}

	return invitation, nil
}

func validateBaseInvitation(baseInvitation domain.BaseInvitation) (errorMessages []string) {
	if !utils.IsWithin(len(baseInvitation.Greeting), GreetingMinLength, GreetingMaxLength) {
		errorMessages = append(errorMessages, fmt.Sprintf("invitation greeting must be between %v to %v characters", GreetingMinLength, GreetingMaxLength))
	}
	if !utils.IsWithin(baseInvitation.MaximumGuestCount, MaximumGuestCountMin, MaximumGuestCountMax) {
		errorMessages = append(errorMessages, fmt.Sprintf("invitation maximum guest count must be between %v to %v", MaximumGuestCountMin, MaximumGuestCountMax))
	}
	if len(baseInvitation.Notes) > NoteMaxLength {
		errorMessages = append(errorMessages, fmt.Sprintf("invitation note must be less than %v characters", NoteMaxLength))
	}
	if len(baseInvitation.MobilePhoneNumber) > MobilePhoneNumberMaxLength {
		errorMessages = append(errorMessages, fmt.Sprintf("invitation mobile phone number must be less than %v in length", MobilePhoneNumberMaxLength))
	}

	return errorMessages
}

func validateInvitationCreateRequest(req *domain.InvitationCreateRequest) (errorMessages []string) {
	return validateBaseInvitation(req.BaseInvitation)
}

func validateInvitationUpdateRequest(req *domain.InvitationUpdateRequest) (errorMessages []string) {
	if req.ID <= 0 {
		errorMessages = append(errorMessages, "invitation id is invalid")
	}
	if !domain.IsValidRSVPStatus(req.Status) {
		errorMessages = append(errorMessages, "status is invalid")
	}

	return append(errorMessages, validateBaseInvitation(req.BaseInvitation)...)
}
