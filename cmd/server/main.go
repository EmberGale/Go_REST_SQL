package main

import (
	"GoRestSQL/internal/handler"
	"GoRestSQL/internal/repository"
	"GoRestSQL/internal/service"
	"fmt"
	"net/http"
	"os"

	_ "github.com/glebarez/go-sqlite"
)

func main() {

	paymentRepo, err := repository.NewSqlitePaymentRepository()
	if err != nil {
		panic(err)
	}

	defer paymentRepo.DB.Close()
	paymentService := service.NewPaymentServiceImpl(paymentRepo)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	router := handler.NewRouter(paymentHandler)

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = fmt.Sprintf(":%s", port)

	if err := http.ListenAndServe(port, router); err != nil {
		panic(err)
	}
}
