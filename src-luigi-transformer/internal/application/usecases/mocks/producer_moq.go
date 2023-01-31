// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/saltpay/go-kafka-driver"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/usecases"
	"sync"
)

// Ensure, that ProducerMock does implement usecases.Producer.
// If this is not the case, regenerate this file with moq.
var _ usecases.Producer = &ProducerMock{}

// ProducerMock is a mock implementation of usecases.Producer.
//
// 	func TestSomethingThatUsesProducer(t *testing.T) {
//
// 		// make and configure a mocked usecases.Producer
// 		mockedProducer := &ProducerMock{
// 			CloseFunc: func()  {
// 				panic("mock out the Close method")
// 			},
// 			WriteMessageFunc: func(ctx context.Context, message kafka.Message) error {
// 				panic("mock out the writeMessage method")
// 			},
// 		}
//
// 		// use mockedProducer in code that requires usecases.Producer
// 		// and then make assertions.
//
// 	}
type ProducerMock struct {
	// CloseFunc mocks the Close method.
	CloseFunc func()

	// WriteMessageFunc mocks the writeMessage method.
	WriteMessageFunc func(ctx context.Context, message kafka.Message) error

	// calls tracks calls to the methods.
	calls struct {
		// Close holds details about calls to the Close method.
		Close []struct {
		}
		// writeMessage holds details about calls to the writeMessage method.
		WriteMessage []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Message is the message argument value.
			Message kafka.Message
		}
	}
	lockClose        sync.RWMutex
	lockWriteMessage sync.RWMutex
}

// Close calls CloseFunc.
func (mock *ProducerMock) Close() {
	if mock.CloseFunc == nil {
		panic("ProducerMock.CloseFunc: method is nil but Producer.Close was just called")
	}
	callInfo := struct {
	}{}
	mock.lockClose.Lock()
	mock.calls.Close = append(mock.calls.Close, callInfo)
	mock.lockClose.Unlock()
	mock.CloseFunc()
}

// CloseCalls gets all the calls that were made to Close.
// Check the length with:
//     len(mockedProducer.CloseCalls())
func (mock *ProducerMock) CloseCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockClose.RLock()
	calls = mock.calls.Close
	mock.lockClose.RUnlock()
	return calls
}

// writeMessage calls WriteMessageFunc.
func (mock *ProducerMock) WriteMessage(ctx context.Context, message kafka.Message) error {
	if mock.WriteMessageFunc == nil {
		panic("ProducerMock.WriteMessageFunc: method is nil but Producer.writeMessage was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Message kafka.Message
	}{
		Ctx:     ctx,
		Message: message,
	}
	mock.lockWriteMessage.Lock()
	mock.calls.WriteMessage = append(mock.calls.WriteMessage, callInfo)
	mock.lockWriteMessage.Unlock()
	return mock.WriteMessageFunc(ctx, message)
}

// WriteMessageCalls gets all the calls that were made to writeMessage.
// Check the length with:
//     len(mockedProducer.WriteMessageCalls())
func (mock *ProducerMock) WriteMessageCalls() []struct {
	Ctx     context.Context
	Message kafka.Message
} {
	var calls []struct {
		Ctx     context.Context
		Message kafka.Message
	}
	mock.lockWriteMessage.RLock()
	calls = mock.calls.WriteMessage
	mock.lockWriteMessage.RUnlock()
	return calls
}
