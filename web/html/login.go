package html

import (
	"encoding/json"
	"net/http"
	"path"
)

const usernameKey = "email"
const passwordKey = "password"
const scopesKey = "scopes"

// Login login user with email and password
func (ar *Router) Login() http.HandlerFunc {
	errorPath := path.Join(ar.PathPrefix, "/misconfiguration")

	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue(usernameKey)
		password := r.FormValue(passwordKey)
		scopesJSON := r.FormValue(scopesKey)
		scopes := []string{}
		app := appFromContext(r.Context())

		err := json.Unmarshal([]byte(scopesJSON), &scopes)
		if err != nil {
			ar.Logger.Printf("Error: Invalid scopes %v", scopesJSON)
			http.Redirect(w, r, errorPath, http.StatusMovedPermanently)
			return
		}

		user, err := ar.UserStorage.UserByNamePassword(username, password)
		if err != nil {
			SetFlash(w, FlashErrorMessageKey, "Invalid Username or Password")

			q := r.URL.Query()
			q.Set(FormKeyAppID, app.ID())
			q.Set(scopesKey, scopesJSON)
			r.URL.RawQuery = q.Encode()

			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusMovedPermanently)
			return
		}

		scopes, err = ar.UserStorage.RequestScopes(user.ID(), scopes)
		if err != nil {
			ar.Logger.Printf("Error: invalid scopes %v for userID: %v", scopes, user.ID())
			http.Redirect(w, r, errorPath, http.StatusMovedPermanently)
			return
		}

		token, err := ar.TokenService.NewToken(user, scopes, app)
		if err != nil {
			ar.Logger.Printf("Error creating token: %v", err)
			ar.Error(w, err, http.StatusInternalServerError, "Server error, try later please")
			return
		}
		tokenString, err := ar.TokenService.String(token)
		if err != nil {
			ar.Logger.Printf("Error stringifying token: %v", err)
			ar.Error(w, err, http.StatusInternalServerError, "Server error, try later please")
			return
		}

		redirectURL := app.RedirectURL() + "#" + tokenString
		http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
	}
}
