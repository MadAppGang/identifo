package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
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
		app := middleware.AppFromContext(r.Context())

		d := challengeRequest{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		agent := r.Header.Get("User-Agent")

		idType := model.AuthIdentityTypePhone
		if len(d.Email) > 0 {
			idType = model.AuthIdentityTypeEmail
		}

		// create uncompleted UserAuthChallenge to challenge controller to create a full challenge
		// if possible
		ch := model.UserAuthChallenge{
			AppID:             app.ID,
			DeviceID:          d.Device,
			UserAgent:         agent,
			CreatedAt:         time.Now(),
			UserCodeChallenge: d.ClientCodeChallenge,
			Strategy: model.AuthStrategy{
				Type: model.AuthStrategyFirstFactor,
				FirstFactor: &model.FirstFactorStrategy{
					Type: model.FirstFactorTypeLocal,
					Local: &model.LocalStrategy{
						Identity:  idType,
						Challenge: model.AuthChallengeType(d.ChallengeType),
					},
				},
			},
		}

		_, err := ar.server.Services().Challenge.RequestChallenge(r.Context(), ch)
		if err != nil && !errors.Is(err, l.ErrorUserNotFound) {
			ar.Error(w, l.ErrorWithLocale(err, locale))
			return
		}

		ar.ServeJSON(w, locale, http.StatusOK, nil)
	}
}

// Solve challenge could be called for requested challenge (for first factor)
// or for user enforced dialog (second factor)
// func (ar *Router) SolveChallenge() http.HandlerFunc {
// 	type challengeRequest struct {
// 		PhoneNumber string   `json:"phone_number"`
// 		Code        string   `json:"code"`
// 		Scopes      []string `json:"scopes"`
// 	}

// 	return func(w http.ResponseWriter, r *http.Request) {
// 		locale := r.Header.Get("Accept-Language")

// 		// check the app is supported for auth challenge requested
// 		// check there are user with this auth challenge factor
// 		// if no user just return "ok" for security purposes
// 		// if yes, create challenge
// 		// save challenge to database
// 		// send SMS/Email or whatever for the challenge
// 		// return ok
// 	}
// }
