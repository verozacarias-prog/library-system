package repository

import "errors"

const pgUniqueViolation = "23505"

var (
	// Database errors
	ErrDatabaseURLNotSet   = errors.New("DATABASE_URL not set")
	ErrConnectionPool      = errors.New("unable to create connection pool")
	ErrDatabaseUnreachable = errors.New("unable to reach database")
	// Loan errors
	ErrLoanNotFound      = errors.New("loan not found or already returned")
	ErrLoanAlreadyActive = errors.New("user already has an active loan for this book")
)
