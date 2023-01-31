package usecases_test

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/matryer/is"

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/testhelpers"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/models"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/ports/mocks"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/usecases"
)

func TestW4Transformer(t *testing.T) {
	is, metricsClient, dbHandlerMock := mockClients(t)

	expectedTxn := LoadCanonicalTransaction()
	w4TransactionJSON := LoadW4Transaction()

	transformer := usecases.NewW4Transformer(dbHandlerMock, nil, "w4schema", metricsClient)

	var w4Transaction models.W4IncomingTransaction
	err := json.Unmarshal(w4TransactionJSON, &w4Transaction)
	is.NoErr(err)
	txn, err := transformer.Translate(context.TODO(), w4Transaction)
	is.NoErr(err)
	// return true if expectedTxn is equal to txn interface
	is.Equal(expectedTxn, txn)
}

func mockClients(t *testing.T) (*is.I, testhelpers.DummyMetricsClient, *mocks.DBHandlerMock) {
	is := is.New(t)
	dummyString := "NULL"
	dummyInt := 100
	dummyTrue := "TRUE"
	metricsClient := testhelpers.DummyMetricsClient{}
	dbHandlerMock := mocks.DBHandlerMock{
		GetFunc: func(ctx context.Context, query string, target []interface{}) error {
			for _, t := range target {
				switch p := t.(type) {
				case *string:
					if strings.Contains(query, "TRUE") {
						*p = dummyTrue
					} else {
						*p = dummyString
					}
				case *int64:
					*p = int64(dummyInt)
				}
			}
			return nil
		},
	}
	return is, metricsClient, &dbHandlerMock
}

