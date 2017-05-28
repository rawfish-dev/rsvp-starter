package domain

import (
	"net/http"
)

type CustomBadRequestError struct {
	Error string `json:"error"`
}

func NewCustomBadRequestError(errorMessage string) (int, interface{}) {
	return http.StatusBadRequest, CustomBadRequestError{errorMessage}
}
