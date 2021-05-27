package api

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/madappgang/identifo/model"
)

// RequestResetPassword requests password reset
// FIXME: We used email as param but search for user by username
func (ar *Router) RequestResetPassword() http.HandlerFunc {
	type resetRequestEmail struct {
		Email string `json:"email,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		d := resetRequestEmail{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}
		if !model.EmailRegexp.MatchString(d.Email) {
			ar.Error(w, ErrorAPIRequestBodyInvalid, http.StatusBadRequest, "", "RequestResetPassword.emailRegexp_MatchString")
			return
		}

		if userExists := ar.userStorage.UserExists(d.Email); !userExists {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusBadRequest, "User with this email does not exist", "RequestResetPassword.UserExists")
			return
		}

		id, err := ar.userStorage.IDByName(d.Email)
		if err != nil {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusBadRequest, err.Error(), "RequestResetPassword.IDByName")
			return
		}

		resetToken, err := ar.tokenService.NewResetToken(id)
		if err != nil {
			ar.Error(w, ErrorAPIAppResetTokenNotCreated, http.StatusInternalServerError, err.Error(), "RequestResetPassword.NewResetToken")
			return
		}

		resetTokenString, err := ar.tokenService.String(resetToken)
		if err != nil {
			ar.Error(w, ErrorAPIAppResetTokenNotCreated, http.StatusInternalServerError, err.Error(), "RequestResetPassword.tokenService_String")
			return
		}

		query := fmt.Sprintf("token=%s", resetTokenString)

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

		if err = ar.emailService.SendResetEmail("Reset Password", d.Email, u.String()); err != nil {
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
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "Token bytes are not in context.", "ResetPassword.TokenBytesFromContext")
			return
		}

		// Get userID from token and update user with this ID.
		userID, err := ar.getTokenSubject(string(accessTokenBytes))
		if err != nil {
			ar.Error(w, ErrorAPIAppCannotExtractTokenSubject, http.StatusInternalServerError, err.Error(), "ResetPassword.getTokenSubject")
			return
		}

		user, err := ar.userStorage.UserByID(userID)
		if err != nil {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusUnauthorized, err.Error(), "ResetPassword.UserByID")
			return
		}

		// Save new password.
		if err := ar.userStorage.ResetPassword(user.ID, d.Password); err != nil {
			ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, "Reset password. Error: "+err.Error(), "ResetPassword.ResetPassword")
			return
		}

		// Refetch user with new password hash.
		if user, err = ar.userStorage.UserByNamePassword(user.Username, d.Password); err != nil {
			ar.Error(w, ErrorAPIRequestBodyOldPasswordInvalid, http.StatusBadRequest, err.Error(), "ResetPassword.RefetchUser")
			return
		}

		result := map[string]string{"result": "ok"}
		ar.ServeJSON(w, http.StatusOK, result)
	}
}
