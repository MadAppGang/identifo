package html

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path"
	"strings"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/middleware"
)

// Register creates user.
func (ar *Router) Register() http.HandlerFunc {
	errorPath := path.Join(ar.PathPrefix, "/misconfiguration")

	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue(usernameKey)
		password := r.FormValue(passwordKey)
		scopesJSON := r.FormValue(scopesKey)
		scopes := []string{}

		if err := json.Unmarshal([]byte(scopesJSON), &scopes); err != nil {
			ar.Logger.Printf("Error: Invalid scopes %v", scopesJSON)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		app := middleware.AppFromContext(r.Context())
		if app == nil {
			ar.Error(w, nil, http.StatusInternalServerError, "Couldn't get app from context")
		}

		redirectToRegister := func() {
			q := r.URL.Query()
			q.Set(FormKeyAppID, app.ID())
			q.Set(scopesKey, scopesJSON)
			r.URL.RawQuery = q.Encode()

			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusFound)
		}

		redirectToLogin := func() {
			r.URL.Path = "login"

			q := r.URL.Query()
			q.Set(FormKeyAppID, app.ID())
			q.Set(scopesKey, scopesJSON)
			r.URL.RawQuery = q.Encode()

			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusFound)
		}

		if app.RegistrationForbidden() {
			SetFlash(w, FlashErrorMessageKey, ErrorRegistrationForbidden.Error())
			redirectToRegister()
			return
		}

		// Authorize user if the app requires authorization.
		azi := authzInfo{
			app:         app,
			userRole:    app.NewUserDefaultRole(),
			resourceURI: r.RequestURI,
			method:      r.Method,
		}

		if err := ar.authorize(azi); err != nil {
			SetFlash(w, FlashErrorMessageKey, err.Error())
			redirectToRegister()
			return
		}

		// Validate password.
		if err := model.StrongPswd(password); err != nil {
			SetFlash(w, FlashErrorMessageKey, err.Error())
			redirectToRegister()
			return
		}

		// Create new user.
		user, err := ar.UserStorage.AddUserByNameAndPassword(username, password, app.NewUserDefaultRole(), nil)
		if err != nil {
			if err == model.ErrorUserExists {
				SetFlash(w, FlashErrorMessageKey, err.Error())
				redirectToRegister()
				return
			}

			ar.Logger.Printf("Error: creating user by name and password %v.", err)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		// Do login flow.
		scopes, err = ar.UserStorage.RequestScopes(user.ID(), scopes)
		if err != nil {
			ar.Logger.Printf("Error: requesting scopes %v.", err)
			http.Redirect(w, r, errorPath, http.StatusFound)
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

		setCookie(w, CookieKeyWebCookieToken, tokenString, int(ar.TokenService.WebCookieTokenLifespan()))
		redirectToLogin()
	}
}

// RegistrationHandler serves registration page.
func (ar *Router) RegistrationHandler() http.HandlerFunc {
	tmpl, err := template.ParseFiles(path.Join(ar.StaticFilesPath.PagesPath, ar.StaticPages.Registration))
	if err != nil {
		ar.Logger.Fatalln("Cannot parse Registration template.", err)
	}
	errorPath := path.Join(ar.PathPrefix, "/misconfiguration")

	return func(w http.ResponseWriter, r *http.Request) {
		app := middleware.AppFromContext(r.Context())
		if app == nil {
			ar.Error(w, nil, http.StatusInternalServerError, "Couldn't get app from context")
		}

		scopesJSON := strings.TrimSpace(r.URL.Query().Get("scopes"))
		scopes := []string{}
		if scopesJSON != "" {
			if err := json.Unmarshal([]byte(scopesJSON), &scopes); err != nil {
				ar.Logger.Printf("Error: Invalid scopes %v. Error: %v", scopesJSON, err)
				http.Redirect(w, r, errorPath, http.StatusFound)
				return
			}
		}

		errorMessage, err := GetFlash(w, r, FlashErrorMessageKey)
		if err != nil {
			ar.Logger.Printf("Error: getting flash message %v", err)
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
}
