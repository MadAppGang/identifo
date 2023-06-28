package model

import (
	"regexp"

	"github.com/madappgang/identifo/v2/localization"
)

var DefaultPasswordPolicy = PasswordPolicy{
	MinPasswordLength:       8,
	RejectCompromised:       false,
	EnforcePasswordStrength: PasswordStrengthNone,
	RequireLowercase:        false,
	RequireUppercase:        false,
	RequireNumber:           false,
	RequireSymbol:           false,
}

type PasswordPolicy struct {
	RestrictMinPasswordLength bool
	MinPasswordLength         int
	RejectCompromised         bool // use HaveBeenPwned passwords to check compromised
	EnforcePasswordStrength   PasswordStrength
	RequireLowercase          bool
	RequireUppercase          bool
	RequireNumber             bool
	RequireSymbol             bool
}

type PasswordPolicyValidationResult struct {
	ValidationRule string
	Valid          bool
}

type PasswordStrength = string

const (
	PasswordStrengthNone    PasswordStrength = "none"
	PasswordStrengthWeak    PasswordStrength = "weak"
	PasswordStrengthAverage PasswordStrength = "average"
	PasswordStrengthStrong  PasswordStrength = "strong"
)

const PasswordSymbols = "!$%^&*()_+{}:@[];'#<>?,./|\\-=?"

// \p{Ll} - lowercase regex
var (
	lowercaseRx = regexp.MustCompile(`\p{Ll}+`)
	uppercaseRx = regexp.MustCompile(`\p{Lu}+`)
	numberRx    = regexp.MustCompile(`\d+`)
	symbolRx    = regexp.MustCompile(`[!\$%\^&\*\(\)_\+{}:@\[\];'#<>\?,\./\|\-=\?]+`)
)

func (pp PasswordPolicy) Validate(p *localization.Printer, pswd string, isCompromised bool) []PasswordPolicyValidationResult {
	result := []PasswordPolicyValidationResult{}

	valid := len(pswd) >= pp.MinPasswordLength
	result = append(result, PasswordPolicyValidationResult{
		p.SD(localization.PasswordLengthPolicy, pp.MinPasswordLength),
		valid,
	})

	if pp.RejectCompromised {
		valid = !isCompromised
		result = append(result, PasswordPolicyValidationResult{
			p.SD(localization.PasswordRejectCompromised),
			valid,
		})
	}

	// TODO: we need port this lib: https://github.com/dwolfhub/zxcvbn-python
	// TODO: or this one: https://github.com/zxcvbn-ts/zxcvbn, to go, because current one is based on deprecated python implementation and now in archived state.
	// if pp.EnforcePasswordStrength {
	// }

	if pp.RequireLowercase {
		result = append(result, PasswordPolicyValidationResult{
			p.SD(localization.PasswordRequireLowercase),
			lowercaseRx.Match([]byte(pswd)),
		})
	}

	if pp.RequireUppercase {
		result = append(result, PasswordPolicyValidationResult{
			p.SD(localization.PasswordRequireUppercase),
			uppercaseRx.Match([]byte(pswd)),
		})
	}

	if pp.RequireNumber {
		result = append(result, PasswordPolicyValidationResult{
			p.SD(localization.PasswordRequireNumber),
			numberRx.Match([]byte(pswd)),
		})
	}

	if pp.RequireSymbol {
		result = append(result, PasswordPolicyValidationResult{
			p.SD(localization.PasswordRequireSymbol),
			symbolRx.Match([]byte(pswd)),
		})
	}

	return result
}
