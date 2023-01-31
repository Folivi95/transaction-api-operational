package middleware

import (
	"net/http"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/logger/zapctx"
)

// ContextMiddleware adds service context to zapcontext.
func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = zapctx.WithCommonFields(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
