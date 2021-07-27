package html

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
)

const emailExpr = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

// SendResetToken POST form request handle with password request handle
func (ar *Router) SendResetToken() http.HandlerFunc {
	tmpl, err := template.New("reset").Parse("Hi! We got a request to reset your password. Click <a href=\"{{.}}\">here</a> to reset your password.")
	emailRegexp, regexpErr := regexp.Compile(emailExpr)

	return func(w http.ResponseWriter, r *http.Request) {
		upath := path.Join(ar.PathPrefix, r.URL.String())
		if err != nil || regexpErr != nil {
			SetFlash(w, FlashErrorMessageKey, "Server Error. Try later please")
			http.Redirect(w, r, upath, http.StatusMovedPermanently)
			return
		}

		err = r.ParseForm()
		if err != nil {
			SetFlash(w, FlashErrorMessageKey, "Invalid request")
			http.Redirect(w, r, upath, http.StatusMovedPermanently)
		}

		name := r.FormValue("email")
		if !emailRegexp.MatchString(name) {
			SetFlash(w, FlashErrorMessageKey, "Invalid email")
			http.Redirect(w, r, upath, http.StatusMovedPermanently)
			return
		}

		if userExists := ar.Server.Storages().User.UserExists(name); !userExists {
			SetFlash(w, FlashErrorMessageKey, "This Email is unregistered")
			http.Redirect(w, r, upath, http.StatusMovedPermanently)
			return
		}

		id, err := ar.Server.Storages().User.IDByName(name)
		if err != nil {
			SetFlash(w, FlashErrorMessageKey, "This Email is unregistered")
			http.Redirect(w, r, upath, http.StatusMovedPermanently)
			return
		}

		t, err := ar.Server.Services().Token.NewResetToken(id)
		if err != nil {
			SetFlash(w, FlashErrorMessageKey, "Server Error. Try later please")
			http.Redirect(w, r, upath, http.StatusMovedPermanently)
			return
		}

		token, err := ar.Server.Services().Token.String(t)
		if err != nil {
			SetFlash(w, FlashErrorMessageKey, "Server Error. Try later please")
			http.Redirect(w, r, upath, http.StatusMovedPermanently)
			return
		}

		query := fmt.Sprintf("token=%s", token)
		host, _ := url.Parse(ar.Host)

		u := &url.URL{
			Scheme:   host.Scheme,
			Host:     host.Host,
			Path:     path.Join(ar.PathPrefix, "password/reset"),
			RawQuery: query,
		}

		var tpl bytes.Buffer
		if err = tmpl.Execute(&tpl, u.String()); err != nil {
			SetFlash(w, FlashErrorMessageKey, "Server Error. Try later please")
			http.Redirect(w, r, upath, http.StatusMovedPermanently)
			return
		}

		err = ar.Server.Services().Email.SendHTML("Reset Password", tpl.String(), name)
		if err != nil {
			SetFlash(w, FlashErrorMessageKey, "Error sending email")
			http.Redirect(w, r, upath, http.StatusMovedPermanently)
			return
		}

		redirectURL := strings.TrimSpace(r.URL.Query().Get(redirectURLParam))
		if redirectURL != "" {
			http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
			return
		}

		url := path.Join(ar.PathPrefix, r.URL.String(), "success")
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	}
}
