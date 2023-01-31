//go:build acceptance
// +build acceptance

package acceptance

import (
	"github.com/saltpay/transaction-api-operational/src-mario-api/black-box-tests/utils"
	"testing"

	"github.com/matryer/is"
)

const fiveRetries = 5

// Sanity check.
func TestGreetingApplication(t *testing.T) {
	client := utils.NewAPIClient(getBaseURL(t), t)

	if err := client.WaitForAPIToBeHealthy(fiveRetries); err != nil {
		t.Fatal(err)
	}
}

func TestTransactionEndpoints(t *testing.T) {
	is := is.New(t)
	client := utils.NewAPIClient(GetBaseURL(t), t)

	acquirer := "some_acquirer"
	transactionId := "some_transactionID"

	_, err := client.GetTransaction(acquirer, transactionId)

	is.NoErr(err)
	//is.Equal(transaction.Transaction.ID, transactionID)
}

func TestTransactionsEndpoints(t *testing.T) {
	is := is.New(t)
	client := utils.NewAPIClient(GetBaseURL(t), t)

	acquirer := "some_acquirer"
	cardAcceptorID := "some_cardAcceptorID"
	limit := 10
	_, err := client.GetTransactions(acquirer, cardAcceptorID, limit)
	is.NoErr(err)
	//is.Equal(len(transaction.Transactions), limit)
	//is.Equal(transaction.Transactions[0].CardAcceptor.ID, merchantID)
}
