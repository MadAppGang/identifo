package api

import (
	"fmt"
	"net/http"

	jwtService "github.com/madappgang/identifo/jwt/service"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/authorization"
	"github.com/madappgang/identifo/web/middleware"
)

type registrationData struct {
	Username  string   `json:"username,omitempty"`
	Password  string   `json:"password,omitempty"`
	Scopes    []string `json:"scopes,omitempty"`
	Anonymous bool     `json:"anonymous,omitempty"`
}

func (rd *registrationData) validate() error {
	usernameLen := len(rd.Username)
	if usernameLen < 6 || usernameLen > 50 {
		return fmt.Errorf("Incorrect username length %d, expected a number between 6 and 50", usernameLen)
	}
	pswdLen := len(rd.Password)
	if pswdLen < 6 || pswdLen > 50 {
		return fmt.Errorf("Incorrect password length %d, expected a number between 6 and 50", pswdLen)
	}
	return nil
}

/*
 * Password rules:
 * at least 6 letters
 * at least 1 upper case
 */

// RegisterWithPassword registers new user with password.
func (ar *Router) RegisterWithPassword() http.HandlerFunc {
	type registrationResponse struct {
		AccessToken  string     `json:"access_token,omitempty"`
		RefreshToken string     `json:"refresh_token,omitempty"`
		User         model.User `json:"user,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		app := middleware.AppFromContext(r.Context())
		if app == nil {
			ar.logger.Println("Error getting App")
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App is not in context.", "RegisterWithPassword.AppFromContext")
			return
		}

		if app.RegistrationForbidden() {
			ar.Error(w, ErrorAPIAppRegistrationForbidden, http.StatusForbidden, "Registration is forbidden in app.", "RegisterWithPassword.RegistrationForbidden")
			return
		}

		// Check if it makes sense to create new user.
		azi := authorization.AuthzInfo{
			App:         app,
			UserRole:    app.NewUserDefaultRole(),
			ResourceURI: r.RequestURI,
			Method:      r.Method,
		}
		if err := ar.Authorizer.Authorize(azi); err != nil {
			ar.Error(w, ErrorAPIAppAccessDenied, http.StatusForbidden, err.Error(), "RegisterWithPassword.Authorizer")
			return
		}

		// Parse registration data.
		rd := registrationData{}
		if ar.MustParseJSON(w, r, &rd) != nil {
			return
		}

		if rd.Anonymous && !app.AnonymousRegistrationAllowed() {
			ar.Error(w, ErrorAPIAppRegistrationForbidden, http.StatusForbidden, "Anonymous login forbidden in the app", "RegisterWithPassword.AnonymousRegistrationForbidden")
			return
		}

		if err := rd.validate(); err != nil {
			ar.Error(w, ErrorAPIRequestBodyParamsInvalid, http.StatusBadRequest, err.Error(), "RegisterWithPassword.validate")
			return
		}

		// Validate password.
		if err := model.StrongPswd(rd.Password); err != nil {
			ar.Error(w, ErrorAPIRequestPasswordWeak, http.StatusBadRequest, err.Error(), "RegisterWithPassword.StrongPswd")
			return
		}

		// Create new user.
		user, err := ar.userStorage.AddUserByNameAndPassword(rd.Username, rd.Password, app.NewUserDefaultRole(), rd.Anonymous)
		if err == model.ErrorUserExists {
			ar.Error(w, ErrorAPIUsernameTaken, http.StatusBadRequest, err.Error(), "RegisterWithPassword.AddUserByNameAndPassword")
			return
		}
		if err != nil {
			ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "RegisterWithPassword.AddUserByNameAndPassword")
			return
		}

		// Do login flow.
		scopes, err := ar.userStorage.RequestScopes(user.ID(), rd.Scopes)
		if err != nil {
			ar.Error(w, ErrorAPIRequestScopesForbidden, http.StatusBadRequest, err.Error(), "RegisterWithPassword.RequestScopes")
			return
		}

		token, err := ar.tokenService.NewAccessToken(user, scopes, app, false)
		if err != nil {
			ar.Error(w, ErrorAPIAppAccessTokenNotCreated, http.StatusForbidden, err.Error(), "RegisterWithPassword.tokenService_NewToken")
			return
		}

		tokenString, err := ar.tokenService.String(token)
		if err != nil {
			ar.Error(w, ErrorAPIAppAccessTokenNotCreated, http.StatusInternalServerError, err.Error(), "RegisterWithPassword.tokenService_String")
			return
		}

		refreshString := ""
		// Requesting offline access?
		if contains(scopes, jwtService.OfflineScope) {
			refresh, err := ar.tokenService.NewRefreshToken(user, scopes, app)
			if err != nil {
				ar.Error(w, ErrorAPIAppRefreshTokenNotCreated, http.StatusInternalServerError, err.Error(), "RegisterWithPassword.tokenService_NewRefreshToken")
				return
			}
			refreshString, err = ar.tokenService.String(refresh)
			if err != nil {
				ar.Error(w, ErrorAPIAppRefreshTokenNotCreated, http.StatusInternalServerError, err.Error(), "RegisterWithPassword.tokenService_String")
				return
			}
		}

		user.Sanitize()
		result := registrationResponse{
			AccessToken:  tokenString,
			RefreshToken: refreshString,
			User:         user,
		}

		ar.ServeJSON(w, http.StatusOK, result)
	}
}
