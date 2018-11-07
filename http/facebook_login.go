package http

import (
	"net/http"

	"github.com/madappgang/identifo/facebook"
	"github.com/madappgang/identifo/model"
)

//FacebookLogin - login/register with facebook
//The user sends the facebookID
//if there is not the user with such facebook_id, function returns 404 (user not found)
//uf register_if_new presents - function creates new user and set username/password (optional to login with email, password)
//access_token is Short-Liven facebook access token, the function will exchange this Short-Lived Token for Long-Lived Token
//refer to facebook docs about that https://developers.facebook.com/docs/facebook-login/access-tokens/refreshing/
//we are not getting user email from facebook account, because it is require to obtain `email` permission
//If your need to ask for additional permissions, you could hardcode it here, optionally it will be acdded to settings.
func (ar *apiRouter) FacebookLogin() http.HandlerFunc {

	type loginData struct {
		Username      string   `json:"username,omitempty"`
		Password      string   `json:"password,omitempty"`
		AccessToken   string   `json:"access_token,omitempty" validate:"required"`
		RegisterIfNew bool     `json:"register_if_new,omitempty"`
		Scopes        []string `json:"scopes,omitempty"`
	}

	type AuthResponse struct {
		AccessToken  string `json:"access_token,omitempty"`
		RefreshToken string `json:"refresh_token,omitempty"`
	}

	registerFacebookUser := func(d loginData, f facebook.User) (model.User, error) {
		//TODO: implement profile
		return ar.userStorage.AddUserWithSocialID(
			model.FacebookIDProvider,
			f.ID,
			d.Username,
			d.Password,
			nil,
		)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		d := loginData{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		app := appFromContext(r.Context())
		if app == nil {
			ar.logger.Println("Error getting App")
			ar.Error(w, ErrorRequestInvalidAppID, http.StatusBadRequest, "")
			return
		}

		fb := facebook.NewClient(d.AccessToken)
		fbProfile, err := fb.MyProfile()
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		//check we had `id` permissions for the access_token
		if len(fbProfile.ID) == 0 {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		user, err := ar.userStorage.UserBySocialID(model.FacebookIDProvider, fbProfile.ID)
		//check error not found, create the new user
		if err == model.ErrorNotFound && d.RegisterIfNew {
			user, err = registerFacebookUser(d, fbProfile)
			if err != nil {
				ar.Error(w, err, http.StatusInternalServerError, "")
				return
			}
		} else if err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		//request the permissions for the user
		scopes, err := ar.userStorage.RequestScopes(user.ID(), d.Scopes)
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		//generate access token
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
