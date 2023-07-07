package model

import (
	_ "embed"
	"encoding/json"
	"time"
)

// TODO: Jack add all missing strategies to json file
//
//go:embed user_auth_strategies.json
var defaultStrategiesBuf []byte

func Strategies() []AuthStrategy {
	var res []AuthStrategy
	json.Unmarshal(defaultStrategiesBuf, &res)
	return res
}

// AuthStrategy - a auth strategy to auth the user
type AuthStrategy interface {
	Info() AuthStrategyInfo
	Type() AuthStrategyType
	Compatible(other AuthStrategy) bool
}

type AuthStrategyInfo struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Score int    `json:"score,omitempty"` // the security score of the strategy, from 1 to 10.
	// Type  AuthStrategyType `json:"type,omitempty"`
}

// FirstFactorInternalStrategy represents first factor internal strategy, like email, phone, etc
type FirstFactorInternalStrategy struct {
	AuthStrategyInfo
	Identity  AuthIdentityType  `json:"identity,omitempty"`
	Challenge AuthChallengeType `json:"challenge,omitempty"`
	Transport AuthTransportType `json:"transport,omitempty"`
}

type FirstFactorFIMStrategy struct {
	AuthStrategyInfo
	FIMType FIMStrategyType `json:"fim_type,omitempty"`
	// TODO: add fields for OIDC and plugin strategy
}

type FirstFactorEnterpriseStrategy struct {
	AuthStrategyInfo
	// TODO: implement support for enterprise strategy
}

type SecondFactorStrategy struct {
	AuthStrategyInfo
	Challenge   AuthChallengeType       `json:"challenge,omitempty"`
	Transport   AuthTransportType       `json:"transport,omitempty"`
	EnrolPolicy SecondFactorEnrolPolicy `json:"enrol_policy,omitempty"`
	Policy      SecondFactorPolicy      `json:"policy,omitempty"`
}

type AnonymousStrategy struct {
	AuthStrategyInfo
}

// FirstFactorInternalStrategy implementation of AuthStrategy interface
func (s FirstFactorInternalStrategy) Info() AuthStrategyInfo {
	return s.AuthStrategyInfo
}

func (s FirstFactorInternalStrategy) Type() AuthStrategyType {
	return AuthStrategyFirstFactorInternal
}

func (s FirstFactorInternalStrategy) Compatible(other AuthStrategy) bool {
	if other.Type() == AuthStrategyFirstFactorInternal {
		ol, ok := other.(FirstFactorInternalStrategy)
		if !ok {
			return false
		}
		if (len(s.Identity) == 0 || s.Identity == ol.Identity) &&
			(len(s.Challenge) == 0 || s.Challenge == ol.Challenge) &&
			(len(s.Transport) == 0 || s.Transport == ol.Transport) {
			return true
		}
	}
	return false
}

// FirstFactorFIMStrategy implementation of AuthStrategy interface
func (s FirstFactorFIMStrategy) Info() AuthStrategyInfo {
	return s.AuthStrategyInfo
}

func (s FirstFactorFIMStrategy) Type() AuthStrategyType {
	return AuthStrategyFirstFactorFIM
}

func (s FirstFactorFIMStrategy) Compatible(other AuthStrategy) bool {
	if other.Type() == AuthStrategyFirstFactorFIM {
		ol, ok := other.(FirstFactorFIMStrategy)
		if !ok {
			return false
		}
		return s.FIMType == ol.FIMType
	}
	return false
}

// FirstFactorFIMStrategy implementation of AuthStrategy interface
func (s FirstFactorEnterpriseStrategy) Info() AuthStrategyInfo {
	return s.AuthStrategyInfo
}

func (s FirstFactorEnterpriseStrategy) Type() AuthStrategyType {
	return AuthStrategyFirstFactorEnterprise
}

func (s FirstFactorEnterpriseStrategy) Compatible(other AuthStrategy) bool {
	return false
}

// SecondFactorStrategy implementation of AuthStrategy interface
func (s SecondFactorStrategy) Info() AuthStrategyInfo {
	return s.AuthStrategyInfo
}

func (s SecondFactorStrategy) Type() AuthStrategyType {
	return AuthStrategySecondFactor
}

func (s SecondFactorStrategy) Compatible(other AuthStrategy) bool {
	if other.Type() == AuthStrategySecondFactor {
		ol, ok := other.(SecondFactorStrategy)
		if !ok {
			return false
		}
		if (len(s.Challenge) == 0 || s.Challenge == ol.Challenge) &&
			(len(s.Transport) == 0 || s.Transport == ol.Transport) {
			return true
		}
	}
	return false
}

// AnonymousStrategy implementation of AuthStrategy interface
func (s AnonymousStrategy) Info() AuthStrategyInfo {
	return s.AuthStrategyInfo
}

func (s AnonymousStrategy) Type() AuthStrategyType {
	return AuthStrategyAnonymous
}

func (s AnonymousStrategy) Compatible(other AuthStrategy) bool {
	return false
}

// AuthStrategyType  show which type of auth this strategy represents
type AuthStrategyType string

const (
	AuthStrategyFirstFactorInternal   AuthStrategyType = "first_factor_internal"
	AuthStrategyFirstFactorFIM        AuthStrategyType = "first_factor_fim"
	AuthStrategyFirstFactorEnterprise AuthStrategyType = "first_factor_enterprise"
	AuthStrategySecondFactor          AuthStrategyType = "second_factor"
	AuthStrategyNone                  AuthStrategyType = "none"
	AuthStrategyAnonymous             AuthStrategyType = "anonymous"
)

