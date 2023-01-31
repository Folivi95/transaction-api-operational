//go:generate moq -out mocks/transactions_moq.go -pkg=mocks . TransactionsSource

package ports

import (
	"context"
	"github.com/saltpay/go-kafka-driver"

	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/application/models"
)

type TransactionsSource interface {
	SaveTransaction(ctx context.Context, transaction models.Transaction) error
	GetTransactionEventsByInternalID(ctx context.Context, request models.GetTransactionEventsByInternalIDRequest) ([]models.Transaction, error)
	GetTransactionsByStoreID(ctx context.Context, request models.GetAllTransactionsRequest) ([]models.Transaction, string, string, error)
}

type TransactionsTransformer interface {
	Execute(ctx context.Context, incomingTransactionMessage kafka.Message) (models.Transaction, error)
}

type KafkaTransactionTransformer struct {
}

func (k *KafkaTransactionTransformer) Execute(ctx context.Context, incomingTransactionMessage kafka.Message) (models.Transaction, error) {
	model := models.Transaction{}
	return model, nil
}
