package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/saltpay/go-kafka-driver"
	"go.uber.org/zap"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/http"
	kafkalistener "github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/kafka"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/logger/zapctx"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/postgres"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/prometheus"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/schemaregistry"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/testhelpers"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/ports"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/usecases"
)

// App holds and creates dependencies for your application.
type App struct {
	ServerConfig       http.ServerConfig
	W4Listener         kafkalistener.Listener
	SolarListener      kafkalistener.Listener
	ListenDimTables    bool
	ListenTransactions bool
}

func newApp(applicationContext context.Context) (*App, error) {
	config := newDefaultConfig()

	// Add common fields to application context
	applicationContext = zapctx.WithCommonFields(applicationContext)

	// Create Prometheus client
	metricsClient := newMetricsClient()

	appConfig, err := loadAppConfig()
	if err != nil {
		return &App{}, fmt.Errorf("failed to load config: %w", err)
	}

	// Create Kafka Clients
	w4ConsumerConfig := kafka.ConsumerConfig{
		Brokers:  appConfig.KafkaConfig.KafkaEndpoint,
		GroupID:  fmt.Sprintf("%s-groupID", appConfig.KafkaConfig.KafkaW4IngressTopic),
		Topic:    appConfig.KafkaConfig.KafkaW4IngressTopic,
		Username: appConfig.KafkaConfig.KafkaUsername,
		Password: appConfig.KafkaConfig.KafkaPassword,
	}
	w4Consumer, err := kafka.NewConsumer(applicationContext, w4ConsumerConfig)
	if err != nil {
		return &App{}, fmt.Errorf("failed to create w4 consumer: %w", err)
	}

	solarConsumerConfig := kafka.ConsumerConfig{
		Brokers:  appConfig.KafkaConfig.KafkaEndpoint,
		GroupID:  fmt.Sprintf("%s-groupID", appConfig.KafkaConfig.KafkaSolarIngressTopic),
		Topic:    appConfig.KafkaConfig.KafkaSolarIngressTopic,
		Username: appConfig.KafkaConfig.KafkaUsername,
		Password: appConfig.KafkaConfig.KafkaPassword,
	}
	solarConsumer, err := kafka.NewConsumer(applicationContext, solarConsumerConfig)
	if err != nil {
		return &App{}, fmt.Errorf("failed to create solar consumer: %w", err)
	}

	producerConfig := kafka.ProducerConfig{
		Addr:     appConfig.KafkaConfig.KafkaEndpoint,
		Topic:    appConfig.KafkaConfig.KafkaEgressTopic,
		Username: appConfig.KafkaConfig.KafkaUsername,
		Password: appConfig.KafkaConfig.KafkaPassword,
	}
	producer, err := kafka.NewProducer(applicationContext, producerConfig)
	if err != nil {
		return &App{}, fmt.Errorf("error creating kafka producer %w", err)
	}

	// create DB handler
	dbHandler, err := postgres.NewDBHandler(appConfig.PostgresConfig.PostgresConnection, metricsClient, 10)
	if err != nil {
		return &App{}, fmt.Errorf("failed to create db handler: %w", err)
	}

	// schema registry client
	srClient := schemaregistry.NewSchemaRegistryClient(appConfig.SchemaRegistryConfig.SchemaRegistryEndpoint, appConfig.SchemaRegistryConfig.SchemaRegistryRefreshTimeSeconds)

	// Create use cases
	var w4Transformer ports.TransactionsTransformer
	var solarTransformer ports.TransactionsTransformer
	if appConfig.KafkaConfig.MockedData {
		w4Transformer = testhelpers.DummyTransformer{}
		solarTransformer = testhelpers.DummyTransformer{}
		go fakeIngestion(applicationContext, producer)
	} else {
		w4Transformer = usecases.NewW4Transformer(dbHandler, srClient, appConfig.SchemaRegistryConfig.W4IngressSchemaKey, metricsClient)
		solarTransformer = usecases.NewSolarTransformer(srClient, appConfig.SchemaRegistryConfig.SolarIngressSchemaKey)
	}

	// Create Transaction Writes
	transactionWriterW4 := usecases.NewTransactionWriter(metricsClient, producer, dbHandler, srClient, appConfig.SchemaRegistryConfig.W4EgressSchemaKey)
	transactionWriterSolar := usecases.NewTransactionWriter(metricsClient, producer, dbHandler, srClient, appConfig.SchemaRegistryConfig.SolarEgressSchemaKey)

	// Create Transaction Listeners
	w4Listener := kafkalistener.NewTransactionsListener(transactionWriterW4, w4Transformer, w4Consumer, metricsClient, "way4")
	solarListener := kafkalistener.NewTransactionsListener(transactionWriterSolar, solarTransformer, solarConsumer, metricsClient, "solar")

	go handleInterrupts(applicationContext)

	return &App{
		ServerConfig:       config,
		W4Listener:         w4Listener,
		SolarListener:      solarListener,
		ListenDimTables:    appConfig.KafkaConfig.ListenDimTables,
		ListenTransactions: appConfig.KafkaConfig.ListenTransactions,
	}, nil
}

func newMetricsClient() ports.MetricsClient {
	prometheusClient := prometheus.New()

	return &prometheusClient
}

// this is just an example of how the services within an app could listen to the
// cancellation signal and stop their work gracefully. So it's probably a decent
// idea to pass the application context to services if you want to do some
// cleanup before finishing.
func handleInterrupts(ctx context.Context) {
	<-ctx.Done()
	zapctx.From(ctx).Info("shutting down")
}

func fakeIngestion(ctx context.Context, producer usecases.Producer) {
	executor := testhelpers.DummyTransformer{}
	for {
		transaction, _ := executor.Execute(ctx, kafka.Message{})
		transactionJSON, _ := json.Marshal(transaction)
		err := producer.WriteMessage(ctx, kafka.Message{Value: transactionJSON})
		if err != nil {
			zapctx.From(ctx).Error("could not write message to topic", zap.Error(err))
		} else {
			zapctx.From(ctx).Debug("fake ingested transaction!")
		}
		time.Sleep(5 * time.Minute)
	}
}
