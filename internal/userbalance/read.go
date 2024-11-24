package userbalance

import (
	"context"
	"encoding/json"
	"time"

	"github.com/MarcusXavierr/api-banco-cooperativa/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
)

type extractedData struct {
	Balance          balance       `json:"saldo"`
	LastTransactions []transaction `json:"ultimas_transacoes"`
}

type balance struct {
	Total       int    `json:"total"`
	ExtractedAt string `json:"data_extrato"`
	Limit       int    `json:"limite"`
}

type transaction struct {
	Value       int    `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
	CreatedAt   string `json:"realizada_em"`
}

func (ub UserBalanceService) HandleExtract() ([]byte, error) {
	transactions, err := ub.getLastTransactions(ub.CTX, ub.User)
	if err != nil {
		return nil, errors.Wrap(err, "error getting last transactions")
	}
	extractedAt := time.Now().Format("2006-01-02T15:04:05.000000Z")

	extractedData := extractedData{
		Balance: balance{
			Total:       int(ub.User.Balance),
			ExtractedAt: extractedAt,
			Limit:       int(ub.User.CreditLimit),
		},
		LastTransactions: transactions,
	}

	response, err := json.Marshal(extractedData)
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling extractedData")
	}

	return response, nil
}

func (ub UserBalanceService) getLastTransactions(ctx context.Context, user *db.User) ([]transaction, error) {
	dbTransactions, err := ub.DB.Conn.GetLastTenTransactions(ctx, pgtype.Int4{Int32: user.ID, Valid: true})
	if err != nil {
		return nil, errors.Wrap(err, "connection error with database")
	}

	transactions := []transaction{}
	for _, dbTransaction := range dbTransactions {
		transactions = append(transactions, transaction{
			Value:       int(dbTransaction.Value),
			Type:        dbTransaction.Type,
			Description: dbTransaction.Description.String,
			CreatedAt:   dbTransaction.CreatedAt.Time.Format("2006-01-02T15:04:05.000000Z"),
		})
	}

	return transactions, nil
}
