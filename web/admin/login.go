package admin

import (
	"net/http"
)

type adminData struct {
	Admin    string `yaml:"admin"`
	Password string `yaml:"password"`
}

// Login logins admin with admin name and password.
func (ar *Router) Login() http.HandlerFunc {
	type loginData struct {
		Email    string `json:"email,omitempty" validate:"required,gte=6,lte=130"`
		Password string `json:"password,omitempty" validate:"required,gte=6,lte=130"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		conf := new(adminData)
		if ar.getConf(w, conf) != nil {
			return
		}

		ld := loginData{}
		if ar.mustParseJSON(w, r, &ld) != nil {
			return
		}

		if (conf.Admin != ld.Email) || (conf.Password != ld.Password) {
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
		return
	}
}
