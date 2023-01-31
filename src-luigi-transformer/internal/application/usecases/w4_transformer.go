package usecases

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/saltpay/go-kafka-driver"
	"go.uber.org/zap"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/logger/zapctx"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/models"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/ports"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/usecases/helpers"
)

const nullValues = "NULL"

type W4Transformer struct {
	dbHandler      ports.DBHandler
	schemaRegistry ports.SchemaRegistry
	schemaKey      string
	metricsClient  ports.MetricsClient
}

func NewW4Transformer(dbHandler ports.DBHandler, schemaRegistry ports.SchemaRegistry, schemaKey string, metricsClient ports.MetricsClient) W4Transformer {
	return W4Transformer{
		dbHandler:      dbHandler,
		schemaRegistry: schemaRegistry,
		schemaKey:      schemaKey,
		metricsClient:  metricsClient,
	}
}

type ErrorNotTransaction struct{}

func (ErrorNotTransaction) Error() string {
	return "Message should not be stored in TAPI"
}

var (
	reDCCField, _           = regexp.Compile(`DCC_IND=(\w+);`)
	reIssuerCountryField, _ = regexp.Compile(`TGTCC=(\w+);`)
	reIssuerProductField, _ = regexp.Compile(`PRODUCT_ID=(\w+);`)
)

func (w W4Transformer) Execute(ctx context.Context, incomingTransactionMessage kafka.Message) (models.Transaction, error) {
	incomingTransaction, err := w.validateAndParse(ctx, incomingTransactionMessage.Value)
	if err != nil {
		return models.Transaction{}, err
	}

	canonicalTransaction, err := w.Translate(ctx, incomingTransaction)
	if err != nil {
		return models.Transaction{}, err
	}

	return canonicalTransaction, nil
}

func (w W4Transformer) validateAndParse(ctx context.Context, incomingTransactionJSON []byte) (models.W4IncomingTransaction, error) {
	// validate and unmarshal incoming message
	_, _, err := w.schemaRegistry.Decode(ctx, incomingTransactionJSON, w.schemaKey)
	if err != nil {
		w.metricsClient.Count("message_validator", 1, []string{"W4-Validator", "failed"})
		zapctx.From(ctx).Error("[W4-Validator] Error validating incoming message: ", zap.Error(err), zap.ByteString("incomingTransactionJSON", incomingTransactionJSON))
		zapctx.From(ctx).Error("[W4-Validator] Ingress message does not comply to the defined schema")
		return models.W4IncomingTransaction{}, err
	}

	incomingTransaction, err := models.W4IncomingTransactionFromBytes(incomingTransactionJSON)
	if err != nil {
		w.metricsClient.Count("message_validator", 1, []string{"W4-Validator", "failed"})
		zapctx.From(ctx).Error("[W4-Validator] Error converting kafka message: ", zap.Error(err))
		return models.W4IncomingTransaction{}, err
	}
	w.metricsClient.Count("message_validator", 1, []string{"W4-Validator", "success"})
	return incomingTransaction, nil
}

func (w W4Transformer) Translate(ctx context.Context, incomingTransaction models.W4IncomingTransaction) (models.Transaction, error) {
	// check if object is a card transaction
	isTransaction, hasTokenData := w.transactionTypeValidator(incomingTransaction)

	// If the message is a transaction, we follow the normal mapping. If not, we update the aux table with the new id (+ token/expiry date)
	if !isTransaction {
		err := w.dbHandler.StoreAuxTable(ctx, incomingTransaction, hasTokenData)
		if err != nil {
			zapctx.From(ctx).Error("[W4Transformer] Error converting kafka message: ", zap.Error(err))
			return models.Transaction{}, err
		}
		return models.Transaction{}, ErrorNotTransaction{}
	}

	transaction := w.mapper(ctx, incomingTransaction)
	return transaction, nil
}

func (w W4Transformer) transactionTypeValidator(transaction models.W4IncomingTransaction) (bool, bool) {
	hasTokenData := len(transaction.After.SaltTokenID) > 0 && len(transaction.After.CardExpire) > 0

	if transaction.After.IsAuthorization == "P" && transaction.After.AmndState == "A" && helpers.IsInTargetChannel(transaction.After.TargetChannel) {
		return true, hasTokenData
	} else {
		return false, hasTokenData
	}
}

