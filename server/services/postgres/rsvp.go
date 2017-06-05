package postgres

import (
	"fmt"
	"strings"
	"time"

	"github.com/rawfish-dev/rsvp-starter/server/domain"
	"github.com/rawfish-dev/rsvp-starter/server/interfaces"
)

type rsvp struct {
	baseModel
	InvitationPrivateID string `db:"invitation_private_id"`
	FullName            string `db:"full_name"`
	Attending           bool   `db:"attending"`
	GuestCount          int    `db:"guest_count"`
	SpecialDiet         bool   `db:"special_diet"`
	Remarks             string `db:"remarks"`
	MobilePhoneNumber   string `db:"mobile_phone_number"`
}

var (
	rsvpColumns = strings.Join([]string{
		"id",
		"invitation_private_id",
		"full_name",
		"attending",
		"guest_count",
		"special_diet",
		"remarks",
		"mobile_phone_number",
	}, ",")
)

func (s *service) InsertRSVP(req *domain.RSVPCreateRequest) (*domain.RSVP, error) {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	rsvp := &rsvp{
		InvitationPrivateID: req.InvitationPrivateID,
		FullName:            req.FullName,
		Attending:           req.Attending,
		GuestCount:          req.GuestCount,
		SpecialDiet:         req.SpecialDiet,
		Remarks:             req.Remarks,
		MobilePhoneNumber:   req.MobilePhoneNumber,
	}

	err := s.gorpDB.Insert(rsvp)
	if err != nil {
		if isRSVPPrivateIDUniqueConstraintError(err) {
			ctxLogger.Warnf("postgres service - unable to insert rsvp with a duplicate private id %v", rsvp.InvitationPrivateID)
			return nil, NewPostgresRSVPPrivateIDUniqueConstraintError()
		}

		ctxLogger.Errorf("postgres service - unable to insert rsvp due to %v", err)
		return nil, NewPostgresOperationError()
	}

	newRSVP := &domain.RSVP{
		BaseRSVP: domain.BaseRSVP{
			FullName:          rsvp.FullName,
			Attending:         rsvp.Attending,
			GuestCount:        rsvp.GuestCount,
			SpecialDiet:       rsvp.SpecialDiet,
			Remarks:           rsvp.Remarks,
			MobilePhoneNumber: rsvp.MobilePhoneNumber,
		},
		ID:                  rsvp.ID,
		InvitationPrivateID: rsvp.InvitationPrivateID,
		UpdatedAt:           rsvp.UpdatedAt.Format(time.RFC3339),
		Completed:           true,
	}

	return newRSVP, nil
}

func (s *service) FindRSVPByID(rsvpID int64) (*domain.RSVP, error) {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	query := fmt.Sprintf(`
		SELECT %v
		FROM rsvps
		WHERE id=$1
	`, rsvpColumns)

	var rsvp rsvp

	err := s.gorpDB.SelectOne(&rsvp, query, rsvpID)
	if err != nil {
		if isNotFoundError(err) {
			ctxLogger.Warnf("postgres service - unable to find rsvp with id %v", rsvpID)
			return nil, NewPostgresRecordNotFoundError()
		}

		ctxLogger.Errorf("postgres service - unable to find rsvp with id %v due to %v", rsvpID, err)
		return nil, NewPostgresOperationError()
	}

	domainRSVP := &domain.RSVP{
		BaseRSVP: domain.BaseRSVP{
			FullName:          rsvp.FullName,
			Attending:         rsvp.Attending,
			GuestCount:        rsvp.GuestCount,
			SpecialDiet:       rsvp.SpecialDiet,
			Remarks:           rsvp.Remarks,
			MobilePhoneNumber: rsvp.MobilePhoneNumber,
		},
		ID:                  rsvp.ID,
		InvitationPrivateID: rsvp.InvitationPrivateID,
		UpdatedAt:           rsvp.UpdatedAt.Format(time.RFC3339),
		Completed:           true,
	}

	return domainRSVP, nil
}

