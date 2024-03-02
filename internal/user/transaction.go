package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MarcusXavierr/rinha-de-backend-2024-q1/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type balanceMovementReq struct {
	Value       int    `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
}

const (
	creditType = "c"
	debitType  = "d"
)

var (
	limitExceededError = errors.New("Credit Limit exceeded")
	invalidInputError  = errors.New("User input is invalid")
)

func (u UserService) HandleBalanceMovement(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params, err := u.MountBalanceTransactionParams(ctx, r)
	if err != nil {
		if errors.Is(invalidInputError, err) {
			http.Error(w, "Invalid input", 422)
			return
		}
		http.Error(w, "Error mounting CreateTransactionParams", 500)
		return
	}

	userBalance, err := u.movementBalance(params, ctx)
	if err != nil {
		if errors.Is(limitExceededError, err) {
			http.Error(w, "You're broke", 422)
			return
		}

		http.Error(w, "Error communicating with database", 502)
	}

	// // HACK: to use later
	// if userBalance.Balance < userBalance.Limit*-1 {
	// 	log.Println(fmt.Sprintf("ao que parece essa pocilga não tinha pego esse bug óbvio|  operação=%s| saldo=%d | limite=%d |\n", params.Type, userBalance.Balance, userBalance.Limit))
	// 	http.Error(w, "You're broke", 422)
	// 	return
	// }

	writeBalanceMovementResponse(userBalance, w)
}

func writeBalanceMovementResponse(userBalance userBalanceData, w http.ResponseWriter) {
	data, err := json.Marshal(userBalance)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
	w.Write(data)
}

func (u UserService) MountBalanceTransactionParams(ctx context.Context, r *http.Request) (db.InsertBalanceTransactionParams, error) {
	userId, err := getUserID(ctx)
	if err != nil {
		return db.InsertBalanceTransactionParams{}, err
	}

	var response balanceMovementReq
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return db.InsertBalanceTransactionParams{}, invalidInputError
	}

	if err := validateBalanceMovementReq(response); err != nil {
		return db.InsertBalanceTransactionParams{}, invalidInputError
	}

	return db.InsertBalanceTransactionParams{
		UserID:      pgtype.Int4{Int32: userId, Valid: true},
		Value:       int32(response.Value),
		Type:        response.Type,
		Description: pgtype.Text{String: response.Description, Valid: true},
	}, nil
}

func validateBalanceMovementReq(response balanceMovementReq) error {
	if response.Type != creditType && response.Type != debitType {
		return invalidInputError
	}

	if len(response.Description) > 10 || len(response.Description) == 0 {
		return invalidInputError
	}

	return nil
}
