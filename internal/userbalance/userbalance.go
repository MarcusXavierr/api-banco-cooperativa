package userbalance

import (
	"context"

	"github.com/MarcusXavierr/api-banco-cooperativa/internal/db"
)

type UserBalanceService struct {
	DB   *db.DBPool
	CTX  context.Context
	User *db.User
}

type TransactionRequest struct {
	Value       int    `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
}

type UserFinancialStatus struct {
	Limit   int32 `json:"limite"`
	Balance int32 `json:"saldo"`
}
