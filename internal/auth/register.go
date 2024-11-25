package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/MarcusXavierr/api-banco-cooperativa/internal/db"
)

func (service AuthService) HandleRegister(w http.ResponseWriter, r *http.Request) {
	conn := service.DB.Conn
	err := r.ParseForm()
	if err != nil {
		log.Print(err)
		http.Error(w, "Could not parse form input", 400)
		return
	}

	hash := sha256.New()
	hash.Write([]byte(r.PostFormValue("password")))

	credit := r.PostFormValue("credit_limit")
	if credit == "" {
		credit = "0"
	}

	limit, err := strconv.Atoi(credit)
	if err != nil {
		log.Print(err)
		http.Error(w, "credit_limit must be a valid number", 500)
	}

	params := db.CreateUserParams{
		Email:       r.PostFormValue("email"),
		Password:    hex.EncodeToString(hash.Sum(nil)),
		Name:        r.PostFormValue("name"),
		CreditLimit: int32(limit),
	}

	if err := conn.CreateUser(r.Context(), params); err != nil {
		log.Print(err)
		http.Error(w, "Error creating the user", 500)
		return
	}

	user, err := conn.RetrieveUserFromEmail(r.Context(), r.PostFormValue("email"))
	if err != nil {
		log.Print(err)
		http.Error(w, "Error retreaving user from database", 500)
		return
	}

	token, err := service.generateTokenJWT(user)
	if err != nil {
		log.Print(err)
		http.Error(w, "Error generating JWT", 500)
		return
	}

	response := map[string]interface{}{
		"user":  user,
		"token": "Bearer " + token,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Print(err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
