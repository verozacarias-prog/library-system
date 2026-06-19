package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/verozacarias-prog/library-system/loans-service/internal/repository"
	"github.com/verozacarias-prog/library-system/loans-service/internal/service"
)

var (
	ErrInvalidRequestBody  = errors.New("invalid request body")
	ErrInvalidLoanID       = errors.New("invalid loan id")
	ErrInvalidUserID       = errors.New("invalid user id")
	ErrInvalidUserOrBookID = errors.New("user_id and book_id must be positive integers")
)

func writeErrorJSON(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrBookNotFound):
		writeErrorJSON(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, service.ErrNoCopiesAvailable):
		writeErrorJSON(w, err.Error(), http.StatusConflict)
	case errors.Is(err, repository.ErrLoanNotFound):
		writeErrorJSON(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, repository.ErrLoanAlreadyActive):
		writeErrorJSON(w, err.Error(), http.StatusConflict)
	default:
		writeErrorJSON(w, err.Error(), http.StatusInternalServerError)
	}
}
