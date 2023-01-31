//go:build unit
// +build unit

package kafkalistener_test

import (
	"context"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/saltpay/go-kafka-driver"

	kafkalistener "github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/kafka"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/kafka/mocks"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/testhelpers"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/models"
	portsMock "github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/ports/mocks"
)

func TestListen(t *testing.T) {
	is := is.New(t)
	metricsClient := testhelpers.DummyMetricsClient{}
	t.Run("should successfully call the writer", func(t *testing.T) {
		done := make(chan bool)

		consumerMock := &mocks.ConsumerMock{
			CloseFunc: func() {},
			CommitMessageFunc: func(_ context.Context, msg kafka.Message) error {
				return nil
			},
			FetchMessageFunc: func(_ context.Context) (kafka.Message, error) {
				// As of now, we're only receiving and processing one payment at a time
				return kafka.Message{}, nil
			},
		}

		mockWriter := portsMock.WriterMock{
			WriteKafkaFunc: func(ctx context.Context, transaction models.Transaction) error {
				done <- true
				return nil
			},
			WriteDBFunc: func(ctx context.Context, transaction models.Transaction) error {
				done <- true
				return nil
			},
		}
		mockTransformer := portsMock.TransactionsTransformerMock{ExecuteFunc: func(ctx context.Context, incomingTransaction kafka.Message) (models.Transaction, error) {
			return models.Transaction{}, nil
		}}
		paymentsListener := kafkalistener.NewTransactionsListener(&mockWriter, &mockTransformer, consumerMock, metricsClient, "")
		go paymentsListener.Listen(context.Background())

		// cleanup
		defer paymentsListener.StopListening()

		select {
		case <-done:
			is.True(len(mockWriter.WriteKafkaCalls()) >= 1)
		case <-time.After(50 * time.Millisecond):
			t.Fatal("time out while waiting for test to finish")
		}
	})
}
