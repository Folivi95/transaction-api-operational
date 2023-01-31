package main

import (
	"context"
	"fmt"

	"github.com/saltpay/transaction-api-operational/src-kamek-cronjob/internal/adapters/logger/zapctx"
	"github.com/saltpay/transaction-api-operational/src-kamek-cronjob/internal/adapters/postgres"
	"github.com/saltpay/transaction-api-operational/src-kamek-cronjob/internal/adapters/prometheus"
	"github.com/saltpay/transaction-api-operational/src-kamek-cronjob/internal/application/ports"
)

// App holds and creates dependencies for your application.
type App struct {
	DBHandler postgres.DBHandler
}

func newApp(applicationContext context.Context) (*App, error) {
	// Add common fields to application context
	applicationContext = zapctx.WithCommonFields(applicationContext)

	// Create Prometheus client
	metricsClient := newMetricsClient()

	appConfig, err := loadAppConfig()
	if err != nil {
		return &App{}, fmt.Errorf("failed to load config: %w", err)
	}

	var dbHandler postgres.DBHandler
	dbHandler, err = postgres.NewDBHandler(appConfig.PostgresConfig.PostgresURL, metricsClient, 10)
	if err != nil {
		return &App{}, fmt.Errorf("failed to create db handler: %w", err)
	}

	go handleInterrupts(applicationContext)

	return &App{
		DBHandler: dbHandler,
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
