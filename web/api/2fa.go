package api

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	jwtService "github.com/madappgang/identifo/jwt/service"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/plugin/shared"
	"github.com/madappgang/identifo/proto"
	"github.com/madappgang/identifo/proto/extensions"
	"github.com/madappgang/identifo/web/middleware"
	"github.com/xlzd/gotp"
)

// EnableTFA enables two-factor authentication for the user.
func (ar *Router) EnableTFA() http.HandlerFunc {
	type tfaSecret struct {
		TFASecret string `json:"tfa_secret"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		app := middleware.AppFromContext(r.Context())
		if app == nil {
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App is not in context.", "EnableTFA.AppFromContext")
			return
		}

		if tfaStatus := app.TFAStatus(); tfaStatus != model.TFAStatusOptional {
			ar.Error(w, ErrorAPIRequestBodyParamsInvalid, http.StatusBadRequest, fmt.Sprintf("App TFA status is '%s', not 'optional'", tfaStatus), "EnableTFA.TFAStatus")
			return
		}

		accessTokenBytes, ok := r.Context().Value(model.TokenRawContextKey).([]byte)
		if !ok {
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "Token bytes are not in context.", "EnableTFA.TokenBytesFromContext")
			return
		}

		// Get userID from token and update user with this ID.
		userID, err := ar.getTokenSubject(string(accessTokenBytes))
		if err != nil {
			ar.Error(w, ErrorAPIAppCannotExtractTokenSubject, http.StatusInternalServerError, err.Error(), "EnableTFA.getTokenSubject")
			return
		}

		user, err := ar.userStorage.UserByID(userID)
		if err != nil {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusBadRequest, err.Error(), "EnableTFA.UserByID")
			return
		}

		if tfaInfo := user.TfaInfo; tfaInfo.IsEnabled && tfaInfo.Secret != "" {
			ar.Error(w, ErrorAPIRequestTFAAlreadyEnabled, http.StatusBadRequest, "TFA already enabled for this user", "EnableTFA.alreadyEnabled")
			return
		}

		if ar.tfaType == model.TFATypeSMS && user.Phone == "" {
			ar.Error(w, ErrorAPIRequestPleaseSetPhoneForTFA, http.StatusBadRequest, "Please specify your phone number to be able to receive one-time passwords", "EnableTFA.setPhone")
			return
		}
		if ar.tfaType == model.TFATypeEmail && user.Email == "" {
			ar.Error(w, ErrorAPIRequestPleaseSetEmailForTFA, http.StatusBadRequest, "Please specify your email address to be able to receive one-time passwords", "EnableTFA.setEmail")
			return
		}

		tfa := proto.TFAInfo{
			IsEnabled: true,
			Secret:    gotp.RandomSecret(16),
		}
		user.TfaInfo = &tfa

		if _, err := ar.userStorage.UpdateUser(userID, user); err != nil {
			ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "EnableTFA.UpdateUser")
			return
		}

		switch ar.tfaType {
		case model.TFATypeApp:
			ar.ServeJSON(w, http.StatusOK, &tfaSecret{TFASecret: tfa.Secret})
			return
		case model.TFATypeSMS, model.TFATypeEmail:
			ar.ServeJSON(w, http.StatusOK, nil)
			return
		}
		ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, fmt.Sprintf("Unknown tfa type '%s'", ar.tfaType), "switch.tfaType")
	}
}

// FinalizeTFA finalizes two-factor authentication.
func (ar *Router) FinalizeTFA() http.HandlerFunc {
	type requestBody struct {
		TFACode string   `json:"tfa_code"`
		Scopes  []string `json:"scopes"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		d := requestBody{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		if len(d.TFACode) == 0 {
			ar.Error(w, ErrorAPIRequestTFACodeEmpty, http.StatusBadRequest, "", "FinalizeTFA.empty")
			return
		}

		oldAccessTokenBytes, ok := r.Context().Value(model.TokenRawContextKey).([]byte)
		if !ok {
			ar.Error(w, ErrorAPIRequestTokenInvalid, http.StatusBadRequest, "Token bytes are not in context.", "FinalizeTFA.TokenBytesFromContext")
			return
		}
		oldAccessTokenString := string(oldAccessTokenBytes)

		userID, err := ar.getTokenSubject(oldAccessTokenString)
		if err != nil {
			ar.Error(w, ErrorAPIAppCannotExtractTokenSubject, http.StatusInternalServerError, err.Error(), "FinalizeTFA.getTokenSubject")
			return
		}

		user, err := ar.userStorage.UserByID(userID)
		if err != nil {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusBadRequest, err.Error(), "FinalizeTFA.UserByID")
			return
		}

		app := middleware.AppFromContext(r.Context())
		if app == nil {
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App is not in context.", "EnableTFA.AppFromContext")
			return
		}

		totp := gotp.NewDefaultTOTP(user.TfaInfo.Secret)
		dontNeedVerification := app.DebugTFACode() != "" && d.TFACode == app.DebugTFACode()

		if verified := totp.Verify(d.TFACode, int(time.Now().Unix())); !(verified || dontNeedVerification) {
			ar.Error(w, ErrorAPIRequestTFACodeInvalid, http.StatusUnauthorized, "", "FinalizeTFA.TOTP_Invalid")
			return
		}

		// Issue new access, and, if requested, refresh token, and then invalidate the old one.
		scopes, err := ar.userStorage.RequestScopes(user.Id, d.Scopes)
		if err != nil {
			ar.Error(w, ErrorAPIRequestScopesForbidden, http.StatusForbidden, err.Error(), "LoginWithPassword.RequestScopes")
			return
		}

		offline := contains(scopes, jwtService.OfflineScope)
		accessToken, refreshToken, err := ar.loginUser(user, d.Scopes, app, offline, false)
		if err != nil {
			ar.Error(w, ErrorAPIAppAccessTokenNotCreated, http.StatusInternalServerError, err.Error(), "LoginWithPassword.loginUser")
			return
		}

		// Blacklist old access token.
		if err := ar.tokenBlacklist.Add(oldAccessTokenString); err != nil {
			ar.logger.Printf("Cannot blacklist old access token: %s\n", err)
		}

		extensions.SanitizeUser(user)
		result := &AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			User:         user,
		}

		ar.userStorage.UpdateLoginMetadata(user.Id)
		ar.ServeJSON(w, http.StatusOK, result)
	}
}

