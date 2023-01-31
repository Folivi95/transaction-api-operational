package http

import (
	"net/http"

	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/application/ports"
)

func NewWebServer(config ServerConfig, transactionSource ports.TransactionsSource) *http.Server {
	return &http.Server{
		Addr:         config.TCPAddress(),
		Handler:      newRouter(transactionSource),
		ReadTimeout:  config.HTTPReadTimeout,
		WriteTimeout: config.HTTPWriteTimeout,
	}
}
