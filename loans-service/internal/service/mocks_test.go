package service_test

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/verozacarias-prog/library-system/loans-service/internal/model"
)

// Mock del BookService
type mockLibraryClient struct {
	mock.Mock
}

func (m *mockLibraryClient) ValidateBook(ctx context.Context, bookID int) (*model.Book, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Book), args.Error(1)
}

func (m *mockLibraryClient) UpdateCopies(ctx context.Context, bookID int, action string) error {
	args := m.Called(ctx, bookID, action)
	return args.Error(0)
}

// Mock del LoanRepository
type mockLoanRepository struct {
	mock.Mock
}

func (m *mockLoanRepository) Create(ctx context.Context, req model.CreateLoanRequest) (*model.Loan, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Loan), args.Error(1)
}

func (m *mockLoanRepository) UpdateStatus(ctx context.Context, loanID int, status string) (*model.Loan, error) {
	args := m.Called(ctx, loanID, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Loan), args.Error(1)
}

func (m *mockLoanRepository) GetActiveByUser(ctx context.Context, userID int) ([]model.Loan, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Loan), args.Error(1)
}

func (m *mockLoanRepository) GetHistoryByUser(ctx context.Context, userID int) ([]model.Loan, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Loan), args.Error(1)
}

func (m *mockLoanRepository) GetByID(ctx context.Context, loanID int) (*model.Loan, error) {
	args := m.Called(ctx, loanID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Loan), args.Error(1)
}