func LoadW4Transaction() []byte {
	return []byte(`{
	"table": "OWS.DOC",
	"op_type": "U",
	"op_ts": "2022-04-14 12:23:53.000",
	"current_ts": "2022-04-14 13:36:34.058",
	"pos": "00000000020000285280",
	"before": {
		"AMND_STATE": "A",
		"AMND_OFFICER": 14400,
		"AMND_DATE": "2022-04-14 12:23:46.000",
		"AMND_PREV": 12413650420,
		"ID": 12413650420,
		"DOC__ORIG__ID": 0,
		"DOC__PREV__ID": 12413650410,
		"DOC__SUMM__ID": 0,
		"NUMBER_OF_SUB_S": 0,
		"ACTION": "0",
		"MESSAGE_CATEGORY": "0",
		"SOURCE_REG_NUM": "M104116TSNQS",
		"ACQ_REF_NUMBER": "0",
		"RET_REF_NUMBER": "210407333288",
		"ISS_REF_NUMBER": "96203234122044541333______0037",
		"PS_REF_NUMBER": "MCC011DOW0414",
		"NW_REF_DATE": "2022-04-14 12:23:45.000",
		"AUTH_CODE": "005988",
		"IS_AUTHORIZATION": "P",
		"REQUEST_CATEGORY": "R",
		"SERVICE_CLASS": "T",
		"SOURCE_CODE": "0",
		"SOURCE_FEE_CODE": "0",
		"TARGET_CODE": "0",
		"TARGET_FEE_CODE": "0",
		"TRANS_TYPE": 5,
		"SOURCE_CHANNEL": "P",
		"S_CAT": "M",
		"SOURCE_IDT_SCHEME": "0",
		"SOURCE_MEMBER_ID": "1501",
		"REC_MEMBER_ID": "0",
		"SOURCE_NUMBER": "44869746",
		"SOURCE_SPC": "0",
		"SOURCE_ACC_TYPE": "0",
		"SOURCE_CONTRACT": 3505810,
		"SOURCE_SERVICE": 340892530,
		"TARGET_CHANNEL": "C",
		"T_CAT": "C",
		"TARGET_IDT_SCHEME": "0",
		"TARGET_MEMBER_ID": "81000000",
		"SEND_MEMBER_ID": "005302",
		"SENDING_BIN": "545302",
		"TARGET_BIN_ID": 0,
		"TARGET_NUMBER": "5555555555554444",
		"TARGET_SPC": "00",
		"TARGET_ACC_TYPE": "0",
		"TARGET_CONTRACT": 659560,
		"TARGET_SERVICE": 330636130,
		"CARD_EXPIRE": "2512",
		"CARD_SEQV_NUMBER": "0",
		"MERCHANT_ID": "9117052",
		"SIC_CODE": "5814",
		"TRANS_CONDITION": "POM",
		"TRANS_COND_ATTR": 93020,
		"SEC_TRANS_COND_ATT": 93030,
		"RECONS_CURR": "352",
		"RECONS_AMOUNT": 10,
		"SETTL_CURR": "352",
		"SETTL_AMOUNT": 10,
		"SOURCE_FEE_CURR": "0",
		"SOURCE_FEE_AMOUNT": 0,
		"TARGET_FEE_CURR": "0",
		"TARGET_FEE_AMOUNT": 0,
		"TRANS_DATE": "2022-04-14 18:03:41.000",
		"SEC_TRANS_DATE": "2022-04-14 12:23:45.000",
		"TRANS_COUNTRY": "ISL",
		"TRANS_STATE": "0",
		"TRANS_CITY": "Reykjavik",
		"TRANS_DETAILS": "Borgun PAX",
		"TRANS_AMOUNT": 10,
		"TRANS_CURR": "352",
		"REASON_CODE": "0",
		"REASON_DETAILS": "0",
		"REQUIREMENT": "0",
		"POSTING_DATE": "0",
		"FX_SETTL_DATE": "2022-04-14 00:00:00.000",
		"POSTING_STATUS": "C",
		"OUTWARD_STATUS": "C",
		"RETURN_CODE": 0,
		"PARTITION_KEY": "X",
		"SYNCH_TAG": "0",
		"BIN_RECORD": 78627,
		"NUMBER_IN_CHAIN": 1,
		"ADD_INFO": "TGTCC=PT;PRODUCT_ID=000;POSTAL_CODE=108;TRANS_LOCATION=Armula 30;TOP_CONTRACT=3505800;MERCH_COUNTRY=ISL;MERCH_SIC=5399;UPD=00;MERCH_CITY=ASDFGH ;VPD=00;MERCH_NAME=BORGUNTESTING ;MSG_AMOUNT=0;ORIG_AMOUNT=10;R...",
		"CHANGE_VERSION": 0,
		"TARGET_COUNTRY": "ISL",
		"REC_DATE": "2022-04-14 12:23:46.000",
		"COMMENT_TEXT": "0",
		"DOC__CHAIN__ID": 0
	},
	"after": {
		"AMND_STATE": "A",
		"AMND_OFFICER": 14400,
		"AMND_DATE": "2022-04-14 12:23:46.000",
		"AMND_PREV": 12413650420,
		"ID": 12413650420,
		"DOC__ORIG__ID": 0,
		"DOC__PREV__ID": 12413650410,
		"DOC__SUMM__ID": 0,
		"NUMBER_OF_SUB_S": 0,
		"ACTION": "0",
		"MESSAGE_CATEGORY": "0",
		"SOURCE_REG_NUM": "M104116TSNQS",
		"ACQ_REF_NUMBER": "0",
		"RET_REF_NUMBER": "210407333288",
		"ISS_REF_NUMBER": "96203234122044541333______0037",
		"PS_REF_NUMBER": "MCC011DOW0414",
		"NW_REF_DATE": "2022-04-14 12:23:45.000",
		"AUTH_CODE": "005988",
		"IS_AUTHORIZATION": "P",
		"REQUEST_CATEGORY": "R",
		"SERVICE_CLASS": "T",
		"SOURCE_CODE": "0",
		"SOURCE_FEE_CODE": "0",
		"TARGET_CODE": "0",
		"TARGET_FEE_CODE": "0",
		"TRANS_TYPE": 5,
		"SOURCE_CHANNEL": "P",
		"S_CAT": "M",
		"SOURCE_IDT_SCHEME": "0",
		"SOURCE_MEMBER_ID": "1501",
		"REC_MEMBER_ID": "0",
		"SOURCE_NUMBER": "44869746",
		"SOURCE_SPC": "0",
		"SOURCE_ACC_TYPE": "0",
		"SOURCE_CONTRACT": 3505810,
		"SOURCE_SERVICE": 340892530,
		"TARGET_CHANNEL": "C",
		"T_CAT": "C",
		"TARGET_IDT_SCHEME": "0",
		"TARGET_MEMBER_ID": "81000000",
		"SEND_MEMBER_ID": "005302",
		"SENDING_BIN": "545302",
		"TARGET_BIN_ID": 0,
		"TARGET_NUMBER": "5555555555554444",
		"TARGET_SPC": "00",
		"TARGET_ACC_TYPE": "0",
		"TARGET_CONTRACT": 659560,
		"TARGET_SERVICE": 330636130,
		"CARD_EXPIRE": "2512",
		"CARD_SEQV_NUMBER": "0",
		"MERCHANT_ID": "9117052",
		"SIC_CODE": "5814",
		"TRANS_CONDITION": "POM",
		"TRANS_COND_ATTR": 93020,
		"SEC_TRANS_COND_ATT": 93030,
		"RECONS_CURR": "352",
		"RECONS_AMOUNT": 10,
		"SETTL_CURR": "352",
		"SETTL_AMOUNT": 10,
		"SOURCE_FEE_CURR": "0",
		"SOURCE_FEE_AMOUNT": 0,
		"TARGET_FEE_CURR": "0",
		"TARGET_FEE_AMOUNT": 0,
		"TRANS_DATE": "2022-04-14 18:03:41.000",
		"SEC_TRANS_DATE": "2022-04-14 12:23:45.000",
		"TRANS_COUNTRY": "ISL",
		"TRANS_STATE": "0",
		"TRANS_CITY": "Reykjavik",
		"TRANS_DETAILS": "Borgun PAX",
		"TRANS_AMOUNT": 10,
		"TRANS_CURR": "352",
		"REASON_CODE": "0",
		"REASON_DETAILS": "0",
		"REQUIREMENT": "0",
		"POSTING_DATE": "0",
		"FX_SETTL_DATE": "2022-04-14 00:00:00.000",
		"POSTING_STATUS": "C",
		"OUTWARD_STATUS": "C",
		"RETURN_CODE": 0,
		"PARTITION_KEY": "X",
		"SYNCH_TAG": "0",
		"BIN_RECORD": 78627,
		"NUMBER_IN_CHAIN": 1,
		"ADD_INFO": "POSTAL_CODE=108;TRANS_LOCATION=Armula 30;TOP_CONTRACT=3505800;MERCH_COUNTRY=ISL;MERCH_SIC=5399;UPD=00;MERCH_CITY=ASDFGH ;VPD=00;MERCH_NAME=BORGUNTESTING ;MSG_AMOUNT=0;ORIG_AMOUNT=10;R...",
		"CHANGE_VERSION": 0,
		"TARGET_COUNTRY": "ISL",
		"REC_DATE": "2022-04-14 12:23:46.000",
		"COMMENT_TEXT": "0",
		"DOC__CHAIN__ID": 0
	}
}`)
}

