//go:build integration
// +build integration

package postgres_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/matryer/is"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/postgres"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/testhelpers"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/models"
)

const (
	connectionString = "postgres://postgres:somePassword@postgres:5432"
)

func TestDBHandler(t *testing.T) {
	metricsClient := testhelpers.DummyMetricsClient{}
	t.Run("should successfully connect when database is available", func(t *testing.T) {
		is := is.New(t)

		_, err := postgres.NewDBHandler(connectionString, metricsClient, 3)
		is.NoErr(err)
	})

	t.Run("should fail when database is not available after polling cycle finishes", func(t *testing.T) {
		is := is.New(t)

		_, err := postgres.NewDBHandler("", metricsClient, 1)
		is.True(err != nil)
	})
}

func TestDBHandler_StoreGet(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	metricsClient := testhelpers.DummyMetricsClient{}
	dbHandler, err := postgres.NewDBHandler(connectionString, metricsClient, 3)
	is.NoErr(err)

	t.Run("successfully stores and get a valid transaction", func(t *testing.T) {
		transaction := models.Transaction{}
		transactionJSON, err := json.Marshal(transaction)
		is.NoErr(err)

		err = dbHandler.Store(ctx, transactionJSON)
		is.NoErr(err)
	})
}

func TestDBHandler_Upsert(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	metricsClient := testhelpers.DummyMetricsClient{}
	dbHandler, err := postgres.NewDBHandler(connectionString, metricsClient, 3)
	is.NoErr(err)

	t.Run("successfully stores and get a valid transaction", func(t *testing.T) {
		entry := map[string]interface{}{
			"ID":                100,
			"AMND_STATE":        "foo",
			"AMND_DATE":         "2022-05-02 15:58:54.000",
			"AMND_OFFICER":      0,
			"AMND_PREV":         0,
			"NAME":              "foo's",
			"CODE":              "foo",
			"TERM_CAT":          "foo",
			"CATEGORY_CODE":     "foo",
			"CONDITION_DETAILS": "foo",
			"DEFAULT_CONDITION": 0,
			"LATE_CONDITION":    0,
			"SECURITY_CODE":     "foo",
			"ADDENDUM":          "foo",
		}

		sample := map[string]interface{}{"after": entry}
		err = dbHandler.Upsert(ctx, sample, "trans_cond")
		is.NoErr(err)

		err = dbHandler.Upsert(ctx, sample, "trans_cond")
		is.NoErr(err)
	})
}

func TestDBHandler_Get(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()

	metricsClient := testhelpers.DummyMetricsClient{}
	dbHandler, err := postgres.NewDBHandler(connectionString, metricsClient, 3)
	is.NoErr(err)

	t.Run("successfully retrieves N elements", func(t *testing.T) {
		entry := map[string]interface{}{
			"ID":                100,
			"AMND_STATE":        "foo",
			"AMND_DATE":         "2022-05-02 15:58:54.000",
			"AMND_OFFICER":      0,
			"AMND_PREV":         0,
			"NAME":              "foo",
			"CODE":              "foo",
			"TERM_CAT":          "foo",
			"CATEGORY_CODE":     "foo",
			"CONDITION_DETAILS": "foo",
			"DEFAULT_CONDITION": 0,
			"LATE_CONDITION":    0,
			"SECURITY_CODE":     "foo",
			"ADDENDUM":          "foo",
		}

		sample := map[string]interface{}{"after": entry}
		err = dbHandler.Upsert(ctx, sample, "trans_cond")
		is.NoErr(err)

		query := "SELECT id, name, code FROM trans_cond where id = 100"
		var name, code string
		var id int
		target := []interface{}{&id, &name, &code}

		err := dbHandler.Get(ctx, query, target)
		is.NoErr(err)
		is.Equal([]interface{}{100, "foo", "foo"}, []interface{}{id, name, code})
	})
}

func TestStoreAuxTable(t *testing.T) {
	var (
		expectedToken = "foo_bar_token"
		expectedDate  = "12/12"
	)
	is := is.New(t)
	ctx := context.Background()

	metricsClient := testhelpers.DummyMetricsClient{}
	dbHandler, err := postgres.NewDBHandler(connectionString, metricsClient, 3)
	is.NoErr(err)

	t.Run("happy path for aux table flow", func(t *testing.T) {
		transaction := models.W4IncomingTransaction{}
		transaction.After.ID = 1000
		transaction.After.CardExpire = expectedDate
		transaction.After.SaltTokenID = expectedToken

		err = dbHandler.StoreAuxTable(ctx, transaction, true)
		is.NoErr(err)

		transaction.After.ID = 1001
		transaction.After.DocPrevID = 1000
		err = dbHandler.StoreAuxTable(ctx, transaction, false)
		is.NoErr(err)

		query := "SELECT expiry_date, salt_token_id FROM aux_transactions_api where id = 1001"
		var expiry, token string
		target := []interface{}{&expiry, &token}

		err = dbHandler.Get(ctx, query, target)
		is.NoErr(err)
		is.Equal(expiry, expectedDate)
		is.Equal(token, expectedToken)
	})
}
