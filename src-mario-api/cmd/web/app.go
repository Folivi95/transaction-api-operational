package main

import (
	"context"
	"fmt"
	kafkaDriver "github.com/saltpay/go-kafka-driver"

	zapctx "github.com/saltpay/go-zap-ctx"

	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/http"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/kafka"
	zctx_hlprs "github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/logger/zapctx_helpers"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/postgres"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/prometheus"
	schemaregistry "github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/schemaRegistry"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/testhelpers"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/application/ports"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/configuration"
	"go.uber.org/zap"
	"time"
)

// App holds and creates dependencies for your application.
type App struct {
	ServerConfig   http.ServerConfig
	DBHandler      ports.TransactionsSource
	AppConfig      configuration.AppConfig
	MetricsClient  ports.MetricsClient
	SchemaRegistry ports.SchemaRegistry
}

func newApp(applicationContext context.Context) (*App, error) {
	config := configuration.NewDefaultConfig()

	// Add common fields to application context
	applicationContext = zctx_hlprs.WithCommonFields(applicationContext)

	// Create Prometheus client
	metricsClient := newMetricsClient()

	appConfig, err := configuration.LoadAppConfig()
	if err != nil {
		return &App{}, fmt.Errorf("failed to load configuration: %w", err)
	}
	initLog(appConfig)

	var dbHandler ports.TransactionsSource
	if appConfig.PostgresConfig.MockedData {
		dbHandler = testhelpers.DummyDBHandler{}
	} else {
		dbHandler, err = postgres.NewDBHandler(appConfig.PostgresConfig.PostgresConnection, metricsClient, 10)
		if err != nil {
			return &App{}, fmt.Errorf("failed to create db handler: %w", err)
		}
	}
	schemaRegistryClient := schemaregistry.NewSchemaRegistryClient(appConfig.SchemaRegistryConfig.SchemaRegistryEndpoint, appConfig.SchemaRegistryConfig.SchemaRegistryRefreshTimeSeconds)

	go handleInterrupts(applicationContext)

	return &App{
		ServerConfig:   config,
		DBHandler:      dbHandler,
		AppConfig:      appConfig,
		MetricsClient:  metricsClient,
		SchemaRegistry: schemaRegistryClient,
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
	zapctx.Info(ctx, "shutting down")
}

// ListenForTransactions Listen to the configured kafka endpoint for transactions and store them
func (a *App) ListenForTransactions(context context.Context) error {
	listener := kafka.NewTransactionListener(a.AppConfig.KafkaConfig)
	err := listener.Listen(context, a.ProcessTransaction)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) ProcessTransaction(ctx context.Context, consumer *kafkaDriver.Consumer, message kafkaDriver.Message) error {
	//validate the schema
	_, valid := a.SchemaRegistry.Decode(ctx, message.Value, a.AppConfig.KafkaConfig.KafkaTransactionSchema)
	if !valid {
		msg := fmt.Sprintf("[TransactionListener] %s %s ", string(message.Value), a.AppConfig.KafkaConfig.KafkaTransactionSchema)
		zapctx.Warn(ctx, msg)
		zapctx.Error(ctx, "[TransactionListener] Received message does not comply to the defined egress schema")
		return consumer.CommitMessage(ctx, message)
	}

	startTime := time.Now()
	a.MetricsClient.Count("invocation_count", 1, []string{"TransactionListener"})

	transformer := ports.KafkaTransactionTransformer{}
	transaction, err := transformer.Execute(ctx, message)
	if err != nil {
		return err
	}

	// store transaction in db
	err = a.DBHandler.SaveTransaction(ctx, transaction)

	if err != nil {
		zapctx.Error(ctx, "[TransactionListener] Unable to process message", zap.Error(err))
		a.MetricsClient.Count("invocation_result", 1, []string{"TransactionListener", "failed"})
		a.MetricsClient.Histogram("execution_time", float64(time.Since(startTime).Milliseconds()), []string{"TransactionListener", "failed"})
		return err
	}
	a.MetricsClient.Count("invocation_result", 1, []string{"TransactionListener", "success"})
	a.MetricsClient.Histogram("execution_time", float64(time.Since(startTime).Milliseconds()), []string{"TransactionListener", "success"})

	err = consumer.CommitMessage(ctx, message)
	if err != nil {
		zapctx.Error(ctx, "[TransactionListener] unable to commit message", zap.Error(err))
		return err
	}
	return nil
}
