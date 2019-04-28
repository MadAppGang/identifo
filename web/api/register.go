package api

import (
	"net/http"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/middleware"
)

/*
 * Password rules:
 * at least 7 letters
 * at least 1 number
 * at least 1 upper case
 * at least 1 special character
 */

//RegisterWithPassword register new user with password
func (ar *Router) RegisterWithPassword() http.HandlerFunc {

	type registrationData struct {
		Username string                 `json:"username,omitempty" validate:"required,gte=6,lte=50"`
		Password string                 `json:"password,omitempty" validate:"required,gte=7,lte=50"`
		Profile  map[string]interface{} `json:"user_profile,omitempty"`
		Scopes   []string               `json:"scopes,omitempty"`
	}

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

		//parse data
		d := registrationData{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		//validate password
		if err := model.StrongPswd(d.Password); err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		//create new user
		user, err := ar.userStorage.AddUserByNameAndPassword(d.Username, d.Password, d.Profile)
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		//do login flow
		scopes, err := ar.userStorage.RequestScopes(user.ID(), d.Scopes)
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
		//requesting offline access ?
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
