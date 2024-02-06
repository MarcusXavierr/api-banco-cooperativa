package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()
	PORT := os.Getenv("HTTP_PORT")

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello darkness my old friend"))
	})

	http.ListenAndServe(":"+PORT, router)
}
