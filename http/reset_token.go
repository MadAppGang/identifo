package http

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

func (ar *apiRouter) SendResetToken() http.HandlerFunc {
	tmpl, err := template.New("reset").Parse("Hi! We got a request to reset your password. Click <a href=\"{{.}}\">here</a> to reset your password.")
	emailRegexp, regexpErr := regexp.Compile(emailExpr)

	return func(w http.ResponseWriter, r *http.Request) {
		if err != nil || regexpErr != nil {
			SetFlash(w, ErrorMessageKey, "Server Error. Try later please")
			http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
			return
		}

		err = r.ParseForm()
		if err != nil {
			SetFlash(w, ErrorMessageKey, "Invalid request")
			http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
		}

		name := r.FormValue("email")
		if !emailRegexp.MatchString(name) {
			SetFlash(w, ErrorMessageKey, "Invalid email")
			http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
			return
		}

		userExists := ar.userStorage.UserExists(name)
		if !userExists {
			SetFlash(w, ErrorMessageKey, "This Email is unregistred")
			http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
			return
		}

		id, err := ar.userStorage.IDByName(name)
		if err != nil {
			SetFlash(w, ErrorMessageKey, "This Email is unregistred")
			http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
			return
		}

		t, err := ar.tokenService.NewResetToken(id)
		if err != nil {
			SetFlash(w, ErrorMessageKey, "Server Error. Try later please")
			http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
			return
		}

		token, err := ar.tokenService.String(t)
		if err != nil {
			SetFlash(w, ErrorMessageKey, "Server Error. Try later please")
			http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
			return
		}

		query := fmt.Sprintf("token=%s", token)
		u := &url.URL{
			Scheme:   "http",
			Host:     r.Host,
			Path:     "password/reset",
			RawQuery: query,
		}

		var tpl bytes.Buffer
		if err = tmpl.Execute(&tpl, u.String()); err != nil {
			SetFlash(w, ErrorMessageKey, "Server Error. Try later please")
			http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
			return
		}

		_, _, err = ar.emailService.SendHTML("Reset Password", tpl.String(), name)
		if err != nil {
			SetFlash(w, ErrorMessageKey, "Error sending email")
			http.Redirect(w, r, r.URL.String(), http.StatusMovedPermanently)
			return
		}

		url := path.Join(r.URL.String(), "success")
		http.Redirect(w, r, url, http.StatusMovedPermanently)
		return
	}
}
