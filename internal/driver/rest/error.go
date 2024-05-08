package rest

import (
	"net/http"
	"time"
)

type Error struct {
	Ok      bool   `json:"ok"`
	Err     string `json:"err"`
	Message string `json:"message"`
	Ts      int64  `json:"ts"`
}

type InternalServerError struct {
	Status  int    `json:"status"`
	Err     string `json:"err"`
	Message string `json:"message"`
	Ts      int64  `json:"ts"`
}

func NewError(isOk bool, err string, msg string) *Error {
	return &Error{
		Ok:      isOk,
		Err:     err,
		Message: msg,
		Ts:      time.Now().Unix(),
	}
}

func NewBadRequest(msg string) *Error {
	return &Error{
		Ok:      false,
		Err:     "ERR_BAD_REQUEST",
		Message: msg,
		Ts:      time.Now().Unix(),
	}
}

func NewReadTimeout() *Error {
	return &Error{
		Ok:      false,
		Err:     "ERR_READ_TIMEOUT",
		Message: "timeout reached when reading payload from client",
		Ts:      time.Now().Unix(),
	}
}

func NewInvalidAccessToken() *Error {
	return &Error{
		Ok:      false,
		Err:     "ERR_INVALID_ACCESS_TOKEN",
		Message: "invalid access token",
		Ts:      time.Now().Unix(),
	}
}

func NewForbiddenAccess() *Error {
	return &Error{
		Ok:      false,
		Err:     "ERR_FORBIDDEN_ACCESS",
		Message: "user doesn't have enough authorization",
		Ts:      time.Now().Unix(),
	}
}

func NewNotFound() *Error {
	return &Error{
		Ok:      false,
		Err:     "ERR_NOT_FOUND",
		Message: "resource is not found",
		Ts:      time.Now().Unix(),
	}
}

func newInternalServerError(msg string) *InternalServerError {
	return &InternalServerError{
		Status:  http.StatusInternalServerError,
		Err:     "ERR_INTERNAL_ERROR",
		Message: msg,
		Ts:      time.Now().Unix(),
	}
}
