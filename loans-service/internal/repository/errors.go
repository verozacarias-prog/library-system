package repository

import "errors"

var (
	// Database errors
	ErrDatabaseURLNotSet   = errors.New("DATABASE_URL not set")
	ErrConnectionPool      = errors.New("unable to create connection pool")
	ErrDatabaseUnreachable = errors.New("unable to reach database")
	// Loan errors
	ErrLoanNotFound = errors.New("loan not found or already returned")
)
