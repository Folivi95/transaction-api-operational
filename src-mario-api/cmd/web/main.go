package main

import (
	"fmt"
	IH "github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/http/auth"
	zap2 "github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/prometheus/zap"
	"go.uber.org/zap"
	"os"

	zapctx "github.com/saltpay/go-zap-ctx"

	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/http"
)

func main() {
	ctx, done := listenForCancellationAndAddToContext()
	defer done()

	//Create Algorithm interface to use for JWT verification
	zapctx.Info(ctx, fmt.Sprintf("JWKS URL: %s", os.Getenv("JWKS")))
	const jwksURI = "https://identity.cloud.saltpay.dev/oauth/v2/oauth-anonymous/jwks"
	IH.SetAlgorithm(jwksURI, ctx)
	app, err := newApp(ctx)
	if err != nil {
		zapctx.Fatal(ctx, "failed to create the app", zap.Error(err))
	}

	// Add error log hook (for alerting)
	zapctx.AddHooks(zap2.NewPrometheusZapHook())

	// Listen for Luigi's exported transactions
	err = app.ListenForTransactions(ctx)
	if err != nil {
		zapctx.Error(ctx, fmt.Sprintf("Unable to start TransactionListener: %s", err))
	}

	// start the web server
	server := http.NewWebServer(app.ServerConfig, app.DBHandler)

	zapctx.Debug(ctx, fmt.Sprintf("Started. Listening on port: %s", app.ServerConfig.Port))
	if err := server.ListenAndServe(); err != nil {
		zapctx.Fatal(ctx, "http server listen failed", zap.Error(err))
	}
}
