package admin

import (
	"fmt"
	"net/http"
)

type adminLoginData struct {
	Login           string `json:"email"`
	LoginEnvName    string `json:"email_env_name"`
	Password        string `json:"password"`
	PasswordEnvName string `json:"password_env_name"`
}

// Login logins admin with admin name and password.
func (ar *Router) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conf := new(adminLoginData)
		if ar.getAdminAccountSettings(w, conf) != nil {
			return
		}

		ld := adminLoginData{}
		if ar.mustParseJSON(w, r, &ld) != nil {
			return
		}

		if (conf.Login != ld.Login) || (conf.Password != ld.Password) {
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
