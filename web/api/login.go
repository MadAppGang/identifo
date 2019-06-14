package api

import (
	"fmt"
	"net/http"

	jwtService "github.com/madappgang/identifo/jwt/service"
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
			ar.Error(w, ErrorAPIRequestBodyParamsInvalid, http.StatusBadRequest, err.Error(), "LoginWithPassword.validate")
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

		offline := contains(scopes, jwtService.OfflineScope)
		accessToken, refreshToken, err := ar.loginUser(user, scopes, app, offline)
		if err != nil {
			ar.Error(w, ErrorAPIAppAccessTokenNotCreated, http.StatusInternalServerError, err.Error(), "LoginWithPassword.loginUser")
			return
		}

		result := AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		go ar.userStorage.UpdateLoginMetadata(user.ID())
		ar.ServeJSON(w, http.StatusOK, result)
	}
}

// loginUser function creates token for user session.
// refreshToken boolean param tells if function should return refresh token too.
func (ar *Router) loginUser(user model.User, scopes []string, app model.AppData, createRefreshToken bool) (accessTokenString, refreshTokenString string, err error) {
	token, err := ar.tokenService.NewToken(user, scopes, app)
	if err != nil {
		return
	}
	accessTokenString, err = ar.tokenService.String(token)
	if err != nil {
		return
	}
	if !createRefreshToken {
		return
	}
	refresh, err := ar.tokenService.NewRefreshToken(user, scopes, app)
	if err != nil {
		return
	}
	refreshTokenString, err = ar.tokenService.String(refresh)
	if err != nil {
		return
	}
	return
}
