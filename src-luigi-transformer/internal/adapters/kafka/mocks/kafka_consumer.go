// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/saltpay/go-kafka-driver"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/kafka"
	"sync"
)

// Ensure, that ConsumerMock does implement kafkalistener.Consumer.
// If this is not the case, regenerate this file with moq.
var _ kafkalistener.Consumer = &ConsumerMock{}

// ConsumerMock is a mock implementation of kafkalistener.Consumer.
//
// 	func TestSomethingThatUsesConsumer(t *testing.T) {
//
// 		// make and configure a mocked kafkalistener.Consumer
// 		mockedConsumer := &ConsumerMock{
// 			CloseFunc: func()  {
// 				panic("mock out the Close method")
// 			},
// 			CommitMessageFunc: func(ctx context.Context, msg kafka.Message) error {
// 				panic("mock out the CommitMessage method")
// 			},
// 			FetchMessageFunc: func(ctx context.Context) (kafka.Message, error) {
// 				panic("mock out the FetchMessage method")
// 			},
// 		}
//
// 		// use mockedConsumer in code that requires kafkalistener.Consumer
// 		// and then make assertions.
//
// 	}
type ConsumerMock struct {
	// CloseFunc mocks the Close method.
	CloseFunc func()

	// CommitMessageFunc mocks the CommitMessage method.
	CommitMessageFunc func(ctx context.Context, msg kafka.Message) error

	// FetchMessageFunc mocks the FetchMessage method.
	FetchMessageFunc func(ctx context.Context) (kafka.Message, error)

	// calls tracks calls to the methods.
	calls struct {
		// Close holds details about calls to the Close method.
		Close []struct {
		}
		// CommitMessage holds details about calls to the CommitMessage method.
		CommitMessage []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Msg is the msg argument value.
			Msg kafka.Message
		}
		// FetchMessage holds details about calls to the FetchMessage method.
		FetchMessage []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
	}
	lockClose         sync.RWMutex
	lockCommitMessage sync.RWMutex
	lockFetchMessage  sync.RWMutex
}

// Close calls CloseFunc.
func (mock *ConsumerMock) Close() {
	if mock.CloseFunc == nil {
		panic("ConsumerMock.CloseFunc: method is nil but Consumer.Close was just called")
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
//     len(mockedConsumer.CloseCalls())
func (mock *ConsumerMock) CloseCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockClose.RLock()
	calls = mock.calls.Close
	mock.lockClose.RUnlock()
	return calls
}

// CommitMessage calls CommitMessageFunc.
func (mock *ConsumerMock) CommitMessage(ctx context.Context, msg kafka.Message) error {
	if mock.CommitMessageFunc == nil {
		panic("ConsumerMock.CommitMessageFunc: method is nil but Consumer.CommitMessage was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Msg kafka.Message
	}{
		Ctx: ctx,
		Msg: msg,
	}
	mock.lockCommitMessage.Lock()
	mock.calls.CommitMessage = append(mock.calls.CommitMessage, callInfo)
	mock.lockCommitMessage.Unlock()
	return mock.CommitMessageFunc(ctx, msg)
}

// CommitMessageCalls gets all the calls that were made to CommitMessage.
// Check the length with:
//     len(mockedConsumer.CommitMessageCalls())
func (mock *ConsumerMock) CommitMessageCalls() []struct {
	Ctx context.Context
	Msg kafka.Message
} {
	var calls []struct {
		Ctx context.Context
		Msg kafka.Message
	}
	mock.lockCommitMessage.RLock()
	calls = mock.calls.CommitMessage
	mock.lockCommitMessage.RUnlock()
	return calls
}

// FetchMessage calls FetchMessageFunc.
func (mock *ConsumerMock) FetchMessage(ctx context.Context) (kafka.Message, error) {
	if mock.FetchMessageFunc == nil {
		panic("ConsumerMock.FetchMessageFunc: method is nil but Consumer.FetchMessage was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockFetchMessage.Lock()
	mock.calls.FetchMessage = append(mock.calls.FetchMessage, callInfo)
	mock.lockFetchMessage.Unlock()
	return mock.FetchMessageFunc(ctx)
}

// FetchMessageCalls gets all the calls that were made to FetchMessage.
// Check the length with:
//     len(mockedConsumer.FetchMessageCalls())
func (mock *ConsumerMock) FetchMessageCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockFetchMessage.RLock()
	calls = mock.calls.FetchMessage
	mock.lockFetchMessage.RUnlock()
	return calls
}
