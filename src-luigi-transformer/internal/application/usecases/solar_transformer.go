package usecases

import (
	"context"
	"fmt"
	"math"

	"github.com/saltpay/go-kafka-driver"
	"go.uber.org/zap"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/logger/zapctx"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/models"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/ports"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/usecases/helpers"
)

type SolarTransformer struct {
	schemaRegistry ports.SchemaRegistry
	schemaKey      string
}

func NewSolarTransformer(schemaRegistry ports.SchemaRegistry, schemaKey string) SolarTransformer {
	return SolarTransformer{
		schemaRegistry: schemaRegistry,
		schemaKey:      schemaKey,
	}
}

func (s SolarTransformer) Execute(ctx context.Context, incomingTransactionMessage kafka.Message) (models.Transaction, error) {
	incomingTransaction, err := s.validateAndParse(ctx, incomingTransactionMessage.Value)
	if err != nil {
		return models.Transaction{}, err
	}

	canonicalTransaction, err := s.Translate(incomingTransaction)
	if err != nil {
		return models.Transaction{}, err
	}

	return canonicalTransaction, nil
}

func (s SolarTransformer) validateAndParse(ctx context.Context, incomingTransactionJSON []byte) (models.SolarIncomingTransaction, error) {
	// validate and unmarshal incoming message
	_, _, err := s.schemaRegistry.Decode(ctx, incomingTransactionJSON, s.schemaKey)
	if err != nil {
		zapctx.From(ctx).Error("[Solar-Validator] Error validating incoming message: ", zap.Error(err), zap.ByteString("incomingTransactionJSON", incomingTransactionJSON))
		zapctx.From(ctx).Error("[Solar-Validator] Ingress message does not comply to the defined schema")
		return models.SolarIncomingTransaction{}, err
	}

	incomingTransaction, err := models.SolarIncomingTransactionFromBytes(incomingTransactionJSON)
	if err != nil {
		zapctx.From(ctx).Error("[SolarTransformer] Error converting kafka message: ", zap.Error(err))
		return models.SolarIncomingTransaction{}, err
	}

	return incomingTransaction, nil
}

