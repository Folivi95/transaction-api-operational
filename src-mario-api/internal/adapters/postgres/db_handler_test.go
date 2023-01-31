//go:build integration
// +build integration

package postgres_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/matryer/is"

	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/application/models"
)

const (
	connectionString = "postgres://postgres:somePassword@postgres:5432"
	storeQuery       = "INSERT INTO transactions_api(transaction) VALUES($1)"
)

func TestDBHandler(t *testing.T) {
	t.Run("should successfully connect when database is available", func(t *testing.T) {
		is := is.New(t)

		_, err := postgres.NewDBHandler(connectionString, nil, 3)
		is.NoErr(err)
	})

	t.Run("should fail when database is not available after polling cycle finishes", func(t *testing.T) {
		is := is.New(t)

		_, err := postgres.NewDBHandler("", nil, 1)
		is.True(err != nil)
	})
}

func TestDBHandler_StoreGet(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	dbHandler, err := postgres.NewDBHandler(connectionString, nil, 3)
	is.NoErr(err)

	err = postgres.CheckDBVersion(dbHandler.DB, 1)
	is.NoErr(err)
	t.Run("successfully stores and get a valid transaction", func(t *testing.T) {
		expectedTransaction := models.Transaction{}

		expectedTransactionJSON, err := json.Marshal(expectedTransaction)
		is.NoErr(err)

		err = dbHandler.Store(ctx, storeQuery, []interface{}{expectedTransactionJSON})
		is.NoErr(err)

		query := "SELECT transaction FROM transactions_api WHERE id = $1"
		response, err := dbHandler.Get(ctx, query, []interface{}{1})
		is.NoErr(err)

		transaction, err := models.NewTransactionFromJSON(response)
		is.NoErr(err)
		is.Equal(transaction, expectedTransaction)
	})

	t.Run("fails if query and arguments are not matching", func(t *testing.T) {
		transaction := models.Transaction{}

		err = dbHandler.Store(ctx, storeQuery, []interface{}{})
		is.True(err != nil)

		query := "INSERT INTO transactions_api(id, transaction) VALUES($1)"
		err = dbHandler.Store(ctx, query, []interface{}{"foobar", transaction})
		is.True(err != nil)

		query = "INSERT INTO transactions_api(id, transaction) VALUES($1)"
		err = dbHandler.Store(ctx, query, []interface{}{transaction, "foobar"})
		is.True(err != nil)
	})
}

func TestDBHandler_GetTransactionByID(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	dbHandler, err := postgres.NewDBHandler(connectionString, nil, 3)
	is.NoErr(err)

	t.Run("happy path for querying by id", func(t *testing.T) {
		expectedTransaction := models.Transaction{ID: "some_id"}

		expectedTransactionJSON, err := json.Marshal(expectedTransaction)
		is.NoErr(err)

		err = dbHandler.Store(ctx, storeQuery, []interface{}{expectedTransactionJSON})
		is.NoErr(err)

		response, err := dbHandler.GetTransaction(ctx, models.GetTransactionEventsByInternalIDRequest{TransactionID: "some_id"})
		is.NoErr(err)

		is.Equal(response, expectedTransaction)
	})
}

func TestDBHandler_GetTransactionByMID(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	dbHandler, err := postgres.NewDBHandler(connectionString, nil, 3)
	is.NoErr(err)

	t.Run("happy path for querying by card acceptor id", func(t *testing.T) {
		expectedTransaction1 := models.Transaction{ID: "some_id", CardAcceptor: models.CardAcceptor{ID: "123"}}
		expectedTransaction2 := models.Transaction{ID: "another_id", CardAcceptor: models.CardAcceptor{ID: "123"}}

		expectedTransactionJSON, err := json.Marshal(expectedTransaction1)
		is.NoErr(err)
		err = dbHandler.Store(ctx, storeQuery, []interface{}{expectedTransactionJSON})
		is.NoErr(err)
		expectedTransactionJSON2, err := json.Marshal(expectedTransaction2)
		is.NoErr(err)
		err = dbHandler.Store(ctx, storeQuery, []interface{}{expectedTransactionJSON2})
		is.NoErr(err)

		response, _, token, err := dbHandler.GetTransactions(ctx, models.GetAllTransactionsRequest{
			CardAcceptorID: "123",
			AcquiringHost:  "",
			Limit:          1,
		})

		is.NoErr(err)
		is.Equal(response, []models.Transaction{expectedTransaction2})

		response, prevToken, _, err := dbHandler.GetTransactions(ctx, models.GetAllTransactionsRequest{
			CardAcceptorID: "123",
			AcquiringHost:  "",
			After:          token,
			Limit:          1,
		})
		is.NoErr(err)
		is.Equal(response, []models.Transaction{expectedTransaction1})

		response, _, _, err = dbHandler.GetTransactions(ctx, models.GetAllTransactionsRequest{
			CardAcceptorID: "123",
			AcquiringHost:  "",
			Before:         prevToken,
			Limit:          1,
		})
		is.NoErr(err)
		is.Equal(response, []models.Transaction{expectedTransaction2})
	})
}
