package service

import "errors"

var (
	ErrLoanInactive = errors.New("Book already returned")
)
