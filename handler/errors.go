package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/dgparker/go-contacts/services"
)

var (
	// ErrNoIDParam message for invalid ID parameter
	ErrNoIDParam = errors.New("Invalid ID parameter")

	// ErrNoID message for no id on req body
	ErrNoID = errors.New("Missing field: _id")

	//ErrInvalidFileType message for POST CSV if file is not of type text/csv
	ErrInvalidFileType = errors.New("Invalid file type")
)

// errorResponse defines our API error messages
type errorResponse struct {
	Err string `json:"err,omitempty"`
}

type postCSVError struct {
	Err            string            `json:"err,omitempty"`
	InvalidEntries []*services.Entry `json:"invalid_entries,omitempty"`
}

// Error writes an API error message to the response and logger.
func Error(w http.ResponseWriter, err error, code int, logger *log.Logger) {
	logger.Printf("http error: %s (code=%d)", err, code)

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&errorResponse{Err: err.Error()})
}
