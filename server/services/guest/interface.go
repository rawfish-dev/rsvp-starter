package guest

import (
	"github.com/rawfish-dev/rsvp-starter/server/domain"
)

type GuestServiceProvider interface {
	CreateCategory(*domain.CategoryCreateRequest) (*domain.Category, error)
	ListCategories() ([]domain.Category, error)
	UpdateCategory(*domain.CategoryUpdateRequest) (*domain.Category, error)
	DeleteCategory(categoryID int64) error

	CreateInvitation(*domain.InvitationCreateRequest) (*domain.Invitation, error)
	ListInvitations() ([]domain.Invitation, error)
	UpdateInvitation(*domain.InvitationUpdateRequest) (*domain.Invitation, error)
	DeleteInvitation(invitationID int64) error
	RetrieveInvitationByPrivateID(privateID string) (*domain.Invitation, error)
	// SendInvitation()

	CreateRSVP(*domain.RSVPCreateRequest) (*domain.RSVP, error)
	ListRSVPs() ([]domain.RSVP, error)
	UpdateRSVP(*domain.RSVPUpdateRequest) (*domain.RSVP, error)
	DeleteRSVP(rsvpID int64) error
	RetrievePrivateRSVP(invitationPrivateID string) (*domain.RSVP, error)
}
