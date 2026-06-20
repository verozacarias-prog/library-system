package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/verozacarias-prog/library-system/loans-service/internal/clients"
	"github.com/verozacarias-prog/library-system/loans-service/internal/repository"
)

var (
	ErrInvalidRequestBody  = errors.New("Invalid request body")
	ErrInvalidLoanID       = errors.New("Invalid loan id")
	ErrInvalidUserID       = errors.New("Invalid user id")
	ErrInvalidUserOrBookID = errors.New("User ID and book ID must be positive integers")
)

func writeErrorJSON(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"statusCode": code,
		"message":    msg,
		"error":      http.StatusText(code),
	})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, clients.ErrBookNotFound):
		writeErrorJSON(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, clients.ErrNoCopiesAvailable):
		writeErrorJSON(w, err.Error(), http.StatusConflict)
	case errors.Is(err, repository.ErrLoanNotFound):
		writeErrorJSON(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, repository.ErrLoanAlreadyActive):
		writeErrorJSON(w, err.Error(), http.StatusConflict)
	default:
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
	}
}
