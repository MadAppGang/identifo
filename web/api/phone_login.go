package api

import (
	"encoding/json"
	"net/http"

	"github.com/madappgang/identifo/v2/l"
)

// Passwordless login, login with code or magic link on email or SMS code.
func (ar *Router) RequestPasswordlessLoginCode() http.HandlerFunc {
	// now RequestChallenge supports only login challenge request.
	// in future it could be extended to 2FA
	return ar.RequestChallenge()
}

// PasswordlessLogin - handles login with code or magic link on email or SMS code.
// If user exists - create new session and return token.
// If user exists and has debug OTP code and app allows debug OTP code and the code provided matched that code - login or register.
// If user does not exist and app allows register passwordless users - register and then login (create session and return token).
// If code is invalid - return error.
func (ar *Router) PasswordlessLogin() http.HandlerFunc {
	type passwordlessLoginData struct {
		Phone  string `json:"phone"`
		Email  string `json:"email"`
		OTP    string `json:"otp"`
		Device string `json:"device"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		var authData passwordlessLoginData
		if err := json.NewDecoder(r.Body).Decode(&authData); err != nil {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyInvalidError, err)
			return
		}

		ar.server.Storages().UCC.ChangeBlockStatus

		ar.server.Storages().User.UpdateLoginMetadata(user.ID)
		// um := model.User{}
		// model.CopyDstFields(rd, um)
		// user, err := ar.server.Storages().UMC.CreateUserWithPassword(r.Context(), um, rd.Password)
		// if err != nil { // this error is already localized.
		// 	ar.Error(w, err)
		// 	return
		// }

		// user = model.CopyFields(user, model.UserFieldsetBasic.Fields())
		// ar.ServeJSON(w, locale, http.StatusOK, user)
	}
}