// RequestDisabledTFA requests link for disabling TFA.
func (ar *Router) RequestDisabledTFA() http.HandlerFunc {
	type requestBody struct {
		Email string `json:"email,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		d := requestBody{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		if !shared.EmailRegexp.MatchString(d.Email) {
			ar.Error(w, ErrorAPIRequestBodyInvalid, http.StatusBadRequest, "", "RequestDisabledTFA.emailRegexp_MatchString")
			return
		}

		if userExists := ar.userStorage.UserExists(d.Email); !userExists {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusBadRequest, "User with this email does not exist", "RequestDisabledTFA.UserExists")
			return
		}

		app := middleware.AppFromContext(r.Context())
		if app == nil {
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App is not in context.", "RequestDisabledTFA.AppFromContext")
			return
		}

		if app.TFAStatus() == model.TFAStatusMandatory {
			ar.Error(w, ErrorAPIRequestMandatoryTFA, http.StatusForbidden, "Two-factor authentication is mandatory for this app", "RequestDisabledTFA.TFAStatusMandatory")
			return
		}

		userID, err := ar.userStorage.IDByName(d.Email)
		if err != nil {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusBadRequest, err.Error(), "RequestDisabledTFA.IDByName")
			return
		}

		resetToken, err := ar.tokenService.NewResetToken(userID)
		if err != nil {
			ar.Error(w, ErrorAPIAppResetTokenNotCreated, http.StatusInternalServerError, err.Error(), "RequestDisabledTFA.NewResetToken")
			return
		}

		resetTokenString, err := ar.tokenService.String(resetToken)
		if err != nil {
			ar.Error(w, ErrorAPIAppResetTokenNotCreated, http.StatusInternalServerError, err.Error(), "RequestDisabledTFA.tokenService_String")
			return
		}

		host, err := url.Parse(ar.Host)
		if err != nil {
			ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "RequestDisabledTFA.URL_parse")
			return
		}

		query := fmt.Sprintf("token=%s", resetTokenString)

		u := &url.URL{
			Scheme:   host.Scheme,
			Host:     host.Host,
			Path:     path.Join(ar.WebRouterPrefix, "tfa/disable"),
			RawQuery: query,
		}

		if err = ar.emailService.SendResetEmail("Disable Two-Factor Authentication", d.Email, u.String()); err != nil {
			ar.Error(w, ErrorAPIEmailNotSent, http.StatusInternalServerError, "Email sending error: "+err.Error(), "RequestDisabledTFA.SendResetEmail")
			return
		}

		result := map[string]string{"result": "ok"}
		ar.ServeJSON(w, http.StatusOK, result)
	}
}

// RequestTFAReset requests link for resetting TFA: deleting old shared secret and establishing the new one.
func (ar *Router) RequestTFAReset() http.HandlerFunc {
	type requestBody struct {
		Email string `json:"email,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		d := requestBody{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		if !shared.EmailRegexp.MatchString(d.Email) {
			ar.Error(w, ErrorAPIRequestBodyInvalid, http.StatusBadRequest, "", "RequestTFAReset.emailRegexp_MatchString")
			return
		}

		if userExists := ar.userStorage.UserExists(d.Email); !userExists {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusBadRequest, "User with this email does not exist", "RequestTFAReset.UserExists")
			return
		}

		userID, err := ar.userStorage.IDByName(d.Email)
		if err != nil {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusBadRequest, err.Error(), "RequestTFAReset.IDByName")
			return
		}

		app := middleware.AppFromContext(r.Context())
		if app == nil {
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App is not in context.", "RequestDisabledTFA.AppFromContext")
			return
		}

		if app.TFAStatus() == model.TFAStatusDisabled {
			ar.Error(w, ErrorAPIRequestDisabledTFA, http.StatusForbidden, "Two-factor authentication is disabled for this app", "RequestTFAReset.TFAStatusDisabled")
			return
		}

		resetToken, err := ar.tokenService.NewResetToken(userID)
		if err != nil {
			ar.Error(w, ErrorAPIAppResetTokenNotCreated, http.StatusInternalServerError, err.Error(), "RequestTFAReset.NewResetToken")
			return
		}

		resetTokenString, err := ar.tokenService.String(resetToken)
		if err != nil {
			ar.Error(w, ErrorAPIAppResetTokenNotCreated, http.StatusInternalServerError, err.Error(), "RequestTFAReset.tokenService_String")
			return
		}

		host, err := url.Parse(ar.Host)
		if err != nil {
			ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "RequestTFAReset.URL_parse")
			return
		}

		query := fmt.Sprintf("token=%s", resetTokenString)

		u := &url.URL{
			Scheme:   host.Scheme,
			Host:     host.Host,
			Path:     path.Join(ar.WebRouterPrefix, "tfa/reset"),
			RawQuery: query,
		}

		if err = ar.emailService.SendResetEmail("Reset Two-Factor Authentication", d.Email, u.String()); err != nil {
			ar.Error(w, ErrorAPIEmailNotSent, http.StatusInternalServerError, "Email sending error: "+err.Error(), "RequestTFAReset.SendResetEmail")
			return
		}

		result := map[string]string{"result": "ok"}
		ar.ServeJSON(w, http.StatusOK, result)
	}
}
