package clients

import "errors"

var (
	ErrBookNotFound              = errors.New("book not found")
	ErrNoCopiesAvailable         = errors.New("no copies available")
	ErrLibraryServiceUnavailable = errors.New("library service unavailable")
)
