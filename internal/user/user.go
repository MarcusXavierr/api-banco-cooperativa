package user

import (
	"context"

	"github.com/MarcusXavierr/rinha-de-backend-2024-q1/internal/db"
	"github.com/jackc/pgx/v5"
)

type DBTransactions interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

type dbTransactionInterface interface {
	Rollback(ctx context.Context) error
	Commit(ctx context.Context) error
}

type UserService struct {
	DBConn         *db.Queries
	DBTransactions DBTransactions
}

type userBalanceData struct {
	Limit   int32 `json:"limite"`
	Balance int32 `json:"saldo"`
}
