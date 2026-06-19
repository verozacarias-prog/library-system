package model

import "time"

type Loan struct {
	ID         int        `json:"id"`
	UserID     int        `json:"user_id"`
	BookID     int        `json:"book_id"`
	LoanedAt   time.Time  `json:"loaned_at"`
	ReturnedAt *time.Time `json:"returned_at,omitempty"`
	Status     string     `json:"status"` // "active" | "returned"
}

type Book struct {
	ID              int `json:"id"`
	AvailableCopies int `json:"availableCopies"`
}

type CreateLoanRequest struct {
	UserID int `json:"user_id"`
	BookID int `json:"book_id"`
}
