package rsvp

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
)

var _ interfaces.RSVPServiceProvider = new(service)

type service struct {
	ctx         context.Context
	rsvpStorage interfaces.RSVPStorage
}

func NewService(ctx context.Context, rsvpStorage interfaces.RSVPStorage) *service {
	return &service{ctx, rsvpStorage}
}

func (s *service) CreateRSVP(req *domain.RSVPCreateRequest) (*domain.RSVP, error) {
	errorMessages := validateRSVPCreateRequest(req)
	if len(errorMessages) > 0 {
		return nil, serviceErrors.NewValidationError(errorMessages)
	}

	newRSVP, err := s.rsvpStorage.InsertRSVP(req)
	if err != nil {
		errorMessage := []string{err.Error()}

		switch err.(type) {
		case postgres.PostgresRSVPPrivateIDUniqueConstraintError:
			return nil, serviceErrors.NewValidationError(errorMessage)
		}

		return nil, serviceErrors.NewGeneralServiceError()
	}

	return newRSVP, nil
}

func (s *service) ListRSVPs() ([]domain.RSVP, error) {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	rsvps, err := s.rsvpStorage.ListRSVPs()
	if err != nil {
		ctxLogger.Error("rsvp service - unable to list all rsvps")
		return nil, serviceErrors.NewGeneralServiceError()
	}

	return rsvps, nil
}

func (s *service) UpdateRSVP(req *domain.RSVPUpdateRequest) (*domain.RSVP, error) {
	errorMessages := validateRSVPUpdateRequest(req)
	if len(errorMessages) > 0 {
		return nil, serviceErrors.NewValidationError(errorMessages)
	}

	rsvp, err := s.rsvpStorage.FindRSVPByID(req.ID)
	if err != nil {
		switch err.(type) {
		case postgres.PostgresRecordNotFoundError:
			return nil, NewRSVPNotFoundError()
		}

		return nil, serviceErrors.NewGeneralServiceError()
	}

	rsvp.FullName = req.FullName
	rsvp.Attending = req.Attending
	rsvp.GuestCount = req.GuestCount
	rsvp.SpecialDiet = req.SpecialDiet
	rsvp.Remarks = req.Remarks
	rsvp.MobilePhoneNumber = req.MobilePhoneNumber

	updatedInvitation, err := s.rsvpStorage.UpdateRSVP(rsvp)
	if err != nil {
		// TODO:: add specific service error handling

		return nil, serviceErrors.NewGeneralServiceError()
	}

	return updatedInvitation, nil
}

func (s *service) DeleteRSVPByID(rsvpID int64) error {
	rsvp, err := s.rsvpStorage.FindRSVPByID(rsvpID)
	if err != nil {
		switch err.(type) {
		case postgres.PostgresRecordNotFoundError:
			return NewRSVPNotFoundError()
		}

		return serviceErrors.NewGeneralServiceError()
	}

	err = s.rsvpStorage.DeleteRSVP(rsvp)
	if err != nil {
		switch err.(type) {
		case postgres.PostgresRecordNotFoundError:
			return NewRSVPNotFoundError()
		}

		return serviceErrors.NewGeneralServiceError()
	}

	return nil
}

func (s *service) RetrievePrivateRSVP(invitationPrivateID string) (*domain.RSVP, error) {
	rsvp, err := s.rsvpStorage.FindRSVPByInvitationPrivateID(invitationPrivateID)
	if err != nil {
		switch err.(type) {
		case postgres.PostgresRecordNotFoundError:
			return nil, NewRSVPNotFoundError()
		}

		return nil, serviceErrors.NewGeneralServiceError()
	}

	return rsvp, nil
}

func validateBaseRSVP(baseRSVP domain.BaseRSVP) (errorMessages []string) {
	if !utils.IsWithin(len(baseRSVP.FullName), GreetingMinLength, GreetingMaxLength) {
		errorMessages = append(errorMessages, fmt.Sprintf("rsvp full name must be between %v to %v characters", GreetingMinLength, GreetingMaxLength))
	}
	// If not attending let's not care about the guest count
	if baseRSVP.Attending {
		if !utils.IsWithin(baseRSVP.GuestCount, MaximumGuestCountMin, MaximumGuestCountMax) {
			errorMessages = append(errorMessages, fmt.Sprintf("rsvp guest count must be between %v to %v", MaximumGuestCountMin, MaximumGuestCountMax))
		}
	}
	if len(baseRSVP.Remarks) > NoteMaxLength {
		errorMessages = append(errorMessages, fmt.Sprintf("rsvp remarks must be less than %v characters", NoteMaxLength))
	}
	if !utils.IsWithin(len(baseRSVP.MobilePhoneNumber), MobilePhoneNumberMinLength, MobilePhoneNumberMaxLength) {
		errorMessages = append(errorMessages, fmt.Sprintf("rsvp mobile phone number must be between %v to %v in length and contain only numbers", MobilePhoneNumberMinLength, MobilePhoneNumberMaxLength))
	}

	return errorMessages
}

func validateRSVPCreateRequest(req *domain.RSVPCreateRequest) (errorMessages []string) {
	return validateBaseRSVP(req.BaseRSVP)
}

func validateRSVPUpdateRequest(req *domain.RSVPUpdateRequest) (errorMessages []string) {
	if req.ID <= 0 {
		errorMessages = append(errorMessages, "rsvp id is invalid")
	}

	return append(errorMessages, validateBaseRSVP(req.BaseRSVP)...)
}
