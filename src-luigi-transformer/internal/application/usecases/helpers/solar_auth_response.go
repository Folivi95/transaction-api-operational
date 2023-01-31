package helpers

var authMap = map[string]string{
	"00": "Approved",
	"01": "RefertoCardIssuer",
	"02": "RefertoCardIssuer, specialcondition",
	"03": "InvalidMerchant",
	"04": "Pickupcard",
	"05": "Donothonor",
	"06": "Error",
	"07": "Pickupcard, specialcondition",
	"09": "Requestinprogress",
	"12": "InvalidTransaction",
	"13": "InvalidAmount",
	"14": "Invalidcardnumber",
	"19": "Re-entertransaction21Noactiontaken",
	"30": "FormatError",
	"41": "Lostcardpickup",
	"43": "Stolencardpickup",
	"51": "Notsufficientfunds",
	"52": "Nocheckingaccount",
	"53": "Nosavingsaccount",
	"54": "Expiredcard",
	"55": "PINincorrect",
	"57": "Transactionnotallowedforcardholder",
	"58": "Transactionnotallowedformerchant",
	"61": "Exceedswithdrawalamountlimit",
	"62": "Restrictedcard",
	"63": "Securityviolation",
	"65": "Activitycountlimitexceeded",
	"75": "PINtriesexceeded",
	"77": "Inconsistentwithoriginal",
	"78": "Noaccount",
	"84": "Pre-authorizationtimetoogreat",
	"85": "Noreasonfordeclinearequestforcardverification",
	"86": "CannotverifyPIN",
	"91": "Issuerunavailable",
	"92": "Invalidreceivinginstitutionid",
	"93": "Transactionviolateslaw",
	"94": "Duplicatetransaction",
	"96": "Systemmalfunction",
}

// GetAuthResponse todo: add a decline_reason logic (involves refactoring the map).
func GetAuthResponse(code string) (string, bool) {
	if reason := authMap[code]; reason == "Approved" || reason == "Noreasonfordeclinearequestforcardverification" {
		return "APPROVED", true
	}

	return "DECLINED", true
}
