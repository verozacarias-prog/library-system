package handler

import "errors"

var (
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrInvalidLoanID      = errors.New("invalid loan id")
	ErrInvalidUserID      = errors.New("invalid user id")
)