func (s SolarTransformer) Translate(incomingTransaction models.SolarIncomingTransaction) (models.Transaction, error) {
	isDCC := isDCCCheck(incomingTransaction.Body.Txn.Amounts.TransactionAmount.Amount, incomingTransaction.Body.Txn.Amounts.OriginatorBillingAmount.Amount,
		incomingTransaction.Body.Txn.Amounts.TransactionAmount.Currency, incomingTransaction.Body.Txn.Amounts.OriginatorBillingAmount.Currency)

	cardLastFourDigits := breakCardNumber(incomingTransaction.Body.Txn.Receiver.Number)

	if len(incomingTransaction.Body.Txn.Originator.Fees.Fee) < 1 {
		incomingTransaction.Body.Txn.Originator.Fees.Fee = append(incomingTransaction.Body.Txn.Originator.Fees.Fee, models.Fee{})
	}

	canonicalTransaction := models.Transaction{
		Institutions: []models.Institution{{
			Name: "SOLAR",
			Type: "ACQUIRER",
		}, {
			Name: incomingTransaction.Body.Txn.Originator.System.Name,
			Type: "PSP",
		}},
		InternalID:              fmt.Sprint(incomingTransaction.Body.Txn.TxnID),
		HostTransactionID:       fmt.Sprint(incomingTransaction.Body.Txn.TxnID),
		HostParentTransactionID: fmt.Sprint(incomingTransaction.Body.Txn.TxnID),
		UID:                     fmt.Sprint(incomingTransaction.Body.Txn.TxnID),
		Type:                    findTransactionType(incomingTransaction.Body.Txn.TransactionType.DirectionClass, incomingTransaction.Body.Txn.Direction, incomingTransaction.Body.Txn.Class),
		CardPresent:             isCardPresent(incomingTransaction.Body.Txn.TxnConditions.CardPresenceType),
		IsDcc:                   isDCC,
		EntryMethod:             incomingTransaction.Body.Txn.TxnConditions.CardDataInputMode,
		ProcessedAt:             incomingTransaction.Body.Txn.TransactionDate,
		AuthCode:                incomingTransaction.Body.Txn.Reference.AuthCode,
		RRN:                     "MISSING_INFO",
		ARN:                     "MISSING_INFO",
		CompanyID:               "MISSING_INFO",
		StoreID:                 incomingTransaction.Body.Txn.MerchantData.CardAcceptorIDCode,
		MerchantID:              incomingTransaction.Body.Txn.MerchantData.CardAcceptorIDCode,
		TransactionSource: models.TransactionSource{
			ID:   "MISSING_INFO",
			Type: "MISSING_INFO",
		},
		Response: models.Response{
			Code:           "MISSING_INFO",
			Classification: "MISSING_INFO",
			Description:    "MISSING_INFO",
		},
		Card: models.Card{
			PanFirstSix:    "MISSING_INFO",
			PanLastFour:    cardLastFourDigits,
			Brand:          incomingTransaction.Body.Txn.Receiver.System.Code,
			ExpiryMonth:    fmt.Sprint(findAttribute(incomingTransaction.Body.Txn.Attributes.Attribute, "card/cardExpiryDate")),
			ExpiryYear:     "MISSING_INFO",
			Issuer:         "MISSING_INFO",
			IssuingCountry: incomingTransaction.Body.Txn.BinTableData.Country,
			IssuingProduct: incomingTransaction.Body.Txn.BinTableData.Product,
			TokenID:        incomingTransaction.Body.Txn.SaltTokenID,
			Currency:       "MISSING_INFO",
		},
		Financial: models.Financial{
			TransactionAmount:   int64(100 * (math.Round(incomingTransaction.Body.Txn.Amounts.OriginatorBillingAmount.Amount*100) / 100)),
			TransactionCurrency: helpers.GetCurrencyIsoCode(incomingTransaction.Body.Txn.Amounts.OriginatorBillingAmount.Currency),
			BillingAmount:       int64(100 * (math.Round(incomingTransaction.Body.Txn.Amounts.TransactionAmount.Amount*100) / 100)),
			BillingCurrency:     helpers.GetCurrencyIsoCode(incomingTransaction.Body.Txn.Amounts.TransactionAmount.Currency),
			Fees: models.Fees{
				Total:    0,
				Currency: "MISSING_INFO",
				Type:     "MISSING_INFO",
				Breakdown: []models.Breakdown{
					{
						Type:               "MISSING_INFO",
						FixedAmount:        0,
						VariablePercentage: 0,
						VariableAmount:     0,
						Currency:           "MISSING_INFO",
					},
				},
			},
		},
	}

	return canonicalTransaction, nil
}

func findAttribute(attributes []models.Attribute, key string) interface{} {
	for _, attr := range attributes {
		if attr.Code == key {
			return attr.Attribute
		}
	}

	return ""
}

func findTransactionType(typeDirection, direction, class string) string {
	if direction == "REVERSE" {
		return "REVERSAL"
	} else if typeDirection == "CREDIT" {
		return "REFUND"
	} else if class == "AUTHORIZATION" {
		return "AUTH"
	}
	return "CAPTURE"
}

func isCardPresent(cardPresenceType string) bool {
	return cardPresenceType == "IS_PRESENT"
}

func breakCardNumber(cardNumber string) string {
	if len(cardNumber) > 8 {
		return cardNumber[len(cardNumber)-4:]
	}
	return ""
}

func isDCCCheck(transactionAmount, originatorBillingAmount float64, transactionCurr, originalBillingCurr string) bool {
	return transactionAmount != originatorBillingAmount || transactionCurr != originalBillingCurr
}