func (s *service) FindRSVPByInvitationPrivateID(invitationPrivateID string) (*domain.RSVP, error) {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	query := fmt.Sprintf(`
		SELECT %v
		FROM rsvps
		WHERE invitation_private_id=$1
	`, rsvpColumns)

	var rsvp rsvp

	err := s.gorpDB.SelectOne(&rsvp, query, invitationPrivateID)
	if err != nil {
		if isNotFoundError(err) {
			ctxLogger.Warnf("postgres service - unable to find rsvp with invitation private id %v", invitationPrivateID)
			return nil, NewPostgresRecordNotFoundError()
		}

		ctxLogger.Errorf("postgres service - unable to find rsvp with invitation private id %v due to %v", invitationPrivateID, err)
		return nil, NewPostgresOperationError()
	}

	domainRSVP := &domain.RSVP{
		BaseRSVP: domain.BaseRSVP{
			FullName:          rsvp.FullName,
			Attending:         rsvp.Attending,
			GuestCount:        rsvp.GuestCount,
			SpecialDiet:       rsvp.SpecialDiet,
			Remarks:           rsvp.Remarks,
			MobilePhoneNumber: rsvp.MobilePhoneNumber,
		},
		// ID:                  rsvp.ID, omit since no operations can be performed against it
		InvitationPrivateID: rsvp.InvitationPrivateID,
		UpdatedAt:           rsvp.UpdatedAt.Format(time.RFC3339),
		Completed:           true,
	}

	return domainRSVP, nil
}

func (s *service) ListRSVPs() ([]domain.RSVP, error) {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	query := fmt.Sprintf(`
		SELECT %v
		FROM rsvps
		ORDER BY updated_at DESC
	`, rsvpColumns)

	var rsvps []rsvp

	_, err := s.gorpDB.Select(&rsvps, query)
	if err != nil {
		ctxLogger.Errorf("postgres service - unable to retrieve all rsvps due to %v", err)
		return nil, NewPostgresOperationError()
	}

	domainRSVPs := make([]domain.RSVP, len(rsvps))
	for idx := range rsvps {
		domainRSVPs[idx] = domain.RSVP{
			BaseRSVP: domain.BaseRSVP{
				FullName:          rsvps[idx].FullName,
				Attending:         rsvps[idx].Attending,
				GuestCount:        rsvps[idx].GuestCount,
				SpecialDiet:       rsvps[idx].SpecialDiet,
				Remarks:           rsvps[idx].Remarks,
				MobilePhoneNumber: rsvps[idx].MobilePhoneNumber,
			},
			ID:                  rsvps[idx].ID,
			InvitationPrivateID: rsvps[idx].InvitationPrivateID,
			UpdatedAt:           rsvps[idx].UpdatedAt.Format(time.RFC3339),
			Completed:           true,
		}
	}

	return domainRSVPs, nil
}

func (s *service) UpdateRSVP(domainRSVP *domain.RSVP) (*domain.RSVP, error) {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	rsvp := &rsvp{
		baseModel: baseModel{
			ID: domainRSVP.ID,
		},
		InvitationPrivateID: domainRSVP.InvitationPrivateID,
		FullName:            domainRSVP.FullName,
		Attending:           domainRSVP.Attending,
		GuestCount:          domainRSVP.GuestCount,
		SpecialDiet:         domainRSVP.SpecialDiet,
		Remarks:             domainRSVP.Remarks,
		MobilePhoneNumber:   domainRSVP.MobilePhoneNumber,
	}

	_, err := s.gorpDB.Update(rsvp)
	if err != nil {
		ctxLogger.Errorf("postgres service - unable to update rsvp %+v due to %v", rsvp, err)
		return nil, NewPostgresOperationError()
	}

	domainRSVP.UpdatedAt = rsvp.UpdatedAt.Format(time.RFC3339)
	domainRSVP.Completed = true

	return domainRSVP, nil
}

func (s *service) DeleteRSVP(rsvp *domain.RSVP) error {
	ctxLogger := s.ctx.Value("logger").(interfaces.Logger)

	_, err := s.gorpDB.Delete(rsvp)
	if err != nil {
		ctxLogger.Errorf("postgres service - unable to delete rsvp with id %v due to %v", rsvp.ID, err)
		return NewPostgresOperationError()
	}

	return nil
}
