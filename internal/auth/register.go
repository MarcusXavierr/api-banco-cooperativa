package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/MarcusXavierr/api-banco-cooperativa/internal/db"
	"github.com/go-playground/validator/v10"
)

func (service AuthService) HandleRegister(w http.ResponseWriter, r *http.Request) {
	conn := service.DB.Conn
	var registerRequest struct {
		Name        string `json:"name" validate:"required"`
		CreditLimit int32  `json:"credit_limit" validate:"numeric,min=0"`
		Email       string `json:"email" validate:"required,email"`
		Senha       string `json:"senha" validate:"required,min=8"`
	}

	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		log.Print(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(registerRequest); err != nil {
		log.Print(err)
		http.Error(w, "Validation failed\n"+err.Error(), http.StatusBadRequest)
		return
	}

	hash := sha256.New()
	hash.Write([]byte(registerRequest.Senha))

	params := db.CreateUserParams{
		Email:       registerRequest.Email,
		Password:    hex.EncodeToString(hash.Sum(nil)),
		Name:        registerRequest.Name,
		CreditLimit: registerRequest.CreditLimit,
	}

	if err := conn.CreateUser(r.Context(), params); err != nil {
		log.Print(err)
		http.Error(w, "Error creating the user", 500)
		return
	}

	user, err := conn.RetrieveUserFromEmail(r.Context(), registerRequest.Email)
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
