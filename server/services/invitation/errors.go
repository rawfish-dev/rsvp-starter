package invitation

var _ error = new(InvitationNotFoundError)

type InvitationNotFoundError struct {
}

func NewInvitationNotFoundError() error {
	return InvitationNotFoundError{}
}

func (i InvitationNotFoundError) Error() string {
	return "invitation not found"
}
