package loans

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/verozacarias-prog/library-system/loans-service/internal/repository"
)

type App struct {
	Pool           *pgxpool.Pool
	LoanRepository *repository.LoanRepository
}

func New(ctx context.Context) (*App, error) {
	pool, err := repository.NewPool(ctx)
	if err != nil {
		return nil, err
	}

	return &App{
		Pool:           pool,
		LoanRepository: repository.NewLoanRepository(pool),
	}, nil
}

func (a *App) Close() {
	a.Pool.Close()
}
