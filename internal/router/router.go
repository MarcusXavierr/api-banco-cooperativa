package router

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/MarcusXavierr/rinha-de-backend-2024-q1/internal/db"
	"github.com/go-chi/chi/v5"
)

func Initialize(dbConn *db.Queries) {
	router := chi.NewRouter()

	initializeRoutes(router, dbConn)

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	address := "0.0.0.0:" + port
	log.Println("Running on port:", address)

	if err := http.ListenAndServe(address, router); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP Server error: %v", err)
	}
}
