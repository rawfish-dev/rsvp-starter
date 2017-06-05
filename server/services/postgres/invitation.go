package postgres

import (
	"fmt"
	"strings"
	"time"

	"github.com/rawfish-dev/rsvp-starter/server/domain"
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"

	"github.com/satori/go.uuid"
)

type invitation struct {
	baseModel
	CategoryID        int64  `db:"category_id"`
	PrivateID         string `db:"private_id"`
	Greeting          string `db:"greeting"`
	MaximumGuestCount int    `db:"maximum_guest_count"`
	Status            string `db:"status"`
	Notes             string `db:"notes"`
	MobilePhoneNumber string `db:"mobile_phone_number"`
}

var (
	invitationColumns = strings.Join([]string{
		"id",
		"category_id",
		"private_id",
		"greeting",
		"maximum_guest_count",
		"status",
		"notes",
		"mobile_phone_number",
		"created_at",
		"updated_at",
	}, ",")
)

func (s *service) InsertInvitation(req *domain.InvitationCreateRequest) (*domain.Invitation, error) {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	invitation := &invitation{
		CategoryID:        req.CategoryID,
		PrivateID:         uuid.NewV4().String(),
		Greeting:          req.Greeting,
		MaximumGuestCount: req.MaximumGuestCount,
		Status:            string(domain.NotSent),
		Notes:             req.Notes,
		MobilePhoneNumber: req.MobilePhoneNumber,
	}

	err := s.gorpDB.Insert(invitation)
	if err != nil {
		if isInvitationGreetingUniqueConstraintError(err) {
			ctxLogger.Warnf("postgres service - unable to insert invitation with a duplicate greeting %v", invitation.Greeting)
			return nil, NewPostgresInvitationGreetingUniqueConstraintError()
		}
		if isInvitationMobilePhoneNumberUniqueConstraintError(err) {
			ctxLogger.Warnf("postgres service - unable to insert invitation with a duplicate mobile phone number %v", invitation.MobilePhoneNumber)
			return nil, NewPostgresInvitationMobilePhoneNumberUniqueConstraintError()
		}

		ctxLogger.Errorf("postgres service - unable to insert invitation due to %v", err)
		return nil, NewPostgresOperationError()
	}

	newInvitation := &domain.Invitation{
		BaseInvitation: domain.BaseInvitation{
			CategoryID:        invitation.CategoryID,
			Greeting:          invitation.Greeting,
			MaximumGuestCount: invitation.MaximumGuestCount,
			Notes:             invitation.Notes,
			MobilePhoneNumber: invitation.MobilePhoneNumber,
		},
		ID:        invitation.ID,
		PrivateID: invitation.PrivateID,
		Status:    domain.RSVPStatus(invitation.Status),
		UpdatedAt: invitation.UpdatedAt.Format(time.RFC3339),
	}

	return newInvitation, nil
}

func (s *service) FindInvitationByID(invitationID int64) (*domain.Invitation, error) {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	query := fmt.Sprintf(`
		SELECT %v
		FROM invitations
		WHERE id=$1
	`, invitationColumns)

	var invitation invitation

	err := s.gorpDB.SelectOne(&invitation, query, invitationID)
	if err != nil {
		if isNotFoundError(err) {
			ctxLogger.Warnf("postgres service - unable to find invitation with id %v", invitationID)
			return nil, NewPostgresRecordNotFoundError()
		}

		ctxLogger.Errorf("postgres service - unable to find invitation with id %v due to %v", invitationID, err)
		return nil, NewPostgresOperationError()
	}

	domainInvitation := &domain.Invitation{
		BaseInvitation: domain.BaseInvitation{
			CategoryID:        invitation.CategoryID,
			Greeting:          invitation.Greeting,
			MaximumGuestCount: invitation.MaximumGuestCount,
			Notes:             invitation.Notes,
			MobilePhoneNumber: invitation.MobilePhoneNumber,
		},
		ID:        invitation.ID,
		PrivateID: invitation.PrivateID,
		Status:    domain.RSVPStatus(invitation.Status),
		UpdatedAt: invitation.UpdatedAt.Format(time.RFC3339),
	}

	return domainInvitation, nil
}

