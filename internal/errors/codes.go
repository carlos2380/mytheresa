package errors

import "net/http"

var (
	ErrInternalServerError = &HttpError{http.StatusInternalServerError, "internal server error", nil}
	ErrNotFound            = &HttpError{http.StatusNotFound, "resource not found", nil}
	ErrBadRequest          = &HttpError{http.StatusBadRequest, "bad request", nil}
	ErrUnauthorized        = &HttpError{http.StatusUnauthorized, "unauthorized", nil}
	ErrForbidden           = &HttpError{http.StatusForbidden, "forbidden", nil}
	ErrIntConvert          = &HttpError{http.StatusInternalServerError, "unsupported int", nil}
	ErrMethodNotAllowed    = &HttpError{http.StatusMethodNotAllowed, "method not allowed", nil}
)
