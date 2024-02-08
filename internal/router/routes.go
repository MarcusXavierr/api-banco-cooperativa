package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func initializeRoutes(router *chi.Mux) {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello friend\n"))
	})
}
