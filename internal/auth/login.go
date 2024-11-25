package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MarcusXavierr/api-banco-cooperativa/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

// TODO: Melhorar a validação de erros a nivel dos handlers
func (service AuthService) HandleLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	conn := service.DB.Conn
	fmt.Print("oi", r.PostFormValue("password"))
	if err != nil {
		log.Print(time.Now(), err)
		http.Error(w, "Could not parse post form", 500)
		return
	}

	hash := sha256.New()
	hash.Write([]byte(r.PostFormValue("password")))

	params := db.VerifyCredentialsParams{
		Email:    r.PostFormValue("email"),
		Password: hex.EncodeToString(hash.Sum(nil)),
	}

	user, err := conn.VerifyCredentials(r.Context(), params)
	if err != nil {
		log.Print(time.Now(), err)
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Credenciais inválidas", 401)
		} else {
			http.Error(w, "Could not connect to database, try again later", 502)
		}
		return
	}

	token, err := service.generateTokenJWT(user)
	if err != nil {
		log.Print(time.Now(), err)
		http.Error(w, "Error generating the jwt token", 500)
		return
	}

	response := fmt.Sprintf("Bearer %s\n", token)
	w.Write([]byte(response))
}
