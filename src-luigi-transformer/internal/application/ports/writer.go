//go:generate moq -out mocks/writer_moq.go -pkg=mocks . Writer

package ports

import (
	"context"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/models"
)

// Writer handles DB and Kafka write-related methods.
type Writer interface {
	// WriteDB(ctx context.Context, transaction models.Transaction) error
	WriteKafka(ctx context.Context, transaction models.Transaction) error
}
