//go:generate moq -out mocks/transactions_transformer_moq.go -pkg=mocks . TransactionsTransformer

package ports

import (
	"context"

	"github.com/saltpay/go-kafka-driver"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/models"
)

type TransactionsTransformer interface {
	Execute(ctx context.Context, incomingTransactionMessage kafka.Message) (models.Transaction, error)
}
