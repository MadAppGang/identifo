package api

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"regexp"

	jwtService "github.com/madappgang/identifo/jwt/service"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/middleware"
)

var (
	phoneRegExp = regexp.MustCompile(`^[\+][0-9]{9,15}$`)
)

const (
	phoneVerificationCodeLength = 6
	smsVerificationCode         = "%v is your SMS verification code!"
)

// RequestVerificationCode requests SMS with verification code.
// To authenticate, user must have a valid phone number.
func (ar *Router) RequestVerificationCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var authData PhoneLogin
		if err := json.NewDecoder(r.Body).Decode(&authData); err != nil {
			ar.Error(w, ErrorAPIRequestBodyInvalid, http.StatusBadRequest, err.Error(), "RequestVerificationCode.Unmarshal")
			return
		}

		if err := authData.validatePhone(); err != nil {
			ar.Error(w, ErrorAPIRequestBodyParamsInvalid, http.StatusBadRequest, err.Error(), "RequestVerificationCode.IsValidPhone")
			return
		}

		// TODO: add limiter here. Check frequency of requests
		code := randStringBytes(phoneVerificationCodeLength)
		err := ar.verificationCodeStorage.CreateVerificationCode(authData.PhoneNumber, code)
		if err != nil {
			ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "RequestVerificationCode.PutVerificationCode")
			return
		}

		err = ar.smsService.SendSMS(authData.PhoneNumber, fmt.Sprintf(smsVerificationCode, code))
		if err != nil {
			ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, "Unable to send sms. Error: "+err.Error(), "RequestVerificationCode.SendSMS")
			return
		}
		ar.ServeJSON(w, http.StatusOK, map[string]string{"message": "SMS code is sent"})
	}
}

// PhoneLogin authenticates user with phone number and verification code.
// If user exists - create new session and return token.
// If user does not exists - register and then login (create session and return token).
// If code is invalid - return error.
func (ar *Router) PhoneLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var authData PhoneLogin
		if err := json.NewDecoder(r.Body).Decode(&authData); err != nil {
			ar.Error(w, ErrorAPIRequestBodyInvalid, http.StatusBadRequest, err.Error(), "PhoneLogin.Unmarshal")
			return
		}
		if err := authData.validateCodeAndPhone(); err != nil {
			ar.Error(w, ErrorAPIRequestBodyParamsInvalid, http.StatusBadRequest, err.Error(), "PhoneLogin.IsValidCodeAndPhone")
			return
		}

		// check verification code
		if exists, err := ar.verificationCodeStorage.IsVerificationCodeFound(authData.PhoneNumber, authData.Code); err != nil {
			ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "PhoneLogin.IsVerificationCodeFound.error")
			return
		} else if !exists {
			ar.Error(w, ErrorAPIVerificationCodeInvalid, http.StatusUnauthorized, "Invalid phone or verification code", "PhoneLogin.IsVerificationCodeFound.not_exists")
			return
		}

		user, err := ar.userStorage.UserByPhone(authData.PhoneNumber)
		if err == model.ErrUserNotFound {
			user, err = ar.userStorage.AddUserByPhone(authData.PhoneNumber)
		}
		if err != nil {
			ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "PhoneLogin.UserByPhone")
			return
		}

		scopes, err := ar.userStorage.RequestScopes(user.ID(), authData.Scopes)
		if err != nil {
			ar.Error(w, ErrorAPIRequestScopesForbidden, http.StatusForbidden, err.Error(), "LoginWithPassword.RequestScopes")
			return
		}
		app := middleware.AppFromContext(r.Context())
		if app == nil {
			ar.logger.Println("Error getting App")
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App is not in context.", "LoginWithPassword.AppFromContext")
			return
		}

		offline := contains(scopes, jwtService.OfflineScope)
		accessToken, refreshToken, err := ar.loginUser(user, scopes, app, offline)
		if err != nil {
			ar.Error(w, ErrorAPIAppAccessTokenNotCreated, http.StatusInternalServerError, err.Error(), "LoginWithPassword.loginUser")
			return
		}
		result := AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		ar.ServeJSON(w, http.StatusOK, result)
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
	if !phoneRegExp.MatchString(l.PhoneNumber) {
		return errors.New("Phone number is not valid. ")
	}
	return nil
}
