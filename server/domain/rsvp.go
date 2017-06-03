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
