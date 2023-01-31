package main

import (
	"os"
	"strconv"

	"go.uber.org/zap"

	"github.com/saltpay/transaction-api-operational/src-kamek-cronjob/internal/adapters/logger/zapctx"
)

func main() {
	ctx, done := listenForCancellationAndAddToContext()

	app, err := newApp(ctx)
	if err != nil {
		zapctx.From(ctx).Fatal("failed to create the app", zap.Error(err))
	}

	interval, err := strconv.Atoi(os.Getenv("DB_CLEAN_INTERVAL"))
	if err != nil {
		zapctx.From(ctx).Warn("no valid value for `DB_CLEAN_INTERVAL`, defaulting to 180")
		interval = 180
	}

	err = app.DBHandler.DeletePastXTransactions(ctx, interval)
	if err != nil {
		zapctx.From(ctx).Error("Failed to delete transactions", zap.Error(err))
	}

	done()
}