func (s *service) FindInvitationByPrivateID(privateID string) (*domain.Invitation, error) {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	query := fmt.Sprintf(`
		SELECT %v
		FROM invitations
		WHERE private_id=$1
	`, invitationColumns)

	var invitation invitation

	err := s.gorpDB.SelectOne(&invitation, query, privateID)
	if err != nil {
		if isNotFoundError(err) {
			ctxLogger.Warnf("postgres service - unable to find invitation with private id %v", privateID)
			return nil, NewPostgresRecordNotFoundError()
		}

		ctxLogger.Errorf("postgres service - unable to find invitation with private id %v due to %v", privateID, err)
		return nil, NewPostgresOperationError()
	}

	domainInvitation := &domain.Invitation{
		BaseInvitation: domain.BaseInvitation{
			CategoryID:        invitation.CategoryID,
			Greeting:          invitation.Greeting,
			MaximumGuestCount: invitation.MaximumGuestCount,
			Notes:             invitation.Notes,
			MobilePhoneNumber: invitation.MobilePhoneNumber,
		},
		ID:        invitation.ID,
		PrivateID: invitation.PrivateID,
		Status:    domain.RSVPStatus(invitation.Status),
		UpdatedAt: invitation.UpdatedAt.Format(time.RFC3339),
	}

	return domainInvitation, nil
}

func (s *service) ListInvitations() ([]domain.Invitation, error) {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	query := fmt.Sprintf(`
		SELECT %v
		FROM invitations
		ORDER BY updated_at DESC
	`, invitationColumns)

	var invitations []invitation

	_, err := s.gorpDB.Select(&invitations, query)
	if err != nil {
		ctxLogger.Errorf("postgres service - unable to retrieve all invitations due to %v", err)
		return nil, NewPostgresOperationError()
	}

	domainInvitations := make([]domain.Invitation, len(invitations))
	for idx := range invitations {
		domainInvitations[idx] = domain.Invitation{
			BaseInvitation: domain.BaseInvitation{
				CategoryID:        invitations[idx].CategoryID,
				Greeting:          invitations[idx].Greeting,
				MaximumGuestCount: invitations[idx].MaximumGuestCount,
				Notes:             invitations[idx].Notes,
				MobilePhoneNumber: invitations[idx].MobilePhoneNumber,
			},
			ID:        invitations[idx].ID,
			PrivateID: invitations[idx].PrivateID,
			Status:    domain.RSVPStatus(invitations[idx].Status),
			UpdatedAt: invitations[idx].UpdatedAt.Format(time.RFC3339),
		}
	}

	return domainInvitations, nil
}

func (s *service) UpdateInvitation(domainInvitation *domain.Invitation) (*domain.Invitation, error) {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	invitation := &invitation{
		baseModel: baseModel{
			ID: domainInvitation.ID,
		},
		CategoryID:        domainInvitation.CategoryID,
		PrivateID:         domainInvitation.PrivateID,
		Greeting:          domainInvitation.Greeting,
		MaximumGuestCount: domainInvitation.MaximumGuestCount,
		Status:            string(domainInvitation.Status),
		Notes:             domainInvitation.Notes,
		MobilePhoneNumber: domainInvitation.MobilePhoneNumber,
	}

	_, err := s.gorpDB.Update(invitation)
	if err != nil {
		ctxLogger.Errorf("postgres service - unable to update invitation %+v due to %v", invitation, err)
		return nil, NewPostgresOperationError()
	}

	domainInvitation.UpdatedAt = invitation.UpdatedAt.Format(time.RFC3339)

	return domainInvitation, nil
}

func (s *service) DeleteInvitation(invitation *domain.Invitation) error {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	_, err := s.gorpDB.Delete(invitation)
	if err != nil {
		ctxLogger.Errorf("postgres service - unable to delete invitation with id %v due to %v", invitation.ID, err)
		return NewPostgresOperationError()
	}

	return nil
}
