package model

import (
	"time"

	"github.com/madappgang/identifo/v2/l"
)

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
	ID                string       `json:"id"`
	UserID            string       `json:"user_id"`
	DeviceID          string       `json:"device_id"`
	AppID             string       `json:"app_id"`
	UserAgent         string       `json:"user_agent"`
	UserCodeChallenge string       `json:"user_code_challenge"`
	Strategy          AuthStrategy `json:"strategy"`
	Solved            bool         `json:"solved"` // is the challenge already solved, could not be solved again. One time challenge.
	CreatedAt         time.Time    `json:"created_at"`
	ExpiresAt         time.Time    `json:"expires_at"`
	SolvedAt          time.Time    `json:"solved_at"`
	SolvedUserAgent   string       `json:"solved_user_agent"`
	SolvedDeviceID    string       `json:"solved_device_id"`
	ExpiresMins       int          `json:"expires_mins"`
	ScopesRequested   []string     `json:"scopes_requested"`
	OTP               string       `json:"value"` // OTP value, it the challenge is OTP code
}

func (u UserAuthChallenge) IsExpired() bool {
	return u.ExpiresAt.Before(time.Now())
}

func (u UserAuthChallenge) Valid() error {
	if u.IsExpired() {
		return l.ErrorOtpExpired
	}
	if u.Solved {
		return l.ErrorOtpAlreadyUsed
	}
	return nil
}
