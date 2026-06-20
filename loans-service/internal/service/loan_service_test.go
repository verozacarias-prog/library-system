package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/verozacarias-prog/library-system/loans-service/internal/clients"
	"github.com/verozacarias-prog/library-system/loans-service/internal/model"
	"github.com/verozacarias-prog/library-system/loans-service/internal/service"
)

func TestCreateLoan_HappyPath(t *testing.T) {
	ctx := context.Background()

	mockClient := &mockLibraryClient{}
	mockRepo := &mockLoanRepository{}
	svc := service.NewLoanService(mockRepo, mockClient)

	req := model.CreateLoanRequest{UserID: 1, BookID: 1}
	expectedLoan := &model.Loan{ID: 1, UserID: 1, BookID: 1, Status: service.StatusActive}
	book := &model.Book{ID: 1, AvailableCopies: 3}

	mockClient.On("ValidateBook", ctx, 1).Return(book, nil)
	mockRepo.On("Create", ctx, req).Return(expectedLoan, nil)
	mockClient.On("UpdateCopies", ctx, 1, service.CopiesActionDecrement).Return(nil)

	loan, err := svc.CreateLoan(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedLoan, loan)
	mockClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestCreateLoan_BookNotAvailable(t *testing.T) {
	ctx := context.Background()

	mockClient := &mockLibraryClient{}
	mockRepo := &mockLoanRepository{}
	svc := service.NewLoanService(mockRepo, mockClient)

	req := model.CreateLoanRequest{UserID: 1, BookID: 1}

	mockClient.On("ValidateBook", ctx, 1).Return(nil, clients.ErrNoCopiesAvailable)

	loan, err := svc.CreateLoan(ctx, req)

	assert.Nil(t, loan)
	assert.ErrorIs(t, err, clients.ErrNoCopiesAvailable)
	mockRepo.AssertNotCalled(t, "Create")
	mockClient.AssertExpectations(t)
}

func TestBookReturned_HappyPath(t *testing.T) {
	ctx := context.Background()

	mockClient := &mockLibraryClient{}
	mockRepo := &mockLoanRepository{}
	svc := service.NewLoanService(mockRepo, mockClient)

	activeLoan := &model.Loan{ID: 1, UserID: 1, BookID: 2, Status: service.StatusActive}
	returnedLoan := &model.Loan{ID: 1, UserID: 1, BookID: 2, Status: service.StatusReturned}

	mockRepo.On("GetByID", ctx, 1).Return(activeLoan, nil)
	mockClient.On("UpdateCopies", ctx, 2, service.CopiesActionIncrement).Return(nil)
	mockRepo.On("UpdateStatus", ctx, 1, service.StatusReturned).Return(returnedLoan, nil)

	loan, err := svc.BookReturned(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, service.StatusReturned, loan.Status)
	mockClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestGetActiveLoans_ReturnsOnlyActiveLoans(t *testing.T) {
	ctx := context.Background()

	mockClient := &mockLibraryClient{}
	mockRepo := &mockLoanRepository{}
	svc := service.NewLoanService(mockRepo, mockClient)

	expectedLoans := []model.Loan{
		{ID: 1, UserID: 5, BookID: 1, Status: service.StatusActive},
		{ID: 2, UserID: 5, BookID: 3, Status: service.StatusActive},
	}

	mockRepo.On("GetActiveByUser", ctx, 5).Return(expectedLoans, nil)

	loans, err := svc.GetActiveLoans(ctx, 5)

	assert.NoError(t, err)
	assert.Len(t, loans, 2)
	for _, l := range loans {
		assert.Equal(t, "active", l.Status)
	}
	mockRepo.AssertExpectations(t)
}
