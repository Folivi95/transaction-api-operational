//go:generate moq -out mocks/kafka_consumer.go -pkg=mocks . Consumer

package kafkalistener

import (
	"context"
	"time"

	kafka "github.com/saltpay/go-kafka-driver"
	"go.uber.org/zap"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/logger/zapctx"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/sync"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/ports"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/usecases"
)

type Consumer interface {
	FetchMessage(ctx context.Context) (kafka.Message, error)
	CommitMessage(ctx context.Context, msg kafka.Message) error
	Close()
}

type Listener interface {
	Listen(ctx context.Context)
	StopListening()
}

type TransactionsListener struct {
	transactionWriter      ports.Writer
	transactionTransformer ports.TransactionsTransformer
	consumer               Consumer
	metricsClient          ports.MetricsClient
	topicName              string
	shouldListen           *sync.AtomicBool
}

const (
	workerPoolSize    = 150
	dimWorkerPoolSize = 10
)

func NewTransactionsListener(
	transactionWriter ports.Writer,
	transactionTransformer ports.TransactionsTransformer,
	consumer Consumer,
	metricsClient ports.MetricsClient,
	topicName string,
) *TransactionsListener {
	shouldListen := sync.New()
	shouldListen.Set()

	return &TransactionsListener{
		transactionWriter:      transactionWriter,
		transactionTransformer: transactionTransformer,
		consumer:               consumer,
		metricsClient:          metricsClient,
		topicName:              topicName,
		shouldListen:           shouldListen,
	}
}

// Listen start listening for new incoming messages for the specified consumer. Since its a blocking long-lived task,
// it should live in a separate goroutine.
func (p *TransactionsListener) Listen(ctx context.Context) {
	messages := make(chan kafka.Message)
	p.startWorkerGroup(ctx, messages)
	for {
		if !p.shouldListen.IsSet() {
			break
		}

		zapctx.From(ctx).Debug("[TransactionsKafkaListener] Trying to fetch new kafka message")
		msg, err := p.consumer.FetchMessage(ctx)
		if err != nil {
			zapctx.From(ctx).Error("[TransactionsKafkaListener] Error fetching new kafka message: ", zap.Error(err))
			continue
		}
		zapctx.From(ctx).Debug("[TransactionsKafkaListener] Got message")

		p.metricsClient.Count("ingress_topic_counter", 1, []string{p.topicName})
		messages <- msg
	}
}

func (p *TransactionsListener) processMessage(ctx context.Context, incomingTransaction kafka.Message) {
	startTime := time.Now()
	transaction, err := p.transactionTransformer.Execute(ctx, incomingTransaction)
	if err != nil {
		// if ErrorNotTransaction, message was not mapped but stored
		if _, isErrorNotTransaction := err.(usecases.ErrorNotTransaction); isErrorNotTransaction {
			p.commitMessage(ctx, incomingTransaction)
			return
		}
		// todo: how to handle bad transaction formatting
		zapctx.From(ctx).Error("[TransactionsKafkaListener] Error transforming transaction: ", zap.Error(err))
		p.metricsClient.Count("egress_topic_counter", 1, []string{p.topicName, "failed"})
		p.metricsClient.Histogram("transaction_transformation_time_ms", float64(time.Since(startTime).Milliseconds()), []string{p.topicName, "failed"})
		p.commitMessage(ctx, incomingTransaction)
		return
	}

	err = p.transactionWriter.WriteKafka(ctx, transaction)
	if err != nil {
		zapctx.From(ctx).Error("[TransactionsKafkaListener] Error processing transaction: ", zap.Error(err))
		p.metricsClient.Histogram("transaction_transformation_time_ms", float64(time.Since(startTime).Milliseconds()), []string{p.topicName, "failed"})
		p.commitMessage(ctx, incomingTransaction)
		return
	}

	p.metricsClient.Count("egress_topic_counter", 1, []string{p.topicName, "success"})
	p.metricsClient.Histogram("transaction_transformation_time_ms", float64(time.Since(startTime).Milliseconds()), []string{p.topicName, "success"})
	p.commitMessage(ctx, incomingTransaction)
	zapctx.From(ctx).Debug("[TransactionsKafkaListener] Successfully processed transaction")
}

// StopListening gracefully request the consumer to stop.
func (p *TransactionsListener) StopListening() {
	p.shouldListen.UnSet()
}

func (p *TransactionsListener) startWorkerGroup(ctx context.Context, messages <-chan kafka.Message) {
	for w := 0; w < workerPoolSize; w++ {
		go p.work(ctx, messages)
	}
}

func (p *TransactionsListener) work(ctx context.Context, messages <-chan kafka.Message) {
	for message := range messages {
		p.processMessage(ctx, message)
	}
}

func (p *TransactionsListener) commitMessage(ctx context.Context, msg kafka.Message) {
	zapctx.From(ctx).Debug("[TransactionsKafkaListener] Committing kafka message")
	if err := p.consumer.CommitMessage(ctx, msg); err != nil {
		zapctx.From(ctx).Error("[TransactionsKafkaListener] Error committing msg: ", zap.Error(err))
	}
}
