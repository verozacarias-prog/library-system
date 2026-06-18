package main

import (
	"context"
	"log"
	"net/http"

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

	log.Println("loans-service running on :8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
