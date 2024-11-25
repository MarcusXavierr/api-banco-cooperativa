package router

import (
	"net/http"

	"github.com/MarcusXavierr/api-banco-cooperativa/internal/auth"
	"github.com/MarcusXavierr/api-banco-cooperativa/internal/db"
	"github.com/MarcusXavierr/api-banco-cooperativa/internal/user"
	"github.com/go-chi/chi/v5"
)

const UserCtxKey = "user"

func initializeRoutes(router *chi.Mux, dbConn *db.DBPool) {
	userService := user.UserService{DB: dbConn}
	authService := auth.AuthService{DB: dbConn}

	router.Route("/clientes/{id}", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return authenticationMiddleware(next, dbConn.Conn)
		})

		r.Get("/extrato", userService.HandleExtract)
		r.Post("/transacoes", userService.HandleTransaction)
	})

	router.Post("/login", authService.HandleLogin)
	router.Post("/cadastro", authService.HandleRegister)

	router.Route("/teste", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("primeiro teste"))
		})
	})
}
