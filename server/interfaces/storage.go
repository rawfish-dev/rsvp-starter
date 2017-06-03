package interfaces

import (
	"github.com/rawfish-dev/rsvp-starter/server/domain"
)

type CategoryStorage interface {
	InsertCategory(*domain.CategoryCreateRequest) (*domain.Category, error)
	FindCategoryByID(categoryID int64) (*domain.Category, error)
	ListCategories() ([]domain.Category, error)
	UpdateCategory(*domain.Category) (*domain.Category, error)
	DeleteCategory(*domain.Category) error
}

type InvitationStorage interface {
	InsertInvitation(*domain.InvitationCreateRequest) (*domain.Invitation, error)
	FindInvitationByID(invitationID int64) (*domain.Invitation, error)
	FindInvitationByPrivateID(privateID string) (*domain.Invitation, error)
	ListInvitations() ([]domain.Invitation, error)
	UpdateInvitation(*domain.Invitation) (*domain.Invitation, error)
	DeleteInvitation(*domain.Invitation) error
}

type RSVPStorage interface {
	InsertRSVP(*domain.RSVPCreateRequest) (*domain.RSVP, error)
	FindRSVPByID(rsvpID int64) (*domain.RSVP, error)
	FindRSVPByInvitationPrivateID(invitationPrivateID string) (*domain.RSVP, error)
	ListRSVPs() ([]domain.RSVP, error)
	UpdateRSVP(*domain.RSVP) (*domain.RSVP, error)
	DeleteRSVP(*domain.RSVP) error
}
