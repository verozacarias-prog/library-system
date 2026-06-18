package service

import (
	"context"
	"fmt"

	"github.com/verozacarias-prog/library-system/loans-service/internal/model"
	"github.com/verozacarias-prog/library-system/loans-service/internal/repository"
)

type BookService interface {
	ValidateBook(ctx context.Context, bookID int) (*model.Book, error)
	UpdateCopies(ctx context.Context, bookID int, action string) error
}

type LoanService struct {
	repository    *repository.LoanRepository
	libraryClient BookService
}

func NewLoanService(repo *repository.LoanRepository, libraryClient BookService) *LoanService {
	return &LoanService{
		repository:    repo,
		libraryClient: libraryClient,
	}
}

func (s *LoanService) CreateLoan(ctx context.Context, req model.CreateLoanRequest) (*model.Loan, error) {
	if _, err := s.libraryClient.ValidateBook(ctx, req.BookID); err != nil {
		return nil, err
	}

	loan, err := s.repository.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := s.libraryClient.UpdateCopies(ctx, req.BookID, CopiesActionDecrement); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUpdateCopies, err)
	}
	return loan, nil
}

func (s *LoanService) BookReturned(ctx context.Context, loanID int) (*model.Loan, error) {
	loan, err := s.repository.UpdateStatus(ctx, loanID, StatusReturned)
	if err != nil {
		return nil, err
	}

	if err := s.libraryClient.UpdateCopies(ctx, loan.BookID, "increment"); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUpdateCopies, err)
	}

	return loan, nil
}

func (s *LoanService) GetActiveLoans(ctx context.Context, userID int) ([]model.Loan, error) {
	return s.repository.GetActiveByUser(ctx, userID)
}

func (s *LoanService) GetLoanHistory(ctx context.Context, userID int) ([]model.Loan, error) {
	return s.repository.GetHistoryByUser(ctx, userID)
}
