package rsvp

var _ error = new(RSVPNotFoundError)

type RSVPNotFoundError struct {
}

func NewRSVPNotFoundError() error {
	return RSVPNotFoundError{}
}

func (r RSVPNotFoundError) Error() string {
	return "rsvp not found"
}
