package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"strings"

	zapctx "github.com/saltpay/go-zap-ctx"
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
			zapctx.Info(
				r.Context(),
				"%s - %s",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
			)
			next.ServeHTTP(w, r)
		})
	}
}
