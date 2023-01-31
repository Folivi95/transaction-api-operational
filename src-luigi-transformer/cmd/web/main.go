package main

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/http"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/logger/zapctx"
)

func main() {
	ctx, done := listenForCancellationAndAddToContext()
	defer done()

	app, err := newApp(ctx)
	if err != nil {
		zapctx.From(ctx).Fatal("failed to create the app", zap.Error(err))
	}

	appConfig, err := loadAppConfig()
	if err != nil {
		zapctx.From(ctx).Fatal("failed to load config for the app", zap.Error(err))
	}
	initLog(appConfig)

	// start listening kafka transaction consumers
	if app.ListenTransactions {
		go app.W4Listener.Listen(ctx)
		go app.SolarListener.Listen(ctx)
	}

	// start the web server
	server := http.NewWebServer(app.ServerConfig)

	zapctx.From(ctx).Info(fmt.Sprintf("Started. Listening on port: %s", app.ServerConfig.Port))
	if err := server.ListenAndServe(); err != nil {
		zapctx.From(ctx).Fatal("http server listen failed", zap.Error(err))
	}
}
