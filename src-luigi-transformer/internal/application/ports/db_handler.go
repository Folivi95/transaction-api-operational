//go:generate moq -out mocks/db_handler_moq.go -pkg mocks . DBHandler

package ports

import (
	"context"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/models"
)

type DBHandler interface {
	Store(ctx context.Context, transactionJSON []byte) error
	Upsert(ctx context.Context, entry map[string]interface{}, tableName string) error
	Get(ctx context.Context, query string, target []interface{}) error
	StoreAuxTable(ctx context.Context, transaction models.W4IncomingTransaction, insert bool) error
}
