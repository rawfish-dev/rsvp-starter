package domain

type RSVPStatus string

const (
	NotSent             RSVPStatus = "NS"
	Sent                RSVPStatus = "ST"
	RepliedAttending    RSVPStatus = "RA"
	RepliedNotAttending RSVPStatus = "RN"
)

func IsValidRSVPStatus(status RSVPStatus) bool {
	for _, validStatus := range []RSVPStatus{NotSent, Sent} {
		if status == validStatus {
			return true
		}
	}

	return false
}

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

type InvitationSMSRequest struct {
	PrivateID string `json:"privateID"`
}

type CategoryCreateRequest struct {
	Tag string `json:"tag"`
}

type CategoryUpdateRequest struct {
	ID  int64  `json:"id"`
	Tag string `json:"tag"`
}

type Category struct {
	ID    int64  `json:"id"`
	Tag   string `json:"tag"`
	Total int    `json:"total"`
}

type BaseRSVP struct {
	FullName          string `json:"fullName"`
	Attending         bool   `json:"attending"`
	GuestCount        int    `json:"guestCount"`
	SpecialDiet       bool   `json:"specialDiet"`
	Remarks           string `json:"remarks"`
	MobilePhoneNumber string `json:"mobilePhoneNumber"`
}

type RSVPCreateRequest struct {
	BaseRSVP
	InvitationPrivateID string `json:"invitationPrivateID"`
	ReCAPTCHAToken      string `json:"reCAPTCHA"`
}

type RSVPUpdateRequest struct {
	BaseRSVP
	ID                  int64  `json:"id"`
	InvitationPrivateID string `json:"invitationPrivateID"`
}

type RSVP struct {
	BaseRSVP
	ID                  int64  `json:"id,omitempty"`
	InvitationPrivateID string `json:"invitationPrivateID,omitempty"`
	Completed           bool   `json:"completed"`
	UpdatedAt           string `json:"updatedAt"`
}
