//go:build integration
// +build integration

package postgres_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/matryer/is"

	"github.com/saltpay/transaction-api-operational/src-kamek-cronjob/internal/adapters/postgres"
	"github.com/saltpay/transaction-api-operational/src-kamek-cronjob/internal/adapters/testhelpers"
)

const connectionString = "postgresql://postgres:postgres@postgresql:5432/postgres"

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

func TestDBHandler_DeletePastXTransactions(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	dbHandler, err := postgres.NewDBHandler(connectionString, nil, 3)
	is.NoErr(err)

	t.Run("successfully stores transaction with specific create_at timestamp and get a valid transaction then delete", func(t *testing.T) {
		expectedTransaction := testhelpers.Transaction{ID: "delete_this_id"}

		expectedTransactionJSON, err := json.Marshal(expectedTransaction)
		is.NoErr(err)

		query := "INSERT INTO transactions_api(transaction, created_at) VALUES($1, $2)"
		err = dbHandler.Store(ctx, query, []interface{}{expectedTransactionJSON, "2020-01-01"})
		is.NoErr(err)
		err = dbHandler.GetTransaction(ctx, "delete_this_id")
		is.NoErr(err)

		err = dbHandler.DeletePastXTransactions(ctx, 30)
		is.NoErr(err)

		err = dbHandler.GetTransaction(ctx, "delete_this_id")
		is.Equal(err, sql.ErrNoRows)
	})

	t.Run("should not fail if there aren't valid transactions", func(t *testing.T) {
		err = dbHandler.DeletePastXTransactions(ctx, 30)
		is.NoErr(err)
	})
}
