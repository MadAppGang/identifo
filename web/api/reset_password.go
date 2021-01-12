package api

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/madappgang/identifo/plugin/shared"
)

// RequestResetPassword requests password reset.
func (ar *Router) RequestResetPassword() http.HandlerFunc {
	type resetRequestEmail struct {
		Email string `json:"email,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		d := resetRequestEmail{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}
		if !shared.EmailRegexp.MatchString(d.Email) {
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
