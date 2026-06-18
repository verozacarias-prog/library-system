package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/verozacarias-prog/library-system/loans-service/internal/model"
	"github.com/verozacarias-prog/library-system/loans-service/internal/handler"
)

func TestCreateLoanHandler_HappyPath(t *testing.T) {
	mockSvc := &mockLoanService{}
	h := handler.NewLoanHandler(mockSvc)

	req := model.CreateLoanRequest{UserID: 1, BookID: 2}
	expectedLoan := &model.Loan{ID: 1, UserID: 1, BookID: 2, Status: "active"}

	mockSvc.On("CreateLoan", mock.Anything, req).Return(expectedLoan, nil)

	body, _ := json.Marshal(req)
	r := httptest.NewRequest(http.MethodPost, "/loans", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.CreateLoan(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCreateLoanHandler_InvalidBody(t *testing.T) {
	mockSvc := &mockLoanService{}
	h := handler.NewLoanHandler(mockSvc)

	r := httptest.NewRequest(http.MethodPost, "/loans", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	h.CreateLoan(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "CreateLoan")
}