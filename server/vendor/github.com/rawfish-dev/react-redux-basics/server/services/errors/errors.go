package errors

import (
	"strings"
)

var _ error = new(GeneralServiceError)

type GeneralServiceError struct {
	errorMessage string
}

func NewGeneralServiceError() GeneralServiceError {
	return GeneralServiceError{}
}

func (g GeneralServiceError) Error() string {
	return ""
}

var _ error = new(ValidationError)

type ValidationError struct {
	errorMessages []string
}

func NewValidationError(errorMessages []string) ValidationError {
	return ValidationError{errorMessages}
}

func (v ValidationError) Error() (fullErrorMessage string) {
	return strings.Join(v.errorMessages, "; ")
}
