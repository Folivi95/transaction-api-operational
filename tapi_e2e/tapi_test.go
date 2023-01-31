//go:build acceptance
// +build acceptance

package tapi_e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/saltpay/go-kafka-driver"
)

const (
	fiveRetries                    = 5
	kafkaTopicIncomingTransactions = "transaction-api-operational-connect-solar-curated-afterProcessTxn"
)

func TestTAPI(t *testing.T) {
	var (
		ctx           = context.Background()
		kafkaEndpoint = os.Getenv("KAFKA_ENDPOINT")
		kafkaUsername = os.Getenv("KAFKA_USERNAME")
		kafkaPassword = os.Getenv("KAFKA_PASSWORD")
		transactionID = 123
	)

	// producer
	producer, err := kafka.NewProducer(ctx, kafka.ProducerConfig{
		Addr:     strings.Split(kafkaEndpoint, ","),
		Topic:    kafkaTopicIncomingTransactions,
		Username: kafkaUsername,
		Password: kafkaPassword,
	})
	if err != nil {
		t.Fatal(err)
	}
	client := NewAPIClient(getBaseURL(t), t)
	if err = client.WaitForAPIToBeHealthy(fiveRetries); err != nil {
		t.Fatal(err)
	}

	is := is.New(t)
	t.Run("An incoming transaction in the kafka topic should generate a transaction in the DB", func(t *testing.T) {
		t.Parallel()
		incomingTransaction := SolarIncomingTransaction{}
		incomingTransaction.Body.Txn.TxnID = transactionID
		transactionRequest := SingleTransactionRequest{
			TransactionID: fmt.Sprint(transactionID),
			AcquiringHost: "",
		}

		// when a message is written
		v, err := json.Marshal(incomingTransaction)
		is.NoErr(err)

		err = producer.WriteMessage(ctx, kafka.Message{
			Value: v,
		})
		is.NoErr(err)

		transactionJSON, err := client.GetTransaction(transactionRequest)
		is.NoErr(err)

		canonicalTransaction, err := NewTransactionFromJSON(transactionJSON)
		is.NoErr(err)
		is.Equal(canonicalTransaction.ID, fmt.Sprint(transactionID))
	})
}
