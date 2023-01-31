//go:build acceptance
// +build acceptance

package acceptance

import (
	"context"
	"github.com/matryer/is"
	"github.com/saltpay/transaction-api-operational/src-mario-api/black-box-tests/utils"
	"testing"
	"time"

	"github.com/saltpay/go-kafka-driver"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/postgres"
)

const (
	kafkaTopic               = "transaction-topic"
	kafkaConnectionString    = "kafka:9092"
	kafkaUsername            = ""
	kafkaPassword            = ""
	postgresConnectionString = "postgres://postgres:somePassword@postgres:5432"
	postgresPingCount        = 10
)

type FakeMetricHandler struct {
}

func (FakeMetricHandler) Histogram(name string, value float64, tags []string) {
}
func (FakeMetricHandler) Count(name string, value int64, tags []string) {
}

// TestKafka add a message onto the kafka transaction topic and verify that the app stores this in the database
func TestKafka(t *testing.T) {
	is := is.New(t)

	ctx := context.Background()
	client := utils.NewAPIClient(getBaseURL(t), t)
	query := "SELECT id  FROM transactions_api LIMIT 1000"

	if err := client.WaitForAPIToBeHealthy(fiveRetries); err != nil {
		t.Fatal(err)
	}

	dbHandler, err := postgres.NewDBHandler(postgresConnectionString, &FakeMetricHandler{}, postgresPingCount)
	if err != nil {
		t.Error(err)
	}

	existingRows := 0
	rows, err := dbHandler.DB.QueryContext(ctx, query)
	if err != nil {
		t.Error(err)
	}
	for rows.Next() {
		existingRows++
	}
	// send some expected transaction?
	sendTransactionToKafka()

	// allow some time for processing of topic
	time.Sleep(time.Second * 5)

	newRows := 0
	rows, err = dbHandler.DB.QueryContext(ctx, query)
	if err != nil {
		t.Error(err)
	}
	for rows.Next() {
		newRows++
	}

	is.Equal(newRows, existingRows+1)
}

func sendTransactionToKafka() {
	ctx := context.Background()
	producer, err := kafka.NewProducer(ctx, kafka.ProducerConfig{
		Addr:     []string{kafkaConnectionString},
		Topic:    kafkaTopic,
		Username: kafkaUsername,
		Password: kafkaPassword,
	})

	defer producer.Close()

	if err = producer.WriteMessage(ctx, kafka.Message{Value: []byte("hello")}); err != nil {
		panic(err)
	}
}