func (w W4Transformer) mapper(ctx context.Context, incomingTransaction models.W4IncomingTransaction) models.Transaction {
	isDCC := checkIfDCC(incomingTransaction.After.AddInfo)
	cardBrand := helpers.GetCardBrand(getCardNumberBrand(incomingTransaction.After.SCat, incomingTransaction.After.TCat, incomingTransaction.After.SourceNumber, incomingTransaction.After.TargetNumber, len(incomingTransaction.After.TargetNumber), false))
	panFirstSix, terminalID := getCardNumber(incomingTransaction.After.SCat, incomingTransaction.After.TCat, incomingTransaction.After.SourceNumber, incomingTransaction.After.TargetNumber, len(incomingTransaction.After.TargetNumber), true)
	panLastFour, terminalID := getCardNumber(incomingTransaction.After.SCat, incomingTransaction.After.TCat, incomingTransaction.After.SourceNumber, incomingTransaction.After.TargetNumber, 4, false)
	transactionType := w.getTransactionType(incomingTransaction.After.IsAuthorization, incomingTransaction.After.RequestCategory, incomingTransaction.After.TransType, incomingTransaction.After.AddInfo, incomingTransaction.After.SourceRegNum)
	cardIssuerCountry, cardIssuerProduct, issuer := w.getCardIssuerData(incomingTransaction.After.AddInfo)

	saltTokenID, expiryDate := w.getTokenData(ctx, incomingTransaction.After.DocPrevID, incomingTransaction.After.SaltTokenID, incomingTransaction.After.CardExpire)
	responseData := w.getResponseAuth(incomingTransaction.After.ReturnCode)
	feesData := w.getFeesData()

	transactionAmount := int64(100 * (math.Round(incomingTransaction.After.TransAmount*100) / 100))
	transactionCurrency := helpers.GetCurrencyIsoCode(incomingTransaction.After.TransCurr)

	var cardHolderAmount int64
	var cardHolderCurrency string
	if isDCC {
		cardHolderCurrency = helpers.GetCurrencyIsoCode(incomingTransaction.After.SettlCurr)
		cardHolderAmount = int64(100 * (math.Round(incomingTransaction.After.SettlAmount*100) / 100))
	} else {
		cardHolderCurrency = transactionCurrency
		cardHolderAmount = transactionAmount
	}
	expiryMonth, expiryYear := "", ""
	if len(expiryDate) >= 4 {
		expiryMonth = expiryDate[0:2]
		expiryYear = expiryDate[2:4]
	} else {
		expiryMonth = nullValues
		expiryYear = nullValues
	}
	canonicalTransaction := models.Transaction{
		Institutions: []models.Institution{{
			Name: "Way4",
			Type: "ACQUIRER",
		}},
		InternalID:              saltTokenID,
		HostTransactionID:       fmt.Sprint(incomingTransaction.After.ID),
		HostParentTransactionID: fmt.Sprint(incomingTransaction.After.DocPrevID),
		UID:                     fmt.Sprint(incomingTransaction.After.ID) + fmt.Sprint(incomingTransaction.After.DocPrevID),
		Type:                    transactionType,
		CardPresent:             false, // TODO: check if this is correct
		IsDcc:                   isDCC,
		EntryMethod:             "MISSING_INFO",
		ProcessedAt:             incomingTransaction.After.TransDate,
		AuthCode:                incomingTransaction.After.AuthCode,
		RRN:                     incomingTransaction.After.RetRefNumber,
		ARN:                     incomingTransaction.After.AcqRefNumber,
		CompanyID:               "MISSING_INFO", // TODO: check if this is correct
		StoreID:                 incomingTransaction.After.MerchantID,
		MerchantID:              incomingTransaction.After.MerchantID,
		Response:                responseData,
		TransactionSource: models.TransactionSource{
			ID:   terminalID,
			Type: "MISSING_INFO",
		},
		Card: models.Card{
			PanFirstSix:    panFirstSix,
			PanLastFour:    panLastFour,
			Brand:          cardBrand,
			ExpiryMonth:    expiryMonth,
			ExpiryYear:     expiryYear,
			Issuer:         issuer,
			IssuingProduct: cardIssuerProduct,
			IssuingCountry: cardIssuerCountry,
			TokenID:        saltTokenID,
			Currency:       cardHolderCurrency,
		},
		Financial: models.Financial{
			TransactionAmount:   transactionAmount,
			TransactionCurrency: transactionCurrency,
			BillingAmount:       cardHolderAmount,
			BillingCurrency:     cardHolderCurrency,
			Fees:                feesData,
		},
		Payout: models.Payout{
			ID:          "MISSING_INFO",
			ScheduledAt: "MISSING_INFO",
		},
	}

	return canonicalTransaction
}

func (w W4Transformer) getFeesData() models.Fees {
	return models.Fees{
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
	}
}

func (w W4Transformer) getTokenData(ctx context.Context, prevID int64, saltTokenID, expiryDate string) (string, string) {
	// if transaction has saltTokenID and expiryDate; do nothing
	if len(saltTokenID) > 0 && len(expiryDate) > 0 {
		return saltTokenID, expiryDate
	}
	// get tokenization data from aux table
	var auxToken, auxExpiryDate string
	query := fmt.Sprintf(`SELECT salt_token_id, expiry_date FROM aux_transactions_api where id = %d`, prevID)
	err := w.dbHandler.Get(ctx, query, []interface{}{&auxToken, &auxExpiryDate})
	if err != nil {
		return nullValues, nullValues
	}

	return auxToken, auxExpiryDate
}

func (w W4Transformer) getResponseAuth(returnCode int64) models.Response {
	return models.Response{
		Code:           fmt.Sprint(returnCode),
		Classification: getTransactionAuthResponse(returnCode),
		Description:    "MISSING_INFO",
	}
}

