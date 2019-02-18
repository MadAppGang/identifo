package html

import (
	"encoding/json"
	"net/http"
	"path"
)

const usernameKey = "email"
const passwordKey = "password"
const scopesKey = "scopes"

// Login logins user with email and password.
func (ar *Router) Login() http.HandlerFunc {
	errorPath := path.Join(ar.PathPrefix, "/misconfiguration")

	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue(usernameKey)
		password := r.FormValue(passwordKey)
		scopesJSON := r.FormValue(scopesKey)
		scopes := []string{}
		app := appFromContext(r.Context())

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

		scopes, err = ar.UserStorage.RequestScopes(user.ID(), scopes)
		if err != nil {
			ar.Logger.Printf("Error: invalid scopes %v for userID: %v", scopes, user.ID())
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		setCookie(w, "identifo-user", user.ID(), 60*60*24*2)
		redirectToLogin()
	}
}
