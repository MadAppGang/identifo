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
	RestrictMinPasswordLength bool             `yaml:"restrictMinPasswordLength" json:"restrict_min_password_length"`
	MinPasswordLength         int              `yaml:"minPasswordLength" json:"min_password_length"`
	RejectCompromised         bool             `yaml:"rejectCompromised" json:"reject_compromised"` // use HaveBeenPwned passwords to check compromised
	EnforcePasswordStrength   PasswordStrength `yaml:"enforcePasswordStrength" json:"enforce_password_strength"`
	RequireLowercase          bool             `yaml:"requireLowercase" json:"require_lowercase"`
	RequireUppercase          bool             `yaml:"requireUppercase" json:"require_uppercase"`
	RequireNumber             bool             `yaml:"requireNumber" json:"require_number"`
	RequireSymbol             bool             `yaml:"requireSymbol" json:"require_symbol"`
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

func (pp PasswordPolicy) Validate(pswd string, isCompromised bool, p *localization.Printer) (bool, []PasswordPolicyValidationResult) {
	result := []PasswordPolicyValidationResult{}
	v := true

	valid := len(pswd) >= pp.MinPasswordLength
	result = append(result, PasswordPolicyValidationResult{
		p.SD(localization.PasswordLengthPolicy, pp.MinPasswordLength),
		valid,
	})
	if !valid {
		v = false
	}

	if pp.RejectCompromised {
		valid = !isCompromised
		result = append(result, PasswordPolicyValidationResult{
			p.SD(localization.PasswordRejectCompromised),
			valid,
		})
		if !valid {
			v = false
		}
	}

	// TODO: we need port this lib: https://github.com/dwolfhub/zxcvbn-python
	// TODO: or this one: https://github.com/zxcvbn-ts/zxcvbn, to go, because current one is based on deprecated python implementation and now in archived state.
	// if pp.EnforcePasswordStrength {
	// }

	if pp.RequireLowercase {
		valid = lowercaseRx.Match([]byte(pswd))
		result = append(result, PasswordPolicyValidationResult{
			p.SD(localization.PasswordRequireLowercase),
			valid,
		})
		if !valid {
			v = false
		}
	}

	if pp.RequireUppercase {
		valid = uppercaseRx.Match([]byte(pswd))
		result = append(result, PasswordPolicyValidationResult{
			p.SD(localization.PasswordRequireUppercase),
			valid,
		})
		if !valid {
			v = false
		}
	}

	if pp.RequireNumber {
		valid = numberRx.Match([]byte(pswd))
		result = append(result, PasswordPolicyValidationResult{
			p.SD(localization.PasswordRequireNumber),
			valid,
		})
		if !valid {
			v = false
		}
	}

	if pp.RequireSymbol {
		valid = symbolRx.Match([]byte(pswd))
		result = append(result, PasswordPolicyValidationResult{
			p.SD(localization.PasswordRequireSymbol),
			valid,
		})
		if !valid {
			v = false
		}
	}

	return v, result
}
