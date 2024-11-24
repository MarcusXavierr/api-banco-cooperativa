package router

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/MarcusXavierr/api-banco-cooperativa/internal/db"
	"github.com/go-chi/chi/v5"
)

func fillUserDataMiddleware(next http.Handler, dbConn *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDRaw := chi.URLParam(r, "id")
		userID, err := strconv.Atoi(userIDRaw)
		if userIDRaw == "" || err != nil {
			log.Printf("Error reading ID: %v", err)
			http.Error(w, "Error reading ID: "+userIDRaw, 500)
		}

		user, err := dbConn.GetUser(r.Context(), int32(userID))
		if err != nil {
			log.Printf("Couldn't find user from ID: %s | Error: %v\n", userIDRaw, err)
			http.Error(w, http.StatusText(404), 404)
			return
		}

		ctx := context.WithValue(r.Context(), UserCtxKey, &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
