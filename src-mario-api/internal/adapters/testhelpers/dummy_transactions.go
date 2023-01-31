package testhelpers

import (
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/application/models"
)

// TODO change method to store multiple transactions in memory and
// be the "data store" for txns for the mock server
func loadTransaction() models.Transaction {
	return models.Transaction{
		Institutions: []models.Institution{
			{
				Name: "SaltPay UK",
				Type: "ACQUIRER",
			},
			{
				Name: "WAY4",
				Type: "ACQUIRING_HOST",
			},
		},
		InternalID:              "31814421-027b-4ac2-8d0d-80d7a9db8007",
		Uid:                     "WAY4116435",
		HostTransactionID:       "116435",
		HostParentTransactionID: "116434",
		Type:                    "AUTH",
		CardPresent:             "MISSING_INFO",
		IsDCC:                   "MISSING_INFO",
		EntryMethod:             "MISSING_INFO",
		ProcessedAt:             "2022-02-21T19:14:44.000Z",
		AuthCode:                "00572Z",
		Rrn:                     "213780213544",
		Arn:                     "BB1010198888882",
		CompanyID:               "639900bd-e2ac-4658-b6b4-99cc87c15204",
		StoreID:                 "4e58198e-6221-41d9-9f38-eb1ceafab46e",
		MerchantID:              "808054690",
		Response: models.Response{
			Code:           "00",
			Description:    "Approved",
			Classification: "APPROVED",
		},
		Card: models.Card{
			PanFirstSix:    "123546",
			PanLastFour:    "1017",
			Brand:          "MC",
			IssuingCountry: "UK",
			IssuingProduct: "MSR",
			ExpiryMonth:    "12",
			ExpiryYear:     "2025",
			Issuer:         "Santander",
			TokenID:        "38c19566-8683-43dc-9dce-e878f0e482be",
			Currency:       "EUR",
		},
		Financial: models.Financial{
			TransactionAmount:   100000000,
			TransactionCurrency: "EUR",
			BillingAmount:       100000000,
			BillingCurrency:     "EUR",
			Fees: models.Fees{
				Total:    1800000,
				Currency: "EUR",
				Breakdown: []models.Breakdown{
					{
						Type:               "MERCHANT_SERVICE_CHARGE",
						FixedAmount:        1800000,
						VariablePercentage: 0,
						VariableAmount:     0,
						Currency:           "EUR",
					},
				},
			},
		},
		Payout: models.Payout{
			ID:          "08355482-1fdf-4522-99c5-8395a9beccdd",
			ScheduledAt: "2022-02-22T19:14:44.000Z",
		},
	}
}

func captureTransaction(internalID string) models.Transaction {
	return models.Transaction{
		Institutions: []models.Institution{
			{
				Name: "SaltPay UK",
				Type: "ACQUIRER",
			},
			{
				Name: "WAY4",
				Type: "ACQUIRING_HOST",
			},
		},
		InternalID:              internalID,
		Uid:                     "WAY4116435",
		HostTransactionID:       "116435",
		HostParentTransactionID: "116434",
		Type:                    "CAPTURE",
		CardPresent:             "MISSING_INFO",
		IsDCC:                   "MISSING_INFO",
		EntryMethod:             "MISSING_INFO",
		ProcessedAt:             "2022-02-21T19:14:45.000Z",
		AuthCode:                "00572Z",
		Rrn:                     "213780213544",
		Arn:                     "BB1010198888882",
		CompanyID:               "639900bd-e2ac-4658-b6b4-99cc87c15204",
		StoreID:                 "4e58198e-6221-41d9-9f38-eb1ceafab46e",
		MerchantID:              "808054690",
		Response: models.Response{
			Code:           "00",
			Description:    "Approved",
			Classification: "APPROVED",
		},
		Card: models.Card{
			PanFirstSix:    "123546",
			PanLastFour:    "1017",
			Brand:          "MC",
			IssuingCountry: "UK",
			IssuingProduct: "MSR",
			ExpiryMonth:    "12",
			ExpiryYear:     "2025",
			Issuer:         "Santander",
			TokenID:        "38c19566-8683-43dc-9dce-e878f0e482be",
			Currency:       "EUR",
		},
		TransactionSource: models.TransactionSource{
			ID:   "2135125",
			Type: "MPOS",
		},
		Financial: models.Financial{
			TransactionAmount:   100000000,
			TransactionCurrency: "EUR",
			BillingAmount:       100000000,
			BillingCurrency:     "EUR",
			Fees: models.Fees{
				Total:    1800000,
				Currency: "EUR",
				Breakdown: []models.Breakdown{
					{
						Type:               "MERCHANT_SERVICE_CHARGE",
						FixedAmount:        1800000,
						VariablePercentage: 0,
						VariableAmount:     0,
						Currency:           "EUR",
					},
				},
			},
		},
		Payout: models.Payout{
			ID:          "08355482-1fdf-4522-99c5-8395a9beccdd",
			ScheduledAt: "2022-02-22T19:14:44.000Z",
		},
	}
}

