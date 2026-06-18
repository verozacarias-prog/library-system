package handler_test

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/verozacarias-prog/library-system/loans-service/internal/model"
)

type mockLoanService struct {
	mock.Mock
}

func (m *mockLoanService) CreateLoan(ctx context.Context, req model.CreateLoanRequest) (*model.Loan, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Loan), args.Error(1)
}

func (m *mockLoanService) BookReturned(ctx context.Context, loanID int) (*model.Loan, error) {
	args := m.Called(ctx, loanID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Loan), args.Error(1)
}

func (m *mockLoanService) GetActiveLoans(ctx context.Context, userID int) ([]model.Loan, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Loan), args.Error(1)
}

func (m *mockLoanService) GetLoanHistory(ctx context.Context, userID int) ([]model.Loan, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Loan), args.Error(1)
}
