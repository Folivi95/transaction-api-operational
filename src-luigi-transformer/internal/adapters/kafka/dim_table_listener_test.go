//go:build unit
// +build unit

package kafkalistener_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/saltpay/go-kafka-driver"

	kafkalistener "github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/kafka"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/kafka/mocks"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/testhelpers"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/ports"
	portsMock "github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/ports/mocks"
)

type fakeValidator struct{}

func (fakeValidator) IsValid(ctx context.Context, schemaRegistry ports.SchemaRegistry, incomingMessage []byte, schemaKey string) bool {
	return true
}

func TestDimTableListen(t *testing.T) {
	is := is.New(t)
	t.Run("should successfully call the writer", func(t *testing.T) {
		done := make(chan bool)
		entry := map[string]string{"key": "value"}
		message, _ := json.Marshal(entry)

		metricsClient := testhelpers.DummyMetricsClient{}
		consumerMock := &mocks.ConsumerMock{
			CloseFunc: func() {},
			CommitMessageFunc: func(_ context.Context, msg kafka.Message) error {
				return nil
			},
			FetchMessageFunc: func(_ context.Context) (kafka.Message, error) {
				// As of now, we're only receiving and processing one payment at a time
				return kafka.Message{Value: message}, nil
			},
		}

		mockDBHandler := portsMock.DBHandlerMock{
			UpsertFunc: func(ctx context.Context, entry map[string]interface{}, tableName string) error {
				done <- true
				return nil
			},
		}

		mockSchemaRegistry := portsMock.SchemaRegistryMock{}

		dimTableListener := kafkalistener.NewDimTableListener(&mockDBHandler, consumerMock, metricsClient, &mockSchemaRegistry, "key", "name")
		dimTableListener.Validator = fakeValidator{}
		go dimTableListener.Listen(context.Background())

		// cleanup
		defer dimTableListener.StopListening()

		select {
		case <-done:
			is.True(len(mockDBHandler.UpsertCalls()) >= 1)
		case <-time.After(50 * time.Millisecond):
			t.Fatal("time out while waiting for test to finish")
		}
	})
}
