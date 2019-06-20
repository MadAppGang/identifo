package admin

import (
	"net/http"
)

type adminLoginData struct {
	Login    string `yaml:"admin" json:"email,omitempty"`
	Password string `yaml:"password" json:"password,omitempty"`
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
			Path:     "/",
			MaxAge:   ar.sessionService.SessionDurationSeconds(),
			HttpOnly: true,
		}
		http.SetCookie(w, c)
		ar.ServeJSON(w, http.StatusOK, nil)
	}
}