func LoadCanonicalTransaction() models.Transaction {
	return models.Transaction{
		Institutions: []models.Institution{{
			Name: "Way4",
			Type: "ACQUIRER",
		}},
		InternalID:              "NULL",
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
			ID:   "44869746",
			Type: "MISSING_INFO"},
		Response: models.Response{
			Code:           "0",
			Classification: "APPROVED",
			Description:    "MISSING_INFO",
		},
		// AuthorisationResponse:   getTransactionAuthResponse(incomingTransaction.After.ReturnCode),
		Card: models.Card{
			PanFirstSix:    "555555",
			PanLastFour:    "4444",
			Brand:          "MC",
			ExpiryMonth:    "NU",
			ExpiryYear:     "LL",
			Issuer:         "MISSING_INFO",
			IssuingProduct: "",
			IssuingCountry: "",
			TokenID:        "NULL",
			Currency:       "ISK",
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
						Type:               "MISSING_INFO",
						FixedAmount:        0,
						VariablePercentage: 0,
						VariableAmount:     0,
						Currency:           "MISSING_INFO",
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

func Test_ProductIDAndTgtcc(t *testing.T) {
	is, metricsClient, dbHandlerMock := mockClients(t)

	type args struct {
		transaction    models.Transaction
		w4Transactions []byte
		IssuingCountry string
		IssuingProduct string
		AddInfo        string
	}
	tests := []struct {
		name string
		args args
		want models.Transaction
	}{
		{
			name: "Should get PRODUCT_ID and TGTCC from the incoming transaction",
			args: args{
				transaction:    LoadCanonicalTransaction(),
				w4Transactions: LoadW4Transaction(),
				IssuingCountry: "PT",
				IssuingProduct: "000",
				AddInfo:        "POSTAL_CODE=108;PRODUCT_ID=000;TGTCC=PT;TRANS_LOCATION=Armula 30;TOP_CONTRACT=3505800;MERCH_COUNTRY=ISL;MERCH_SIC=5399;UPD=00;MERCH_CITY=ASDFGH ;VPD=00;MERCH_NAME=BORGUNTESTING ;MSG_AMOUNT=0;ORIG_AMOUNT=10;R...",
			},
			want: LoadCanonicalTransaction(),
		},
		{
			name: "test the case where PRODUCT_ID and TGTCC are not present in the incoming transaction",
			args: args{
				transaction:    LoadCanonicalTransaction(),
				w4Transactions: LoadW4Transaction(),
				IssuingCountry: "",
				IssuingProduct: "",
				AddInfo:        "POSTAL_CODE=108;TRANS_LOCATION=Armula 30;TOP_CONTRACT=3505800;MERCH_COUNTRY=ISL;MERCH_SIC=5399;UPD=00;MERCH_CITY=ASDFGH ;VPD=00;MERCH_NAME=BORGUNTESTING ;MSG_AMOUNT=0;ORIG_AMOUNT=10;R...",
			},
			want: LoadCanonicalTransaction(),
		},
		{
			name: "test the case where PRODUCT_ID and TGTCC are present but the value is empty",
			args: args{
				transaction:    LoadCanonicalTransaction(),
				w4Transactions: LoadW4Transaction(),
				IssuingCountry: "",
				IssuingProduct: "",
				AddInfo:        "POSTAL_CODE=108;PRODUCT_ID=;TGTCC=;TRANS_LOCATION=Armula 30;TOP_CONTRACT=3505800;MERCH_COUNTRY=ISL;MERCH_SIC=5399;UPD=00;MERCH_CITY=ASDFGH ;VPD=00;MERCH_NAME=BORGUNTESTING ;MSG_AMOUNT=0;ORIG_AMOUNT=10;R...",
			},
			want: LoadCanonicalTransaction(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := usecases.NewW4Transformer(dbHandlerMock, nil, "w4schema", metricsClient)
			var w4Transaction models.W4IncomingTransaction
			err := json.Unmarshal(tt.args.w4Transactions, &w4Transaction)
			is.NoErr(err)
			tt.want.Card.IssuingCountry = tt.args.IssuingCountry
			tt.want.Card.IssuingProduct = tt.args.IssuingProduct

			w4Transaction.After.AddInfo = tt.args.AddInfo
			w4Transaction.Before.AddInfo = tt.args.AddInfo

			got, _ := transformer.Translate(context.TODO(), w4Transaction)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductIDAndTgtcc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_W4TransformerCardBrand(t *testing.T) {
	is, metricsClient, dbHandlerMock := mockClients(t)

	type args struct {
		transaction    models.Transaction
		w4Transactions []byte
		cardBrand      string
		targetNumber   string
	}
	tests := []struct {
		name string
		args args
		want models.Transaction
	}{
		{
			name: "targetNumber(Credit) should return VISA",
			args: args{
				transaction:    LoadCanonicalTransaction(),
				w4Transactions: LoadW4Transaction(),
				cardBrand:      "VISA",
				targetNumber:   "4242424242424242",
			},
			want: LoadCanonicalTransaction(),
		},
		{
			name: "targetNumber(Debit) should return VISA",
			args: args{
				transaction:    LoadCanonicalTransaction(),
				w4Transactions: LoadW4Transaction(),
				cardBrand:      "VISA",
				targetNumber:   "4000056655665556",
			},
			want: LoadCanonicalTransaction(),
		},
		{
			name: "targetNumber(Credit) should return MC",
			args: args{
				transaction:    LoadCanonicalTransaction(),
				w4Transactions: LoadW4Transaction(),
				cardBrand:      "MC",
				targetNumber:   "5555555555554444",
			},
			want: LoadCanonicalTransaction(),
		},
		{
			name: "targetNumber(Debit) should return MC",
			args: args{
				transaction:    LoadCanonicalTransaction(),
				w4Transactions: LoadW4Transaction(),
				cardBrand:      "MC",
				targetNumber:   "5200828282828210",
			},
			want: LoadCanonicalTransaction(),
		},
		{
			name: "targetNumber(credit) should return AMEX",
			args: args{
				transaction:    LoadCanonicalTransaction(),
				w4Transactions: LoadW4Transaction(),
				cardBrand:      "AMEX",
				targetNumber:   "378282246310005",
			},
			want: LoadCanonicalTransaction(),
		},
		{
			name: "targetNumber(debit) should return AMEX",
			args: args{
				transaction:    LoadCanonicalTransaction(),
				w4Transactions: LoadW4Transaction(),
				cardBrand:      "AMEX",
				targetNumber:   "371449635398431",
			},
			want: LoadCanonicalTransaction(),
		},
		// {
		// 	name: "targetNumber(credit) should return UPI",
		// 	args: args{
		// 		transaction:    LoadCanonicalTransaction(),
		// 		w4Transactions: LoadW4Transaction(),
		// 		cardBrand:      "UPI",
		// 		targetNumber:   "6200000000000005",
		// 	},
		// 	want: LoadCanonicalTransaction(),
		// },
		{
			name: "targetNumber(credit) should return JCB",
			args: args{
				transaction:    LoadCanonicalTransaction(),
				w4Transactions: LoadW4Transaction(),
				cardBrand:      "JCB",
				targetNumber:   "3566002020360505",
			},
			want: LoadCanonicalTransaction(),
		},
		{
			name: "targetNumber(credit) DCI should return OTHER",
			args: args{
				transaction:    LoadCanonicalTransaction(),
				w4Transactions: LoadW4Transaction(),
				cardBrand:      "OTHER",
				targetNumber:   "3056930009020004",
			},
			want: LoadCanonicalTransaction(),
		},
		{
			name: "targetNumber(credit) should return MAESTRO",
			args: args{
				transaction:    LoadCanonicalTransaction(),
				w4Transactions: LoadW4Transaction(),
				cardBrand:      "MAESTRO",
				targetNumber:   "6423543159906829",
			},
			want: LoadCanonicalTransaction(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := usecases.NewW4Transformer(dbHandlerMock, nil, "w4schema", metricsClient)
			var w4Transaction models.W4IncomingTransaction
			err := json.Unmarshal(tt.args.w4Transactions, &w4Transaction)
			is.NoErr(err)
			tt.want.Card.Brand = ""
			tt.want.Card.Brand = tt.args.cardBrand
			w4Transaction.After.TargetNumber = tt.args.targetNumber
			w4Transaction.Before.TargetNumber = tt.args.targetNumber

			if got, _ := transformer.Translate(context.TODO(), w4Transaction); got.Card.Brand != tt.want.Card.Brand {
				t.Errorf("Test_W4TransformerCardBrand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsDCCField(t *testing.T) {
	is, metricsClient, dbHandlerMock := mockClients(t)

	type args struct {
		transaction    models.Transaction
		w4Transactions []byte
		isDCC          bool
		addInfo        string
	}
	tests := []struct {
		name string
		args args
		want models.Transaction
	}{
		{
			name: "isDCC should be true if addInfo contains DCC",
			args: args{
				transaction:    LoadCanonicalTransaction(),
				w4Transactions: LoadW4Transaction(),
				isDCC:          true,
				addInfo:        "DCC_IND=Y;POSTAL_CODE=108;TRANS_LOCATION=Armula 30;TOP_CONTRACT=3505800;MERCH_COUNTRY=ISL;MERCH_SIC=5399;UPD=00;MERCH_CITY=ASDFGH ;VPD=00;MERCH_NAME=BORGUNTESTING ;MSG_AMOUNT=0;ORIG_AMOUNT=10;R...",
			},
			want: LoadCanonicalTransaction(),
		},
		{
			name: "isDCC should be false if addInfo contains DCC",
			args: args{
				transaction:    LoadCanonicalTransaction(),
				w4Transactions: LoadW4Transaction(),
				isDCC:          false,
				addInfo:        "DCC_IND=A;POSTAL_CODE=108;TRANS_LOCATION=Armula 30;TOP_CONTRACT=3505800;MERCH_COUNTRY=ISL;MERCH_SIC=5399;UPD=00;MERCH_CITY=ASDFGH ;VPD=00;MERCH_NAME=BORGUNTESTING ;MSG_AMOUNT=0;ORIG_AMOUNT=10;R...",
			},
			want: LoadCanonicalTransaction(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := usecases.NewW4Transformer(dbHandlerMock, nil, "w4schema", metricsClient)
			var w4Transaction models.W4IncomingTransaction
			err := json.Unmarshal(tt.args.w4Transactions, &w4Transaction)
			is.NoErr(err)
			tt.want.IsDcc = tt.args.isDCC
			w4Transaction.After.AddInfo = tt.args.addInfo
			w4Transaction.Before.AddInfo = tt.args.addInfo

			// Compare the two interface values
			got, _ := transformer.Translate(context.TODO(), w4Transaction)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Test_IsDCCField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCardNumber(t *testing.T) {
	is, metricsClient, dbHandlerMock := mockClients(t)

	type args struct {
		transaction    models.Transaction
		w4Transactions []byte
		PanLastFour    string
		terminalID     string
		sCat           string
		tCat           string
		SourceNumber   string
		TargetNumber   string
	}
	tests := []struct {
		name string
		args args
		want models.Transaction
	}{
		{
			name: "getCardNumber should return the last 4 card number",
			args: args{
				transaction:    LoadCanonicalTransaction(),
				w4Transactions: LoadW4Transaction(),
				PanLastFour:    "4545",
				sCat:           "M",
				tCat:           "C",
				SourceNumber:   "5555555555554545",
				terminalID:     "5555555555554545",
				TargetNumber:   "5555555555554545",
			},
			want: LoadCanonicalTransaction(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := usecases.NewW4Transformer(dbHandlerMock, nil, "w4schema", metricsClient)
			var w4Transaction models.W4IncomingTransaction
			err := json.Unmarshal(tt.args.w4Transactions, &w4Transaction)
			is.NoErr(err)
			tt.want.TransactionSource.ID = "5555555555554545"
			tt.want.Card.PanLastFour = tt.args.PanLastFour
			tt.want.Card.Brand = "MC"
			w4Transaction.After.SCat = tt.args.sCat
			w4Transaction.After.TCat = tt.args.tCat
			w4Transaction.After.SourceNumber = tt.args.SourceNumber
			w4Transaction.After.TargetNumber = tt.args.TargetNumber

			// Compare the two interface values
			got, _ := transformer.Translate(context.TODO(), w4Transaction)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Test_getCardNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
