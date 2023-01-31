package models

import "encoding/json"

type Transaction struct {
	Institutions            []Institution     `json:"institutions"`
	InternalID              string            `json:"internal_id"`
	HostTransactionID       string            `json:"host_transaction_id"`
	HostParentTransactionID string            `json:"host_parent_transaction_id"`
	Uid                     string            `json:"uid"`
	Type                    string            `json:"type"`
	CardPresent             string            `json:"card_present"`
	IsDCC                   string            `json:"is_dcc"`
	EntryMethod             string            `json:"entry_method"`
	ProcessedAt             string            `json:"processed_at"`
	AuthCode                string            `json:"auth_code"`
	Rrn                     string            `json:"rrn"`
	Arn                     string            `json:"arn"`
	CompanyID               string            `json:"company_id"`
	StoreID                 string            `json:"store_id"`
	MerchantID              string            `json:"merchant_id"`
	TransactionSource       TransactionSource `json:"transaction_source"`
	Response                Response          `json:"response"`
	Card                    Card              `json:"card"`
	Financial               Financial         `json:"financial"`
	Payout                  Payout            `json:"payout"`
}
type TransactionSource struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
type Card struct {
	PanFirstSix    string `json:"pan_first_six"`
	PanLastFour    string `json:"pan_last_four"`
	Brand          string `json:"brand"`
	ExpiryMonth    string `json:"expiry_month"`
	ExpiryYear     string `json:"expiry_year"`
	IssuingProduct string `json:"issuing_product"`
	IssuingCountry string `json:"issuing_country"`
	Issuer         string `json:"issuer"`
	TokenID        string `json:"token_id"`
	Currency       string `json:"currency"`
}

type Financial struct {
	TransactionAmount   int64  `json:"transaction_amount"`
	TransactionCurrency string `json:"transaction_currency"`
	BillingAmount       int64  `json:"billing_amount"`
	BillingCurrency     string `json:"billing_currency"`
	Fees                Fees   `json:"fees"`
}

type Fees struct {
	Total     int64       `json:"total"`
	Currency  string      `json:"currency"`
	Type      string      `json:"type"`
	Breakdown []Breakdown `json:"breakdown"`
}

type Breakdown struct {
	Type               string `json:"type"`
	FixedAmount        int64  `json:"fixed_amount"`
	VariablePercentage int64  `json:"variable_percentage"`
	VariableAmount     int64  `json:"variable_amount"`
	Currency           string `json:"currency"`
}

type Institution struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Payout struct {
	ID          string `json:"id"`
	ScheduledAt string `json:"scheduled_at"`
}

type Response struct {
	Code           string `json:"code"`
	Description    string `json:"description"`
	Classification string `json:"classification"`
}

func NewTransactionFromJSON(in []byte) (Transaction, error) {
	var transaction Transaction
	err := json.Unmarshal(in, &transaction)
	return transaction, err
}
