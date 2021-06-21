package html

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/authorization"
	"github.com/madappgang/identifo/web/middleware"
)

// Register creates user.
func (ar *Router) Register() http.HandlerFunc {
	errorPath := path.Join(ar.PathPrefix, "/misconfiguration")

	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue(usernameKey)
		password := r.FormValue(passwordKey)
		scopesJSON := r.FormValue(scopesKey)
		callbackURL := r.FormValue(callbackURLKey)
		inviteToken := r.FormValue(inviteTokenKey)
		scopes := []string{}

		if err := json.Unmarshal([]byte(scopesJSON), &scopes); err != nil {
			ar.Logger.Printf("Error: Invalid scopes %v", scopesJSON)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		var isAnonymous bool
		var err error
		if isAnonymousStr := r.FormValue(isAnonymousKey); len(isAnonymousStr) > 0 {
			isAnonymous, err = strconv.ParseBool(isAnonymousStr)
			if err != nil {
				ar.Logger.Printf("Error: Invalid anonymous parameter %s", isAnonymousStr)
				http.Redirect(w, r, errorPath, http.StatusFound)
				return
			}
		}

		app := middleware.AppFromContext(r.Context())
		if app.ID == "" {
			ar.Error(w, nil, http.StatusInternalServerError, "Couldn't get app from context")
		}

		redirectToRegister := func() {
			q := r.URL.Query()
			q.Set(FormKeyAppID, app.ID)
			q.Set(scopesKey, scopesJSON)
			q.Set(callbackURLKey, callbackURL)
			r.URL.RawQuery = q.Encode()

			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusFound)
		}

		redirectToLogin := func() {
			r.URL.Path = "login"

			q := r.URL.Query()
			q.Set(FormKeyAppID, app.ID)
			q.Set(scopesKey, scopesJSON)
			q.Set(callbackURLKey, callbackURL)
			r.URL.RawQuery = q.Encode()

			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusFound)
		}

		if app.RegistrationForbidden {
			SetFlash(w, FlashErrorMessageKey, ErrorRegistrationForbidden.Error())
			redirectToRegister()
			return
		}

		if isAnonymous && !app.AnonymousRegistrationAllowed {
			SetFlash(w, FlashErrorMessageKey, ErrorRegistrationForbidden.Error())
			redirectToRegister()
			return
		}

		userRole := app.NewUserDefaultRole
		if inviteToken != "" {
			parsedInviteToken, err := ar.Server.Services().Token.Parse(inviteToken)
			if err != nil {
				ar.Logger.Printf("Error: Invalid invite token %s", inviteToken)
				http.Redirect(w, r, errorPath, http.StatusFound)
				return
			}

			role, ok := parsedInviteToken.Payload()["role"].(string)
			if ok {
				userRole = role
			}
		}

		// Authorize user if the app requires authorization.
		azi := authorization.AuthzInfo{
			App:         app,
			UserRole:    userRole,
			ResourceURI: r.RequestURI,
			Method:      r.Method,
		}

		if err := ar.Authorizer.Authorize(azi); err != nil {
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
		user, err := ar.Server.Storages().User.AddUserByNameAndPassword(username, password, userRole, isAnonymous)
		if err != nil {
			if err == model.ErrorUserExists {
				SetFlash(w, FlashErrorMessageKey, err.Error())
				redirectToRegister()
				return
			}

			ar.Logger.Printf("error creating user by name and password %v", err)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		// Do login flow.
		scopes, err = ar.Server.Storages().User.RequestScopes(user.ID, scopes)
		if err != nil {
			ar.Logger.Printf("error requesting scopes %v", err)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		token, err := ar.Server.Services().Token.NewWebCookieToken(user)
		if err != nil {
			ar.Logger.Printf("error creating auth token %v", err)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		tokenString, err := ar.Server.Services().Token.String(token)
		if err != nil {
			ar.Logger.Printf("error while making a call token stringify: %v", err)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		setCookie(w, CookieKeyWebCookieToken, tokenString, int(ar.Server.Services().Token.WebCookieTokenLifespan()))
		redirectToLogin()
	}
}

// RegistrationHandler serves registration page.
func (ar *Router) RegistrationHandler() http.HandlerFunc {
	tmpl, err := ar.Server.Storages().Static.ParseTemplate(model.StaticPagesNames.Registration)
	if err != nil {
		ar.Logger.Fatalln("cannot parse registration template", err)
	}
	errorPath := path.Join(ar.PathPrefix, "/misconfiguration")

	return func(w http.ResponseWriter, r *http.Request) {
		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
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

		inviteToken := strings.TrimSpace(r.URL.Query().Get("token"))

		errorMessage, err := GetFlash(w, r, FlashErrorMessageKey)
		if err != nil {
			ar.Logger.Printf("Error: getting flash message %v", err)
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		data := map[string]interface{}{
			"Error":       errorMessage,
			"Prefix":      ar.PathPrefix,
			"Scopes":      scopesJSON,
			"CallbackUrl": strings.TrimSpace(r.URL.Query().Get(callbackURLKey)),
			"AppId":       app.ID,
			"InviteToken": inviteToken,
		}

		if err = tmpl.Execute(w, data); err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
		}
	}
}
