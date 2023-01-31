//go:build acceptance_mock
// +build acceptance_mock

package acceptance_mock

import (
	"github.com/saltpay/transaction-api-operational/src-mario-api/black-box-tests/utils"
	"testing"

	"github.com/matryer/is"
)

const fiveRetries = 5

// Sanity check.
func TestGreetingApplication(t *testing.T) {
	client := utils.NewAPIClient(utils.GetBaseURL(t), t)

	if err := client.WaitForAPIToBeHealthy(fiveRetries); err != nil {
		t.Fatal(err)
	}
}

func TestTransactionEndpoints(t *testing.T) {
	is := is.New(t)
	client := utils.NewAPIClient(utils.GetBaseURL(t), t)

	internalID := "some_transactionID"
	transaction, err := client.GetTransaction(internalID)

	is.NoErr(err)
	is.Equal(transaction.Transactions[0].InternalID, internalID)
}

func TestTransactionsEndpoints(t *testing.T) {
	is := is.New(t)
	client := utils.NewAPIClient(utils.GetBaseURL(t), t)

	after := "after_string"
	before := "before_string"
	endDate := "end_date_string"
	startDate := "start_date_string"
	storeID := "store_id_string"
	limit := 10

	transactions, err := client.GetTransactions(after, before, endDate, startDate, storeID, limit)
	is.Equal(err, nil)
	is.Equal(len(transactions.Transactions), limit)
	is.Equal(transactions.Transactions[0].StoreID, storeID)
}
