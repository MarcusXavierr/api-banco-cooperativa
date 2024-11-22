package user

import (
	"log"
	"net/http"

	"github.com/MarcusXavierr/rinha-de-backend-2024-q1/internal/userbalance"
)

func (u *UserService) HandleExtract(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := GetUser(ctx)
	if err != nil {
		http.Error(w, "Cannot get user from context", 502)
	}

	data, err := userbalance.UserBalanceService{DB: u.DB, User: user, CTX: ctx}.HandleExtract()
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(502), 502)
		return
	}

	w.Write(data)
}
