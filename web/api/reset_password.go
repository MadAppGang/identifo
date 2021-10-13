package api

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/middleware"
)

// RequestResetPassword requests password reset
func (ar *Router) RequestResetPassword() http.HandlerFunc {
	type resetRequest struct {
		login
		TFACode string `json:"tfa_code,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		d := resetRequest{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		if err := d.login.validate(); err != nil {
			ar.Error(w, ErrorAPIRequestBodyParamsInvalid, http.StatusBadRequest, err.Error(), "RequestResetPassword.validate")
			return
		}

		if err := ar.checkSupportedWays(d.login); err != nil {
			ar.Error(w, ErrorAPIAppLoginWithUsernameNotSupported, http.StatusBadRequest, err.Error(), "RequestResetPassword.unsupportedLoginWays")
			return
		}

		user, err := ar.server.Storages().User.UserByEmail(d.Email)
		if err == model.ErrUserNotFound {
			result := map[string]string{"result": "ok"}
			ar.ServeJSON(w, http.StatusOK, result)
			return
		} else if err != nil {
			ar.Error(w, ErrorAPIInternalServerError, http.StatusBadRequest, "Unable to get user with email", "RequestResetPassword.ErrorGettingUser")
			return
		}

		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusInternalServerError, "App is not in context.", "RequestResetPassword.AppFromContext")
			return
		}

		if user.TFAInfo.IsEnabled && ar.tfaType != model.TFATypeEmail {
			if d.TFACode != "" {
				otpVerified, err := ar.verifyOTPCode(user, d.TFACode)
				if err != nil {
					ar.Error(w, ErrorAPIRequestScopesForbidden, http.StatusForbidden, err.Error(), "FinalizeTFA.OTP_Invalid")
					return
				}

				dontNeedVerification := app.DebugTFACode != "" && d.TFACode == app.DebugTFACode

				if !(otpVerified || dontNeedVerification) {
					ar.Error(w, ErrorAPIRequestTFACodeInvalid, http.StatusUnauthorized, "", "FinalizeTFA.OTP_Invalid")
					return
				}
			} else {
				if err := ar.sendOTPCode(user); err != nil {
					ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "RequestResetPassword.SendOTPCode")
					return
				}
				result := map[string]string{"result": "tfa-required"}
				ar.ServeJSON(w, http.StatusOK, result)
				return
			}
		}

		resetToken, err := ar.server.Services().Token.NewResetToken(user.ID)
		if err != nil {
			ar.Error(w, ErrorAPIAppResetTokenNotCreated, http.StatusInternalServerError, err.Error(), "RequestResetPassword.NewResetToken")
			return
		}

		resetTokenString, err := ar.server.Services().Token.String(resetToken)
		if err != nil {
			ar.Error(w, ErrorAPIAppResetTokenNotCreated, http.StatusInternalServerError, err.Error(), "RequestResetPassword.tokenService_String")
			return
		}

		query := fmt.Sprintf("appId=%s&token=%s", app.ID, resetTokenString)

		host, err := url.Parse(ar.Host)
		if err != nil {
			ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "RequestResetPassword.URL_parse")
			return
		}

		u := &url.URL{
			Scheme:   host.Scheme,
			Host:     host.Host,
			Path:     path.Join(ar.WebRouterPrefix, "password/reset"),
			RawQuery: query,
		}
		uu := &url.URL{Scheme: host.Scheme, Host: host.Host, Path: path.Join(ar.WebRouterPrefix, "password/reset")}

		resetEmailData := model.ResetEmailData{
			User:  user,
			Token: resetTokenString,
			URL:   u.String(),
			Host:  uu.String(),
		}

		if err = ar.server.Services().Email.SendResetEmail("Reset Password", d.Email, resetEmailData); err != nil {
			ar.Error(w, ErrorAPIEmailNotSent, http.StatusInternalServerError, "Email sending error: "+err.Error(), "RequestResetPassword.SendResetEmail")
			return
		}

		result := map[string]string{"result": "ok"}
		ar.ServeJSON(w, http.StatusOK, result)
	}
}

// ResetPassword handles password reset form submission (POST request).
func (ar *Router) ResetPassword() http.HandlerFunc {
	type newPassword struct {
		Password string `json:"password,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		d := newPassword{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}
		if err := model.StrongPswd(d.Password); err != nil {
			ar.Error(w, ErrorAPIRequestPasswordWeak, http.StatusBadRequest, err.Error(), "ResetPassword.StrongPswd")
			return
		}

		accessTokenBytes, ok := r.Context().Value(model.TokenRawContextKey).([]byte)
		if !ok {
			ar.Error(w, ErrorAPIRequestTokenInvalid, http.StatusBadRequest, "Token bytes are not in context.", "ResetPassword.TokenBytesFromContext")
			return
		}

		// Get userID from token and update user with this ID.
		userID, err := ar.getTokenSubject(string(accessTokenBytes))
		if err != nil {
			ar.Error(w, ErrorAPIAppCannotExtractTokenSubject, http.StatusInternalServerError, err.Error(), "ResetPassword.getTokenSubject")
			return
		}

		user, err := ar.server.Storages().User.UserByID(userID)
		if err != nil {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusUnauthorized, err.Error(), "ResetPassword.UserByID")
			return
		}

		// Save new password.
		if err := ar.server.Storages().User.ResetPassword(user.ID, d.Password); err != nil {
			ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, "Reset password. Error: "+err.Error(), "ResetPassword.ResetPassword")
			return
		}

		// Refetch user with new password hash.
		if user, err = ar.server.Storages().User.UserByUsername(user.Username); err != nil {
			ar.Error(w, ErrorAPIRequestBodyOldPasswordInvalid, http.StatusBadRequest, err.Error(), "ResetPassword.RefetchUser")
			return
		}

		result := map[string]string{"result": "ok"}
		ar.ServeJSON(w, http.StatusOK, result)
	}
}
