package interfaces

import (
	"time"

	"github.com/rawfish-dev/rsvp-starter/server/domain"
)

type SessionServiceProvider interface {
	CreateWithExpiry(username string) (authToken string, err error)
	IsSessionValid(authToken string) (valid bool, err error)
	Destroy(authToken string) (err error)
}

type JWTServiceProvider interface {
	GenerateAuthToken(additionalClaims map[string]string, duration time.Duration) (authToken string, err error)
	ParseToken(token string) (claims map[string]interface{}, err error)
	IsAuthTokenValid(authToken string) (valid bool)
}

type SecurityServiceProvider interface {
	ValidateCredentials(username, password string) (valid bool)
	VerifyReCAPTCHA(token string) (valid bool)
}

type CategoryServiceProvider interface {
	CreateCategory(*domain.CategoryCreateRequest) (*domain.Category, error)
	ListCategories() ([]domain.Category, error)
	UpdateCategory(*domain.CategoryUpdateRequest) (*domain.Category, error)
	DeleteCategoryByID(categoryID int64) error
}

type InvitationServiceProvider interface {
	CreateInvitation(*domain.InvitationCreateRequest) (*domain.Invitation, error)
	ListInvitations([]domain.RSVP) ([]domain.Invitation, error)
	UpdateInvitation(*domain.InvitationUpdateRequest) (*domain.Invitation, error)
	DeleteInvitationByID(invitationID int64) error
	RetrieveInvitationByPrivateID(privateID string) (*domain.Invitation, error)
	// SendInvitation()
}

type RSVPServiceProvider interface {
	CreateRSVP(*domain.RSVPCreateRequest) (*domain.RSVP, error)
	ListRSVPs() ([]domain.RSVP, error)
	UpdateRSVP(*domain.RSVPUpdateRequest) (*domain.RSVP, error)
	DeleteRSVPByID(rsvpID int64) error
	RetrievePrivateRSVP(invitationPrivateID string) (*domain.RSVP, error)
}
