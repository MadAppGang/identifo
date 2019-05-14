package model

import "unicode"

// StrongPswd validates password
func StrongPswd(pswd string) error {
	seven, uppper, _, _, invalid := verifyPassword(pswd)
	if invalid {
		return ErrorPasswordWrongSymbols
	} else if !seven {
		return ErrorPasswordShouldHave6Letters
	} else if !uppper {
		return ErrorPasswordNoUppercase
	}
	return nil
}

func verifyPassword(s string) (sixLettersOrMore, upper, number, special, invalid bool) {
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
	sixLettersOrMore = letters >= 6
	return
}
