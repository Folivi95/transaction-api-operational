package models

type Transaction struct {
	Institutions            []Institution     `json:"institutions"`
	InternalID              string            `json:"internal_id"`
	HostTransactionID       string            `json:"host_transaction_id"`
	HostParentTransactionID string            `json:"host_parent_transaction_id"`
	UID                     string            `json:"uid"`
	Type                    string            `json:"type"`
	CardPresent             bool              `json:"card_present"`
	IsDcc                   bool              `json:"is_dcc"`
	EntryMethod             string            `json:"entry_method"`
	ProcessedAt             string            `json:"processed_at"`
	AuthCode                string            `json:"auth_code"`
	RRN                     string            `json:"rrn"`
	ARN                     string            `json:"arn"`
	CompanyID               string            `json:"company_id"`
	StoreID                 string            `json:"store_id"`
	MerchantID              string            `json:"merchant_id"`
	TransactionSource       TransactionSource `json:"transaction_source"`
	Response                Response          `json:"response"`
	Card                    Card              `json:"card"`
	Financial               Financial         `json:"financial"`
	Payout                  Payout            `json:"payout"`
}

type Response struct {
	Code           string `json:"code"`
	Classification string `json:"classification"`
	Description    string `json:"description"`
}

type TransactionSource struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
type Institution struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type InterchangeFee struct {
	FixedAmount      string `json:"fixed_amount,omitempty"`
	PercentageAmount string `json:"percentage_amount,omitempty"`
}

type Store struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	City            string `json:"city"`
	Country         string `json:"country"`
	ServiceCategory string `json:"service_category"`
	InternalGmdID   string `json:"internal_gmd_id"`
	CompanyID       string `json:"company_id"`
}

type CardReader struct {
	TerminalID string `json:"terminal_id"`
	Type       string `json:"type"`
}

type Card struct {
	PanFirstSix    string `json:"pan_first_six"`
	PanLastFour    string `json:"pan_last_four"`
	Brand          string `json:"brand"`
	ExpiryMonth    string `json:"expiry_month"`
	ExpiryYear     string `json:"expiry_year"`
	Issuer         string `json:"issuer"`
	IssuingCountry string `json:"issuing_country"`
	IssuingProduct string `json:"issuing_product"`
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
type Payout struct {
	ID          string `json:"id"`
	ScheduledAt string `json:"scheduled_at"`
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
type MerchantServiceCharge struct {
	FixedAmount      int    `json:"fixed_amount"`
	PercentageAmount string `json:"percentage_amount"`
	Currency         string `json:"currency"`
}
