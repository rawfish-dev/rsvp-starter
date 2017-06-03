package domain

type BaseInvitation struct {
	CategoryID        int64  `json:"categoryID"`
	Greeting          string `json:"greeting"`
	MaximumGuestCount int    `json:"maximumGuestCount"`
	Notes             string `json:"notes"`
	MobilePhoneNumber string `json:"mobilePhoneNumber"`
}

type InvitationCreateRequest struct {
	BaseInvitation
}

type InvitationUpdateRequest struct {
	BaseInvitation
	ID     int64      `json:"id"`
	Status RSVPStatus `json:"status"`
}

type Invitation struct {
	BaseInvitation
	ID        int64      `json:"id"`
	PrivateID string     `json:"privateID"`
	Status    RSVPStatus `json:"status"`
	UpdatedAt string     `json:"updatedAt"`
}
