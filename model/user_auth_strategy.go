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
