package userbalance

import (
	"context"
	"log"

	"github.com/MarcusXavierr/rinha-de-backend-2024-q1/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
)

func (ub UserBalanceService) HandleTransaction(request TransactionRequest) (UserFinancialStatus, error) {
	dbParams := db.RegisterTransactionParams{
		UserID:      pgtype.Int4{Int32: ub.User.ID, Valid: true},
		Value:       int32(request.Value),
		Type:        request.Type,
		Description: pgtype.Text{String: request.Description, Valid: true},
	}

	return ub.registerTransaction(dbParams, ub.CTX)
}

func (ub UserBalanceService) registerTransaction(params db.RegisterTransactionParams, ctx context.Context) (UserFinancialStatus, error) {
	emptyStatus := UserFinancialStatus{}
	tx, err := ub.DB.Transactions.Begin(ctx)
	if err != nil {
		return emptyStatus, errors.Wrap(err, "starting transaction failed")
	}
	defer tx.Rollback(ctx)
	connWithTx := ub.DB.Conn.WithTx(tx)

	if err = connWithTx.RegisterTransaction(ctx, params); err != nil {
		return emptyStatus, errors.Wrap(err, "error registering transaction")
	}

	financialStatus, err := ub.updateBalance(params, connWithTx)
	log.Printf("financialStatus: %v", financialStatus)
	if err != nil {
		return emptyStatus, errors.Wrap(err, "updating user balance failed")
	}
	err = tx.Commit(ctx)

	return financialStatus, err
}

func (ub UserBalanceService) updateBalance(params db.RegisterTransactionParams, qtx *db.Queries) (UserFinancialStatus, error) {
	financialStatus := UserFinancialStatus{Limit: ub.User.CreditLimit, Balance: ub.User.Balance + params.Value}
	log.Printf("financialStatus inside: %v", financialStatus)

	if financialStatus.Balance < financialStatus.Limit*-1 {
		return UserFinancialStatus{}, financialStatusError{}
	}

	updateBalanceParams := db.UpdateUserBalanceParams{Balance: financialStatus.Balance, ID: params.UserID.Int32}
	return financialStatus, qtx.UpdateUserBalance(ub.CTX, updateBalanceParams)
}
