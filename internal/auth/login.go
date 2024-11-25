package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/MarcusXavierr/api-banco-cooperativa/internal/db"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

// TODO: Melhorar a validação de erros a nivel dos handlers
func (service AuthService) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Email string `json:"email" validate:"required,email"`
		Senha string `json:"senha" validate:"required,min=8"`
	}
	conn := service.DB.Conn
	hash := sha256.New()

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		log.Print(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(loginRequest); err != nil {
		log.Print(err)
		http.Error(w, "Validation failed\n"+err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	fmt.Println(loginRequest)

	hash.Write([]byte(loginRequest.Senha))

	params := db.VerifyCredentialsParams{
		Email:    loginRequest.Email,
		Password: hex.EncodeToString(hash.Sum(nil)),
	}

	user, err := conn.VerifyCredentials(r.Context(), params)
	if err != nil {
		log.Print(err)
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Credenciais inválidas", 401)
		} else {
			http.Error(w, "Could not connect to database, try again later", 502)
		}
		return
	}

	token, err := service.generateTokenJWT(user)
	if err != nil {
		log.Print(err)
		http.Error(w, "Error generating the jwt token", 500)
		return
	}

	response := fmt.Sprintf("Bearer %s\n", token)
	w.Write([]byte(response))
}
