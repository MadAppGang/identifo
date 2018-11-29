package http

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
)

func (ar *apiRouter) SendResetToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		name := r.FormValue("email")
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

		mb := messageBody(u.String())

		_, _, err = ar.emailService.SendHTML("Reset Password", mb, name)
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

func messageBody(link string) string {
	message := "Hi! we got a request to reset your password. Click <a href=\"%s\">here</a> to reset your password"
	return fmt.Sprintf(message, link)
}
