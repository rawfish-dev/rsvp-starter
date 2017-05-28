package postgres

import (
	"database/sql"
	"strings"
)

func isNotFoundError(err error) bool {
	return err != nil && err == sql.ErrNoRows
}

func isCategoryTagUniqueConstraintError(err error) bool {
	return strings.Contains(err.Error(), `duplicate key value violates unique constraint "unique_tag"`)
}

func isInvitationGreetingUniqueConstraintError(err error) bool {
	return strings.Contains(err.Error(), `duplicate key value violates unique constraint "unique_greeting"`)
}

func isInvitationMobilePhoneNumberUniqueConstraintError(err error) bool {
	return strings.Contains(err.Error(), `duplicate key value violates unique constraint "unique_mobile_phone_number"`)
}

func isRSVPPrivateIDUniqueConstraintError(err error) bool {
	return strings.Contains(err.Error(), `duplicate key value violates unique constraint "unique_invitation_private_id"`)
}
