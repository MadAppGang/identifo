package admin

import (
	"net/http"
)

const (
	adminFormKey    = "email" // TODO: change to 'admin' when admin-login.html template is provided.
	passwordFormKey = "password"
)

type adminData struct {
	Admin    string `yaml:"admin"`
	Password string `yaml:"password"`
}

// Login logins admin with admin name and password.
func (ar *Router) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conf := new(adminData)
		if ar.getConf(w, conf) != nil {
			return
		}

		username := r.FormValue(adminFormKey)
		password := r.FormValue(passwordFormKey)

		if (conf.Admin != username) || (conf.Password != password) {
			ar.Error(w, ErrorIncorrectLogin, http.StatusBadRequest, "")
			return
		}

		session, err := ar.sessionService.NewSession()
		if err != nil {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, "")
			return
		}

		if err = ar.sessionStorage.InsertSession(session); err != nil {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, "")
			return
		}

		c := &http.Cookie{
			Name:     cookieName,
			Value:    encode(session.ID),
			MaxAge:   ar.sessionService.SessionDurationSeconds(),
			HttpOnly: true,
		}
		http.SetCookie(w, c)
		// TODO: redirect to success page.
		return
	}
}
