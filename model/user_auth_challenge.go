package model

import "time"

// UserAuthChallenge is a list of Auth challenges user have to solve.
// Specific for an auth strategy.
// For example we use SMS code to login with phone number.
// For this reason we are using the following strategy:
// Type: AuthStrategyFirstFactor
// FirstFactor: {
//		Type: local
//		Local: {
//			Identity: AuthIdentityTypePhone,
//			Challenge: AuthChallengeTypeOTP,
//			Transport: AuthTransportTypeSMS
// 		}
// }
// The expire will be 10 mins for example,
// and when user enters OTP code in web app we call API:
// AuthenticateWithChallenge(OTP: 1234)
// We look for the latest Challenge with specific strategy and
// validate the expire date
// issue the JWT tokens and mark the strategy as Solved

type UserAuthChallenge struct {
	ID        string
	UserID    string
	DeviceID  string
	Strategy  AuthStrategy
	Solved    bool // is the challenge already solved, could not be solved again. One time challenge.
	CreatedAt time.Time
	ExpiresAt time.Time
}
