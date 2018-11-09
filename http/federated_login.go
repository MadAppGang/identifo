package http

import (
	"net/http"
	"strings"

	"github.com/madappgang/identifo/model"
)

//FederatedLoginData represents federated login input data
type FederatedLoginData struct {
	FederatedIDProvider string   `json:"provider,omitempty" validate:"required"`
	AccessToken         string   `json:"access_token,omitempty" validate:"required"`
	RegisterIfNew       bool     `json:"register_if_new,omitempty"`
	Scopes              []string `json:"scopes,omitempty"`
}

//AuthResponse is response with successful auth data
type AuthResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

//FederatedLogin - login/register with federated Identity
//The user sends the Identity provider access_token
//Identifo is sendind request to identity provider to get user_profile and identity user ID
//and then searching the user with this federated identity ID in the user pool
//if there is not the user with such identity, function returns 404 (user not found)
//uf register_if_new presents - function creates new user with no  username/password, there is dedicated endpoint to link username/password with federated account
func (ar *apiRouter) FederatedLogin() http.HandlerFunc {

	var federatedProviders = map[string]bool{
		strings.ToLower(string(model.FacebookIDProvider)): true,
		strings.ToLower(string(model.GoogleIDProvider)):   false, //TODO: add later
		strings.ToLower(string(model.TwitterIDProvider)):  false, //TODO: add later
	}

	return func(w http.ResponseWriter, r *http.Request) {
		d := FederatedLoginData{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		if !federatedProviders[strings.ToLower(d.FederatedIDProvider)] {
			ar.logger.Println("Federated provider is not supported:", d.FederatedIDProvider)
			ar.Error(w, ErrorFederatedProviderIsNotSupported, http.StatusBadRequest, "")
			return
		}

		app := appFromContext(r.Context())
		if app == nil {
			ar.logger.Println("Error getting App")
			ar.Error(w, ErrorRequestInvalidAppID, http.StatusBadRequest, "")
			return
		}

		var federatedID string
		var err error
		fid := model.FederatedIdentityProvider(strings.ToUpper(d.FederatedIDProvider))
		switch fid {
		case model.FacebookIDProvider:
			federatedID, err = ar.FacebookUserID(d.AccessToken)
		default:
			err = ErrorFederatedProviderIsNotSupported
		}

		if err != nil {
			ar.logger.Println("Error getting federated user ID:", err)
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		user, err := ar.userStorage.UserByFederatedID(fid, federatedID)
		//check error not found, create the new user
		if err == model.ErrorNotFound && d.RegisterIfNew {
			user, err = ar.userStorage.AddUserWithFederatedID(fid, federatedID)
			if err != nil {
				ar.Error(w, err, http.StatusInternalServerError, "")
				return
			}
		} else if err == model.ErrorNotFound && !d.RegisterIfNew {
			ar.Error(w, err, http.StatusNotFound, "")
			return
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
