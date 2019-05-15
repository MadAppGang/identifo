package html

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path"
	"strings"

	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/tokensrvc"
	"github.com/madappgang/identifo/web/middleware"
)

// LoginHandler serves login page or redirects to the callback_url if user is already authenticated.
func (ar *Router) LoginHandler(pathComponents ...string) http.HandlerFunc {
	tmpl, err := template.ParseFiles(path.Join(pathComponents...))
	errorPath := path.Join(ar.PathPrefix, "/misconfiguration")
	tokenValidator := jwt.NewDefaultValidator("identifo", ar.TokenService.Issuer(), "", tokensrvc.WebCookieTokenType)

	return func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		app := middleware.AppFromContext(r.Context())
		if app == nil {
			ar.Error(w, nil, http.StatusInternalServerError, "Couldn't get app from context")
		}

		scopesJSON := strings.TrimSpace(r.URL.Query().Get("scopes"))
		scopes := []string{}
		if err := json.Unmarshal([]byte(scopesJSON), &scopes); err != nil {
			ar.Logger.Printf("Error: Invalid scopes %v", scopesJSON)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		serveTemplate := func() {
			errorMessage, err := GetFlash(w, r, FlashErrorMessageKey)
			if err != nil {
				ar.Error(w, err, http.StatusInternalServerError, "")
				return
			}

			data := map[string]interface{}{
				"Error":  errorMessage,
				"Prefix": ar.PathPrefix,
				"Scopes": scopesJSON,
				"AppId":  app.ID(),
			}

			if err = tmpl.Execute(w, data); err != nil {
				ar.Error(w, err, http.StatusInternalServerError, "")
			}
		}

		tstr, err := getCookie(r, CookieKeyWebCookieToken)
		if err != nil || tstr == "" {
			ar.Logger.Printf("Error getting auth token cookie: %v", err)
			deleteCookie(w, CookieKeyWebCookieToken)
			serveTemplate()
			return
		}

		webCookieToken, err := ar.TokenService.Parse(tstr)
		if err != nil {
			ar.Logger.Printf("Error invalid token %v", err)
			deleteCookie(w, CookieKeyWebCookieToken)
			serveTemplate()
			return
		}

		if err = tokenValidator.Validate(webCookieToken); err != nil {
			ar.Logger.Printf("Error invalid token %v", err)
			deleteCookie(w, CookieKeyWebCookieToken)
			serveTemplate()
			return
		}

		userID := webCookieToken.UserID()
		user, err := ar.UserStorage.UserByID(userID)
		if err != nil {
			ar.Logger.Printf("Error: getting UserByID: %v, userID: %v", err, userID)
			serveTemplate()
			return
		}

		scopes, err = ar.UserStorage.RequestScopes(userID, scopes)
		if err != nil {
			ar.Logger.Printf("Error: invalid scopes %v for userID: %v", scopes, userID)
			serveTemplate()
			return
		}

		token, err := ar.TokenService.NewToken(user, scopes, app)
		if err != nil {
			ar.Logger.Printf("Error creating token: %v", err)
			serveTemplate()
			return
		}

		tokenString, err := ar.TokenService.String(token)
		if err != nil {
			ar.Logger.Printf("Error stringifying token: %v", err)
			serveTemplate()
			return
		}

		redirectURL := app.RedirectURL() + "#" + tokenString
		http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
	}
}
