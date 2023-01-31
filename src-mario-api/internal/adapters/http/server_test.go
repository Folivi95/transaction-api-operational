//go:build unit
// +build unit

package http_test

import (
	"github.com/saltpay/transaction-api-operational/src-mario-api/black-box-tests/utils"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/http"
	"os"
)

func TestNewWebServer(t *testing.T) {
	os.Setenv("jwks", "https://identity.cloud.saltpay.dev/oauth/v2/oauth-anonymous/jwks")
	is := is.New(t)
	webServer := http.NewWebServer(http.ServerConfig{}, nil)

	svr := httptest.NewServer(webServer.Handler)
	defer svr.Close()

	t.Run("Server starts up healthy", func(t *testing.T) {
		client := utils.NewAPIClient(svr.URL, t)
		err := client.CheckIfHealthy()
		is.NoErr(err)
	})
}
