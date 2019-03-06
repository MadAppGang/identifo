package api

import (
	"net/http"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/middleware"
)

//LoginWithPassword - login user with email and password
func (ar *Router) LoginWithPassword() http.HandlerFunc {

	type loginData struct {
		Username    string   `json:"username,omitempty" validate:"required,gte=6,lte=130"`
		Password    string   `json:"password,omitempty" validate:"required,gte=6,lte=130"`
		DeviceToken string   `json:"device_token,omitempty"`
		Scopes      []string `json:"scopes,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		d := loginData{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		user, err := ar.userStorage.UserByNamePassword(d.Username, d.Password)
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		scopes, err := ar.userStorage.RequestScopes(user.ID(), d.Scopes)
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		app := middleware.AppFromContext(r.Context())
		if app == nil {
			ar.logger.Println("Error getting App")
			ar.Error(w, ErrorRequestInvalidAppID, http.StatusBadRequest, "")
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

		result := AuthResponse{
			AccessToken:  tokenString,
			RefreshToken: refreshString,
		}
		ar.ServeJSON(w, http.StatusOK, result)
	}
}
