package router

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func Initialize() {
	router := chi.NewRouter()
	// router.Use(middleware.Logger)

	initializeRoutes(router)

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
