package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Transactions interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

type DBPool struct {
	Conn         *Queries
	Transactions Transactions
}
