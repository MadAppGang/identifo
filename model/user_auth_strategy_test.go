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
	s1 := model.FirstFactorFIMStrategy{FIMType: model.UserIdentityTypeApple}
	s2 := model.FirstFactorFIMStrategy{FIMType: model.UserIdentityTypeApple}

	assert.True(t, s1.Compatible(s2))

	s1.FIMType = ""
	assert.False(t, s1.Compatible(s2))

	s1 = model.FirstFactorFIMStrategy{FIMType: model.UserIdentityTypeApple}
	s2.FIMType = model.UserIdentityTypeOIDC
	assert.False(t, s1.Compatible(s2))

	s2.FIMType = ""
	assert.False(t, s1.Compatible(s2))
}

func TestLocalStrategyCompatible(t *testing.T) {
	s1 := model.FirstFactorInternalStrategy{
		Identity: model.AuthIdentityTypePhone,
	}
	s2 := model.FirstFactorInternalStrategy{
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
	s1 := model.FirstFactorInternalStrategy{}
	s2 := model.FirstFactorInternalStrategy{
		Identity:  model.AuthIdentityTypeEmail,
		Challenge: model.AuthChallengeTypeOTP,
		Transport: model.AuthTransportTypePush,
	}
	assert.True(t, s1.Compatible(s2))

	s1.Identity = model.AuthIdentityTypeAnonymous
	assert.False(t, s1.Compatible(s2))

	s1.Identity = ""
	assert.True(t, s1.Compatible(s2))

	s1.Identity = model.AuthIdentityTypeEmail
	assert.True(t, s1.Compatible(s2))
}

func TestFilterCompatible(t *testing.T) {
	a := model.FirstFactorInternalStrategy{}
	a0 := model.FirstFactorInternalStrategy{
		Identity:  model.AuthIdentityTypeEmail,
		Challenge: model.AuthChallengeTypeOTP,
		Transport: model.AuthTransportTypePush,
	}
	a1 := model.AnonymousStrategy{}
	a2 := model.FirstFactorEnterpriseStrategy{}
	a3 := model.FirstFactorFIMStrategy{
		FIMType: model.UserIdentityTypeApple,
	}
	a4 := model.FirstFactorInternalStrategy{
		Identity:  model.AuthIdentityTypePhone,
		Challenge: model.AuthChallengeTypeMagicLink,
	}
	al := []model.AuthStrategy{a0, a1, a2, a3, a4}

	f := model.FilterCompatible(a, al)
	assert.Len(t, f, 2)
	assert.Contains(t, f, a0)
	assert.Contains(t, f, a4)

	a.Identity = model.AuthIdentityTypePhone
	f = model.FilterCompatible(a, al)
	assert.Len(t, f, 1)
	assert.Contains(t, f, a4)
}
