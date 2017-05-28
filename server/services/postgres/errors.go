package postgres

type PostgresOperationError struct {
}

func NewPostgresOperationError() error {
	return PostgresOperationError{}
}

func (p PostgresOperationError) Error() string {
	return ""
}

type PostgresRecordNotFoundError struct {
}

func NewPostgresRecordNotFoundError() error {
	return PostgresRecordNotFoundError{}
}

func (p PostgresRecordNotFoundError) Error() string {
	return ""
}

type PostgresCategoryTagUniqueConstraintError struct {
}

func NewPostgresCategoryTagUniqueConstraintError() error {
	return PostgresCategoryTagUniqueConstraintError{}
}

func (p PostgresCategoryTagUniqueConstraintError) Error() string {
	return ""
}

type PostgresInvitationGreetingUniqueConstraintError struct {
}

func NewPostgresInvitationGreetingUniqueConstraintError() error {
	return PostgresInvitationGreetingUniqueConstraintError{}
}

func (p PostgresInvitationGreetingUniqueConstraintError) Error() string {
	return "greeting already exists"
}

type PostgresInvitationMobilePhoneNumberUniqueConstraintError struct {
}

func NewPostgresInvitationMobilePhoneNumberUniqueConstraintError() error {
	return PostgresInvitationMobilePhoneNumberUniqueConstraintError{}
}

func (p PostgresInvitationMobilePhoneNumberUniqueConstraintError) Error() string {
	return "mobile phone number already exists"
}

type PostgresRSVPPrivateIDUniqueConstraintError struct {
}

func NewPostgresRSVPPrivateIDUniqueConstraintError() error {
	return PostgresRSVPPrivateIDUniqueConstraintError{}
}

func (p PostgresRSVPPrivateIDUniqueConstraintError) Error() string {
	return "rsvp already exists for invitation"
}
