package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MarcusXavierr/rinha-de-backend-2024-q1/internal/userbalance"
	"github.com/pkg/errors"
)

var (
	invalidInputError = errors.New("User input is invalid")
)

type limitExceededError interface {
	HasExceededLimit() bool
}

const (
	creditType = "c"
	debitType  = "d"
)

func (u UserService) HandleTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params, err := u.parseTransactionRequest(r)
	user, userErr := GetUser(ctx)
	if err != nil {
		log.Printf("%v\n%v\n", err, userErr)
		http.Error(w, "Invalid input", 422)
		return
	}

	balance, err := userbalance.UserBalanceService{DB: u.DB, CTX: ctx, User: user}.HandleTransaction(params)
	log.Printf("balance: %v", balance)
	if err != nil {
		log.Printf("%v+", err)
		if l, ok := errors.Cause(err).(limitExceededError); ok && l.HasExceededLimit() {
			http.Error(w, "You're broke", 422)
		} else {
			http.Error(w, "Error communicating with database", 502)
		}
		return
	}

	writeBalanceMovementResponse(balance, w)
}

func writeBalanceMovementResponse(userBalance userbalance.UserFinancialStatus, w http.ResponseWriter) {
	data, err := json.Marshal(userBalance)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
	w.Write(data)
}

func (u UserService) parseTransactionRequest(r *http.Request) (userbalance.TransactionRequest, error) {
	var request userbalance.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("%v+\n", err)
		return userbalance.TransactionRequest{}, errors.Wrap(invalidInputError, "Error decoding request body")
	}

	if err := validateTransactionRequest(request); err != nil {
		return userbalance.TransactionRequest{}, errors.Wrap(err, "Error validating request body")
	}

	return request, nil
}

func validateTransactionRequest(response userbalance.TransactionRequest) error {
	if response.Type != creditType && response.Type != debitType {
		return invalidInputError
	}

	if response.Type == debitType {
		response.Value = response.Value * -1
	}

	if len(response.Description) > 10 || len(response.Description) == 0 {
		return invalidInputError
	}

	return nil
}
