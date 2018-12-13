package api

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"regexp"
)

const emailExpr = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

//RequestResetPassword - request reset password
func (ar *Router) RequestResetPassword() http.HandlerFunc {

	tmpl, _ := template.New("reset").Parse("Hi! We got a request to reset your password. Click <a href=\"{{.}}\">here</a> to reset your password.")
	emailRegexp, _ := regexp.Compile(emailExpr)

	type resetRequestEmail struct {
		Email string `json:"email,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		d := resetRequestEmail{}
		if ar.MustParseJSON(w, r, &d) != nil {
			ar.Error(w, ErrorWrongInput, http.StatusBadRequest, "Invalid input data")
			return
		}
		if !emailRegexp.MatchString(d.Email) {
			ar.Error(w, ErrorWrongInput, http.StatusBadRequest, "Invalid Email")
			return
		}

		userExists := ar.userStorage.UserExists(d.Email)
		if !userExists {
			ar.Error(w, ErrorWrongInput, http.StatusBadRequest, "User with with email is not registered")
			return
		}

		id, err := ar.userStorage.IDByName(d.Email)
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, "User with with email is not registered")
			return
		}

		t, err := ar.tokenService.NewResetToken(id)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "TokenService error:"+err.Error())
			return
		}

		token, err := ar.tokenService.String(t)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "TokenService error:"+err.Error())
			return
		}

		query := fmt.Sprintf("token=%s", token)
		host, _ := url.Parse(ar.Host)

		u := &url.URL{
			Scheme:   host.Scheme,
			Host:     host.Host,
			Path:     path.Join(ar.WebRouterPrefix, "password/reset"),
			RawQuery: query,
		}

		var tpl bytes.Buffer
		if err = tmpl.Execute(&tpl, u.String()); err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "Email template compile error:"+err.Error())
			return
		}

		err = ar.emailService.SendHTML("Reset Password", tpl.String(), d.Email)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "Email sending error:"+err.Error())
			return
		}

		result := map[string]string{"Result": "ok"}
		ar.ServeJSON(w, http.StatusOK, result)
	}
}
