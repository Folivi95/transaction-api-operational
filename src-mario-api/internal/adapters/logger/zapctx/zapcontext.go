package zapctx

import (
	"context"
	"os"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type key struct{}

// hooks holds the functions to be added as hooks for the logger.
var hooks []func(entry zapcore.Entry) error

// From returns a logger if it is present in the context, otherwise it returns the default logger.
func From(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(key{}).(*zap.Logger); ok {
		return l
	}

	// set log level from environment
	atom := zap.NewAtomicLevel()
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		atom.SetLevel(zap.DebugLevel)
	case "info":
		atom.SetLevel(zap.InfoLevel)
	case "error":
		atom.SetLevel(zap.ErrorLevel)
	default:
		atom.SetLevel(zap.InfoLevel)
	}

	return zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zapcore.EncoderConfig{
				TimeKey:        "@timestamp",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "msg",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.NanosDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			}),
			zapcore.AddSync(os.Stdout),
			atom,
		),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.DPanicLevel))
}

// WithFields will return a context with a new logger with fields.
func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	return With(ctx, From(ctx).With(fields...))
}

// With will add l to the context.
func With(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, key{}, l)
}

func WithCommonFields(ctx context.Context) context.Context {
	env := os.Getenv("ENV")
	service := os.Getenv("SERVICE")

	return WithFields(ctx,
		zap.String("env", env),
		zap.String("service", service),
		zap.Stringer("trace.id", uuid.New()),
	)
}

// Hook defines the required signature for your log entries hooks.
type Hook func(entry zapcore.Entry) error

// AddHooks allows you to add zap hooks to be executed when a logging operation happens. Keep in mind that
// your hooks can impact performance and that this is not safe to be used concurrently(go routine), nor when already logging, so use
// this when launching the application and configuring the logger.
func AddHooks(hks ...Hook) {
	for _, h := range hks {
		hooks = append(hooks, h)
	}
}