func (w W4Transformer) getCardIssuerData(info string) (string, string, string) {
	var country, cardBrand, issuer string
	// check if reIssuerCountryField regex was matched
	if reIssuerCountryField.MatchString(info) {
		country = reIssuerCountryField.FindStringSubmatch(info)[1]
	}
	// check if reIssuerProductField regex was matched
	if reIssuerProductField.MatchString(info) {
		cardBrand = reIssuerProductField.FindStringSubmatch(info)[1]
	}
	issuer = "MISSING_INFO"
	return country, cardBrand, issuer
}

func checkIfDCC(info string) bool {
	// check if reDCCField regex was matched
	if reDCCField.MatchString(info) {
		if reDCCField.FindStringSubmatch(info)[1] == "Y" {
			return true
		}
	}
	return false
}

func getTransactionAuthResponse(returnCode int64) string {
	if returnCode == 0o0 || returnCode == 85 || returnCode == 0 {
		return "APPROVED"
	}
	if returnCode == 10 {
		return "PARTIAL_APPROVAL"
	} else {
		return "DECLINED"
	}
}

func (w W4Transformer) getTransactionType(isAuth, requestCategory string, transType int64, addInfo string, sourceRegNum string) string {
	if transType == 5 && isAuth == "P" && requestCategory == "Q" {
		return "AUTH"
	} else if transType == 5 && isAuth == "P" && requestCategory == "Q" && strings.Contains(addInfo, "PREAUTH=Y;") {
		return "PREAUTH"
	} else if transType == 5 && isAuth == "P" && requestCategory == "P" {
		return "PREAUTH_COMPLETION_ADVICE"
	} else if transType == 5 && isAuth == "P" && requestCategory == "R" {
		return "AUTH_REVERSAL"
	} else if transType == 5 && isAuth == "P" && requestCategory == "R" && strings.Contains(addInfo, "PREAUTH=Y;") {
		return "PREAUTH_REVERSAL"
	} else if transType == 5 && isAuth == "P" && requestCategory == "J" {
		return "AUTH_ADJUSTMENT"
	} else if transType == 5 && isAuth == "P" && requestCategory == "J" && strings.Contains(addInfo, "PREAUTH=Y;") {
		return "PREAUTH_ADJUSTMENT"
	} else if transType == 5 && isAuth == "N" && requestCategory == "P" && sourceRegNum == "" {
		return "CAPTURE"
	} else if transType == 5 && isAuth == "N" && requestCategory == "P" && sourceRegNum != "" {
		return "AUTH_AND_CAPTURE"
	} else if transType == 15 && isAuth == "N" && requestCategory == "P" {
		return "REFUND"
	} else if transType == 15 && isAuth == "N" && requestCategory == "R" {
		return "REFUND_REVERSAL"
	} else if transType == 5 && isAuth == "N" && requestCategory == "R" {
		return "CAPTURE_REVERSAL"
	} else if transType == 5 && isAuth == "N" && requestCategory == "J" {
		return "CAPTURE_ADJUSTMENT"
	} else if transType == 17 && isAuth == "N" && requestCategory == "P" {
		return "CHARGEBACK"
	} else if transType == 7 && isAuth == "N" && requestCategory == "P" {
		return "2ND_PRESENTMENT"
	} else if transType == 7 && isAuth == "N" && requestCategory == "R" {
		return "2ND_PRESENTMENT_REVERSAL"
	} else if transType == 20 && isAuth == "N" && requestCategory == "P" {
		return "2ND_CHARGEBACK"
	} else if transType == 43 && isAuth == "N" && requestCategory == "P" {
		return "FEE_COLLECTION INC"
	} else if transType == 40 && isAuth == "N" && requestCategory == "P" {
		return "FEE_COLLECTION OUT"
	} else {
		return "DECLINED"
	}
}

func getCardNumber(sCat, tCat, sourceNumber, targetNumber string, length int, firstSixDigits bool) (string, string) {
	if sCat == "C" && len(sourceNumber) > 13 {
		if !firstSixDigits {
			return sourceNumber[len(sourceNumber)-length:], targetNumber
		} else {
			return sourceNumber[0:6], targetNumber
		}
	} else if tCat == "C" && len(targetNumber) > 13 {
		if !firstSixDigits {
			return targetNumber[len(targetNumber)-length:], sourceNumber
		} else {
			return targetNumber[0:6], sourceNumber
		}
	}

	return "", ""
}

func getCardNumberBrand(sCat, tCat, sourceNumber, targetNumber string, length int, firstSixDigits bool) string {
	if sCat == "C" && len(sourceNumber) > 13 {
		if !firstSixDigits {
			return sourceNumber[len(sourceNumber)-length:]
		} else {
			return sourceNumber[0:6]
		}
	} else if tCat == "C" && len(targetNumber) > 13 {
		if !firstSixDigits {
			return targetNumber[len(targetNumber)-length:]
		} else {
			return targetNumber[0:6]
		}
	}

	return ""
}
