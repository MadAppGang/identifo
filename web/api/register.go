package api

import (
	"fmt"
	"net/http"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/middleware"
)

type registrationData struct {
	Username string                 `json:"username,omitempty"`
	Password string                 `json:"password,omitempty"`
	Profile  map[string]interface{} `json:"user_profile,omitempty"`
	Scopes   []string               `json:"scopes,omitempty"`
}

func (rd *registrationData) validate() error {
	usernameLen := len(rd.Username)
	if usernameLen < 6 || usernameLen > 50 {
		return fmt.Errorf("Incorrect email length %d, expected a number between 6 and 50", usernameLen)
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

// RegisterWithPassword register new user with password
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
			ar.Error(w, ErrorRequestInvalidAppID, http.StatusBadRequest, "")
			return
		}

		if app.RegistrationForbidden() {
			ar.Error(w, ErrorRegistrationForbidden, http.StatusForbidden, "")
			return
		}

		// Parse registration data.
		rd := registrationData{}
		if ar.MustParseJSON(w, r, &rd) != nil {
			return
		}

		if err := rd.validate(); err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		// Validate password.
		if err := model.StrongPswd(rd.Password); err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		// Create new user.
		user, err := ar.userStorage.AddUserByNameAndPassword(rd.Username, rd.Password, rd.Profile)
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		// Do login flow.
		scopes, err := ar.userStorage.RequestScopes(user.ID(), rd.Scopes)
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		token, err := ar.tokenService.NewToken(user, scopes, app)
		if err != nil {
			ar.Error(w, err, http.StatusUnauthorized, "")
			return
		}

		tokenString, err := ar.tokenService.String(token)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		refreshString := ""
		// Requesting offline access?
		if contains(scopes, model.OfflineScope) {
			refresh, err := ar.tokenService.NewRefreshToken(user, scopes, app)
			if err != nil {
				ar.Error(w, err, http.StatusInternalServerError, "")
				return
			}
			refreshString, err = ar.tokenService.String(refresh)
			if err != nil {
				ar.Error(w, err, http.StatusInternalServerError, "")
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
