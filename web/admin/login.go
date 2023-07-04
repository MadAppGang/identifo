package admin

import (
	"net/http"
	"os"

	"github.com/madappgang/identifo/v2/l"
)

type adminLoginData struct {
	Login    string `json:"email"`
	Password string `json:"password"`
}

// Login logins admin with admin name and password.
func (ar *Router) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		login, pswd, err := ar.getAdminAccountSettings()
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAdminPanelAdminCredentialsError, err.Error())
			return
		}

		ld := adminLoginData{}
		if err = ar.mustParseJSON(w, r, &ld); err != nil {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPIJsonParseError, err.Error())
			return
		}

		if (login != ld.Login) || (pswd != ld.Password) {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAdminPanelAdminCredentialsMismatch)
			return
		}

		session, err := ar.server.Services().Session.NewSession()
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAdminPanelCreateSession, err)
			return
		}

		if err = ar.server.Storages().Session.InsertSession(session); err != nil {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAdminPanelCreateSession, err)
			return
		}

		c := &http.Cookie{
			Name:     cookieName,
			Value:    encode(session.ID),
			Path:     "/",
			MaxAge:   ar.server.Services().Session.SessionDurationSeconds(),
			HttpOnly: true,
		}
		http.SetCookie(w, c)
		ar.ServeJSON(w, locale, http.StatusOK, nil)
	}
}

func (ar *Router) getAdminAccountSettings() (string, string, error) {
	loginEnvName := ar.server.Settings().AdminAccount.LoginEnvName
	pswdEnvName := ar.server.Settings().AdminAccount.PasswordEnvName

	if len(loginEnvName) == 0 || len(pswdEnvName) == 0 {
		return "", "", l.ErrorAdminPanelAdminCredentialsNotSet
	}
	login := os.Getenv(loginEnvName)
	password := os.Getenv(pswdEnvName)

	if len(login) == 0 || len(password) == 0 {
		return "", "", l.ErrorAdminPanelAdminCredentialsNotSet
	}

	return login, password, nil
}
