package api

import (
	"fmt"
	"net/http"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/middleware"
)

type loginData struct {
	Username    string   `json:"username,omitempty"`
	Password    string   `json:"password,omitempty"`
	DeviceToken string   `json:"device_token,omitempty"`
	Scopes      []string `json:"scopes,omitempty"`
}

func (ld *loginData) validate() error {
	usernameLen := len(ld.Username)
	if usernameLen < 6 || usernameLen > 130 {
		return fmt.Errorf("Incorrect username length %d, expected a number between 6 and 130", usernameLen)
	}
	pswdLen := len(ld.Password)
	if pswdLen < 6 || pswdLen > 130 {
		return fmt.Errorf("Incorrect password length %d, expected a number between 6 and 130", pswdLen)
	}
	return nil
}

// LoginWithPassword logs user in with email and password.
func (ar *Router) LoginWithPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ld := loginData{}
		if ar.MustParseJSON(w, r, &ld) != nil {
			return
		}

		if err := ld.validate(); err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		user, err := ar.userStorage.UserByNamePassword(ld.Username, ld.Password)
		if err != nil {
			ar.Error(w, ErrorAPIRequestIncorrectEmailOrPassword, http.StatusUnauthorized, err.Error(), "LoginWithPassword.UserByNamePassword")
			return
		}

		scopes, err := ar.userStorage.RequestScopes(user.ID(), ld.Scopes)
		if err != nil {
			ar.Error(w, ErrorAPIRequestScopesForbidden, http.StatusForbidden, err.Error(), "LoginWithPassword.RequestScopes")
			return
		}

		app := middleware.AppFromContext(r.Context())
		if app == nil {
			ar.logger.Println("Error getting App")
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App is not in context.", "LoginWithPassword.AppFromContext")
			return
		}

		token, err := ar.tokenService.NewToken(user, scopes, app)
		if err != nil {
			ar.Error(w, ErrorAPIAppAccessTokenNotCreated, http.StatusForbidden, err.Error(), "LoginWithPassword.tokenService_NewToken")
			return
		}
		tokenString, err := ar.tokenService.String(token)
		if err != nil {
			ar.Error(w, ErrorAPIAppAccessTokenNotCreated, http.StatusInternalServerError, err.Error(), "LoginWithPassword.tokenService_String")
			return
		}

		refreshString := ""
		// requesting offline access ?
		if contains(scopes, model.OfflineScope) {
			refresh, err := ar.tokenService.NewRefreshToken(user, scopes, app)
			if err != nil {
				ar.Error(w, ErrorAPIAppRefreshTokenNotCreated, http.StatusInternalServerError, err.Error(), "LoginWithPassword.tokenService_NewRefreshToken")
				return
			}
			refreshString, err = ar.tokenService.String(refresh)
			if err != nil {
				ar.Error(w, ErrorAPIAppRefreshTokenNotCreated, http.StatusInternalServerError, err.Error(), "LoginWithPassword.tokenService_String")
				return
			}
		}

		result := AuthResponse{
			AccessToken:  tokenString,
			RefreshToken: refreshString,
		}
		ar.ServeJSON(w, http.StatusOK, result)
	}
}
