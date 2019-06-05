package admin

import (
	"fmt"
	"net/http"
)

type adminLoginData struct {
	Login    string `yaml:"admin" json:"email,omitempty"`
	Password string `yaml:"password" json:"password,omitempty"`
}

func (ld *adminLoginData) validate() error {
	loginLen := len(ld.Login)
	if loginLen < 6 || loginLen > 130 {
		return fmt.Errorf("Incorrect login length %d, expected a number between 6 and 130", loginLen)
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
		conf := new(adminLoginData)
		if ar.getAccountConf(w, conf) != nil {
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
			MaxAge:   ar.sessionService.SessionDurationSeconds(),
			HttpOnly: true,
		}
		http.SetCookie(w, c)
	}
}
