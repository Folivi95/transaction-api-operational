package http

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/http/handlers"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/http/middleware"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/application/ports"
)

func newRouter(transactionSource ports.TransactionsSource) *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.Prometheus)
	router.Use(middleware.ContextMiddleware)
	router.Use(middleware.LoggingMiddleware())
	monitoring(router)

	transactionHandler := handlers.TransactionHandler{TransactionSource: transactionSource}
	serveTransactions(router, transactionHandler)

	return router
}

func monitoring(router *mux.Router) {
	router.HandleFunc("/internal/health_check", handlers.HealthCheck)
	router.Handle("/internal/metrics", promhttp.Handler())
}

func serveTransactions(router *mux.Router, transactionHandler handlers.TransactionHandler) {
	router.HandleFunc("/transactions", transactionHandler.GetTransactions)
	router.HandleFunc("/transactions/{internal_id}", transactionHandler.GetTransaction)
}