func authTransaction(internalID string) models.Transaction {
	return models.Transaction{
		Institutions: []models.Institution{
			{
				Name: "SaltPay UK",
				Type: "ACQUIRER",
			},
			{
				Name: "WAY4",
				Type: "ACQUIRING_HOST",
			},
		},
		InternalID:              internalID,
		Uid:                     "WAY4116435",
		HostTransactionID:       "116435",
		HostParentTransactionID: "116434",
		Type:                    "AUTH",
		CardPresent:             "MISSING_INFO",
		IsDCC:                   "MISSING_INFO",
		EntryMethod:             "MISSING_INFO",
		ProcessedAt:             "2022-02-21T19:14:44.000Z",
		AuthCode:                "00572Z",
		Rrn:                     "213780213544",
		Arn:                     "BB1010198888882",
		CompanyID:               "639900bd-e2ac-4658-b6b4-99cc87c15204",
		StoreID:                 "4e58198e-6221-41d9-9f38-eb1ceafab46e",
		MerchantID:              "808054690",
		Response: models.Response{
			Code:           "00",
			Description:    "Approved",
			Classification: "APPROVED",
		},
		Card: models.Card{
			PanFirstSix:    "123546",
			PanLastFour:    "1017",
			Brand:          "MC",
			IssuingCountry: "UK",
			IssuingProduct: "MSR",
			ExpiryMonth:    "12",
			ExpiryYear:     "2025",
			Issuer:         "Santander",
			TokenID:        "38c19566-8683-43dc-9dce-e878f0e482be",
			Currency:       "EUR",
		},
		TransactionSource: models.TransactionSource{
			ID:   "2135125",
			Type: "MPOS",
		},
		Financial: models.Financial{
			TransactionAmount:   100000000,
			TransactionCurrency: "EUR",
			BillingAmount:       100000000,
			BillingCurrency:     "EUR",
			Fees: models.Fees{
				Total:    1800000,
				Currency: "EUR",
				Breakdown: []models.Breakdown{
					{
						Type:               "MERCHANT_SERVICE_CHARGE",
						FixedAmount:        1800000,
						VariablePercentage: 0,
						VariableAmount:     0,
						Currency:           "EUR",
					},
				},
			},
		},
		Payout: models.Payout{
			ID:          "08355482-1fdf-4522-99c5-8395a9beccdd",
			ScheduledAt: "2022-02-22T19:14:44.000Z",
		},
	}
}

func reversalTransaction(internalID string) models.Transaction {
	return models.Transaction{
		Institutions: []models.Institution{
			{
				Name: "SaltPay UK",
				Type: "ACQUIRER",
			},
			{
				Name: "WAY4",
				Type: "ACQUIRING_HOST",
			},
		},
		InternalID:              internalID,
		Uid:                     "WAY4116435",
		HostTransactionID:       "116435",
		HostParentTransactionID: "116434",
		Type:                    "REVERSAL",
		CardPresent:             "MISSING_INFO",
		IsDCC:                   "MISSING_INFO",
		EntryMethod:             "MISSING_INFO",
		ProcessedAt:             "2022-02-21T19:14:44.000Z",
		AuthCode:                "00572Z",
		Rrn:                     "213780213544",
		Arn:                     "BB1010198888882",
		CompanyID:               "639900bd-e2ac-4658-b6b4-99cc87c15204",
		StoreID:                 "4e58198e-6221-41d9-9f38-eb1ceafab46e",
		MerchantID:              "808054690",
		Response: models.Response{
			Code:           "00",
			Description:    "Approved",
			Classification: "APPROVED",
		},
		Card: models.Card{
			PanFirstSix:    "123546",
			PanLastFour:    "1017",
			Brand:          "MC",
			IssuingCountry: "UK",
			IssuingProduct: "MSR",
			ExpiryMonth:    "12",
			ExpiryYear:     "2025",
			Issuer:         "Santander",
			TokenID:        "38c19566-8683-43dc-9dce-e878f0e482be",
			Currency:       "EUR",
		},
		TransactionSource: models.TransactionSource{
			ID:   "2135125",
			Type: "MPOS",
		},
		Financial: models.Financial{
			TransactionAmount:   100000000,
			TransactionCurrency: "EUR",
			BillingAmount:       100000000,
			BillingCurrency:     "EUR",
			Fees: models.Fees{
				Total:    1800000,
				Currency: "EUR",
				Breakdown: []models.Breakdown{
					{
						Type:               "MERCHANT_SERVICE_CHARGE",
						FixedAmount:        1800000,
						VariablePercentage: 0,
						VariableAmount:     0,
						Currency:           "EUR",
					},
				},
			},
		},
		Payout: models.Payout{
			ID:          "08355482-1fdf-4522-99c5-8395a9beccdd",
			ScheduledAt: "2022-02-22T19:14:44.000Z",
		},
	}
}
