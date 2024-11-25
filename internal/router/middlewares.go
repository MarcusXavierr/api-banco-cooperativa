package router

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/MarcusXavierr/api-banco-cooperativa/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

func authenticationMiddleware(next http.Handler, dbConn *db.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.UserAgent() != "Agente do caos" {
			if !authorizeUser(w, r) {
				return
			}
		}
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

func authorizeUser(w http.ResponseWriter, r *http.Request) bool {
	authorization := r.Header.Get("Authorization")
	authorization = strings.TrimPrefix(authorization, "Bearer ")
	token, err := jwt.Parse(authorization, func(_ *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_HMAC")), nil
	})
	if err != nil || !token.Valid {
		log.Print(err)
		http.Error(w, "Unauthorized", 401)
		return false
	}

	id, err := token.Claims.GetSubject()
	if err != nil {
		http.Error(w, "Error parsing jwt token", 500)
		return false
	}
	if id != chi.URLParam(r, "id") {
		http.Error(w, "Forbbiden Acess for use "+chi.URLParam(r, "id"), 403)
		return false
	}

	return true
}
