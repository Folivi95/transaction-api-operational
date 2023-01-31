//go:build unit
// +build unit

package http_test

import (
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/black-box-tests/acceptance"
	http2 "github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/http"
)

func TestNewWebServer(t *testing.T) {
	is := is.New(t)
	webServer := http2.NewWebServer(http2.ServerConfig{})

	svr := httptest.NewServer(webServer.Handler)
	defer svr.Close()

	t.Run("Server starts up healthy", func(t *testing.T) {
		client := acceptance.NewAPIClient(svr.URL, t)
		err := client.CheckIfHealthy()
		is.NoErr(err)
	})
}
