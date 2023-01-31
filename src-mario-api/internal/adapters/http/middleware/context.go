package middleware

import (
	"net/http"

	zctx_hlprs "github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/logger/zapctx_helpers"
)

// ContextMiddleware adds service context to zapcontext.
func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = zctx_hlprs.WithCommonFields(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
