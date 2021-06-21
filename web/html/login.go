package html

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"

	jwtValidator "github.com/madappgang/identifo/jwt/validator"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/authorization"
	"github.com/madappgang/identifo/web/middleware"
)

const (
	usernameKey    = "email"
	passwordKey    = "password"
	scopesKey      = "scopes"
	isAnonymousKey = "anonymous"
	callbackURLKey = "callbackUrl"
	redirectURIKey = "redirectUri"
	inviteTokenKey = "inviteToken"
)

// Login logs user in with email and password.
func (ar *Router) Login() http.HandlerFunc {
	errorPath := path.Join(ar.PathPrefix, "/misconfiguration")

	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue(usernameKey)
		password := r.FormValue(passwordKey)
		scopesJSON := r.FormValue(scopesKey)
		callbackURL := r.FormValue(callbackURLKey)
		scopes := []string{}
		app := middleware.AppFromContext(r.Context())

		redirectToLogin := func() {
			q := r.URL.Query()
			q.Set(FormKeyAppID, app.ID)
			q.Set(scopesKey, scopesJSON)
			q.Set(callbackURLKey, callbackURL)
			r.URL.RawQuery = q.Encode()

			http.Redirect(w, r, path.Join(ar.PathPrefix, r.URL.String()), http.StatusFound)
		}

		if err := json.Unmarshal([]byte(scopesJSON), &scopes); err != nil {
			ar.Logger.Printf("invalid scopes %v", scopesJSON)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		user, err := ar.UserStorage.UserByNamePassword(username, password)
		if err != nil {
			SetFlash(w, FlashErrorMessageKey, "invalid Username or Password")
			redirectToLogin()
			return
		}

		if _, err = ar.UserStorage.RequestScopes(user.ID, scopes); err != nil {
			ar.Logger.Printf("invalid scopes %v for userID: %v", scopes, user.ID)
			http.Redirect(w, r, errorPath, http.StatusFound)
			redirectToLogin()
			return
		}

		// Authorize user if the app requires authorization.
		azi := authorization.AuthzInfo{
			App:         app,
			UserRole:    user.AccessRole,
			ResourceURI: r.RequestURI,
			Method:      r.Method,
		}

		if err := ar.Authorizer.Authorize(azi); err != nil {
			SetFlash(w, FlashErrorMessageKey, err.Error())
			redirectToLogin()
			return
		}

		token, err := ar.TokenService.NewWebCookieToken(user)
		if err != nil {
			ar.Logger.Printf("error creating auth token %v", err)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		tokenString, err := ar.TokenService.String(token)
		if err != nil {
			ar.Logger.Printf("error making a call to stringify the token: %v", err)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		ar.UserStorage.UpdateLoginMetadata(user.ID)
		setCookie(w, CookieKeyWebCookieToken, tokenString, int(ar.TokenService.WebCookieTokenLifespan()))
		redirectToLogin()
	}
}

// LoginHandler serves login page or redirects to the callback_url if user is already authenticated.
func (ar *Router) LoginHandler() http.HandlerFunc {
	tmpl, err := ar.staticFilesStorage.ParseTemplate(model.StaticPagesNames.Login)
	if err != nil {
		ar.Logger.Fatalln("Cannot parse Login template.", err)
	}
	errorPath := path.Join(ar.PathPrefix, "/misconfiguration")
	tokenValidator := jwtValidator.NewValidator(
		[]string{"identifo"},
		[]string{ar.TokenService.Issuer()},
		[]string{},
		[]string{model.TokenTypeWebCookie},
	)

	return func(w http.ResponseWriter, r *http.Request) {
		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.Logger.Printf("Error: App not found.")
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		scopesJSON := strings.TrimSpace(r.URL.Query().Get(scopesKey))
		scopes := []string{}
		if err := json.Unmarshal([]byte(scopesJSON), &scopes); err != nil {
			ar.Logger.Printf("Error: Invalid scopes %v", scopesJSON)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		callbackURL := strings.TrimSpace(r.URL.Query().Get(callbackURLKey))
		if !contains(app.RedirectURLs, callbackURL) {
			ar.Logger.Printf("Unauthorized redirect url %v for app %v", callbackURL, app.ID)
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
				"Error":       errorMessage,
				"Prefix":      ar.PathPrefix,
				"Scopes":      scopesJSON,
				"CallbackURL": callbackURL,
				"AppId":       app.ID,
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

		// TODO: Add TFA support.
		token, err := ar.TokenService.NewAccessToken(user, scopes, app, false, nil)
		if err != nil {
			ar.Logger.Printf("Error creating token: %v", err)
			serveTemplate()
			return
		}

		tokenString, err := ar.TokenService.String(token)
		if err != nil {
			ar.Logger.Printf("Error stringify token: %v", err)
			serveTemplate()
			return
		}

		refreshString := ""

		// If requesting offline access then generate and set refreshString
		if contains(scopes, model.OfflineScope) {
			refresh, err := ar.TokenService.NewRefreshToken(user, scopes, app)
			if err != nil {
				ar.Logger.Printf("Error creating refresh token: %v", err)
				serveTemplate()
				return
			}
			refreshString, err = ar.TokenService.String(refresh)
			if err != nil {
				ar.Logger.Printf("Error stringify refresh token: %v", err)
				serveTemplate()
				return
			}

		}

		redirectURL := callbackURL + "?token=" + tokenString

		if refreshString != "" {
			redirectURL = redirectURL + "&refresh_token=" + refreshString
		}

		http.Redirect(w, r, redirectURL, http.StatusFound)
	}
}
