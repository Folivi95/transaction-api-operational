package usecases

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/saltpay/go-kafka-driver"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/logger/zapctx"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/models"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/ports"
)

type Producer interface {
	WriteMessage(ctx context.Context, message kafka.Message) error
	Close()
}

type TransactionWriter struct {
	metricsClient ports.MetricsClient
	producer      Producer
	dbHandler     ports.DBHandler
	srClient      ports.SchemaRegistry
	schemaKey     string
}

func NewTransactionWriter(metricsClient ports.MetricsClient, producer Producer, handler ports.DBHandler, srClient ports.SchemaRegistry, schemaKey string) *TransactionWriter {
	return &TransactionWriter{
		metricsClient: metricsClient,
		producer:      producer,
		dbHandler:     handler,
		srClient:      srClient,
		schemaKey:     schemaKey,
	}
}

func (t TransactionWriter) WriteKafka(ctx context.Context, transaction models.Transaction) error {
	// send to kafka
	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		zapctx.From(ctx).Error("[MessageProcessorKafka] Unable to unmarshal message")
	}

	// validate against egress schema
	schema := t.checkSchema(ctx)
	_, _, err = t.srClient.Decode(ctx, transactionJSON, schema)
	if err != nil {
		msg := fmt.Sprintf("[MessageProcessorKafka] %s %s ", transactionJSON, schema)
		zapctx.From(ctx).Warn(msg)
		zapctx.From(ctx).Error("[MessageProcessorKafka] Produced message does not comply to the defined egress schema")
		return err
	}
	err = t.producer.WriteMessage(ctx, kafka.Message{Value: transactionJSON})
	if err != nil {
		return err
	}

	return nil
}

func (t TransactionWriter) checkSchema(ctx context.Context) string {
	schema := ""
	// check if schemaKey is empty
	if len(t.schemaKey) == 0 {
		schema = "ent-canonical_transaction-v2"
		zapctx.From(ctx).Warn("[MessageProcessorKafka] Schema key is empty")
	} else {
		schema = t.schemaKey
	}
	return schema
}
