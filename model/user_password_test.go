package model_test

import (
	"fmt"
	"testing"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestValidPassword(t *testing.T) {
	policy := model.PasswordPolicy{
		MinPasswordLength:       10,
		RejectCompromised:       true,
		EnforcePasswordStrength: model.PasswordStrengthWeak,
		RequireLowercase:        true,
		RequireUppercase:        true,
		RequireNumber:           true,
		RequireSymbol:           true,
	}
	_, results := policy.Validate("Abcdefg1", true)

	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: model.ValidationRule{
			Description: l.PasswordLengthPolicy,
			Params:      []any{policy.MinPasswordLength},
		},
		Valid: false,
	})
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: model.ValidationRule{
			Description: l.PasswordRejectCompromised,
			Params:      []any{},
		},
		Valid: false,
	})
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: model.ValidationRule{
			Description: l.PasswordRequireLowercase,
			Params:      []any{},
		},
		Valid: true,
	})
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: model.ValidationRule{
			Description: l.PasswordRequireUppercase,
			Params:      []any{},
		},
		Valid: true,
	})
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: model.ValidationRule{
			Description: l.PasswordRequireNumber,
			Params:      []any{},
		},
		Valid: true,
	})
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: model.ValidationRule{
			Description: l.PasswordRequireSymbol,
			Params:      []any{},
		},
		Valid: false,
	})

	fmt.Println(results)
}

func TestSymbolPassword(t *testing.T) {
	policy := model.PasswordPolicy{
		MinPasswordLength:       10,
		RejectCompromised:       true,
		EnforcePasswordStrength: model.PasswordStrengthWeak,
		RequireLowercase:        true,
		RequireUppercase:        true,
		RequireNumber:           true,
		RequireSymbol:           true,
	}

	_, results := policy.Validate("Abcdefg1", true)
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: model.ValidationRule{
			Description: l.PasswordRequireSymbol,
			Params:      []any{},
		},
		Valid: false,
	})

	_, results = policy.Validate("Abcdef!<>g1", true)
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: model.ValidationRule{
			Description: l.PasswordRequireSymbol,
			Params:      []any{},
		},
		Valid: true,
	})

	fmt.Println(results)
}

func TestLengthPassword(t *testing.T) {
	policy := model.PasswordPolicy{
		MinPasswordLength:       10,
		RejectCompromised:       true,
		EnforcePasswordStrength: model.PasswordStrengthWeak,
		RequireLowercase:        true,
		RequireUppercase:        true,
		RequireNumber:           true,
		RequireSymbol:           true,
	}

	_, results := policy.Validate("Abcdefg1", true)
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: model.ValidationRule{
			Description: l.PasswordLengthPolicy,
			Params:      []any{policy.MinPasswordLength},
		},
		Valid: false,
	})

	_, results = policy.Validate("Abcdef!<>g1fffdd", true)
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: model.ValidationRule{
			Description: l.PasswordLengthPolicy,
			Params:      []any{policy.MinPasswordLength},
		},

		Valid: true,
	})

	fmt.Println(results)
}

func TestCompromised(t *testing.T) {
	policy := model.PasswordPolicy{
		MinPasswordLength:       10,
		RejectCompromised:       true,
		EnforcePasswordStrength: model.PasswordStrengthWeak,
		RequireLowercase:        true,
		RequireUppercase:        true,
		RequireNumber:           true,
		RequireSymbol:           true,
	}

	_, results := policy.Validate("Abcdefg1", true)
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: model.ValidationRule{
			Description: l.PasswordRejectCompromised,
			Params:      []any{},
		},
		Valid: false,
	})
	_, results = policy.Validate("Abcdefg1", false)
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: model.ValidationRule{
			Description: l.PasswordRejectCompromised,
			Params:      []any{},
		},
		Valid: true,
	})

	policy.RejectCompromised = false
	_, results = policy.Validate("Abcdefg1", true)
	assert.NotContains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: model.ValidationRule{
			Description: l.PasswordRejectCompromised,
			Params:      []any{},
		},
		Valid: true,
	})
	assert.NotContains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: model.ValidationRule{
			Description: l.PasswordRejectCompromised,
			Params:      []any{},
		},
		Valid: false,
	})

	fmt.Println(results)
}
