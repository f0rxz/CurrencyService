package main

import (
	"currencyservice/internal/controller/httpservice"
	"currencyservice/internal/repo"
	"currencyservice/internal/repo/currencies"
	"currencyservice/internal/usecase/exchangerate"
	"fmt"
	"log"
)

func main() {
	db, err := repo.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	repo := currencies.NewRepo(db)

	exchangeUsecase := exchangerate.NewUsecase(repo)

	server := httpservice.NewServer(exchangeUsecase)

	server.SetupRoutes()

	fmt.Println("Server is running on port 8080")
	if err := server.Start(8080); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
