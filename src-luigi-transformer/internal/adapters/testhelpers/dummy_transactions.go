package testhelpers

import (
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/models"
)

func LoadCanonicalTransaction() models.Transaction {
	return models.Transaction{
		Institutions: []models.Institution{{
			Name: "Way4",
			Type: "ACQUIRER",
		}},
		InternalID:              "",
		HostTransactionID:       "12413650420",
		HostParentTransactionID: "12413650410",
		UID:                     "12413650420" + "12413650410",
		Type:                    "AUTH_REVERSAL",
		CardPresent:             false, // TODO: check if this is correct
		IsDcc:                   false,
		EntryMethod:             "MISSING_INFO",
		ProcessedAt:             "2022-04-14 18:03:41.000",
		AuthCode:                "005988",
		RRN:                     "210407333288",
		ARN:                     "0",
		CompanyID:               "MISSING_INFO", // TODO: check if this is correct
		StoreID:                 "9117052",
		MerchantID:              "9117052",
		TransactionSource: models.TransactionSource{
			ID:   "MISSING_INFO",
			Type: "MISSING_INFO",
		},
		Response: models.Response{
			Code:           "0",
			Classification: "APPROVED",
			Description:    "MISSING_INFO",
		},
		Card: models.Card{
			PanFirstSix:    "555555",
			PanLastFour:    "4444",
			Brand:          "MC",
			ExpiryMonth:    "NULL",
			ExpiryYear:     "NULL",
			Issuer:         "",
			IssuingProduct: "",
			IssuingCountry: "",
			TokenID:        "",
			Currency:       "",
		},
		Financial: models.Financial{
			TransactionAmount:   1000,
			TransactionCurrency: "ISK",
			BillingAmount:       1000,
			BillingCurrency:     "ISK",
			Fees: models.Fees{
				Total:    0,
				Currency: "MISSING_INFO",
				Type:     "MISSING_INFO",
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
			ID:          "MISSING_INFO",
			ScheduledAt: "MISSING_INFO",
		},
	}
}

func LoadSolarJSON() []byte {
	return []byte(`{
  "header": {
    "protocol": {
      "name": "solar-events",
      "version": "1.0"
    },
    "messageId": "2c26064a-2501-4960-9f6d-f4b6264314c3",
    "messageDate": "2022-03-22T10:32:04",
    "originator": {
      "system": "SOLAR"
    },
    "receiver": null,
    "responseParams": {
      "collectByReference": true,
      "includes": {
        "include": [
          {
            "data": "balances",
            "include": true
          },
          {
            "data": "binTableData",
            "include": true
          }
        ]
      }
    },
    "locale": null,
    "workflow": null,
    "preciseTime": null
  },
  "body": {
    "txn": {
      "txnId": 116435,
      "transactionType": {
        "id": 1,
        "code": "purchase",
        "name": "Purchase",
        "directionClass": "DEBIT"
      },
      "class": "AUTHORIZATION",
      "category": "REQUEST",
      "direction": "ORIGINAL",
      "transactionDate": "2022-03-22T10:32:00",
      "settlementDate": "2022-03-23",
      "saltTokenId": "38c19566-8683-43dc-9dce-e878f0e482be",
      "systemDate": "2022-03-22T10:32:00.804201",
      "reference": {
        "authCode": "00555Z",
        "retrievalRefNumber": "208110036464"
      },
      "originator": {
        "system": {
          "id": 100007,
          "code": "solar-slq-h2h-switch-payments",
          "name": "SwitchPayments (PSP)",
          "category": "ON_US"
        },
        "agreement": {
          "id": 123174,
          "number": "5311-620"
        },
        "accessor": {
          "id": 113156,
          "number": "55312620",
          "type": {
            "id": 100001,
            "code": "TERM",
            "name": "Terminal"
          }
        },
        "number": "55312620",
        "feeTotalAmount": {
          "amount": -0.06,
          "currency": "978"
        },
        "fees": {
          "fee": [
            {
              "id": 158332,
              "feeType": {
                "id": 100023,
                "code": "purchaseTxnFee",
                "name": "Purchase Fee",
                "feeClass": "TXN_FEE",
                "directionClass": "DEBIT"
              },
              "tariff": {
                "id": 100001,
                "tariffType": "FEE",
                "tariffClass": "PRODUCT",
                "tariffGroupId": 100000,
                "feeTariffValue": {
                  "id": 100002,
                  "amountType": "ORIGINATOR",
                  "base": {
                    "value": 0,
                    "currency": "978"
                  },
                  "min": {
                    "value": 0,
                    "currency": "978"
                  },
                  "max": {
                    "value": 0,
                    "currency": "978"
                  },
                  "percentValue": 1.5,
                  "valueCurrency": "978"
                }
              },
              "direction": "ORIGINAL",
              "amounts": {
                "feeAmount": {
                  "amount": 0.06,
                  "currency": "978"
                },
                "billingAmount": {
                  "amount": 0.06,
                  "currency": "978"
                }
              }
            }
          ]
        },
        "balances": {
          "balance": [
            {
              "balanceType": {
                "id": 100000,
                "code": "ON_HOLD",
                "name": "On-hold balance"
              },
              "currency": "978",
              "own": 0,
              "available": 0,
              "blocked": 0,
              "loan": 0,
              "overlimit": 0,
              "overdue": 0,
              "creditLimit": 0,
              "finBlocking": 0,
              "interests": 0,
              "penalty": 0
            },
            {
              "balanceType": {
                "id": 100001,
                "code": "CREDIT_BALANCE",
                "name": "Credit balance"
              },
              "currency": "978",
              "own": 0,
              "available": 0,
              "blocked": 0,
              "loan": 0,
              "overlimit": 0,
              "overdue": 0,
              "creditLimit": 0,
              "finBlocking": 0,
              "interests": 0,
              "penalty": 0
            },
            {
              "balanceType": {
                "id": 100002,
                "code": "ACQ_BALANCE",
                "name": "Acquiring balance"
              },
              "currency": "978",
              "own": 12546.66,
              "available": 12546.66,
              "blocked": 0,
              "loan": 0,
              "overlimit": 0,
              "overdue": 0,
              "creditLimit": 0,
              "finBlocking": 0,
              "interests": 0,
              "penalty": 0
            }
          ]
        }
      },
      "receiver": {
        "system": {
          "id": 100000,
          "code": "MC",
          "name": "MasterCard",
          "category": "EXTERNAL",
          "systemGroup": {
            "id": 100001,
            "code": "MC",
            "name": "MasterCard"
          }
        },
        "agreement": {
          "id": 123267,
          "number": "ISSUER_MC"
        },
        "accessor": {},
        "number": "545721****1040",
        "feeTotalAmount": {
          "amount": 0,
          "currency": "978"
        }
      },
      "merchantData": {
        "merchantName": "35 West Street",
        "merchantCity": "New York",
        "merchantCountry": "PT",
        "mcc": "5311",
        "cardAcceptorIdCode": "5311-620"
      },
      "binTableData": {
        "product": "MDS",
        "country": "BE"
      },
      "amounts": {
        "transactionAmount": {
          "amount": 4.2,
          "currency": "978"
        },
        "settlementAmount": {
          "amount": 4.2,
          "currency": "978"
        },
        "originatorAmount": {
          "amount": 4.2,
          "currency": "978"
        },
        "originatorBillingAmount": {
          "amount": 4.2,
          "currency": "978"
        },
        "receiverAmount": {
          "amount": 4.2,
          "currency": "978"
        },
        "receiverBillingAmount": {
          "amount": 4.2,
          "currency": "978"
        }
      },
      "attributes": {
        "attribute": [
          {
            "code": "card/cardSequenceNumber",
            "attribute": "0"
          },
          {
            "code": "card/cardExpiryDate",
            "attribute": "12/25"
          },
          {
            "code": "origAcsrType",
            "attribute": "TERM"
          },
          {
            "code": "slrRefData/sw5.chain.id",
            "attribute": 50363709
          },
          {
            "code": "posData/cardDataInputMode",
            "attribute": "ICC"
          },
          {
            "code": "rcvrSystemGroup",
            "attribute": "MC"
          },
          {
            "code": "rcvrCountryCode",
            "attribute": "056"
          },
          {
            "code": "productCode",
            "attribute": "MDS"
          },
          {
            "code": "ORIGINATOR",
            "attribute": 4.2
          },
          {
            "code": "originatorLimiterIds",
            "attribute": []
          },
          {
            "code": "receiverLimiterIds",
            "attribute": []
          }
        ]
      },
      "txnConditions": {
        "token": false,
        "isECommerceIndicator": false,
        "cardDataInputMode": "ICC",
        "chipDataReadFailed": false,
        "cardholderPresenceType": "IS_PRESENT",
        "cardPresenceType": "IS_PRESENT",
        "catLevel": "TERM_MOBILE"
      },
      "txnDetails": "null",
      "response": {
        "code": "00",
        "message": "00 – Approved"
      },
      "txnStatus": "FINISHED"
    }
  }
}
`)
}
