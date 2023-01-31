package testhelpers

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"

	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/application/models"
)

type DummyDBHandler struct{}

func (m DummyDBHandler) GetTransactionEventsByInternalID(ctx context.Context, body models.GetTransactionEventsByInternalIDRequest) ([]models.Transaction, error) {
	t := loadTransaction()
	t.InternalID = body.InternalID
	return []models.Transaction{
		authTransaction("31814421-027b-4ac2-8d0d-80d7a9db8008"),
		captureTransaction("31814421-027b-4ac2-8d0d-80d7a9db8008"),
	}, nil
}

func (m DummyDBHandler) GetTransactionsByStoreID(ctx context.Context, body models.GetAllTransactionsRequest) ([]models.Transaction, string, string, error) {
	t := loadTransaction()
	t.StoreID = body.StoreID
	var transactions []models.Transaction
	countOfTransactions := 20
	for i := 0; i <= 20; i++ {
		uid, _ := uuid.NewUUID()
		transactions = append(transactions, authTransaction(uid.String()))
		transactions = append(transactions, captureTransaction(uid.String()))
		transactions = append(transactions, reversalTransaction(uid.String()))
	}
	if body.Limit > countOfTransactions {
		body.Limit = countOfTransactions
	}
	finalIDToken := base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(body.Limit)))
	prevIDToken := base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(1)))
	return transactions[0:body.Limit], prevIDToken, finalIDToken, nil
}

func (m DummyDBHandler) SaveTransaction(ctx context.Context, transaction models.Transaction) error {
	return nil
}
