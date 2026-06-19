package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	loans "github.com/verozacarias-prog/library-system/loans-service/internal"
	"github.com/verozacarias-prog/library-system/loans-service/internal/clients"
	"github.com/verozacarias-prog/library-system/loans-service/internal/handler"
)

func main() {
	ctx := context.Background()

	libraryClient := clients.NewLibraryClient()

	app, err := loans.New(ctx, libraryClient)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	defer app.Close()

	loanHandler := handler.NewLoanHandler(app.LoanService)
	router := handler.NewRouter(loanHandler)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("loans-service running on :8081")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("shutdown complete")
}
