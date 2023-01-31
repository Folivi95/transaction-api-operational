//go:build unit
// +build unit

package usecases_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/matryer/is"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/testhelpers"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/models"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/usecases"
)

func TestSolarTransformer(t *testing.T) {
	// Skip test
	t.Skip("Skipping test")
	is := is.New(t)
	expectedTxn := testhelpers.LoadCanonicalTransaction()

	solarTxnJSON := testhelpers.LoadSolarJSON()

	var solarTxn models.SolarIncomingTransaction
	err := json.Unmarshal(solarTxnJSON, &solarTxn)
	is.NoErr(err)
	txn, err := usecases.SolarTransformer{}.Translate(solarTxn)
	is.NoErr(err)
	fmt.Println(expectedTxn)
	fmt.Println(txn)
	is.Equal(expectedTxn, txn)
}
