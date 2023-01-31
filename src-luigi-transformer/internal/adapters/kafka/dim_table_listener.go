package kafkalistener

import (
	"context"
	"encoding/json"
	"time"

	kafka "github.com/saltpay/go-kafka-driver"
	"go.uber.org/zap"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/logger/zapctx"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/sync"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/ports"
)

type Validator interface {
	IsValid(ctx context.Context, schemaRegistry ports.SchemaRegistry, incomingMessage []byte, schemaKey string) bool
}

type ValidatorStruct struct{}

func (ValidatorStruct) IsValid(ctx context.Context, schemaRegistry ports.SchemaRegistry, incomingMessage []byte, schemaKey string) bool {
	// validate and unmarshal incoming message
	_, _, err := schemaRegistry.Decode(ctx, incomingMessage, schemaKey)
	if err != nil {
		zapctx.From(ctx).Error("[Dim-Validator] Ingress message does not comply to the defined schema")
		return false
	}

	return true
}

type DimTableListener struct {
	dbHandler      ports.DBHandler
	consumer       Consumer
	metricsClient  ports.MetricsClient
	tableName      string
	schemaRegistry ports.SchemaRegistry
	schemaKey      string
	shouldListen   *sync.AtomicBool
	Validator      Validator
}

func NewDimTableListener(
	dbHandler ports.DBHandler,
	consumer Consumer,
	metricsClient ports.MetricsClient,
	schemaRegistry ports.SchemaRegistry,
	schemaKey string,
	tableName string,
) *DimTableListener {
	shouldListen := sync.New()
	shouldListen.Set()

	return &DimTableListener{
		dbHandler:      dbHandler,
		consumer:       consumer,
		metricsClient:  metricsClient,
		tableName:      tableName,
		schemaRegistry: schemaRegistry,
		schemaKey:      schemaKey,
		shouldListen:   shouldListen,
		Validator:      ValidatorStruct{},
	}
}

// Listen start listening for new incoming messages for the specified consumer. Since its a blocking long-lived task,
// it should live in a separate goroutine.
func (p *DimTableListener) Listen(ctx context.Context) {
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

		p.metricsClient.Count("ingress_topic_counter", 1, []string{p.tableName})
		messages <- msg
	}
}

func (p *DimTableListener) processMessage(ctx context.Context, incomingMessage kafka.Message) {
	startTime := time.Now()
	var entry map[string]interface{}

	if p.Validator.IsValid(ctx, p.schemaRegistry, incomingMessage.Value, p.schemaKey) {
		err := json.Unmarshal(incomingMessage.Value, &entry)
		if err != nil {
			zapctx.From(ctx).Error("[TransactionsKafkaListener] Error unmarshalling the message: ", zap.Error(err))
			p.metricsClient.Count("egress_topic_counter", 1, []string{p.tableName, "failed"})
			p.metricsClient.Histogram("transaction_transformation_time_ms", float64(time.Since(startTime).Milliseconds()), []string{p.tableName, "failed"})
			p.commitMessage(ctx, incomingMessage)
			return
		}

		err = p.dbHandler.Upsert(ctx, entry, p.tableName)
		if err != nil {
			zapctx.From(ctx).Error("[TransactionsKafkaListener] Error processing message: ", zap.Error(err))
			p.metricsClient.Count("egress_topic_counter", 1, []string{p.tableName, "failed"})
			p.metricsClient.Histogram("transaction_transformation_time_ms", float64(time.Since(startTime).Milliseconds()), []string{p.tableName, "failed"})
			p.commitMessage(ctx, incomingMessage)
			return
		}
	} else {
		p.metricsClient.Count("egress_topic_counter", 1, []string{p.tableName, "failed"})
		p.metricsClient.Histogram("transaction_transformation_time_ms", float64(time.Since(startTime).Milliseconds()), []string{p.tableName, "failed"})
		return
	}

	p.commitMessage(ctx, incomingMessage)
	p.metricsClient.Count("egress_topic_counter", 1, []string{p.tableName, "success"})
	p.metricsClient.Histogram("transaction_transformation_time_ms", float64(time.Since(startTime).Milliseconds()), []string{p.tableName, "success"})
	zapctx.From(ctx).Debug("[TransactionsKafkaListener] Successfully processed message")
}

// StopListening gracefully request the consumer to stop.
func (p *DimTableListener) StopListening() {
	p.shouldListen.UnSet()
}

func (p *DimTableListener) startWorkerGroup(ctx context.Context, messages <-chan kafka.Message) {
	for w := 0; w < dimWorkerPoolSize; w++ {
		go p.work(ctx, messages)
	}
}

func (p *DimTableListener) work(ctx context.Context, messages <-chan kafka.Message) {
	for message := range messages {
		p.processMessage(ctx, message)
	}
}

func (p *DimTableListener) commitMessage(ctx context.Context, msg kafka.Message) {
	zapctx.From(ctx).Debug("[TransactionsKafkaListener] Committing kafka message")
	if err := p.consumer.CommitMessage(ctx, msg); err != nil {
		zapctx.From(ctx).Error("[TransactionsKafkaListener] Error committing msg: ", zap.Error(err))
	}
}
