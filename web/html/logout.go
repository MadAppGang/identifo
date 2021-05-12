package html

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
		if len(app.ID) == 0 {
			ar.Error(w, nil, http.StatusInternalServerError, "couldn't get app from context")
			return
		}

		scopesJSON := strings.TrimSpace(r.URL.Query().Get("scopes"))
		scopes := []string{}
		if err := json.Unmarshal([]byte(scopesJSON), &scopes); err != nil {
			ar.Logger.Printf("Invalid scopes %v", scopesJSON)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		callbackURL := strings.TrimSpace(r.URL.Query().Get(callbackURLKey))
		if !contains(app.RedirectURLs, callbackURL) {
			ar.Logger.Printf("Unauthorized callback url %v for app %v", callbackURL, app.ID)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		redirectURI := strings.TrimSpace(r.URL.Query().Get(redirectURIKey))
		if redirectURI == "" {
			http.Redirect(w, r, callbackURL, http.StatusFound)
			return
		}

		redirectURIParsed, err := url.Parse(redirectURI)
		if err != nil {
			ar.Logger.Printf("cannot parse redirect url %v", redirectURI)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		if redirectURIParsed.Host != r.Host {
			ar.Logger.Printf("provided redirect url host %v is not allowed", redirectURIParsed.Host)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		redirectURI = fmt.Sprintf("%s?callbackUrl=%s", redirectURI, callbackURL)
		http.Redirect(w, r, redirectURI, http.StatusFound)
	}
}
