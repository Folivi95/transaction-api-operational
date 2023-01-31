package helpers

import "regexp"

type CardBrand string

const (
	CardBrandAmericanExpress CardBrand = "American Express"
	CardBrandDinersClub      CardBrand = "Diners Club"
	CardBrandDiscover        CardBrand = "Discover"
	CardBrandJCB             CardBrand = "JCB"
	CardBrandMasterCard      CardBrand = "MasterCard"
	CardBrandUnionPay        CardBrand = "UnionPay"
	CardBrandVisa            CardBrand = "Visa"
	CardBrandUnknown         CardBrand = "Unknown"
	CardMaestroBrand         CardBrand = "Maestro"
	CardMaestroBrandUK       CardBrand = "Maestro UK"
	defaultFilter                      = "OTHER"
)

var cardBrands = map[CardBrand]string{
	CardBrandAmericanExpress: "American Express",
	CardBrandDinersClub:      "Diners Club",
	CardBrandDiscover:        "Discover",
	CardBrandJCB:             "JCB",
	CardBrandMasterCard:      "MasterCard",
	CardBrandUnionPay:        "UnionPay",
	CardBrandVisa:            "Visa",
	CardMaestroBrand:         "Maestro",
	CardMaestroBrandUK:       "Maestro UK",
	CardBrandUnknown:         "Unknown",
}

// CardBrandGetCardBrands returns a list of CardBrands that match the cardNumber's prefix,
// if any are found; otherwise, CardBrandUnknown.
func CardBrandGetCardBrands(cardNumber string) []CardBrand {
	if cardNumber == "" {
		return []CardBrand{CardBrandUnknown}
	}

	matchingCards := CardBrandGetMatchingCards(cardNumber)
	if len(matchingCards) > 0 {
		return matchingCards
	}

	return []CardBrand{CardBrandUnknown}
}

// CardBrandGetMatchingCards returns a list of CardBrands that match the cardNumber's prefix.
func CardBrandGetMatchingCards(cardNumber string) []CardBrand {
	var matchingCards []CardBrand
	for cardBrand := range cardBrands {
		if cardBrand.GetPatternForLength().MatchString(cardNumber) {
			matchingCards = append(matchingCards, cardBrand)
		}
	}
	return matchingCards
}

// GetPatternForLength returns a regular expression that matches the card number prefix
// for a given card brand.
func (cardBrand CardBrand) GetPatternForLength() *regexp.Regexp {
	switch cardBrand {
	case CardBrandAmericanExpress:
		return regexp.MustCompile("^3[47][0-9]{0,13}$")
	case CardBrandDinersClub:
		return regexp.MustCompile("^(36|30|38|39)[0-9]*$")
	case CardBrandDiscover:
		return regexp.MustCompile("^6(?:011|5[0-9]{2})[0-9]{0,12}$")
	case CardBrandJCB:
		return regexp.MustCompile("^(?:2131|1800|35[0-9]{3})[0-9]{0,11}$")
	case CardBrandMasterCard:
		return regexp.MustCompile("^5[1-5][0-9]{0,14}$")
	case CardBrandUnionPay:
		return regexp.MustCompile("^(62|81)[0-9]*$")
	case CardBrandVisa:
		return regexp.MustCompile("^4[0-9]{0,15}$")
	case CardMaestroBrand:
		return regexp.MustCompile(`^(?:50|5[6-9]|6[0-9])\d+$`)
	case CardMaestroBrandUK:
		return regexp.MustCompile(`^(?:5020|5038|6304|6759|676[1-3])\d{0,8}$`)
	case CardBrandUnknown:
		return regexp.MustCompile("^$")
	}
	return nil
}

func GetCardBrand(cardNumber string) string {
	cardBrand := CardBrandGetCardBrands(cardNumber)[0]

	switch cardBrand {
	case CardBrandVisa:
		return "VISA"
	case CardBrandMasterCard:
		return "MC"
	case CardBrandAmericanExpress:
		return "AMEX"
	case CardBrandUnionPay:
		return "UPI"
	case CardBrandJCB:
		return "JCB"
	case CardMaestroBrand:
		return "MAESTRO"
	case CardMaestroBrandUK:
		return "MAESTRO"
	case CardBrandDinersClub:
		return defaultFilter
	case CardBrandDiscover:
		return defaultFilter
	case CardBrandUnknown:
		return defaultFilter
	}
	return defaultFilter
}

// mapping possible names for O(1) check.
var transCondName = map[string]bool{
	"El Comm Non-Secure w CVV2 Token":   true,
	"El Comm Non-Secure w CVV2":         true,
	"El Comm Non-Secure":                true,
	"Mail/Phone Order Recurring w CVV2": true,
	"El Comm Secure w/o cert Token":     true,
	"El Comm SSL Token":                 true,
	"Mail/Phone Order Single":           true,
	"El Comm SSL w CVV2":                true,
	"El Comm SSL":                       true,
	"Mail/Phone Order Recurring":        true,
	"Mail/Phone Order Single w CVV2":    true,
	"El Comm Secure w/cert Token":       true,
	"El Comm Secure w/cert":             true,
	"El Comm Secure w/o cert":           true,
	"El Comm SSL w CVV2 Token":          true,
	"El Comm Non-Secure Token":          true,
}

func IsCardPresent(name string) bool {
	_, isPresent := transCondName[name]
	return isPresent
}

var targetChannelTransaction = map[string]bool{
	"V": true,
	"E": true,
	"X": true,
	"H": true,
	"C": true,
	"J": true,
	"a": true,
	"b": true,
	"K": true,
	"o": true,
}

func IsInTargetChannel(channel string) bool {
	_, isPresent := targetChannelTransaction[channel]
	return isPresent
}
