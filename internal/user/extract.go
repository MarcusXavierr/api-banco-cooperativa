package user

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/MarcusXavierr/rinha-de-backend-2024-q1/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type extractResponse struct {
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

func (u *UserService) HandleExtract(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value("cliente").(*db.User)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	response, err := u.mountExtractData(r.Context(), user)
	if err != nil {
		log.Printf("Error while mounting extractedData: %v", err)
		http.Error(w, http.StatusText(502), 502)
		return
	}

	w.Write(response)
}

func (u *UserService) mountExtractData(ctx context.Context, user *db.User) ([]byte, error) {
	dbTransactions, err := u.DBConn.GetLastTenTransactions(ctx, pgtype.Int4{Int32: user.ID, Valid: true})
	if err != nil {
		return nil, err
	}

	transactions := mountTransactions(dbTransactions)
	// TODO: Depois ver como tratar a formatação de datas pra ficar certinho. Agora tá falando que o ano é 80841
	extractedAt := time.Now().Format("2024-01-17T02:34:38.543030Z")

	extractedData := extractResponse{
		Balance: balance{
			Total:       int(user.Balance),
			ExtractedAt: extractedAt,
			Limit:       int(user.CreditLimit),
		},
		LastTransactions: transactions,
	}

	response, err := json.Marshal(extractedData)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func mountTransactions(dbTransactions []db.Transaction) []transaction {
	transactions := []transaction{}
	for _, dbTransaction := range dbTransactions {
		transactions = append(transactions, transaction{
			Value:       int(dbTransaction.Value),
			Type:        dbTransaction.Type,
			Description: dbTransaction.Description.String,
			CreatedAt:   dbTransaction.CreatedAt.Time.Format("2024-01-17T02:34:38.543030Z"),
		})
	}

	return transactions
}
