package api

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/authorization"
	"github.com/madappgang/identifo/v2/web/middleware"
)

const (
	phoneVerificationCodeLength = 6
	smsVerificationCode         = "%v is your SMS verification code!"
)

// RequestVerificationCode requests SMS with verification code.
// To authenticate, user must have a valid phone number.
// check the app is supported for phone login
// check there are user with this phone number
// if no user just return "ok" for security purposes
// if yes, create challenge
// save challenge to database
// send SMS with challenge
// return ok

func (ar *Router) RequestVerificationCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		if !ar.SupportedLoginWays.Phone {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.APIAPPPhoneLoginNotSupported)
			return
		}

		var authData PhoneLogin
		if err := json.NewDecoder(r.Body).Decode(&authData); err != nil {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyInvalidError, err)
			return
		}

		if err := authData.validatePhone(); err != nil {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyInvalidError, err)
			return
		}

		// TODO: add limiter here. Check frequency of requests.
		code := randStringBytes(phoneVerificationCodeLength)
		if err := ar.server.Storages().Verification.CreateVerificationCode(authData.PhoneNumber, code); err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorStorageVerificationCreateError, err)
			return
		}

		if err := ar.server.Services().SMS.SendSMS(authData.PhoneNumber, fmt.Sprintf(smsVerificationCode, code)); err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorServiceSmsSendError, err)
			return
		}
		result := map[string]string{"result": "ok", "message": "SMS code is sent"}
		ar.ServeJSON(w, locale, http.StatusOK, result)
	}
}

// PhoneLogin authenticates user with phone number and verification code.
// If user exists - create new session and return token.
// If user does not exist - register and then login (create session and return token).
// If code is invalid - return error.
func (ar *Router) PhoneLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		var authData PhoneLogin
		if err := json.NewDecoder(r.Body).Decode(&authData); err != nil {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyInvalidError, err)
			return
		}
		if err := authData.validateCodeAndPhone(); err != nil {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyInvalidError, err)
			return
		}

		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPIAPPNoAPPInContext)
			return
		}

		needVerification := app.DebugTFACode == "" || authData.Code != app.DebugTFACode
		if needVerification { // check verification code
			if exists, err := ar.server.Storages().Verification.IsVerificationCodeFound(authData.PhoneNumber, authData.Code); err != nil {
				ar.Error(w, locale, http.StatusInternalServerError, l.ErrorStorageVerificationFindError, err)
				return
			} else if !exists {
				ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAPILoginCodeInvalid)
				return
			}
		}

		user, err := ar.server.Storages().User.UserByPhone(authData.PhoneNumber)
		if err == l.ErrorUserNotFound {
			if !ar.server.Settings().Login.AllowRegisterMissing {
				ar.LocalizedError(w, locale, http.StatusUnauthorized, l.ErrorAPIAPPRegistrationForbidden)
				return
			}

			// Generate random password for feature reset if needed
			user, err = ar.server.Storages().User.AddUserWithPassword(
				model.User{Phone: authData.PhoneNumber, Scopes: model.SliceIntersect(app.Scopes, authData.Scopes)},
				model.RandomPassword(15),
				app.NewUserDefaultRole,
				false)
		}
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorStorageFindUserPhoneError, err)
			return
		}

		// Authorize user if the app requires authorization.
		azi := authorization.AuthzInfo{
			App:         app,
			UserRole:    user.AccessRole,
			ResourceURI: r.RequestURI,
			Method:      r.Method,
		}
		if err := ar.Authorizer.Authorize(azi); err != nil {
			ar.LocalizedError(w, locale, http.StatusForbidden, l.APIAccessDenied)
			return
		}

		// if app requires scope, we need to check user has at leas one scope
		if len(app.Scopes) > 0 && len(model.SliceIntersect(app.Scopes, user.Scopes)) == 0 {
			ar.LocalizedError(w, locale, http.StatusForbidden, l.ErrorAPPLoginNoScope)
			return
		}

		// Do login flow.
		scopes := []string{}
		// if we requested any scope, let's provide all the scopes user has and requested
		if len(authData.Scopes) > 0 {
			scopes = model.SliceIntersect(authData.Scopes, user.Scopes)
		}
		if model.SliceContains(authData.Scopes, "offline") && app.Offline {
			scopes = append(scopes, "offline")
		}

		tokenPayload, err := ar.getTokenPayloadForApp(app, user)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAPIAPPUnableToTokenPayloadForAPPError, app.ID, err)
			return
		}

		offline := contains(scopes, model.OfflineScope)
		accessToken, refreshToken, err := ar.loginUser(user, scopes, app, offline, false, tokenPayload)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAPILoginError, err)
			return
		}

		user = user.Sanitized()
		result := AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			User:         user,
		}

		ar.server.Storages().User.UpdateLoginMetadata(user.ID)
		ar.ServeJSON(w, locale, http.StatusOK, result)
	}
}

// Generate user code
func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		code := "1234567890"
		if k, err := rand.Int(rand.Reader, big.NewInt(int64(len(code)))); err != nil {
			panic(err)
		} else {
			b[i] = code[int(k.Int64())]
		}
	}
	return string(b)
}

// PhoneLogin is used to parse input data from the client during phone login.
type PhoneLogin struct {
	PhoneNumber string   `json:"phone_number"`
	Code        string   `json:"code"`
	Scopes      []string `json:"scopes"`
}

func (l *PhoneLogin) validateCodeAndPhone() error {
	if len(l.Code) == 0 {
		return errors.New("Verification code is too short or missing. ")
	}
	return l.validatePhone()
}

func (l *PhoneLogin) validatePhone() error {
	if !model.PhoneRegexp.MatchString(l.PhoneNumber) {
		return errors.New("Phone number is not valid. ")
	}
	return nil
}
