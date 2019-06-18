package html

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/madappgang/identifo/web/middleware"
)

const (
	usernameKey = "email"
	passwordKey = "password"
	scopesKey   = "scopes"
)

// Login logins user with email and password.
func (ar *Router) Login() http.HandlerFunc {
	errorPath := path.Join(ar.PathPrefix, "/misconfiguration")

	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue(usernameKey)
		password := r.FormValue(passwordKey)
		scopesJSON := r.FormValue(scopesKey)
		scopes := []string{}
		app := middleware.AppFromContext(r.Context())

		redirectToLogin := func() {
			q := r.URL.Query()
			q.Set(FormKeyAppID, app.ID())
			q.Set(scopesKey, scopesJSON)
			r.URL.RawQuery = q.Encode()

			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusFound)
		}

		if err := json.Unmarshal([]byte(scopesJSON), &scopes); err != nil {
			ar.Logger.Printf("Error: Invalid scopes %v", scopesJSON)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		user, err := ar.UserStorage.UserByNamePassword(username, password)
		if err != nil {
			SetFlash(w, FlashErrorMessageKey, "Invalid Username or Password")
			redirectToLogin()
			return
		}

		if _, err = ar.UserStorage.RequestScopes(user.ID(), scopes); err != nil {
			ar.Logger.Printf("Error: invalid scopes %v for userID: %v", scopes, user.ID())
			SetFlash(w, FlashErrorMessageKey, err.Error())
			redirectToLogin()
			return
		}

		token, err := ar.TokenService.NewWebCookieToken(user)
		if err != nil {
			ar.Logger.Printf("Error creating auth token %v", err)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		tokenString, err := ar.TokenService.String(token)
		if err != nil {
			ar.Logger.Printf("Error stringifying token: %v", err)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		go ar.UserStorage.UpdateLoginMetadata(user.ID())
		setCookie(w, CookieKeyWebCookieToken, tokenString, int(ar.TokenService.WebCookieTokenLifespan()))
		redirectToLogin()
	}
}
