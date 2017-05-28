package guest

import (
	"fmt"

	"bitbucket.org/rawfish-dev/wedding-rsvp/server/domain"
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/services/postgres"
	"bitbucket.org/rawfish-dev/wedding-rsvp/server/utils"
	serviceErrors "github.com/rawfish-dev/react-redux-basics/server/services/errors"
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

func (s *service) CreateInvitation(req *domain.InvitationCreateRequest) (*domain.Invitation, error) {
	errorMessages := validateInvitationCreateRequest(req)
	if len(errorMessages) > 0 {
		return nil, serviceErrors.NewValidationError(errorMessages)
	}

	// Populate default extension to mobile phone number
	if req.MobilePhoneNumber == "" {
		req.MobilePhoneNumber = defaultPhoneExtension
	}

	newInvitation, err := s.guestStorage.InsertInvitation(req)
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

func (s *service) ListInvitations() ([]domain.Invitation, error) {
	invitations, err := s.guestStorage.FindAllInvitations()
	if err != nil {
		s.baseService.Error("guest service - unable to list all invitations")
		return nil, serviceErrors.NewGeneralServiceError()
	}

	rsvps, err := s.guestStorage.FindAllRSVPs()
	if err != nil {
		s.baseService.Error("guest service - unable to retrieve all rsvps")
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

	invitation, err := s.guestStorage.FindInvitationByID(req.ID)
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

	updatedInvitation, err := s.guestStorage.UpdateInvitation(invitation)
	if err != nil {
		errorMessage := []string{err.Error()}

		switch err.(type) {
		case postgres.PostgresInvitationGreetingUniqueConstraintError, postgres.PostgresInvitationMobilePhoneNumberUniqueConstraintError:
			return nil, serviceErrors.NewValidationError(errorMessage)
		}

		return nil, serviceErrors.NewGeneralServiceError()
	}

	return updatedInvitation, nil
}

func (s *service) DeleteInvitation(invitationID int64) error {
	err := s.guestStorage.DeleteInvitationByID(invitationID)
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
	invitation, err := s.guestStorage.FindInvitationByPrivateID(privateID)
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
