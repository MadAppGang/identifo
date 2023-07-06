package api

import (
	"net/http"
)

// check the app is supported for auth challenge requested
// check there are user with this auth challenge factor
// check the user is not blocked
// check the user is enrolled to this auth type
// check the user email verified if needed
// check the user phone number verified if needed
// if no user just return "ok" for security purposes
// if yes, create challenge
// save challenge to database
// send SMS/Email or whatever for the challenge
// return ok
// now we can request only for first factor:
// - phone number
// - email
// - passkey?(TODO)
// - guardian?(TODO)
// the auth type for those are:
// - OTP
// - magic link
// - passkey key?(TODO)
// - guardian response with guardian SDK?(TODO)
func (ar *Router) RequestChallenge() http.HandlerFunc {
	type challengeRequest struct {
		PhoneNumber         string   `json:"phone"` // email or phone should be filled
		Email               string   `json:"email"`
		ChallengeType       string   `json:"challenge_type"`   // preferable challenge type: otp or magic link for now
		ClientCodeChallenge string   `json:"client_challenge"` // random verification code from user
		Scopes              []string `json:"scopes"`           // requested scopes
		Device              string   `json:"device"`           // device info  from client
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		d := challengeRequest{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}
		// AuthChallengeType
		err := ar.server.Services().Challenge.RequestChallenge(r.Context())
	}
}

// Solve challenge could be called for requested challenge (for first factor)
// or for user enforced dialog (second factor)
func (ar *Router) SolveChallenge() http.HandlerFunc {
	type challengeRequest struct {
		PhoneNumber string   `json:"phone_number"`
		Code        string   `json:"code"`
		Scopes      []string `json:"scopes"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		// check the app is supported for auth challenge requested
		// check there are user with this auth challenge factor
		// if no user just return "ok" for security purposes
		// if yes, create challenge
		// save challenge to database
		// send SMS/Email or whatever for the challenge
		// return ok
	}
}
