package errors

import (
	"fmt"
	"net/http"
)

type HttpError struct {
	Code    int
	Message string
	Err     error
}

func (e *HttpError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("status %d: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("status %d: %s", e.Code, e.Message)
}

func (e *HttpError) Respond(w http.ResponseWriter) {
	http.Error(w, e.Error(), e.Code)
}

func Wrap(err error, httpError HttpError) *HttpError {
	httpError.Err = err
	return &httpError
}
