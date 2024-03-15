package util

// Can add more currencies in the future. No need for a Database to keep track of what's supported and what isn't
const (
	USD  = "USD"
	CAD  = "CAD"
	EUR  = "EUR"
	LBP  = "LBP"
	DLBP = "DLBP"
)

// To check if its supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, CAD, EUR, LBP, DLBP:
		return true
	}
	return false
}
