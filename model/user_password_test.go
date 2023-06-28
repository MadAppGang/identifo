package model_test

import (
	"fmt"
	"testing"

	"github.com/madappgang/identifo/v2/localization"
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
	p, _ := localization.NewPrinter("en")
	results := policy.Validate(p, "Abcdefg1", true)

	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: p.SD(localization.PasswordLengthPolicy, policy.MinPasswordLength),
		Valid:          false,
	})
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: p.SD(localization.PasswordRejectCompromised),
		Valid:          false,
	})
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: p.SD(localization.PasswordRequireLowercase),
		Valid:          true,
	})
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: p.SD(localization.PasswordRequireUppercase),
		Valid:          true,
	})
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: p.SD(localization.PasswordRequireNumber),
		Valid:          true,
	})
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: p.SD(localization.PasswordRequireSymbol),
		Valid:          false,
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
	p, _ := localization.NewPrinter("en")

	results := policy.Validate(p, "Abcdefg1", true)
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: p.SD(localization.PasswordRequireSymbol),
		Valid:          false,
	})

	results = policy.Validate(p, "Abcdef!<>g1", true)
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: p.SD(localization.PasswordRequireSymbol),
		Valid:          true,
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
	p, _ := localization.NewPrinter("en")

	results := policy.Validate(p, "Abcdefg1", true)
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: p.SD(localization.PasswordLengthPolicy, policy.MinPasswordLength),
		Valid:          false,
	})

	results = policy.Validate(p, "Abcdef!<>g1fffdd", true)
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: p.SD(localization.PasswordLengthPolicy, policy.MinPasswordLength),
		Valid:          true,
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
	p, _ := localization.NewPrinter("en")

	results := policy.Validate(p, "Abcdefg1", true)
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: p.SD(localization.PasswordRejectCompromised),
		Valid:          false,
	})
	results = policy.Validate(p, "Abcdefg1", false)
	assert.Contains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: p.SD(localization.PasswordRejectCompromised),
		Valid:          true,
	})

	policy.RejectCompromised = false
	results = policy.Validate(p, "Abcdefg1", true)
	assert.NotContains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: p.SD(localization.PasswordRejectCompromised),
		Valid:          true,
	})
	assert.NotContains(t, results, model.PasswordPolicyValidationResult{
		ValidationRule: p.SD(localization.PasswordRejectCompromised),
		Valid:          false,
	})

	fmt.Println(results)
}
