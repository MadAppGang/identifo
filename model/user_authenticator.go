package model

// UserAuthenticator is a way how user can authenticate itself
type UserAuthenticator struct {
	Strategy AuthStrategy
}

type AuthStrategy string

// ID: username, email, phone, anonymous
// What: password, magic link, guardian, otp
// How: email, sms, push, oob

const (
	AuthStrategyUsernameAndPassword AuthStrategy = "username_password"
	AuthStrategyEmailAndPassword    AuthStrategy = "email_password"
	AuthStrategyEmailAndMagicLink   AuthStrategy = "email_link"
	AuthStrategyEmailAndGuardian    AuthStrategy = "email_guardian"
	AuthStrategyEmailAndOPT         AuthStrategy = "email_otp"
	AuthStrategyPhoneAndPassword    AuthStrategy = "phone_password"
	AuthStrategyPhoneAndMagicLink   AuthStrategy = "phone_link"
	AuthStrategyPhoneAndGuardian    AuthStrategy = "phone_guardian"
	AuthStrategyPhoneAndOTP         AuthStrategy = "phone_otp"
	AuthStrategyPushOTP             AuthStrategy = "push_otp"
	AuthStrategyPushMagicLink       AuthStrategy = "push_link"
	AuthStrategyAnonymous           AuthStrategy = "anonymous"
	AuthStrategyFIDO                AuthStrategy = "fido"
)
