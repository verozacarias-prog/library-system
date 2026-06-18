package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(loanHandler *LoanHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", Health)

	r.Route("/loans", func(r chi.Router) {
		r.Post("/", loanHandler.CreateLoan)
		r.Patch("/{id}", loanHandler.ReturnLoan)
		r.Get("/users/{userID}", loanHandler.GetActiveLoans)
		r.Get("/users/{userID}/history", loanHandler.GetLoanHistory)
	})

	return r
}