package domain

import (
	"net/http"
)

const (
	invalidJSONBodyMessage = "JSON request was invalid"
)

type CustomBadRequestError struct {
	Error string `json:"error"`
}

func NewCustomBadRequestError(errorMessage string) (int, interface{}) {
	return http.StatusBadRequest, CustomBadRequestError{errorMessage}
}

func NewInvalidJSONBodyError() (int, interface{}) {
	return http.StatusBadRequest, CustomBadRequestError{invalidJSONBodyMessage}
}
