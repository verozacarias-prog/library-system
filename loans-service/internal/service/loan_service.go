package service

import (
	"context"
	"errors"

	"github.com/verozacarias-prog/library-system/loans-service/internal/clients"
	"github.com/verozacarias-prog/library-system/loans-service/internal/model"
)

type BookService interface {
	ValidateBook(ctx context.Context, bookID int) (*model.Book, error)
	UpdateCopies(ctx context.Context, bookID int, action string) error
}

type LoanRepository interface {
	Create(ctx context.Context, req model.CreateLoanRequest) (*model.Loan, error)
	GetByID(ctx context.Context, loanID int) (*model.Loan, error)
	UpdateStatus(ctx context.Context, loanID int, status string) (*model.Loan, error)
	GetActiveByUser(ctx context.Context, userID int) ([]model.Loan, error)
	GetHistoryByUser(ctx context.Context, userID int) ([]model.Loan, error)
}

type LoanService struct {
	repository    LoanRepository
	libraryClient BookService
}

func NewLoanService(repo LoanRepository, libraryClient BookService) *LoanService {
	return &LoanService{
		repository:    repo,
		libraryClient: libraryClient,
	}
}

func (s *LoanService) CreateLoan(ctx context.Context, req model.CreateLoanRequest) (*model.Loan, error) {
	if _, err := s.libraryClient.ValidateBook(ctx, req.BookID); err != nil {
		switch {
		case errors.Is(err, clients.ErrBookNotFound):
			return nil, ErrBookNotFound
		case errors.Is(err, clients.ErrNoCopiesAvailable):
			return nil, ErrNoCopiesAvailable
		default:
			return nil, ErrLibraryServiceUnavailable
		}
	}

	if err := s.libraryClient.UpdateCopies(ctx, req.BookID, CopiesActionDecrement); err != nil {
		return nil, ErrLibraryServiceUnavailable
	}

	loan, err := s.repository.Create(ctx, req)
	if err != nil {
		s.libraryClient.UpdateCopies(context.Background(), req.BookID, CopiesActionIncrement)
		return nil, err
	}
	return loan, nil
}

func (s *LoanService) BookReturned(ctx context.Context, loanID int) (*model.Loan, error) {
	loan, err := s.repository.GetByID(ctx, loanID)
	if err != nil {
		return nil, err
	}

	if loan.Status != StatusActive {
		return nil, ErrLoanInactive
	}

	if err := s.libraryClient.UpdateCopies(ctx, loan.BookID, CopiesActionIncrement); err != nil {
		return nil, ErrLibraryServiceUnavailable
	}

	returned, err := s.repository.UpdateStatus(ctx, loanID, StatusReturned)
	if err != nil {
		s.libraryClient.UpdateCopies(context.Background(), loan.BookID, CopiesActionDecrement)
		return nil, err
	}

	return returned, nil
}

func (s *LoanService) GetActiveLoans(ctx context.Context, userID int) ([]model.Loan, error) {
	return s.repository.GetActiveByUser(ctx, userID)
}

func (s *LoanService) GetLoanHistory(ctx context.Context, userID int) ([]model.Loan, error) {
	return s.repository.GetHistoryByUser(ctx, userID)
}
