package repository

import "errors"

const pgUniqueViolation = "23505"

var (
	// Database errors
	ErrDatabaseURLNotSet   = errors.New("DATABASE_URL not set")
	ErrConnectionPool      = errors.New("Unable to create connection pool")
	ErrDatabaseUnreachable = errors.New("Unable to reach database")
	// Loan errors
	ErrLoanNotFound      = errors.New("Loan not found or already returned")
	ErrLoanAlreadyActive = errors.New("User already has an active loan for this book")
)
