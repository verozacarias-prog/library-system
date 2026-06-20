package service

import "errors"

var (
	ErrUpdateCopies              = errors.New("Failed to update copies")
	ErrLoanInactive              = errors.New("Book already returned")
)
