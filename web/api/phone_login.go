package api

import (
	"encoding/json"
	"net/http"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
)

// Passwordless login, login with code or magic link on email or SMS code.
func (ar *Router) RequestPasswordlessLoginCode() http.HandlerFunc {
	// now RequestChallenge supports only login challenge request.
	// in future it could be extended to 2FA
	return ar.RequestChallenge()
}

// PasswordlessLogin - handles login with code or magic link on email or SMS code.
// verify challenge
// If user exists - create new session and return token.
// If user exists and has debug OTP code and app allows debug OTP code and the code provided matched that code - login or register.
// If user does not exist and app allows register passwordless users - register and then login (create session and return token).
// If code is invalid - return error.
func (ar *Router) PasswordlessLogin() http.HandlerFunc {
	type passwordlessLoginData struct {
		Phone  string   `json:"phone"`
		Email  string   `json:"email"`
		OTP    string   `json:"otp"`
		Device string   `json:"device"`
		Scopes []string `json:"scopes"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")
		agent := r.Header.Get("User-Agent")

		var d passwordlessLoginData
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyInvalidError, err)
			return
		}

		app := middleware.AppFromContext(r.Context())

		idType := model.AuthIdentityTypePhone
		idValue := d.Phone
		transport := model.AuthTransportTypeSMS
		if len(d.Email) > 0 {
			idType = model.AuthIdentityTypeEmail
			idValue = d.Email
			transport = model.AuthTransportTypeEmail
		}

		// create uncompleted UserAuthChallenge to verify challenge we have created
		challenge := model.UserAuthChallenge{
			AppID:     app.ID,
			DeviceID:  d.Device,
			UserAgent: agent,
			Strategy: model.FirstFactorInternalStrategy{
				Identity:  idType,
				Transport: transport,
			},
		}
		// let's check if the challenge has been solved
		u, app, challenge, err := ar.server.Storages().UCC.VerifyChallenge(r.Context(), challenge, idValue)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.LocalizedString(err.Error()))
		}

		// let's handle new user
		if model.IsNewUserID(u.ID) {
			u.ID = ""
			u, err = ar.server.Storages().UMC.CreateUser(r.Context(), u)
			if err != nil { // this error is already localized.
				ar.LocalizedError(w, locale, http.StatusInternalServerError, l.LocalizedString(err.Error()))
				return
			}
		}

		tokens, err := ar.server.Storages().UC.GetJWTTokens(r.Context(), app, u, d.Scopes)

		// if user send us challenge code while requesting the challenge from us - return it
		// trust is two way street 😉
		if len(challenge.UserCodeChallenge) > 0 {
			tokens.ClientChallenge = &challenge.UserCodeChallenge
		}

		ar.ServeJSON(w, locale, http.StatusOK, tokens)
	}
}
