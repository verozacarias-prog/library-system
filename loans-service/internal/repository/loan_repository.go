package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/verozacarias-prog/library-system/loans-service/internal/model"
)

type LoanRepository struct {
	pool *pgxpool.Pool
}

func NewLoanRepository(pool *pgxpool.Pool) *LoanRepository {
	return &LoanRepository{pool: pool}
}

func (r *LoanRepository) Create(ctx context.Context, req model.CreateLoanRequest) (*model.Loan, error) {
	loan := &model.Loan{}
	err := r.pool.QueryRow(ctx,
		QueryCreateLoan,
		req.UserID, req.BookID, time.Now().UTC(),
	).Scan(&loan.ID, &loan.UserID, &loan.BookID, &loan.LoanedAt, &loan.ReturnedAt, &loan.Status)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return nil, ErrLoanAlreadyActive
		}
		return nil, err
	}
	return loan, err
}

func (r *LoanRepository) UpdateStatus(ctx context.Context, loanID int, status string) (*model.Loan, error) {
	loan := &model.Loan{}
	err := r.pool.QueryRow(ctx,
		QueryUpdateStatus,
		time.Now(), status, loanID,
	).Scan(&loan.ID, &loan.UserID, &loan.BookID, &loan.LoanedAt, &loan.ReturnedAt, &loan.Status)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrLoanNotFound
	}
	return loan, err
}

func (r *LoanRepository) GetActiveByUser(ctx context.Context, userID int) ([]model.Loan, error) {
	rows, err := r.pool.Query(ctx, QueryGetActiveByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	loans := make([]model.Loan, 0)
	for rows.Next() {
		var l model.Loan
		if err := rows.Scan(&l.ID, &l.UserID, &l.BookID, &l.LoanedAt, &l.ReturnedAt, &l.Status); err != nil {
			return nil, err
		}
		loans = append(loans, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return loans, nil
}

func (r *LoanRepository) GetHistoryByUser(ctx context.Context, userID int) ([]model.Loan, error) {
	rows, err := r.pool.Query(ctx, QueryGetHistoryByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	loans := make([]model.Loan, 0)
	for rows.Next() {
		var l model.Loan
		if err := rows.Scan(&l.ID, &l.UserID, &l.BookID, &l.LoanedAt, &l.ReturnedAt, &l.Status); err != nil {
			return nil, err
		}
		loans = append(loans, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return loans, nil
}

func (r *LoanRepository) GetByID(ctx context.Context, loanID int) (*model.Loan, error) {
	loan := &model.Loan{}
	err := r.pool.QueryRow(ctx, QueryGetByID, loanID).
		Scan(&loan.ID, &loan.UserID, &loan.BookID, &loan.LoanedAt, &loan.ReturnedAt, &loan.Status)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrLoanNotFound
	}
	return loan, err
}
