package models

type GetTransactionEventsByInternalIDRequest struct {
	AcquiringHost string `json:"acquiring_host"`
	InternalID    string `json:"internal_id"`
}

func NewGetTransactionEventsByInternalIDRequestFromParameters(internalID string) GetTransactionEventsByInternalIDRequest {
	return GetTransactionEventsByInternalIDRequest{
		InternalID: internalID,
	}
}

type GetTransactionEventsByInternalIDResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type GetAllTransactionsRequest struct {
	StoreID   string `json:"store_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	After     string `json:"after"`
	Before    string `json:"before"`
	Limit     int    `json:"limit"`
}

func NewGetAllTransactionsRequestFromParameters(storeID, startDate, endDate, after, before string, limit int) GetAllTransactionsRequest {
	return GetAllTransactionsRequest{
		StoreID:   storeID,
		After:     after,
		Before:    before,
		Limit:     limit,
		StartDate: startDate,
		EndDate:   endDate,
	}
}

type GetAllTransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
	Pagination   Pagination    `json:"pagination"`
}

type Pagination struct {
	Cursors Cursors `json:"cursors"`
	Paths   Paths   `json:"paths"`
}

type Cursors struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

type Paths struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

type TransactionErrorResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type InvalidCursors struct{}

func (InvalidCursors) Error() string {
	return "Request has two different cursors"
}

type HTTPResponse struct {
	Code              int               `json:"code"`
	Status            string            `json:"status"`
	RequestParameters RequestParameters `json:"request_parameters,omitempty"`
}

type RequestParameters struct {
	StoreID    string `json:"store_id,omitempty"`
	StartDate  string `json:"start_date,omitempty"`
	EndDate    string `json:"end_date,omitempty"`
	Before     string `json:"before,omitempty"`
	After      string `json:"after,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	InternalID string `json:"internal_id,omitempty"`
}
