// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/application/models"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/application/ports"
	"sync"
)

// Ensure, that TransactionsSourceMock does implement ports.TransactionsSource.
// If this is not the case, regenerate this file with moq.
var _ ports.TransactionsSource = &TransactionsSourceMock{}

// TransactionsSourceMock is a mock implementation of ports.TransactionsSource.
//
// 	func TestSomethingThatUsesTransactionsSource(t *testing.T) {
//
// 		// make and configure a mocked ports.TransactionsSource
// 		mockedTransactionsSource := &TransactionsSourceMock{
// 			GetTransactionEventsByInternalIDFunc: func(ctx context.Context, request models.GetTransactionEventsByInternalIDRequest) ([]models.Transaction, error) {
// 				panic("mock out the GetTransactionEventsByInternalID method")
// 			},
// 			GetTransactionsByStoreIDFunc: func(ctx context.Context, request models.GetAllTransactionsRequest) ([]models.Transaction, string, string, error) {
// 				panic("mock out the GetTransactionsByStoreID method")
// 			},
// 			SaveTransactionFunc: func(ctx context.Context, transaction models.Transaction) error {
// 				panic("mock out the SaveTransaction method")
// 			},
// 		}
//
// 		// use mockedTransactionsSource in code that requires ports.TransactionsSource
// 		// and then make assertions.
//
// 	}
type TransactionsSourceMock struct {
	// GetTransactionEventsByInternalIDFunc mocks the GetTransactionEventsByInternalID method.
	GetTransactionEventsByInternalIDFunc func(ctx context.Context, request models.GetTransactionEventsByInternalIDRequest) ([]models.Transaction, error)

	// GetTransactionsByStoreIDFunc mocks the GetTransactionsByStoreID method.
	GetTransactionsByStoreIDFunc func(ctx context.Context, request models.GetAllTransactionsRequest) ([]models.Transaction, string, string, error)

	// SaveTransactionFunc mocks the SaveTransaction method.
	SaveTransactionFunc func(ctx context.Context, transaction models.Transaction) error

	// calls tracks calls to the methods.
	calls struct {
		// GetTransactionEventsByInternalID holds details about calls to the GetTransactionEventsByInternalID method.
		GetTransactionEventsByInternalID []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Request is the request argument value.
			Request models.GetTransactionEventsByInternalIDRequest
		}
		// GetTransactionsByStoreID holds details about calls to the GetTransactionsByStoreID method.
		GetTransactionsByStoreID []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Request is the request argument value.
			Request models.GetAllTransactionsRequest
		}
		// SaveTransaction holds details about calls to the SaveTransaction method.
		SaveTransaction []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Transaction is the transaction argument value.
			Transaction models.Transaction
		}
	}
	lockGetTransactionEventsByInternalID sync.RWMutex
	lockGetTransactionsByStoreID         sync.RWMutex
	lockSaveTransaction                  sync.RWMutex
}

// GetTransactionEventsByInternalID calls GetTransactionEventsByInternalIDFunc.
func (mock *TransactionsSourceMock) GetTransactionEventsByInternalID(ctx context.Context, request models.GetTransactionEventsByInternalIDRequest) ([]models.Transaction, error) {
	if mock.GetTransactionEventsByInternalIDFunc == nil {
		panic("TransactionsSourceMock.GetTransactionEventsByInternalIDFunc: method is nil but TransactionsSource.GetTransactionEventsByInternalID was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Request models.GetTransactionEventsByInternalIDRequest
	}{
		Ctx:     ctx,
		Request: request,
	}
	mock.lockGetTransactionEventsByInternalID.Lock()
	mock.calls.GetTransactionEventsByInternalID = append(mock.calls.GetTransactionEventsByInternalID, callInfo)
	mock.lockGetTransactionEventsByInternalID.Unlock()
	return mock.GetTransactionEventsByInternalIDFunc(ctx, request)
}

// GetTransactionEventsByInternalIDCalls gets all the calls that were made to GetTransactionEventsByInternalID.
// Check the length with:
//     len(mockedTransactionsSource.GetTransactionEventsByInternalIDCalls())
func (mock *TransactionsSourceMock) GetTransactionEventsByInternalIDCalls() []struct {
	Ctx     context.Context
	Request models.GetTransactionEventsByInternalIDRequest
} {
	var calls []struct {
		Ctx     context.Context
		Request models.GetTransactionEventsByInternalIDRequest
	}
	mock.lockGetTransactionEventsByInternalID.RLock()
	calls = mock.calls.GetTransactionEventsByInternalID
	mock.lockGetTransactionEventsByInternalID.RUnlock()
	return calls
}

// GetTransactionsByStoreID calls GetTransactionsByStoreIDFunc.
func (mock *TransactionsSourceMock) GetTransactionsByStoreID(ctx context.Context, request models.GetAllTransactionsRequest) ([]models.Transaction, string, string, error) {
	if mock.GetTransactionsByStoreIDFunc == nil {
		panic("TransactionsSourceMock.GetTransactionsByStoreIDFunc: method is nil but TransactionsSource.GetTransactionsByStoreID was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Request models.GetAllTransactionsRequest
	}{
		Ctx:     ctx,
		Request: request,
	}
	mock.lockGetTransactionsByStoreID.Lock()
	mock.calls.GetTransactionsByStoreID = append(mock.calls.GetTransactionsByStoreID, callInfo)
	mock.lockGetTransactionsByStoreID.Unlock()
	return mock.GetTransactionsByStoreIDFunc(ctx, request)
}

// GetTransactionsByStoreIDCalls gets all the calls that were made to GetTransactionsByStoreID.
// Check the length with:
//     len(mockedTransactionsSource.GetTransactionsByStoreIDCalls())
func (mock *TransactionsSourceMock) GetTransactionsByStoreIDCalls() []struct {
	Ctx     context.Context
	Request models.GetAllTransactionsRequest
} {
	var calls []struct {
		Ctx     context.Context
		Request models.GetAllTransactionsRequest
	}
	mock.lockGetTransactionsByStoreID.RLock()
	calls = mock.calls.GetTransactionsByStoreID
	mock.lockGetTransactionsByStoreID.RUnlock()
	return calls
}

// SaveTransaction calls SaveTransactionFunc.
func (mock *TransactionsSourceMock) SaveTransaction(ctx context.Context, transaction models.Transaction) error {
	if mock.SaveTransactionFunc == nil {
		panic("TransactionsSourceMock.SaveTransactionFunc: method is nil but TransactionsSource.SaveTransaction was just called")
	}
	callInfo := struct {
		Ctx         context.Context
		Transaction models.Transaction
	}{
		Ctx:         ctx,
		Transaction: transaction,
	}
	mock.lockSaveTransaction.Lock()
	mock.calls.SaveTransaction = append(mock.calls.SaveTransaction, callInfo)
	mock.lockSaveTransaction.Unlock()
	return mock.SaveTransactionFunc(ctx, transaction)
}

// SaveTransactionCalls gets all the calls that were made to SaveTransaction.
// Check the length with:
//     len(mockedTransactionsSource.SaveTransactionCalls())
func (mock *TransactionsSourceMock) SaveTransactionCalls() []struct {
	Ctx         context.Context
	Transaction models.Transaction
} {
	var calls []struct {
		Ctx         context.Context
		Transaction models.Transaction
	}
	mock.lockSaveTransaction.RLock()
	calls = mock.calls.SaveTransaction
	mock.lockSaveTransaction.RUnlock()
	return calls
}