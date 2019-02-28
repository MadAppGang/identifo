package html

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"

	"github.com/madappgang/identifo/web/shared"

	"github.com/madappgang/identifo/model"
)

// RenewIDToken creates new id_token if user is already authenticated.
func (ar *Router) RenewIDToken(pathComponents ...string) http.HandlerFunc {
	tmpl, err := template.ParseFiles(path.Join(pathComponents...))

	return func(w http.ResponseWriter, r *http.Request) {
		serveTemplate := func(errorMessage, IDToken, callbackURL string) {
			if err != nil {
				ar.Logger.Printf("Error parsing template: %v", err)
				ar.Error(w, err, 500, "Error parsing template")
				return
			}

			data := map[string]interface{}{
				"Error":       errorMessage,
				"IDToken":     IDToken,
				"CallbackURL": callbackURL,
			}

			if err := tmpl.Execute(w, data); err != nil {
				ar.Logger.Printf("Error executing template: %v", err)
				ar.Error(w, err, 500, "Error executing template")
				return
			}
		}

		appID := strings.TrimSpace(r.URL.Query().Get(FormKeyAppID))
		app, err := shared.AppByID(ar.AppStorage, appID)
		if err != nil {
			message := fmt.Sprintf("Error getting App by ID: %v", err)
			serveTemplate(message, "", "")
			return
		}

		tstr, err := getCookie(r, CookieKeyWebCookieToken)
		if err != nil || tstr == "" {
			ar.Logger.Printf("Error getting toke from cookie: %v", err)
			deleteCookie(w, CookieKeyWebCookieToken)
			serveTemplate("not authorized", "", app.RedirectURL())
			return
		}

		webCookieToken, err := shared.ParseToken(tstr, ar.TokenService, model.WebCookieTokenType)
		if err != nil {
			ar.Logger.Printf("Error invalid token: %v", err)
			deleteCookie(w, CookieKeyWebCookieToken)
			serveTemplate("not authorized", "", app.RedirectURL())
			return
		}

		userID := webCookieToken.UserID()

		user, err := ar.UserStorage.UserByID(userID)
		if err != nil {
			ar.Logger.Printf("Error: getting UserByID: %v, userID: %v", err, userID)
			deleteCookie(w, CookieKeyWebCookieToken)
			serveTemplate("invalid user token", "", app.RedirectURL())
			return
		}

		scopesJSON := strings.TrimSpace(r.URL.Query().Get("scopes"))
		scopes := []string{}
		if err := json.Unmarshal([]byte(scopesJSON), &scopes); err != nil {
			ar.Logger.Printf("Error: Invalid scopes %v", scopesJSON)
			serveTemplate("invalid scopes", "", app.RedirectURL())
			return
		}

		scopes, err = ar.UserStorage.RequestScopes(userID, scopes)
		if err != nil {
			ar.Logger.Printf("Error: invalid scopes %v for userID: %v", scopes, userID)
			message := fmt.Sprintf("user not allowed to access this scopes %v", scopes)
			serveTemplate(message, "", app.RedirectURL())
			return
		}

		token, err := ar.TokenService.NewToken(user, scopes, app)
		if err != nil {
			ar.Logger.Printf("Error creating token: %v", err)
			serveTemplate("server error", "", app.RedirectURL())
			return
		}

		tokenString, err := ar.TokenService.String(token)
		if err != nil {
			ar.Logger.Printf("Error stringifying token: %v", err)
			serveTemplate("server error", "", app.RedirectURL())
			return
		}

		serveTemplate("", tokenString, app.RedirectURL())
		return
	}
}
