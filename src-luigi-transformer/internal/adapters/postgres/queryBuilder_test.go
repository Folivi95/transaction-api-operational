//go:build unit

package postgres_test

import (
	"testing"

	"github.com/matryer/is"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/postgres"
)

const expectedQuery = `INSERT INTO tableName (ADDENDUM) VALUES ($escape$foo's$escape$) ON CONFLICT (ID) DO UPDATE SET ADDENDUM = $escape$foo's$escape$;`

func TestQueryBuilder(t *testing.T) {
	is := is.New(t)
	t.Run("test single quote input", func(t *testing.T) {
		entry := map[string]interface{}{
			"ADDENDUM": "foo's",
		}
		sample := map[string]interface{}{"after": entry}
		query, err := postgres.QueryBuilder("tableName", "ID", sample)
		is.NoErr(err)
		is.Equal(query, expectedQuery) // todo: validate query with more parameters (builder changes order every run)
	})
}
