package html

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"

	"github.com/madappgang/identifo/web/middleware"
)

// Logout removes user's session.
func (ar *Router) Logout() http.HandlerFunc {
	errorPath := path.Join(ar.PathPrefix, "/misconfiguration")

	return func(w http.ResponseWriter, r *http.Request) {
		deleteCookie(w, CookieKeyWebCookieToken)

		app := middleware.AppFromContext(r.Context())
		if app == nil {
			ar.Error(w, nil, http.StatusInternalServerError, "Couldn't get app from context")
			return
		}

		scopesJSON := strings.TrimSpace(r.URL.Query().Get("scopes"))
		scopes := []string{}
		if err := json.Unmarshal([]byte(scopesJSON), &scopes); err != nil {
			ar.Logger.Printf("Invalid scopes %v", scopesJSON)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		r.URL.Path = path.Join(ar.PathPrefix, "/login")
		http.Redirect(w, r, r.URL.String(), http.StatusFound)
	}
}
