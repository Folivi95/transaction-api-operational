//go:build unit
// +build unit

package middleware_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	is2 "github.com/matryer/is"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/http/middleware"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/logger/zapctx"
)

func TestLoggingMiddleware(t *testing.T) {
	loggedHTTPRequest := httptest.NewRequest(http.MethodDelete, "/some/not/blacklisted/url", http.NoBody)

	testCases := []struct {
		name         string
		blacklist    []string
		req          *http.Request
		wantLogEntry []observer.LoggedEntry
	}{
		{
			name: "should log request that does not match any of the blacklist routes",
			req:  loggedHTTPRequest,
			wantLogEntry: []observer.LoggedEntry{
				{
					Entry: zapcore.Entry{
						Level:   zap.DebugLevel,
						Message: fmt.Sprintf("%s - %s", loggedHTTPRequest.Method, loggedHTTPRequest.URL),
					},
					Context: []zapcore.Field{},
				},
			},
		},
		{
			name:         "should not log /internal routes",
			req:          httptest.NewRequest(http.MethodGet, "/internal/something", http.NoBody),
			wantLogEntry: []observer.LoggedEntry{},
		},
		{
			name:         "should not log configured blacklist routes",
			blacklist:    []string{"/blacklisted"},
			req:          httptest.NewRequest(http.MethodGet, "/blacklisted/something", http.NoBody),
			wantLogEntry: []observer.LoggedEntry{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// setup test logger
			core, logs := observer.New(zap.DebugLevel)
			logger := zap.New(core)

			// add test logger to context
			ctx := zapctx.With(context.Background(), logger)
			req := tc.req.WithContext(ctx)

			is := is2.New(t)
			h := middleware.LoggingMiddleware(tc.blacklist...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)
			is.Equal(logs.AllUntimed(), tc.wantLogEntry)
		})
	}
}
