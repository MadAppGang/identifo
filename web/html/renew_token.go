package html

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	jwtValidator "github.com/madappgang/identifo/jwt/validator"
	"github.com/madappgang/identifo/model"
)

// RenewToken creates new id_token if user is already authenticated.
func (ar *Router) RenewToken() http.HandlerFunc {
	tmpl, err := ar.Server.Storages().Static.ParseTemplate(model.StaticPagesNames.WebMessage)
	if err != nil {
		ar.Logger.Fatalln("Cannot parse WebMessage template.", err)
	}
	tokenValidator := jwtValidator.NewValidator(
		[]string{"identifo"},
		[]string{ar.Server.Services().Token.Issuer()},
		[]string{},
		[]string{model.TokenTypeWebCookie},
	)

	return func(w http.ResponseWriter, r *http.Request) {
		serveTemplate := func(errorMessage, AccessToken, redirectURI string) {
			if err != nil {
				ar.Logger.Printf("Error parsing template: %v", err)
				ar.Error(w, err, http.StatusInternalServerError, "Error parsing template")
				return
			}

			data := map[string]interface{}{
				"Error":       errorMessage,
				"AccessToken": AccessToken,
				"RedirectUri": redirectURI,
			}

			if err := tmpl.Execute(w, data); err != nil {
				ar.Logger.Printf("Error executing template: %v", err)
				ar.Error(w, err, http.StatusInternalServerError, "Error executing template")
				return
			}
		}

		appID := strings.TrimSpace(r.URL.Query().Get(FormKeyAppID))
		app, err := ar.Server.Storages().App.ActiveAppByID(appID)
		if err != nil {
			message := fmt.Sprintf("Error getting App by ID: %v", err)
			ar.Logger.Printf(message)
			serveTemplate(message, "", "")
			return
		}

		redirectURI := strings.TrimSpace(r.URL.Query().Get(redirectURIKey))

		tstr, err := getCookie(r, CookieKeyWebCookieToken)
		if err != nil || tstr == "" {
			ar.Logger.Printf("Error getting token from cookie: %v", err)
			deleteCookie(w, CookieKeyWebCookieToken)
			serveTemplate("not authorized", "", redirectURI)
			return
		}
		webCookieToken, err := ar.Server.Services().Token.Parse(tstr)
		if err != nil {
			ar.Logger.Printf("Error invalid token: %v", err)
			deleteCookie(w, CookieKeyWebCookieToken)
			serveTemplate("not authorized", "", redirectURI)
			return
		}

		if err = tokenValidator.Validate(webCookieToken); err != nil {
			ar.Logger.Printf("Error invalid token: %v", err)
			deleteCookie(w, CookieKeyWebCookieToken)
			serveTemplate("not authorized", "", redirectURI)
			return
		}

		userID := webCookieToken.UserID()

		user, err := ar.Server.Storages().User.UserByID(userID)
		if err != nil {
			ar.Logger.Printf("Error: getting UserByID: %v, userID: %v", err, userID)
			deleteCookie(w, CookieKeyWebCookieToken)
			serveTemplate("invalid user token", "", redirectURI)
			return
		}

		scopesJSON := strings.TrimSpace(r.URL.Query().Get(scopesKey))
		scopes := []string{}
		if err := json.Unmarshal([]byte(scopesJSON), &scopes); err != nil {
			ar.Logger.Printf("Error: Invalid scopes %v", scopesJSON)
			serveTemplate("invalid scopes", "", redirectURI)
			return
		}

		scopes, err = ar.Server.Storages().User.RequestScopes(userID, scopes)
		if err != nil {
			ar.Logger.Printf("Error: invalid scopes %v for userID: %v", scopes, userID)
			message := fmt.Sprintf("user not allowed to access this scopes %v", scopes)
			serveTemplate(message, "", redirectURI)
			return
		}

		token, err := ar.Server.Services().Token.NewAccessToken(user, scopes, app, false, nil)
		if err != nil {
			ar.Logger.Printf("Error creating token: %v", err)
			serveTemplate("server error", "", redirectURI)
			return
		}

		tokenString, err := ar.Server.Services().Token.String(token)
		if err != nil {
			ar.Logger.Printf("Error stringifying token: %v", err)
			serveTemplate("server error", "", redirectURI)
			return
		}

		serveTemplate("", tokenString, redirectURI)
	}
}
