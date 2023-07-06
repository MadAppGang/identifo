package model

import (
	_ "embed"
	"encoding/json"
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
type AuthStrategy struct {
	ID           string                `json:"id,omitempty"`
	Name         string                `json:"name,omitempty"`
	Type         AuthStrategyType      `json:"type,omitempty"`
	FirstFactor  *FirstFactorStrategy  `json:"first_factor,omitempty"`
	SecondFactor *SecondFactorStrategy `json:"second_factor,omitempty"`
	Score        int                   `json:"score,omitempty"` // the security score of the strategy, from 1 to 10.
}

// AuthStrategyType  show which type of auth this strategy represents
type AuthStrategyType string

const (
	AuthStrategyFirstFactor  AuthStrategyType = "first_factor"
	AuthStrategySecondFactor AuthStrategyType = "second_factor"
	AuthStrategyNone         AuthStrategyType = "none"
	AuthStrategyAnonymous    AuthStrategyType = "anonymous"
)

// FirstFactorType - the type of first factor authentication. Now local is supported only.
type FirstFactorType string

const (
	FirstFactorTypeLocal      FirstFactorType = "local"
	FirstFactorTypeFIM        FirstFactorType = "fim" // FIM if federated identity, also called as social login
	FirstFactorTypeEnterprise FirstFactorType = "enterprise"
)

// FirstFactorStrategy - the strategy for the first factor authentication.
type FirstFactorStrategy struct {
	Type  FirstFactorType `json:"type,omitempty"`
	Local *LocalStrategy  `json:"local,omitempty"` // User Identity is managed directly by Identifo
	FIM   *FIMStrategy    `json:"fim,omitempty"`   // FIM if federated identity, also called as social login
	// TODO:
	// SSO: SAML // user identity managed by someone else, delegated
}

// LocalStrategy - local first factor strategy, the user can authenticate with local account.
type LocalStrategy struct {
	Identity  AuthIdentityType  `json:"identity,omitempty"`
	Challenge AuthChallengeType `json:"challenge,omitempty"`
	Transport AuthTransportType `json:"transport,omitempty"`
}

type FIMStrategy struct {
	Type FIMStrategyType `json:"type,omitempty"`
	// TODO: add fields for OIDC and plugin strategy
}

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
	AuthTransportTypeEmail  AuthTransportType = "email"
	AuthTransportTypeSMS    AuthTransportType = "sms"
	AuthTransportTypePush   AuthTransportType = "push"
	AuthTransportTypeNone   AuthTransportType = "none"
	AuthTransportTypeSocket AuthTransportType = "socket"

	// FIMStrategyType - the FIM (Federated Identity Management) provider strategy name.
	FIMStrategyTypeNone     FIMStrategyType = "none"
	FIMStrategyTypeOIDC     FIMStrategyType = "oidc"
	FIMStrategyTypeApple    FIMStrategyType = "apple"
	FIMStrategyTypeFirebase FIMStrategyType = "firebase"
	FIMStrategyTypeGoogle   FIMStrategyType = "google"
)

// SecondFactorStrategy - the strategy for the second factor authentication.
type SecondFactorStrategy struct {
	Challenge   AuthChallengeType       `json:"challenge,omitempty"`
	Transport   AuthTransportType       `json:"transport,omitempty"`
	EnrolPolicy SecondFactorEnrolPolicy `json:"enrol_policy,omitempty"`
	Policy      SecondFactorPolicy      `json:"policy,omitempty"`
}

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
func (s AuthStrategy) FilterCompatible(strats []AuthStrategy) []AuthStrategy {
	res := []AuthStrategy{}
	for _, strat := range strats {
		if s.Compatible(strat) {
			res = append(res, strat)
		}
	}
	return res
}

// Compatible return true if other strategy is compatible with s.
func (s AuthStrategy) Compatible(other AuthStrategy) bool {
	// should be the same strategy type
	if s.Type == other.Type {
		// if the strategy is not none or anonymous, it is compatible, no other values to check.
		if s.Type == AuthStrategyNone || s.Type == AuthStrategyAnonymous {
			return true
		}
		// let's check the first factor
		if s.Type == AuthStrategyFirstFactor {
			// no other requirements, we are looking for the first factor strategy only
			if s.FirstFactor == nil {
				return true
			}
			// the strategy from the list is incomplete, it has not first factor strategy data
			if other.FirstFactor == nil {
				return false
			}
			return s.FirstFactor.Compatible(*other.FirstFactor)
		}
		if s.Type == AuthStrategySecondFactor {
			if s.SecondFactor == nil {
				return true
			}
			if other.SecondFactor == nil {
				return false
			}
			return s.SecondFactor.Compatible(*other.SecondFactor)
		}
	}
	return false
}

// Compatible returns true other strategy is compatible with s
func (s FirstFactorStrategy) Compatible(other FirstFactorStrategy) bool {
	if s.Type == other.Type {
		// check local strategy
		if s.Type == FirstFactorTypeLocal {
			// we are happy for any local strategy
			if s.Local == nil {
				return true
			}
			// s has requirements, but other strategy has no local strategy data
			if other.Local == nil {
				return false
			}
			return s.Local.Compatible(*other.Local)
		}
		if s.Type == FirstFactorTypeFIM {
			// we are happy for any federated identity management(FIM) strategy
			if s.FIM == nil {
				return true
			}
			if other.FIM == nil {
				return false
			}
			return s.FIM.Compatible(*other.FIM)
		}
	}
	return false
}

// Compatible returns true other strategy is compatible with s
func (s LocalStrategy) Compatible(other LocalStrategy) bool {
	if (len(s.Identity) == 0 || s.Identity == other.Identity) &&
		(len(s.Challenge) == 0 || s.Challenge == other.Challenge) &&
		(len(s.Transport) == 0 || s.Transport == other.Transport) {
		return true
	}
	return false
}

// Compatible returns true other strategy is compatible with s
func (s FIMStrategy) Compatible(other FIMStrategy) bool {
	return s.Type == other.Type || len(s.Type) == 0
}

// Compatible returns true other strategy is compatible with s
func (s SecondFactorStrategy) Compatible(other SecondFactorStrategy) bool {
	if (len(s.Challenge) == 0 || s.Challenge == other.Challenge) &&
		(len(s.Transport) == 0 || s.Transport == other.Transport) {
		return true
	}
	return false
}


func (a AuthStrategy)LoginType