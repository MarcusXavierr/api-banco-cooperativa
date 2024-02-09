package user

import (
	"net/http"

	"github.com/MarcusXavierr/rinha-de-backend-2024-q1/internal/db"
)

type UserService struct {
	DBConn *db.Queries
}

func HandleTransaction(w http.ResponseWriter, r *http.Request) {
	// code here
}
