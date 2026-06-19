package service

import "errors"

var (
	ErrInvalidStatus             = errors.New("invalid status")
	ErrBookNotFound              = errors.New("book not found")
	ErrNoCopiesAvailable         = errors.New("no copies available")
	ErrLibraryServiceUnavailable = errors.New("library service unavailable")
	ErrUpdateCopies              = errors.New("failed to update copies")
	ErrLoanInactive              = errors.New("loan inactive")
)
