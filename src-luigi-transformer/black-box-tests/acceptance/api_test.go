//go:build acceptance
// +build acceptance

package acceptance

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/saltpay/go-kafka-driver"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/models"
)

const fiveRetries = 5

// Sanity check.
func TestGreetingApplication(t *testing.T) {
	client := NewAPIClient(getBaseURL(t), t)

	if err := client.WaitForAPIToBeHealthy(fiveRetries); err != nil {
		t.Fatal(err)
	}
}

// produces to the incoming topic and consumes from the outgoing.
func TestProducer(t *testing.T) {
	var (
		ctx           = context.Background()
		kafkaEndpoint = os.Getenv("KAFKA_ENDPOINT")
		kafkaUsername = os.Getenv("KAFKA_USERNAME")
		kafkaPassword = os.Getenv("KAFKA_PASSWORD")
	)
	// create consumer and producer
	consumer, err := kafka.NewConsumer(ctx, kafka.ConsumerConfig{
		Brokers:  strings.Split(kafkaEndpoint, ","),
		GroupID:  fmt.Sprintf("%s-groupId", kafkaTopicOutgoingTransactions),
		Topic:    kafkaTopicOutgoingTransactions,
		Username: kafkaUsername,
		Password: kafkaPassword,
	})
	if err != nil {
		t.Fatal(err)
	}

	// make sure server is healthy
	client := NewAPIClient(getBaseURL(t), t)
	if err := client.WaitForAPIToBeHealthy(fiveRetries); err != nil {
		t.Fatal(err)
	}

	t.Run("happy path for kafka topics", func(t *testing.T) {
		// given a transaction, kafka consumer, and kafka producer
		is := is.New(t)
		messageReceived := make(chan kafka.Message)

		// listener to consume message. Channel will receive message to be asserted later
		go func() {
			time.Sleep(time.Second)
			msg, err := consumer.FetchMessage(ctx)
			is.NoErr(err)
			messageReceived <- msg
		}()

		select {
		case kafkaMessage := <-messageReceived:
			// if channel receives a message, then the message should be the expected one
			var transaction models.Transaction
			err := json.Unmarshal(kafkaMessage.Value, &transaction)
			is.NoErr(err)
		case <-time.After(30 * time.Second):
			t.Fatal("time out while waiting for test to finish")
		}
	})
}
