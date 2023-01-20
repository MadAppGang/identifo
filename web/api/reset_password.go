package api

import (
	"fmt"
	"net/http"
	"net/url"

	l "github.com/madappgang/identifo/v2/localization"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
)

// RequestResetPassword requests password reset
func (ar *Router) RequestResetPassword() http.HandlerFunc {
	type resetRequest struct {
		login
		TFACode      string `json:"tfa_code,omitempty"`
		ResetPageURL string `json:"reset_page_url,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		d := resetRequest{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		if err := d.login.validate(); err != nil {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyInvalidError, err)
			return
		}

		if err := ar.checkSupportedWays(d.login); err != nil {
			ar.Error(w, locale, http.StatusBadRequest, l.APIAPPUsernameLoginNotSupported)
			return
		}

		user, err := ar.server.Storages().User.UserByEmail(d.Email)
		if err == model.ErrUserNotFound {
			// return ok, but there is no user
			// TODO: add logging for for reset password for user, who is not exist
			result := map[string]string{"result": "ok"}
			ar.ServeJSON(w, locale, http.StatusOK, result)
			return
		} else if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorStorageFindUserEmailError, d.Email, err)
			return
		}

		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIAPPNoAPPInContext)
			return
		}

		_, enabled2FA, _ := ar.check2FA(app.TFAStatus, ar.tfaType, user)

		if enabled2FA && ar.tfaType != model.TFATypeEmail {
			if d.TFACode != "" {
				otpVerified, err := ar.verifyOTPCode(user, d.TFACode)
				if err != nil {
					ar.Error(w, locale, http.StatusForbidden, l.Error2FAVerifyFailError, err)
					return
				}

				dontNeedVerification := app.DebugTFACode != "" && d.TFACode == app.DebugTFACode

				if !(otpVerified || dontNeedVerification) {
					ar.Error(w, locale, http.StatusUnauthorized, l.ErrorAPILoginCodeInvalid)
					return
				}
			} else {
				if err := ar.sendOTPCode(app, user); err != nil {
					ar.Error(w, locale, http.StatusInternalServerError, l.ErrorServiceOtpSendError, err)
					return
				}
				result := map[string]string{"result": "tfa-required"}
				ar.ServeJSON(w, locale, http.StatusOK, result)
				return
			}
		}

		resetToken, err := ar.server.Services().Token.NewResetToken(user.ID)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorTokenUnableToCreateResetTokenError, err)
			return
		}

		resetTokenString, err := ar.server.Services().Token.String(resetToken)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorTokenUnableToCreateResetTokenError, err)
			return
		}

		query := fmt.Sprintf("appId=%s&token=%s", app.ID, resetTokenString)
		u := &url.URL{
			Scheme:   ar.Host.Scheme,
			Host:     ar.Host.Host,
			RawQuery: query,
		}

		resetPath := model.DefaultLoginWebAppSettings.ResetPasswordURL

		// if app requested reset password custom page, use it.
		if len(d.ResetPageURL) > 0 {
			resetPath = d.ResetPageURL
		} else if app.LoginAppSettings != nil && len(app.LoginAppSettings.ResetPasswordURL) > 0 {
			// rewrite path for app, if app has specific web app login settings
			resetPath = app.LoginAppSettings.ResetPasswordURL
		}

		resetPathURL, err := url.Parse(resetPath)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorAPPResetUrlError, resetPath, app.ID, err)
			return
		}

		// app settings could rewrite host or just path, if path is absolute - it rewrites host as well
		if resetPathURL.IsAbs() {
			u.Scheme = resetPathURL.Scheme
			u.Host = resetPathURL.Host
		}

		u.Path = resetPathURL.Path

		uu := &url.URL{Scheme: u.Scheme, Host: u.Host, Path: u.Path}

		resetEmailData := ResetEmailData{
			User:  user,
			Token: resetTokenString,
			URL:   u.String(),
			Host:  uu.String(),
		}

		if err = ar.server.Services().Email.SendTemplateEmail(
			model.EmailTemplateTypeResetPassword,
			app.GetCustomEmailTemplatePath(),
			"Reset Password",
			d.Email,
			model.EmailData{
				User: user,
				Data: resetEmailData,
			},
		); err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorServiceEmailSendError, err)
			return
		}

		result := map[string]string{"result": "ok"}
		ar.ServeJSON(w, locale, http.StatusOK, result)
	}
}

// ResetPassword handles password reset form submission (POST request).
func (ar *Router) ResetPassword() http.HandlerFunc {
	type newPassword struct {
		Password string `json:"password,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		d := newPassword{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}
		if err := model.StrongPswd(d.Password); err != nil {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequestPasswordWeak, err)
			return
		}

		accessTokenBytes, ok := r.Context().Value(model.TokenRawContextKey).([]byte)
		if !ok {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIContextNoToken)
			return
		}

		// Get userID from token and update user with this ID.
		userID, err := ar.getTokenSubject(string(accessTokenBytes))
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorAPIRequestTokenSubError, err)
			return
		}

		user, err := ar.server.Storages().User.UserByID(userID)
		if err != nil {
			ar.Error(w, locale, http.StatusUnauthorized, l.ErrorStorageFindUserIDError, userID, err)
			return
		}

		// Save new password.
		if err := ar.server.Storages().User.ResetPassword(user.ID, d.Password); err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorStorageResetPasswordUserError, user.ID, err)
			return
		}

		result := map[string]string{"result": "ok"}
		ar.ServeJSON(w, locale, http.StatusOK, result)
	}
}
