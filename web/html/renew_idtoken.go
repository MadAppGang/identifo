package html

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

// RenewIDToken creates new id_token if user is already authenticated.
func (ar *Router) RenewIDToken(pathComponents ...string) http.HandlerFunc {
	tmpl, err := template.ParseFiles(pathComponents...)

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
		if appID == "" {
			serveTemplate("Empty appId param", "", "")
			return
		}

		app, err := ar.AppStorage.AppByID(appID)
		if err != nil {
			message := fmt.Sprintf("Error getting App by ID: %v", err)
			serveTemplate(message, "", "")
			return
		}

		if !app.Active() {
			message := fmt.Sprintf("App with ID: %v is inactive", app.ID())
			serveTemplate(message, "", "")
			return
		}

		encryptedID, err := getCookie(r, CookieKeyUserID)
		if err != nil || encryptedID == "" {
			serveTemplate("not authorized", "", app.RedirectURL())
			deleteCookie(w, CookieKeyUserID)
			return
		}

		uID, err := ar.Encryptor.Decrypt([]byte(encryptedID))
		if err != nil {
			serveTemplate("not authorized", "", app.RedirectURL())
			deleteCookie(w, CookieKeyUserID)
			return
		}

		userID := string(uID)
		user, err := ar.UserStorage.UserByID(userID)
		if err != nil {
			ar.Logger.Printf("Error: getting UserByID: %v, userID: %v", err, userID)
			deleteCookie(w, CookieKeyUserID)
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
