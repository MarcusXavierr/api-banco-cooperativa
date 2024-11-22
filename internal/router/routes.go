package router

import (
	"net/http"

	"github.com/MarcusXavierr/rinha-de-backend-2024-q1/internal/db"
	"github.com/MarcusXavierr/rinha-de-backend-2024-q1/internal/user"
	"github.com/go-chi/chi/v5"
)

const UserCtxKey = "user"

func initializeRoutes(router *chi.Mux, dbConn *db.DBPool) {
	userService := user.UserService{DB: dbConn}

	router.Route("/clientes/{id}", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return fillUserDataMiddleware(next, dbConn.Conn)
		})

		r.Get("/extrato", userService.HandleExtract)
		r.Post("/transacoes", userService.HandleTransaction)
	})
}
