package testhelpers

import "time"

type Transaction struct {
	Institutions          []Institution `json:"institutions"`
	ID                    string        `json:"id"`
	Type                  string        `json:"type"`
	AuthorisationResponse string        `json:"authorisation_response"`
	ProcessedAt           time.Time     `json:"processed_at"`
	IsCardPresent         bool          `json:"is_card_present"`
	AuthorisationCode     string        `json:"authorisation_code"`
	IsDcc                 bool          `json:"is_dcc"`
	EntryMethod           string        `json:"entry_method"`
	CardAcceptor          CardAcceptor  `json:"card_acceptor"`
	CardReader            CardReader    `json:"card_reader"`
	Card                  Card          `json:"card"`
	Financial             Financial     `json:"financial"`
}

type Institution struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CardAcceptor struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	City            string `json:"city"`
	Country         string `json:"country"`
	ServiceCategory string `json:"service_category"`
}

type CardReader struct {
	TerminalID string `json:"terminal_id"`
	Type       string `json:"type"`
}

type Card struct {
	PanFirstFour  string `json:"pan_first_four"`
	PanLastFour   string `json:"pan_last_four"`
	Brand         string `json:"brand"`
	ExpiryDate    string `json:"expiry_date"`
	SaltTokenID   string `json:"salt_token_id"`
	IssuerCountry string `json:"issuer_country"`
	IssuerProduct string `json:"issuer_product"`
}

type Financial struct {
	TransactionAmount     int                   `json:"transaction_amount"`
	TransactionCurrency   string                `json:"transaction_currency"`
	CardholderAmount      int                   `json:"cardholder_amount"`
	CardholderCurrency    string                `json:"cardholder_currency"`
	MerchantServiceCharge MerchantServiceCharge `json:"merchant_service_charge"`
}

type MerchantServiceCharge struct {
	FixedAmount      int    `json:"fixed_amount"`
	PercentageAmount string `json:"percentage_amount"`
	CurrencyCode     string `json:"currency_code"`
}