// types for Identity, Challenge and Transport
type (
	AuthIdentityType  string
	AuthChallengeType string
	AuthTransportType string
	FIMStrategyType   string
)

const (
	// AuthIdentityType - what kind of user identity we are using to authenticate a user.
	AuthIdentityTypeID        AuthIdentityType = "id"
	AuthIdentityTypeEmail     AuthIdentityType = "email"
	AuthIdentityTypePhone     AuthIdentityType = "phone"
	AuthIdentityTypeUsername  AuthIdentityType = "username"
	AuthIdentityTypeAnonymous AuthIdentityType = "anonymous"

	// AuthChallengeType - the challenge type we are using to auth the sure.
	AuthChallengeTypePassword      AuthChallengeType = "password"
	AuthChallengeTypeOTP           AuthChallengeType = "otp"
	AuthChallengeTypeMagicLink     AuthChallengeType = "magic_link"
	AuthChallengeTypeNone          AuthChallengeType = "none"
	AuthChallengeTypeRecoveryCodes AuthChallengeType = "recovery_codes"
	AuthChallengeTypeGuardian      AuthChallengeType = "guardian"
	AuthChallengeTypeWebauthn      AuthChallengeType = "webauthn"

	// AuthTransportType - the transport we are using to deliver the authentication challenge.
	AuthTransportTypeEmail         AuthTransportType = "email"
	AuthTransportTypeSMS           AuthTransportType = "sms"
	AuthTransportTypePush          AuthTransportType = "push"
	AuthTransportTypeNone          AuthTransportType = "none"
	AuthTransportTypeSocket        AuthTransportType = "socket"
	AuthTransportTypeAuthenticator AuthTransportType = "authenticator"

	// FIMStrategyType - the FIM (Federated Identity Management) provider strategy name.
	FIMStrategyTypeNone     FIMStrategyType = "none"
	FIMStrategyTypeOIDC     FIMStrategyType = "oidc"
	FIMStrategyTypeApple    FIMStrategyType = "apple"
	FIMStrategyTypeFirebase FIMStrategyType = "firebase"
	FIMStrategyTypeGoogle   FIMStrategyType = "google"
)

// SecondFactorEnrolPolicy - the policy for the second factor enrolment, the way is user can enrol second factor.
type SecondFactorEnrolPolicy string

const (
	// SecondFactorEnrolPolicyAlways - the user can never enrol second factor.
	SecondFactorEnrolPolicyNever SecondFactorEnrolPolicy = "never"
	// SecondFactorEnrolPolicyAlways - the user can always enrol second factor by himself.
	SecondFactorEnrolPolicySelfEnrol SecondFactorEnrolPolicy = "self_enrol"
	// SecondFactorEnrolPolicyAlways - the user can be enrolled by admin or during sign-up process.
	SecondFactorEnrolPolicyDeny SecondFactorEnrolPolicy = "deny"
)

// SecondFactorPolicy - the policy for the second factor authentication, when are challenge for second factor from the user.
type SecondFactorPolicy string

const (
	// SecondFactorPolicyNever - we never ask for second factor.
	SecondFactorPolicyNever SecondFactorPolicy = "never"
	// SecondFactorPolicyNaive - we ask for second factor only if client app requests for that.
	SecondFactorPolicyNaive SecondFactorPolicy = "naive"
	// SecondFactorPolicyAlways - we always ask for second factor.
	SecondFactorPolicyAlways SecondFactorPolicy = "always"
	// SecondFactorPolicyAdaptive - we ask for second factor only when user is trying to login from risky sources. MFA not required by default.
	SecondFactorPolicyAdaptive SecondFactorPolicy = "adaptive"
	// SecondFactorPolicyAdaptive - we ask for second factor only when user is trying to login from risky sources. MFA required in this case. - we ask for second factor only when user is trying to login from untrusted sources. MFA not required by default.
	SecondFactorPolicyAdaptiveMFARequired SecondFactorPolicy = "adaptive_mfa_required"
	// SecondFactorPolicyCustom - we delegate the policy to a plugin.
	SecondFactorPolicyCustom SecondFactorPolicy = "custom"
)

func (a AuthIdentityType) Field() string {
	switch a {
	case AuthIdentityTypeID:
		return UserFieldID
	case AuthIdentityTypeEmail:
		return UserFieldEmail
	case AuthIdentityTypeUsername:
		return UserFieldUsername
	case AuthIdentityTypePhone:
		return UserFieldPhone
	default:
		return ""
	}
}

// FilterCompatible returns all compatible strategies from stats slice.
func FilterCompatible(s AuthStrategy, strats []AuthStrategy) []AuthStrategy {
	res := []AuthStrategy{}
	for _, strat := range strats {
		if s.Compatible(strat) {
			res = append(res, strat)
		}
	}
	return res
}

// ExpireChallengeDuration returns a duration for the challenge to expire for this specific strategy type.
// If this type does not support challenge expiration, it returns 0.
func ExpireChallengeDuration(s AuthStrategy) time.Duration {
	switch s.Type() {
	case AuthStrategyFirstFactorInternal:
		f, ok := s.(FirstFactorInternalStrategy)
		if ok {
			if f.Transport == AuthTransportTypeAuthenticator {
				return time.Minute * 1
			}
			if f.Challenge == AuthChallengeTypeOTP {
				return time.Minute * 5
			}
			if f.Challenge == AuthChallengeTypeMagicLink {
				return time.Minute * 30
			}
		}
	}
	return time.Hour * 24 * 7 * 4 * 12 * 10 // something around ten years, never expire
}
