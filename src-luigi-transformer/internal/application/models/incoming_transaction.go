package models

import (
	"encoding/json"
)

type W4IncomingTransaction struct {
	Table     string `json:"table"`
	OpType    string `json:"op_type"`
	OpTS      string `json:"op_ts"`
	CurrentTS string `json:"current_ts"`
	Pos       string `json:"pos"`
	Before    struct {
		AcqRefNumber    string  `json:"ACQ_REF_NUMBER"`
		Action          string  `json:"ACTION"`
		AddInfo         string  `json:"ADD_INFO"`
		AmndDate        string  `json:"AMND_DATE"`
		AmndOfficer     int64   `json:"AMND_OFFICER"`
		AmndPrev        int64   `json:"AMND_PREV"`
		AmndState       string  `json:"AMND_STATE"`
		AuthCode        string  `json:"AUTH_CODE"`
		BinRecord       int64   `json:"BIN_RECORD"`
		CardExpire      string  `json:"CARD_EXPIRE"`
		CardSeqvNumber  string  `json:"CARD_SEQV_NUMBER"`
		ChangeVersion   int64   `json:"CHANGE_VERSION"`
		CommentText     string  `json:"COMMENT_TEXT"`
		DocChainID      int64   `json:"DOC__CHAIN__ID"`
		DocOrigID       int64   `json:"DOC__ORIG__ID"`
		DocPrevID       int64   `json:"DOC__PREV__ID"`
		DocSummID       int64   `json:"DOC__SUMM__ID"`
		FxSettlDate     string  `json:"FX_SETTL_DATE"`
		ID              int64   `json:"ID"`
		IsAuthorization string  `json:"IS_AUTHORIZATION"`
		IssRefNumber    string  `json:"ISS_REF_NUMBER"`
		MerchantID      string  `json:"MERCHANT_ID"`
		MessageCategory string  `json:"MESSAGE_CATEGORY"`
		NumberInChain   int64   `json:"NUMBER_IN_CHAIN"`
		NumberOfSubS    int64   `json:"NUMBER_OF_SUB_S"`
		NwRefDate       string  `json:"NW_REF_DATE"`
		OutwardStatus   string  `json:"OUTWARD_STATUS"`
		PartitionKey    string  `json:"PARTITION_KEY"`
		PostingDate     string  `json:"POSTING_DATE"`
		PostingStatus   string  `json:"POSTING_STATUS"`
		PsRefNumber     string  `json:"PS_REF_NUMBER"`
		ReasonCode      string  `json:"REASON_CODE"`
		ReasonDetails   string  `json:"REASON_DETAILS"`
		RecDate         string  `json:"REC_DATE"`
		RecMemberID     string  `json:"REC_MEMBER_ID"`
		ReconsAmount    float64 `json:"RECONS_AMOUNT"`
		ReconsCurr      string  `json:"RECONS_CURR"`
		RequestCategory string  `json:"REQUEST_CATEGORY"`
		Requirement     string  `json:"REQUIREMENT"`
		RetRefNumber    string  `json:"RET_REF_NUMBER"`
		ReturnCode      int64   `json:"RETURN_CODE"`
		SCat            string  `json:"S_CAT"`
		SecTransCondAtt int64   `json:"SEC_TRANS_COND_ATT"`
		SecTransDate    string  `json:"SEC_TRANS_DATE"`
		SendMemberID    string  `json:"SEND_MEMBER_ID"`
		SendingBin      string  `json:"SENDING_BIN"`
		ServiceClass    string  `json:"SERVICE_CLASS"`
		SettlAmount     float64 `json:"SETTL_AMOUNT"`
		SettlCurr       string  `json:"SETTL_CURR"`
		SicCode         string  `json:"SIC_CODE"`
		SourceAccType   string  `json:"SOURCE_ACC_TYPE"`
		SourceChannel   string  `json:"SOURCE_CHANNEL"`
		SourceCode      string  `json:"SOURCE_CODE"`
		SourceContract  int64   `json:"SOURCE_CONTRACT"`
		SourceFeeAmount float64 `json:"SOURCE_FEE_AMOUNT"`
		SourceFeeCode   string  `json:"SOURCE_FEE_CODE"`
		SourceFeeCurr   string  `json:"SOURCE_FEE_CURR"`
		SourceIdtScheme string  `json:"SOURCE_IDT_SCHEME"`
		SourceMemberID  string  `json:"SOURCE_MEMBER_ID"`
		SourceNumber    string  `json:"SOURCE_NUMBER"`
		SourceRegNum    string  `json:"SOURCE_REG_NUM"`
		SourceService   int64   `json:"SOURCE_SERVICE"`
		SourceSpc       string  `json:"SOURCE_SPC"`
		SynchTag        string  `json:"SYNCH_TAG"`
		TCat            string  `json:"T_CAT"`
		TargetAccType   string  `json:"TARGET_ACC_TYPE"`
		TargetBinID     int64   `json:"TARGET_BIN_ID"`
		TargetChannel   string  `json:"TARGET_CHANNEL"`
		TargetCode      string  `json:"TARGET_CODE"`
		TargetContract  int64   `json:"TARGET_CONTRACT"`
		TargetCountry   string  `json:"TARGET_COUNTRY"`
		TargetFeeAmount float64 `json:"TARGET_FEE_AMOUNT"`
		TargetFeeCode   string  `json:"TARGET_FEE_CODE"`
		TargetFeeCurr   string  `json:"TARGET_FEE_CURR"`
		TargetIdtScheme string  `json:"TARGET_IDT_SCHEME"`
		TargetMemberID  string  `json:"TARGET_MEMBER_ID"`
		TargetNumber    string  `json:"TARGET_NUMBER"`
		TargetService   int64   `json:"TARGET_SERVICE"`
		TargetSpc       string  `json:"TARGET_SPC"`
		TransAmount     float64 `json:"TRANS_AMOUNT"`
		TransCity       string  `json:"TRANS_CITY"`
		TransCondAttr   int64   `json:"TRANS_COND_ATTR"`
		TransCondition  string  `json:"TRANS_CONDITION"`
		TransCountry    string  `json:"TRANS_COUNTRY"`
		TransCurr       string  `json:"TRANS_CURR"`
		TransDate       string  `json:"TRANS_DATE"`
		TransDetails    string  `json:"TRANS_DETAILS"`
		TransState      string  `json:"TRANS_STATE"`
		TransType       int64   `json:"TRANS_TYPE"`
	} `json:"before"`
	After struct {
		AcqRefNumber    string  `json:"ACQ_REF_NUMBER"`
		Action          string  `json:"ACTION"`
		AddInfo         string  `json:"ADD_INFO"`
		AmndDate        string  `json:"AMND_DATE"`
		AmndOfficer     int64   `json:"AMND_OFFICER"`
		AmndPrev        int64   `json:"AMND_PREV"`
		AmndState       string  `json:"AMND_STATE"`
		AuthCode        string  `json:"AUTH_CODE"`
		BinRecord       int64   `json:"BIN_RECORD"`
		CardExpire      string  `json:"CARD_EXPIRE"`
		CardSeqvNumber  string  `json:"CARD_SEQV_NUMBER"`
		ChangeVersion   int64   `json:"CHANGE_VERSION"`
		CommentText     string  `json:"COMMENT_TEXT"`
		DocChainID      int64   `json:"DOC__CHAIN__ID"`
		DocOrigID       int64   `json:"DOC__ORIG__ID"`
		DocPrevID       int64   `json:"DOC__PREV__ID"`
		DocSummID       int64   `json:"DOC__SUMM__ID"`
		FxSettlDate     string  `json:"FX_SETTL_DATE"`
		ID              int64   `json:"ID"`
		IsAuthorization string  `json:"IS_AUTHORIZATION"`
		IssRefNumber    string  `json:"ISS_REF_NUMBER"`
		MerchantID      string  `json:"MERCHANT_ID"`
		MessageCategory string  `json:"MESSAGE_CATEGORY"`
		NumberInChain   int64   `json:"NUMBER_IN_CHAIN"`
		NumberOfSubS    int64   `json:"NUMBER_OF_SUB_S"`
		NwRefDate       string  `json:"NW_REF_DATE"`
		OutwardStatus   string  `json:"OUTWARD_STATUS"`
		PartitionKey    string  `json:"PARTITION_KEY"`
		PostingDate     string  `json:"POSTING_DATE"`
		PostingStatus   string  `json:"POSTING_STATUS"`
		PsRefNumber     string  `json:"PS_REF_NUMBER"`
		ReasonCode      string  `json:"REASON_CODE"`
		ReasonDetails   string  `json:"REASON_DETAILS"`
		RecDate         string  `json:"REC_DATE"`
		RecMemberID     string  `json:"REC_MEMBER_ID"`
		ReconsAmount    float64 `json:"RECONS_AMOUNT"`
		ReconsCurr      string  `json:"RECONS_CURR"`
		RequestCategory string  `json:"REQUEST_CATEGORY"`
		Requirement     string  `json:"REQUIREMENT"`
		RetRefNumber    string  `json:"RET_REF_NUMBER"`
		ReturnCode      int64   `json:"RETURN_CODE"`
		SaltTokenID     string  `json:"SALT_TOKEN_ID"`
		SCat            string  `json:"S_CAT"`
		SecTransCondAtt int64   `json:"SEC_TRANS_COND_ATT"`
		SecTransDate    string  `json:"SEC_TRANS_DATE"`
		SendMemberID    string  `json:"SEND_MEMBER_ID"`
		SendingBin      string  `json:"SENDING_BIN"`
		ServiceClass    string  `json:"SERVICE_CLASS"`
		SettlAmount     float64 `json:"SETTL_AMOUNT"`
		SettlCurr       string  `json:"SETTL_CURR"`
		SicCode         string  `json:"SIC_CODE"`
		SourceAccType   string  `json:"SOURCE_ACC_TYPE"`
		SourceChannel   string  `json:"SOURCE_CHANNEL"`
		SourceCode      string  `json:"SOURCE_CODE"`
		SourceContract  int64   `json:"SOURCE_CONTRACT"`
		SourceFeeAmount float64 `json:"SOURCE_FEE_AMOUNT"`
		SourceFeeCode   string  `json:"SOURCE_FEE_CODE"`
		SourceFeeCurr   string  `json:"SOURCE_FEE_CURR"`
		SourceIdtScheme string  `json:"SOURCE_IDT_SCHEME"`
		SourceMemberID  string  `json:"SOURCE_MEMBER_ID"`
		SourceNumber    string  `json:"SOURCE_NUMBER"`
		SourceRegNum    string  `json:"SOURCE_REG_NUM"`
		SourceService   int64   `json:"SOURCE_SERVICE"`
		SourceSpc       string  `json:"SOURCE_SPC"`
		SynchTag        string  `json:"SYNCH_TAG"`
		TCat            string  `json:"T_CAT"`
		TargetAccType   string  `json:"TARGET_ACC_TYPE"`
		TargetBinID     int64   `json:"TARGET_BIN_ID"`
		TargetChannel   string  `json:"TARGET_CHANNEL"`
		TargetCode      string  `json:"TARGET_CODE"`
		TargetContract  int64   `json:"TARGET_CONTRACT"`
		TargetCountry   string  `json:"TARGET_COUNTRY"`
		TargetFeeAmount float64 `json:"TARGET_FEE_AMOUNT"`
		TargetFeeCode   string  `json:"TARGET_FEE_CODE"`
		TargetFeeCurr   string  `json:"TARGET_FEE_CURR"`
		TargetIdtScheme string  `json:"TARGET_IDT_SCHEME"`
		TargetMemberID  string  `json:"TARGET_MEMBER_ID"`
		TargetNumber    string  `json:"TARGET_NUMBER"`
		TargetService   int64   `json:"TARGET_SERVICE"`
		TargetSpc       string  `json:"TARGET_SPC"`
		TransAmount     float64 `json:"TRANS_AMOUNT"`
		TransCity       string  `json:"TRANS_CITY"`
		TransCondAttr   int64   `json:"TRANS_COND_ATTR"`
		TransCondition  string  `json:"TRANS_CONDITION"`
		TransCountry    string  `json:"TRANS_COUNTRY"`
		TransCurr       string  `json:"TRANS_CURR"`
		TransDate       string  `json:"TRANS_DATE"`
		TransDetails    string  `json:"TRANS_DETAILS"`
		TransState      string  `json:"TRANS_STATE"`
		TransType       int64   `json:"TRANS_TYPE"`
	} `json:"after"`
}

func W4IncomingTransactionFromBytes(in []byte) (W4IncomingTransaction, error) {
	var incomingInstruction W4IncomingTransaction
	err := json.Unmarshal(in, &incomingInstruction)

	return incomingInstruction, err
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
			SaltTokenID     string `json:"saltTokenId"`
			Class           string `json:"class"`
			Category        string `json:"category"`
			Direction       string `json:"direction"`
			TransactionDate string `json:"transactionDate"`
			SettlementDate  string `json:"settlementDate"`
			SystemDate      string `json:"systemDate"`
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
					Fee []Fee `json:"fee"`
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

type Fee struct {
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
}

func SolarIncomingTransactionFromBytes(in []byte) (SolarIncomingTransaction, error) {
	var incomingInstruction SolarIncomingTransaction
	err := json.Unmarshal(in, &incomingInstruction)

	return incomingInstruction, err
}
