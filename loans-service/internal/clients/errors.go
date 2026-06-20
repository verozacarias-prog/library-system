package clients

import "errors"

var (
	ErrBookNotFound              = errors.New("Book not found")
	ErrNoCopiesAvailable         = errors.New("No copies available")
	ErrLibraryServiceUnavailable = errors.New("Library service unavailable")
)
