package api

import (
	"fmt"
	"net/http"
	"time"

	jwtService "github.com/madappgang/identifo/jwt/service"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/middleware"
	"github.com/xlzd/gotp"
)

type loginData struct {
	Username    string   `json:"username,omitempty"`
	Password    string   `json:"password,omitempty"`
	DeviceToken string   `json:"device_token,omitempty"`
	TFACode     string   `json:"tfa_code,omitempty"`
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

// LoginWithPassword logs user in with username and password.
func (ar *Router) LoginWithPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !ar.SupportedLoginWays.Username {
			ar.Error(w, ErrorAPIAppLoginWithUsernameNotSupported, http.StatusBadRequest, "Application does not support login with username", "LoginWithPassword.supportedLoginWays")
			return
		}

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

		// Execute two-factor auth if user enabled it.
		if user.TFAInfo().IsEnabled {
			if err = ar.execute2FA(w, ld, user.TFAInfo().Secret); err != nil {
				return
			}
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

		ar.userStorage.UpdateLoginMetadata(user.ID())
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

func (ar *Router) execute2FA(w http.ResponseWriter, ld loginData, secret string) error {
	if len(ld.TFACode) == 0 {
		err := fmt.Errorf("Empty 2FA code")
		ar.Error(w, ErrorAPIRequestTFACodeEmpty, http.StatusBadRequest, err.Error(), "LoginWithPassword.execute2FA")
		return err
	}

	totp := gotp.NewDefaultTOTP(secret)
	if verified := totp.Verify(ld.TFACode, int(time.Now().Unix())); !verified {
		err := fmt.Errorf("Invalid one-time password")
		ar.Error(w, ErrorAPIRequestTFACodeInvalid, http.StatusUnauthorized, err.Error(), "LoginWithPassword.execute2FA")
		return err
	}
	return nil
}

