package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/verozacarias-prog/library-system/loans-service/internal/model"
	"github.com/verozacarias-prog/library-system/loans-service/internal/repository"
	"github.com/verozacarias-prog/library-system/loans-service/internal/service"
)

type LoanService interface {
	CreateLoan(ctx context.Context, req model.CreateLoanRequest) (*model.Loan, error)
	BookReturned(ctx context.Context, loanID int) (*model.Loan, error)
	GetActiveLoans(ctx context.Context, userID int) ([]model.Loan, error)
	GetLoanHistory(ctx context.Context, userID int) ([]model.Loan, error)
}

type LoanHandler struct {
	service LoanService
}

func NewLoanHandler(svc LoanService) *LoanHandler {
	return &LoanHandler{service: svc}
}

func (h *LoanHandler) CreateLoan(w http.ResponseWriter, r *http.Request) {
	var req model.CreateLoanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	if req.UserID <= 0 || req.BookID <= 0 {
		http.Error(w, "user_id and book_id must be positive integers", http.StatusBadRequest)
		return
	}

	loan, err := h.service.CreateLoan(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(loan)
}

func (h *LoanHandler) ReturnLoan(w http.ResponseWriter, r *http.Request) {
	loanID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, ErrInvalidLoanID.Error(), http.StatusBadRequest)
		return
	}

	loan, err := h.service.BookReturned(r.Context(), loanID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loan)
}

func (h *LoanHandler) GetActiveLoans(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, ErrInvalidUserID.Error(), http.StatusBadRequest)
		return
	}

	loans, err := h.service.GetActiveLoans(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loans)
}

func (h *LoanHandler) GetLoanHistory(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, ErrInvalidUserID.Error(), http.StatusBadRequest)
		return
	}

	loans, err := h.service.GetLoanHistory(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loans)
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrBookNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, service.ErrNoCopiesAvailable):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, repository.ErrLoanNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
