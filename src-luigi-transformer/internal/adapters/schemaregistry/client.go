//go:generate moq -out mocks/refresh_scheduler_moq.go -pkg mocks . RefreshScheduler
//go:generate moq -out mocks/schema_registry_moq.go -pkg mocks . SchemaRegistry

package schemaregistry

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"strings"
	"time"

	"github.com/linkedin/goavro/v2"
	"github.com/riferrei/srclient"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/logger/zapctx"
)

type RefreshScheduler interface {
	AfterFunc(t time.Duration, f func()) *time.Timer
}

type SchedulerFunc func(d time.Duration, f func()) *time.Timer

func (s SchedulerFunc) AfterFunc(d time.Duration, f func()) *time.Timer {
	return s(d, f)
}

type SchemaRegistry interface {
	GetLatestSchema(subject string) (*srclient.Schema, error)
	GetSchemaByVersion(subject string, version int) (*srclient.Schema, error)
}

type Client struct {
	schemaTable map[string][]string
	Session     SchemaRegistry
	Scheduler   RefreshScheduler
}

func NewSchemaRegistryClient(endpoint string, refreshIntervalSeconds int) *Client {
	schemaTable := make(map[string][]string)
	session := srclient.CreateSchemaRegistryClient(endpoint)
	scheduler := SchedulerFunc(time.AfterFunc)

	client := &Client{
		schemaTable: schemaTable,
		Session:     session,
		Scheduler:   scheduler,
	}
	client.RefreshSchemasTrigger(refreshIntervalSeconds)

	return client
}

func (c *Client) Decode(ctx context.Context, buf []byte, schemaKey string) (interface{}, []byte, error) {
	codec, err := c.GetNewCodec(ctx, schemaKey)
	if err != nil {
		zapctx.From(ctx).Warn("[SchemaRegistryClient] Failed to get new codec.", zap.Error(err))
		return nil, nil, err
	}

	native, data, err := codec.NativeFromTextual(buf)
	if err != nil {
		zapctx.From(ctx).Warn("[SchemaRegistryClient] Failed to convert avro data to Go native type.", zap.Error(err))
		return nil, data, err
	}

	return native, data, nil
}

func (c *Client) GetNewCodec(ctx context.Context, schemaKey string) (goavro.Codec, error) {
	schema, err := c.getLatestSchema(ctx, schemaKey)
	if err != nil {
		zapctx.From(ctx).Warn("[SchemaRegistryClient] Failed to get latest schema.", zap.Error(err))
		return goavro.Codec{}, err
	}

	// Convert schema slice to string in Avro format
	schemaStr := "[" + strings.Join(schema, ",") + "]"

	// log schemaStr for debug purpose
	msg := fmt.Sprintf("[GetNewCodec] Schema string: %s", schemaStr)
	zapctx.From(ctx).Warn(msg)

	codec, err := goavro.NewCodecForStandardJSON(schemaStr)
	if err != nil {
		zapctx.From(ctx).Error("[registryHandler] Failed to create Codec for given schema.", zap.Error(err))
		return goavro.Codec{}, err
	}

	return *codec, err
}

func (c *Client) getLatestSchema(ctx context.Context, schemaKey string) ([]string, error) {
	schema, schemaRegistered := c.schemaTable[schemaKey]

	var err error
	if !schemaRegistered {
		schema, err = c.fetchSchema(ctx, schemaKey)
		if err != nil {
			zapctx.From(ctx).Warn("[SchemaRegistryClient] Failed to fetch schema.", zap.Error(err))
			return []string{}, err
		}
	}

	return schema, nil
}

func (c *Client) fetchSchema(ctx context.Context, schemaKey string) ([]string, error) {
	schema, err := c.Session.GetLatestSchema(schemaKey)
	if err != nil {
		return []string{}, err
	}

	var newSchema []string
	// Handling References
	if schema.References() != nil {
		for _, reference := range schema.References() {
			// Fetch reference schema by Version ID and appends to new constructor
			refSchema, err := c.Session.GetSchemaByVersion(reference.Subject, reference.Version)
			if err != nil {
				zapctx.From(ctx).Error("[registryHandler] Failed to get reference schema", zap.Error(err))
				continue
			}

			newSchema = append(newSchema, refSchema.Schema())
		}
	}

	// Appends envelope schema to entities schema
	newSchema = append(newSchema, schema.Schema())

	if len(newSchema) == 0 {
		zapctx.From(ctx).Error("[registryHandler] Empty schema for given subject name. Aborting.", zap.Error(err))
		return []string{}, err
	}

	c.schemaTable[schemaKey] = newSchema
	return newSchema, err
}

func (c *Client) RefreshSchemasTrigger(refreshIntervalSeconds int) {
	c.refreshSchemas(refreshIntervalSeconds)
}

func (c *Client) refreshSchemas(refreshIntervalSeconds int) {
	ctx := context.Background()
	for schemaKey := range c.schemaTable {
		newSchema, err := c.fetchSchema(ctx, schemaKey)
		if err != nil {
			zapctx.From(ctx).Error("[refreshSchemas] Failed to retrieve schema", zap.Error(err))
			continue
		}
		c.schemaTable[schemaKey] = newSchema
	}

	c.Scheduler.AfterFunc(time.Duration(refreshIntervalSeconds)*time.Second, func() {
		c.RefreshSchemasTrigger(refreshIntervalSeconds)
	})
}
