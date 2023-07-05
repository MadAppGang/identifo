package management

import (
	"net/http"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
)

func (ar *Router) getResetPasswordToken(w http.ResponseWriter, r *http.Request) {
	locale := r.Header.Get("Accept-Language")

	var d ResetPasswordTokenRequest
	if ar.MustParseJSON(w, r, &d) != nil {
		return
	}

	if d.Email != "" && !model.EmailRegexp.MatchString(d.Email) {
		ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyEmailInvalid)
		return
	}

	if len(d.Email) == 0 {
		ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyEmailInvalid)
		return
	}

	user, err := ar.server.Storages().UC.UserBySecondaryID(r.Context(), model.AuthIdentityTypeEmail, d.Email)
	if err == l.ErrorUserNotFound {
		// return ok, but there is no user
		ar.logger.Printf("Trying to reset password for the user, which is not exists: %s. Sending back ok to user for security reason.", d.Email)
		result := map[string]string{"result": "ok"}
		ar.ServeJSON(w, locale, http.StatusOK, result)
		return
	} else if err != nil {
		ar.Error(w, locale, http.StatusInternalServerError, l.ErrorStorageFindUserEmailError, d.Email, err)
		return
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

	result := map[string]string{"result": "ok", "token": resetTokenString}
	ar.ServeJSON(w, locale, http.StatusOK, result)
}
