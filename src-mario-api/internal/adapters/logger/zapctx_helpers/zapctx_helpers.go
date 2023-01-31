package zapctxhelpers

import (
	"context"
	"os"

	"github.com/google/uuid"
	"go.uber.org/zap"

	zapctx "github.com/saltpay/go-zap-ctx"
)

func WithCommonFields(ctx context.Context) context.Context {
	env := os.Getenv("ENV")
	service := os.Getenv("SERVICE")

	return zapctx.WithFields(ctx,
		zap.String("env", env),
		zap.String("service", service),
		zap.Stringer("trace.id", uuid.New()),
	)
}
