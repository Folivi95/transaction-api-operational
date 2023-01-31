package schemaregistry_test

import (
	"context"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/riferrei/srclient"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/schemaregistry"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/schemaregistry/mocks"
)

const validSchema = `{
	"type": "record",
	"namespace": "namespace", 
	"name": "name",
	"doc": "Some data",
	"fields": [
		{"name": "f1", "type": "string"}
	]
}`

func TestClient_GetNewCodec(t *testing.T) {
	is := is.New(t)

	t.Run("client should fetch schema only once if for identical requests", func(t *testing.T) {
		registryClient := schemaregistry.NewSchemaRegistryClient("some_endpoint", 5)
		mockedRegistry := &mocks.SchemaRegistryMock{
			GetLatestSchemaFunc: func(subject string) (*srclient.Schema, error) {
				schema, _ := srclient.NewSchema(1, validSchema, "PROTOBUF", 2, nil, nil, nil)
				return schema, nil
			},
		}
		registryClient.Session = mockedRegistry

		_, err := registryClient.GetNewCodec(context.TODO(), "some_key")
		is.NoErr(err)
		is.Equal(len(mockedRegistry.GetLatestSchemaCalls()), 1)

		_, err = registryClient.GetNewCodec(context.TODO(), "some_key")
		is.NoErr(err)
		is.Equal(len(mockedRegistry.GetLatestSchemaCalls()), 1)
	})

	t.Run("client should refresh schema after some time", func(t *testing.T) {
		authCalled := make(chan bool)
		registryClient := schemaregistry.NewSchemaRegistryClient("some_endpoint", 1)
		mockedRegistry := &mocks.SchemaRegistryMock{
			GetLatestSchemaFunc: func(subject string) (*srclient.Schema, error) {
				schema, _ := srclient.NewSchema(1, validSchema, "PROTOBUF", 2, nil, nil, nil)
				return schema, nil
			},
		}
		mockedScheduler := &mocks.RefreshSchedulerMock{
			AfterFuncFunc: func(t time.Duration, f func()) *time.Timer {
				go func() {
					authCalled <- true
				}()
				return time.NewTimer(2 * time.Second)
			},
		}
		registryClient.Session = mockedRegistry
		registryClient.Scheduler = mockedScheduler

		_, err := registryClient.GetNewCodec(context.TODO(), "some_key")
		is.NoErr(err)
		is.Equal(len(mockedRegistry.GetLatestSchemaCalls()), 1)

		select {
		case <-authCalled:
			is.Equal(len(mockedRegistry.GetLatestSchemaCalls()), 2)
		case <-time.After(2000 * time.Millisecond):
			t.Fatal("timed out before refreshing tokens")
		}
	})
}
