package interfaces

import (
	"github.com/rawfish-dev/rsvp-starter/server/domain"
)

type Storage interface {
	Close() error
	CategoryStorage
	InvitationStorage
	RSVPStorage
}

type CategoryStorage interface {
	InsertCategory(*domain.CategoryCreateRequest) (*domain.Category, error)
	FindCategoryByID(categoryID int64) (*domain.Category, error)
	FindAllCategories() ([]domain.Category, error)
	UpdateCategory(*domain.Category) (*domain.Category, error)
	DeleteCategoryByID(categoryID int64) error
}

type InvitationStorage interface {
	InsertInvitation(*domain.InvitationCreateRequest) (*domain.Invitation, error)
	FindInvitationByID(invitationID int64) (*domain.Invitation, error)
	FindInvitationByPrivateID(privateID string) (*domain.Invitation, error)
	FindAllInvitations() ([]domain.Invitation, error)
	UpdateInvitation(*domain.Invitation) (*domain.Invitation, error)
	DeleteInvitationByID(invitationID int64) error
}

type RSVPStorage interface {
	InsertRSVP(*domain.RSVPCreateRequest) (*domain.RSVP, error)
	FindRSVPByID(rsvpID int64) (*domain.RSVP, error)
	FindRSVPByInvitationPrivateID(invitationPrivateID string) (*domain.RSVP, error)
	FindAllRSVPs() ([]domain.RSVP, error)
	UpdateRSVP(*domain.RSVP) (*domain.RSVP, error)
	DeleteRSVPByID(rsvpID int64) error
}
