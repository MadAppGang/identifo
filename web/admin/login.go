package admin

import (
	"fmt"
	"net/http"
)

type adminData struct {
	Admin    string `yaml:"admin" json:"admin"`
	Password string `yaml:"password" json:"password"`
}

type loginData struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (ld *loginData) validate() error {
	emailLen := len(ld.Email)
	if emailLen < 6 || emailLen > 130 {
		return fmt.Errorf("Incorrect email length %d, expected a number between 6 and 130", emailLen)
	}
	pswdLen := len(ld.Password)
	if pswdLen < 6 || pswdLen > 130 {
		return fmt.Errorf("Incorrect password length %d, expected a number between 6 and 130", pswdLen)
	}
	return nil
}

// Login logins admin with admin name and password.
func (ar *Router) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conf := new(adminData)
		if ar.getAccountConf(w, conf) != nil {
			return
		}

		ld := loginData{}
		if ar.mustParseJSON(w, r, &ld) != nil {
			return
		}

		if err := ld.validate(); err != nil {
			ar.Error(w, ErrorIncorrectLogin, http.StatusBadRequest, err.Error())
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
	}
}
