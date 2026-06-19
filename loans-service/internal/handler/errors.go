package handler

import (
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

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrBookNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, service.ErrNoCopiesAvailable):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, repository.ErrLoanNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, repository.ErrLoanAlreadyActive):
		http.Error(w, err.Error(), http.StatusConflict)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
