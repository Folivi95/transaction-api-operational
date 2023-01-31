package main

import (
	logger "github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/logger/prometheus"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/configuration"
	"strings"

	zapctx "github.com/saltpay/go-zap-ctx"
	"go.uber.org/zap/zapcore"
)

func initLog(appConfig configuration.AppConfig) {
	switch strings.ToLower(appConfig.LogLevel.LogLevel) {
	case "debug":
		zapctx.SetLogLevel(zapcore.DebugLevel)
	case "info":
		zapctx.SetLogLevel(zapcore.InfoLevel)
	case "error":
		zapctx.SetLogLevel(zapcore.ErrorLevel)
	case "warn":
		zapctx.SetLogLevel(zapcore.WarnLevel)
	}

	hook := logger.NewPrometheusHook()
	zapctx.AddHooks(hook)
}
