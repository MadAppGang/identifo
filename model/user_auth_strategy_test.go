package model_test

import (
	_ "embed"
	"testing"

	"github.com/madappgang/identifo/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestStrategies(t *testing.T) {
	strategies := model.Strategies()
	assert.Len(t, strategies, 3)
}

func TestSecondFactorStrategyCompatible(t *testing.T) {
	s1 := model.SecondFactorStrategy{
		Challenge: model.AuthChallengeTypeOTP,
		Transport: model.AuthTransportTypeEmail,
	}

	s2 := model.SecondFactorStrategy{
		Challenge:   model.AuthChallengeTypeOTP,
		Transport:   model.AuthTransportTypeEmail,
		EnrolPolicy: model.SecondFactorEnrolPolicySelfEnrol,
		Policy:      model.SecondFactorPolicyCustom,
	}

	assert.True(t, s1.Compatible(s2))

	s1.Transport = ""
	s1.Policy = model.SecondFactorPolicyAdaptive
	assert.True(t, s1.Compatible(s2))

	s1.Challenge = ""
	assert.True(t, s1.Compatible(s2))

	s1 = model.SecondFactorStrategy{
		Challenge: model.AuthChallengeTypeOTP,
		Transport: model.AuthTransportTypeEmail,
	}
	s2.Challenge = model.AuthChallengeTypeGuardian
	assert.False(t, s1.Compatible(s2))

	s2.Challenge = ""
	assert.False(t, s1.Compatible(s2))
}

func TestFIMStrategyCompatible(t *testing.T) {
	s1 := model.FIMStrategy{Type: model.FIMStrategyTypeApple}
	s2 := model.FIMStrategy{Type: model.FIMStrategyTypeApple}

	assert.True(t, s1.Compatible(s2))

	s1.Type = ""
	assert.True(t, s1.Compatible(s2))

	s1 = model.FIMStrategy{Type: model.FIMStrategyTypeApple}
	s2.Type = model.FIMStrategyTypeFirebase
	assert.False(t, s1.Compatible(s2))

	s2.Type = ""
	assert.False(t, s1.Compatible(s2))
}

func TestLocalStrategyCompatible(t *testing.T) {
	s1 := model.LocalStrategy{
		Identity: model.AuthIdentityTypePhone,
	}
	s2 := model.LocalStrategy{
		Identity:  model.AuthIdentityTypePhone,
		Challenge: model.AuthChallengeTypeOTP,
		Transport: model.AuthTransportTypePush,
	}

	assert.True(t, s1.Compatible(s2))

	s1.Identity = ""
	assert.True(t, s1.Compatible(s2))

	s1.Identity = model.AuthIdentityTypeUsername
	assert.False(t, s1.Compatible(s2))
}

func TestFirstFactorStrategyCompatible(t *testing.T) {
	s1 := model.FirstFactorStrategy{Type: model.FirstFactorTypeLocal}
	s2 := model.FirstFactorStrategy{
		Type: model.FirstFactorTypeLocal,
		Local: &model.LocalStrategy{
			Identity:  model.AuthIdentityTypeEmail,
			Challenge: model.AuthChallengeTypeOTP,
			Transport: model.AuthTransportTypePush,
		},
	}
	assert.True(t, s1.Compatible(s2))

	s1.Local = &model.LocalStrategy{
		Identity: model.AuthIdentityTypeAnonymous,
	}
	assert.False(t, s1.Compatible(s2))

	s1.Local.Identity = ""
	assert.True(t, s1.Compatible(s2))

	s1.Local.Identity = model.AuthIdentityTypeEmail
	assert.True(t, s1.Compatible(s2))

	s2.Local = nil
	assert.False(t, s1.Compatible(s2))
}

func TestFilterCompatible(t *testing.T) {
	a := model.AuthStrategy{
		Type:        model.AuthStrategyFirstFactor,
		FirstFactor: &model.FirstFactorStrategy{Type: model.FirstFactorTypeLocal},
	}
	a0 := model.AuthStrategy{
		Type: model.AuthStrategyFirstFactor,
		FirstFactor: &model.FirstFactorStrategy{
			Type: model.FirstFactorTypeLocal,
			Local: &model.LocalStrategy{
				Identity:  model.AuthIdentityTypeEmail,
				Challenge: model.AuthChallengeTypeOTP,
				Transport: model.AuthTransportTypePush,
			},
		},
	}
	a1 := model.AuthStrategy{Type: model.AuthStrategyNone}
	a2 := model.AuthStrategy{Type: model.AuthStrategyAnonymous}
	a3 := model.AuthStrategy{
		Type: model.AuthStrategyFirstFactor,
		FirstFactor: &model.FirstFactorStrategy{
			Type: model.FirstFactorTypeFIM,
			FIM: &model.FIMStrategy{
				Type: model.FIMStrategyTypeApple,
			},
		},
	}
	a4 := model.AuthStrategy{
		Type: model.AuthStrategyFirstFactor,
		FirstFactor: &model.FirstFactorStrategy{
			Type: model.FirstFactorTypeLocal,
			Local: &model.LocalStrategy{
				Identity:  model.AuthIdentityTypePhone,
				Challenge: model.AuthChallengeTypeMagicLink,
			},
		},
	}
	al := []model.AuthStrategy{a0, a1, a2, a3, a4}

	f := a.FilterCompatible(al)
	assert.Len(t, f, 2)
	assert.Contains(t, f, a0)
	assert.Contains(t, f, a4)

	a.FirstFactor.Local = &model.LocalStrategy{
		Identity: model.AuthIdentityTypePhone,
	}
	f = a.FilterCompatible(al)
	assert.Len(t, f, 1)
	assert.Contains(t, f, a4)
}
