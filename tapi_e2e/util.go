package tapi_e2e

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
)

const (
	LocalURL = "http://localhost:8080"
)

type SingleTransactionRequest struct {
	TransactionID string `json:"transaction_id"`
	AcquiringHost string `json:"acquiring_host"`
}

func getBaseURL(t *testing.T) string {
	url := os.Getenv("BASE_URL")
	if url == "" {
		url = LocalURL
		startWebserver(t)
	}

	return url
}

func startWebserver(t *testing.T) {
	t.Helper()

	compose := testcontainers.NewLocalDockerCompose(
		[]string{"../../docker-compose.yaml"},
		strings.ToLower(uuid.New().String()),
	)
	webContainer := compose.WithCommand([]string{"up", "-d", "web"})
	invokeErr := webContainer.Invoke()

	if invokeErr.Error != nil {
		t.Fatal(invokeErr)
	}

	t.Cleanup(func() {
		compose.Down()
	})
}

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

func NewTransactionFromJSON(in []byte) (Transaction, error) {
	var transaction Transaction
	err := json.Unmarshal(in, &transaction)
	return transaction, err
}

type SolarIncomingTransaction struct {
	Header struct {
		Protocol struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"protocol"`
		MessageID   string `json:"messageId"`
		MessageDate string `json:"messageDate"`
		Originator  struct {
			System string `json:"system"`
		} `json:"originator"`
		Receiver       interface{} `json:"receiver"`
		ResponseParams struct {
			CollectByReference bool `json:"collectByReference"`
			Includes           struct {
				Include []struct {
					Data    string `json:"data"`
					Include bool   `json:"include"`
				} `json:"include"`
			} `json:"includes"`
		} `json:"responseParams"`
		Locale      interface{} `json:"locale"`
		Workflow    interface{} `json:"workflow"`
		PreciseTime interface{} `json:"preciseTime"`
	} `json:"header"`
	Body struct {
		Txn struct {
			TxnID           int `json:"txnId"`
			TransactionType struct {
				ID             int    `json:"id"`
				Code           string `json:"code"`
				Name           string `json:"name"`
				DirectionClass string `json:"directionClass"`
			} `json:"transactionType"`
			SaltTokenID     string    `json:"salt_token_id"`
			Class           string    `json:"class"`
			Category        string    `json:"category"`
			Direction       string    `json:"direction"`
			TransactionDate time.Time `json:"transactionDate"`
			SettlementDate  string    `json:"settlementDate"`
			SystemDate      string    `json:"systemDate"`
			Reference       struct {
				AuthCode           string `json:"authCode"`
				RetrievalRefNumber string `json:"retrievalRefNumber"`
			} `json:"reference"`
			Originator struct {
				System struct {
					ID       int    `json:"id"`
					Code     string `json:"code"`
					Name     string `json:"name"`
					Category string `json:"category"`
				} `json:"system"`
				Agreement struct {
					ID     int    `json:"id"`
					Number string `json:"number"`
				} `json:"agreement"`
				Accessor struct {
					ID     int    `json:"id"`
					Number string `json:"number"`
					Type   struct {
						ID   int    `json:"id"`
						Code string `json:"code"`
						Name string `json:"name"`
					} `json:"type"`
				} `json:"accessor"`
				Number         string `json:"number"`
				FeeTotalAmount struct {
					Amount   float64 `json:"amount"`
					Currency string  `json:"currency"`
				} `json:"feeTotalAmount"`
				Fees struct {
					Fee []struct {
						ID      int `json:"id"`
						FeeType struct {
							ID             int    `json:"id"`
							Code           string `json:"code"`
							Name           string `json:"name"`
							FeeClass       string `json:"feeClass"`
							DirectionClass string `json:"directionClass"`
						} `json:"feeType"`
						Tariff struct {
							ID             int    `json:"id"`
							TariffType     string `json:"tariffType"`
							TariffClass    string `json:"tariffClass"`
							TariffGroupID  int    `json:"tariffGroupId"`
							FeeTariffValue struct {
								ID         int    `json:"id"`
								AmountType string `json:"amountType"`
								Base       struct {
									Value    int    `json:"value"`
									Currency string `json:"currency"`
								} `json:"base"`
								Min struct {
									Value    int    `json:"value"`
									Currency string `json:"currency"`
								} `json:"min"`
								Max struct {
									Value    int    `json:"value"`
									Currency string `json:"currency"`
								} `json:"max"`
								PercentValue  float64 `json:"percentValue"`
								ValueCurrency string  `json:"valueCurrency"`
							} `json:"feeTariffValue"`
						} `json:"tariff"`
						Direction string `json:"direction"`
						Amounts   struct {
							FeeAmount struct {
								Amount   float64 `json:"amount"`
								Currency string  `json:"currency"`
							} `json:"feeAmount"`
							BillingAmount struct {
								Amount   float64 `json:"amount"`
								Currency string  `json:"currency"`
							} `json:"billingAmount"`
						} `json:"amounts"`
					} `json:"fee"`
				} `json:"fees"`
				Balances struct {
					Balance []struct {
						BalanceType struct {
							ID   int    `json:"id"`
							Code string `json:"code"`
							Name string `json:"name"`
						} `json:"balanceType"`
						Currency    string  `json:"currency"`
						Own         float64 `json:"own"`
						Available   float64 `json:"available"`
						Blocked     int     `json:"blocked"`
						Loan        int     `json:"loan"`
						Overlimit   int     `json:"overlimit"`
						Overdue     int     `json:"overdue"`
						CreditLimit int     `json:"creditLimit"`
						FinBlocking int     `json:"finBlocking"`
						Interests   int     `json:"interests"`
						Penalty     int     `json:"penalty"`
					} `json:"balance"`
				} `json:"balances"`
			} `json:"originator"`
			Receiver struct {
				System struct {
					ID          int    `json:"id"`
					Code        string `json:"code"`
					Name        string `json:"name"`
					Category    string `json:"category"`
					SystemGroup struct {
						ID   int    `json:"id"`
						Code string `json:"code"`
						Name string `json:"name"`
					} `json:"systemGroup"`
				} `json:"system"`
				Agreement struct {
					ID     int    `json:"id"`
					Number string `json:"number"`
				} `json:"agreement"`
				Accessor       struct{} `json:"accessor"`
				Number         string   `json:"number"`
				FeeTotalAmount struct {
					Amount   int    `json:"amount"`
					Currency string `json:"currency"`
				} `json:"feeTotalAmount"`
			} `json:"receiver"`
			MerchantData struct {
				MerchantName       string `json:"merchantName"`
				MerchantCity       string `json:"merchantCity"`
				MerchantCountry    string `json:"merchantCountry"`
				Mcc                string `json:"mcc"`
				CardAcceptorIDCode string `json:"cardAcceptorIdCode"`
			} `json:"merchantData"`
			BinTableData struct {
				Product string `json:"product"`
				Country string `json:"country"`
			} `json:"binTableData"`
			Amounts struct {
				TransactionAmount struct {
					Amount   float64 `json:"amount"`
					Currency string  `json:"currency"`
				} `json:"transactionAmount"`
				SettlementAmount struct {
					Amount   float64 `json:"amount"`
					Currency string  `json:"currency"`
				} `json:"settlementAmount"`
				OriginatorAmount struct {
					Amount   float64 `json:"amount"`
					Currency string  `json:"currency"`
				} `json:"originatorAmount"`
				OriginatorBillingAmount struct {
					Amount   float64 `json:"amount"`
					Currency string  `json:"currency"`
				} `json:"originatorBillingAmount"`
				ReceiverAmount struct {
					Amount   float64 `json:"amount"`
					Currency string  `json:"currency"`
				} `json:"receiverAmount"`
				ReceiverBillingAmount struct {
					Amount   float64 `json:"amount"`
					Currency string  `json:"currency"`
				} `json:"receiverBillingAmount"`
			} `json:"amounts"`
			Attributes struct {
				Attribute []Attribute `json:"attribute"`
			} `json:"attributes"`
			TxnConditions struct {
				Token                  bool   `json:"token"`
				IsECommerceIndicator   bool   `json:"isECommerceIndicator"`
				CardDataInputMode      string `json:"cardDataInputMode"`
				ChipDataReadFailed     bool   `json:"chipDataReadFailed"`
				CardholderPresenceType string `json:"cardholderPresenceType"`
				CardPresenceType       string `json:"cardPresenceType"`
				CatLevel               string `json:"catLevel"`
			} `json:"txnConditions"`
			TxnDetails string `json:"txnDetails"`
			Response   struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"response"`
			TxnStatus string `json:"txnStatus"`
		} `json:"txn"`
	} `json:"body"`
}

type Attribute struct {
	Code      string      `json:"code"`
	Attribute interface{} `json:"attribute"`
}
