package router

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/MarcusXavierr/rinha-de-backend-2024-q1/internal/db"
	"github.com/MarcusXavierr/rinha-de-backend-2024-q1/internal/user"
	"github.com/go-chi/chi/v5"
)

const UserCtxKey = "user"

func initializeRoutes(router *chi.Mux, dbConn *db.Queries, dbTx user.DBTransactions) {
	userService := user.UserService{DBConn: dbConn, DBTransactions: dbTx}

	router.Route("/clientes/{id}", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return userCtx(next, dbConn)
		})

		r.Get("/extrato", userService.HandleExtract)
		r.Post("/transacoes", userService.HandleBalanceMovement)
	})
}

func userCtx(next http.Handler, dbConn *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clienteIDStr := chi.URLParam(r, "id")
		clienteID, err := strconv.Atoi(clienteIDStr)
		if clienteIDStr == "" || err != nil {
			log.Printf("Error reading ID: %v", err)
			http.Error(w, http.StatusText(422), 404)
		}

		user, err := dbConn.GetUser(r.Context(), int32(clienteID))
		if err != nil {
			log.Printf("Couldn't find user from ID: %s | Error: %v\n", clienteIDStr, err)
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), UserCtxKey, &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
