package user

import (
	"context"
	"errors"
	"log"

	"github.com/MarcusXavierr/rinha-de-backend-2024-q1/internal/db"
)

func (u UserService) movementBalance(params db.InsertBalanceTransactionParams, ctx context.Context) (userBalanceData, error) {
	// userBalance, err := mountUserNewBalance(params, ctx, u.DBConn)
	userBalance, err := mountUserNewBalanceAndValidate(params, ctx, u.DBConn)

	if err != nil {
		return userBalanceData{}, err
	}

	tx, err := u.DBTransactions.Begin(ctx)
	if err != nil {
		return userBalanceData{}, err
	}

	defer tx.Rollback(ctx)
	connWithTx := u.DBConn.WithTx(tx)

	if err = connWithTx.InsertBalanceTransaction(ctx, params); err != nil {
		return userBalanceData{}, err
	}

	err = connWithTx.UpdateUserBalance(ctx, db.UpdateUserBalanceParams{Balance: userBalance.Balance, ID: params.UserID.Int32})

	if err != nil {
		return userBalanceData{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return userBalanceData{}, err
	}

	return userBalance, nil
}

func mountUserNewBalance(params db.InsertBalanceTransactionParams, ctx context.Context, conn *db.Queries) (userBalanceData, error) {
	user, err := conn.GetUser(ctx, params.UserID.Int32)
	if err != nil {
		return userBalanceData{}, err
	}
	log.Printf("userBalance: %d | ParamsValue: %d \n", user.Balance, params.Value)

	if params.Type == creditType {
		return userBalanceData{Limit: user.CreditLimit, Balance: user.Balance + params.Value}, nil
	}

	limitExceeded := (user.Balance-params.Value)*-1 > user.CreditLimit

	if limitExceeded {
		return userBalanceData{}, limitExceededError
	}

	return userBalanceData{Limit: user.CreditLimit, Balance: user.Balance - params.Value}, nil
}

func mountUserNewBalanceAndValidate(params db.InsertBalanceTransactionParams, ctx context.Context, qtx *db.Queries) (userBalanceData, error) {
	user, err := qtx.GetUser(ctx, params.UserID.Int32)
	if err != nil {
		return userBalanceData{}, err
	}

	var balanceData userBalanceData

	if params.Type == creditType {
		balanceData = userBalanceData{Limit: user.CreditLimit, Balance: user.Balance + params.Value}
	} else {
		balanceData = userBalanceData{Limit: user.CreditLimit, Balance: user.Balance - params.Value}
	}

	if balanceData.Balance < balanceData.Limit*-1 {
		return userBalanceData{}, limitExceededError
	}

	return balanceData, nil
}

func getUserID(ctx context.Context) (int32, error) {
	user, ok := ctx.Value("user").(*db.User)
	if !ok {
		log.Printf("Error getting user")
		return 0, errors.New("invalid data")
	}

	return user.ID, nil
}
