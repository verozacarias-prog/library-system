package loans

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/verozacarias-prog/library-system/loans-service/internal/repository"
	"github.com/verozacarias-prog/library-system/loans-service/internal/service"
)

type App struct {
	Pool        *pgxpool.Pool
	LoanService *service.LoanService
}

func New(ctx context.Context, libraryClient service.BookService) (*App, error) {
	pool, err := repository.NewPool(ctx)
	if err != nil {
		return nil, err
	}

	repository := repository.NewLoanRepository(pool)
	svc := service.NewLoanService(repository, libraryClient)

	return &App{
		Pool:        pool,
		LoanService: svc,
	}, nil
}

func (a *App) Close() {
	a.Pool.Close()
}
