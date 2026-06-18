package service

import (
	"context"

	"github.com/verozacarias-prog/library-system/loans-service/internal/model"
	"github.com/verozacarias-prog/library-system/loans-service/internal/repository"
)

type BookService interface {
	ValidateBook(ctx context.Context, bookID int) (*model.Book, error)
	UpdateCopies(ctx context.Context, bookID int, action string) error
}

type LoanService struct {
	repo *repository.LoanRepository
}

func NewLoanService(repo *repository.LoanRepository) *LoanService {
	return &LoanService{repo: repo}
}

func (s *LoanService) CreateLoan(ctx context.Context, req model.CreateLoanRequest) (*model.Loan, error) {
	return s.repo.Create(ctx, req)
}

func (s *LoanService) BookReturned(ctx context.Context, loanID int) (*model.Loan, error) {
	return s.repo.UpdateStatus(ctx, loanID, "returned")
}

func (s *LoanService) GetActiveLoans(ctx context.Context, userID int) ([]model.Loan, error) {
	return s.repo.GetActiveByUser(ctx, userID)
}

func (s *LoanService) GetLoanHistory(ctx context.Context, userID int) ([]model.Loan, error) {
	return s.repo.GetHistoryByUser(ctx, userID)
}
