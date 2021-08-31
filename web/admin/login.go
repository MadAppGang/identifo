package admin

import (
	"fmt"
	"net/http"
	"os"
)

type adminLoginData struct {
	Login    string `json:"email"`
	Password string `json:"password"`
}

// Login logins admin with admin name and password.
func (ar *Router) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, pswd, err := ar.getAdminAccountSettings()
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		ld := adminLoginData{}
		if err = ar.mustParseJSON(w, r, &ld); err != nil {
			ar.Error(w, fmt.Errorf("unable to parse login and pssword: %s", err.Error), http.StatusBadRequest, "")
			return
		}

		if (login != ld.Login) || (pswd != ld.Password) {
			ar.Error(w, ErrorIncorrectLogin, http.StatusBadRequest, "")
			return
		}

		session, err := ar.server.Services().Session.NewSession()
		if err != nil {
			ar.Error(w, fmt.Errorf("Cannot create session: %s", err), http.StatusInternalServerError, "")
			return
		}

		if err = ar.server.Storages().Session.InsertSession(session); err != nil {
			ar.Error(w, fmt.Errorf("Cannot insert session: %s", err), http.StatusInternalServerError, "")
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
		ar.ServeJSON(w, http.StatusOK, nil)
	}
}

func (ar *Router) getAdminAccountSettings() (string, string, error) {
	loginEnvName := ar.server.Settings().AdminAccount.LoginEnvName
	pswdEnvName := ar.server.Settings().AdminAccount.PasswordEnvName

	if len(loginEnvName) == 0 || len(pswdEnvName) == 0 {
		return "", "", ErrorAdminAccountIsNotSet
	}
	login := os.Getenv(loginEnvName)
	password := os.Getenv(pswdEnvName)

	if len(login) == 0 || len(password) == 0 {
		return "", "", ErrorAdminAccountNoEmailAndPassword
	}

	return login, password, nil
}
