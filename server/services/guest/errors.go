package guest

type CategoryNotFoundError struct {
}

func NewCategoryNotFoundError() error {
	return CategoryNotFoundError{}
}

func (c CategoryNotFoundError) Error() string {
	return "category not found"
}

type InvitationNotFoundError struct {
}

func NewInvitationNotFoundError() error {
	return InvitationNotFoundError{}
}

func (i InvitationNotFoundError) Error() string {
	return "invitation not found"
}

type RSVPNotFoundError struct {
}

func NewRSVPNotFoundError() error {
	return RSVPNotFoundError{}
}

func (r RSVPNotFoundError) Error() string {
	return "rsvp not found"
}
