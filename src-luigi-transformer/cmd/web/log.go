package main

import (
	logger "github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/logger/prometheus"

	zapctx "github.com/saltpay/go-zap-ctx"
	"go.uber.org/zap/zapcore"
)

func initLog(appConfig AppConfig) {
	switch appConfig.LogLevel {
	case "debug":
		zapctx.SetLogLevel(zapcore.DebugLevel)
	case "info":
		zapctx.SetLogLevel(zapcore.InfoLevel)
	case "error":
		zapctx.SetLogLevel(zapcore.ErrorLevel)
	}

	hook := logger.NewPrometheusHook()
	zapctx.AddHooks(hook)
}
