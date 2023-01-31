//go:generate moq -out mocks/schema_registry_moq.go -pkg=mocks . SchemaRegistry

package ports

import (
	"context"
)

type SchemaRegistry interface {
	Decode(ctx context.Context, buf []byte, schemaKey string) (interface{}, []byte, error)
}
