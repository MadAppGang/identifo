package model

import "unicode"

// StrongPswd validates password
func StrongPswd(pswd string) error {
	seven, number, uppper, _, invalid := verifyPassword(pswd)
	if invalid {
		return ErrorPasswordWrongSymbols
	} else if !seven {
		return ErrorPasswordShouldHave7Letter
	} else if !number {
		return ErrorPasswordNoNumbers
	} else if !uppper {
		return ErrorPasswordNoUppercase
	}
	return nil
}

func verifyPassword(s string) (sevenOrMore, number, upper, special, invalid bool) {
	letters := 0
	for _, s := range s {
		switch {
		case unicode.IsNumber(s):
			number = true
		case unicode.IsUpper(s):
			upper = true
			letters++
		case unicode.IsPunct(s) || unicode.IsSymbol(s):
			special = true
		case unicode.IsLetter(s) || s == ' ':
			letters++
		default:
			return false, false, false, false, true
		}
	}
	invalid = false
	sevenOrMore = letters >= 7
	return
}
