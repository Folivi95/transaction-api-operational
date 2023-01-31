package testhelpers

import (
	"context"

	"github.com/saltpay/go-kafka-driver"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/models"
)

type DummyTransformer struct{}

func (d DummyTransformer) Execute(context context.Context, incomingTransaction kafka.Message) (models.Transaction, error) {
	return LoadCanonicalTransaction(), nil
}
