//go:build unit
// +build unit

package sync_test

import (
	"testing"

	"github.com/matryer/is"

	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/sync"
)

func TestAtomicBool_IsSet(t *testing.T) {
	is := is.New(t)
	t.Run("should be true", func(t *testing.T) {
		b := sync.New()
		b.Set()
		is.True(b.IsSet())
	})

	t.Run("should be false", func(t *testing.T) {
		b := sync.New()
		b.UnSet()
		is.True(!b.IsSet())
	})
}
