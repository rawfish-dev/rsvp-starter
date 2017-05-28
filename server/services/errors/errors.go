package errors

import (
	"strings"
)

type GeneralServiceError struct {
	errorMessage string
}

func NewGeneralServiceError() error {
	return GeneralServiceError{}
}

func (g GeneralServiceError) Error() string {
	return ""
}

type ValidationError struct {
	errorMessages []string
}

func NewValidationError(errorMessages []string) error {
	return ValidationError{errorMessages}
}

func (v ValidationError) Error() (fullErrorMessage string) {
	return strings.Join(v.errorMessages, "; ")
}
