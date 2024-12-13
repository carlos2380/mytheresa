package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)

	response := map[string]interface{}{
		"code":  e.Code,
		"error": e.Message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func Unwrap(err error) {
	log.Println("Error occurred:")
	unwrappedErr := err
	for unwrappedErr != nil {
		log.Printf(" - %v", unwrappedErr)
		unwrappedErr = errors.Unwrap(unwrappedErr)
	}
}
