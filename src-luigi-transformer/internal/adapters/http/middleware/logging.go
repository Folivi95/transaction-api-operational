package middleware

import (
	"net/http"
	"strings"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/logger/zapctx"
)

// LoggingMiddleware logs the method and the url of the incoming request
// if the url does not begin with any of the blacklistPrefixes.
// A "/internal" blacklist is added by default.
func LoggingMiddleware(blacklistPrefixes ...string) func(http.Handler) http.Handler {
	blacklist := map[string]struct{}{
		"/internal": {},
	}
	for _, blp := range blacklistPrefixes {
		blacklist[blp] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for route := range blacklist {
				if strings.HasPrefix(r.URL.Path, route) {
					next.ServeHTTP(w, r)
					return
				}
			}
			zapctx.From(r.Context()).Sugar().Debugf("%s - %s", r.Method, r.URL)
			next.ServeHTTP(w, r)
		})
	}
}
