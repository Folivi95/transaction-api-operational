// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/application/ports"
	"sync"
)

// Ensure, that SchemaRegistryMock does implement ports.SchemaRegistry.
// If this is not the case, regenerate this file with moq.
var _ ports.SchemaRegistry = &SchemaRegistryMock{}

// SchemaRegistryMock is a mock implementation of ports.SchemaRegistry.
//
// 	func TestSomethingThatUsesSchemaRegistry(t *testing.T) {
//
// 		// make and configure a mocked ports.SchemaRegistry
// 		mockedSchemaRegistry := &SchemaRegistryMock{
// 			DecodeFunc: func(ctx context.Context, msg []byte, schemaKey string) (interface{}, bool) {
// 				panic("mock out the Decode method")
// 			},
// 		}
//
// 		// use mockedSchemaRegistry in code that requires ports.SchemaRegistry
// 		// and then make assertions.
//
// 	}
type SchemaRegistryMock struct {
	// DecodeFunc mocks the Decode method.
	DecodeFunc func(ctx context.Context, msg []byte, schemaKey string) (interface{}, bool)

	// calls tracks calls to the methods.
	calls struct {
		// Decode holds details about calls to the Decode method.
		Decode []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Msg is the msg argument value.
			Msg []byte
			// SchemaKey is the schemaKey argument value.
			SchemaKey string
		}
	}
	lockDecode sync.RWMutex
}

// Decode calls DecodeFunc.
func (mock *SchemaRegistryMock) Decode(ctx context.Context, msg []byte, schemaKey string) (interface{}, bool) {
	if mock.DecodeFunc == nil {
		panic("SchemaRegistryMock.DecodeFunc: method is nil but SchemaRegistry.Decode was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		Msg       []byte
		SchemaKey string
	}{
		Ctx:       ctx,
		Msg:       msg,
		SchemaKey: schemaKey,
	}
	mock.lockDecode.Lock()
	mock.calls.Decode = append(mock.calls.Decode, callInfo)
	mock.lockDecode.Unlock()
	return mock.DecodeFunc(ctx, msg, schemaKey)
}

// DecodeCalls gets all the calls that were made to Decode.
// Check the length with:
//     len(mockedSchemaRegistry.DecodeCalls())
func (mock *SchemaRegistryMock) DecodeCalls() []struct {
	Ctx       context.Context
	Msg       []byte
	SchemaKey string
} {
	var calls []struct {
		Ctx       context.Context
		Msg       []byte
		SchemaKey string
	}
	mock.lockDecode.RLock()
	calls = mock.calls.Decode
	mock.lockDecode.RUnlock()
	return calls
}