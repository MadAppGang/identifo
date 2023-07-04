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
	ValidationRule ValidationRule
	Valid          bool
}

// Validation Rule is localizable human-readable rule description with params to be rendered.
type ValidationRule struct {
	Description localization.LocalizedString
	Params      []any
}

func (vr PasswordPolicyValidationResult) Error() error {
	if vr.Valid == true {
		return nil // no error
	}
	return localization.LocalizedError{
		ErrID:   vr.ValidationRule.Description,
		Details: vr.ValidationRule.Params,
	}
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

func (pp PasswordPolicy) Validate(pswd string, isCompromised bool) (bool, []PasswordPolicyValidationResult) {
	result := []PasswordPolicyValidationResult{}
	v := true

	valid := len(pswd) >= pp.MinPasswordLength
	result = append(result, PasswordPolicyValidationResult{
		ValidationRule{
			Description: localization.PasswordLengthPolicy,
			Params:      []any{pp.MinPasswordLength},
		},
		valid,
	})
	if !valid {
		v = false
	}

	if pp.RejectCompromised {
		valid = !isCompromised
		result = append(result, PasswordPolicyValidationResult{
			ValidationRule{
				Description: localization.PasswordRejectCompromised,
				Params:      []any{},
			},
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
			ValidationRule{
				Description: localization.PasswordRequireLowercase,
				Params:      []any{},
			},
			valid,
		})
		if !valid {
			v = false
		}
	}

	if pp.RequireUppercase {
		valid = uppercaseRx.Match([]byte(pswd))
		result = append(result, PasswordPolicyValidationResult{
			ValidationRule{
				Description: localization.PasswordRequireUppercase,
				Params:      []any{},
			},
			valid,
		})
		if !valid {
			v = false
		}
	}

	if pp.RequireNumber {
		valid = numberRx.Match([]byte(pswd))
		result = append(result, PasswordPolicyValidationResult{
			ValidationRule{
				Description: localization.PasswordRequireNumber,
				Params:      []any{},
			},
			valid,
		})
		if !valid {
			v = false
		}
	}

	if pp.RequireSymbol {
		valid = symbolRx.Match([]byte(pswd))
		result = append(result, PasswordPolicyValidationResult{
			ValidationRule{
				Description: localization.PasswordRequireSymbol,
				Params:      []any{},
			},
			valid,
		})
		if !valid {
			v = false
		}
	}

	return v, result
}
