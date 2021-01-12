package api

import (
	"fmt"
	"net/http"
	"strings"

	jwtService "github.com/madappgang/identifo/jwt/service"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/plugin/shared"
	"github.com/madappgang/identifo/proto"
	"github.com/madappgang/identifo/proto/extensions"
	"github.com/madappgang/identifo/web/authorization"
	"github.com/madappgang/identifo/web/middleware"
)

// FederatedLoginData represents federated login input data.
type FederatedLoginData struct {
	FederatedIDProvider string   `json:"provider,omitempty" validate:"required"`
	AccessToken         string   `json:"access_token,omitempty"`
	RegisterIfNew       bool     `json:"register_if_new,omitempty"`
	Scopes              []string `json:"scopes,omitempty"`
	AuthorizationCode   string   `json:"authorization_code,omitempty"` // Specific for Sign In with Apple.
}

// FederatedLogin provides login/registration with federated identity.
// First, user sends the identity provider access token to Identifo.
// Then, Identifo sends request to identity provider to get user profile and identity user ID,
// and then search for the user with this federated identity ID in the user pool.
// If there is no user with such identity, function returns 404 (user not found).
// If register_if_new presents - function creates new user without username/password,
// there is a dedicated endpoint to link username/password to federated account.
func (ar *Router) FederatedLogin() http.HandlerFunc {
	var federatedProviders = map[string]bool{
		strings.ToLower(string(model.FacebookIDProvider)): true,
		strings.ToLower(string(model.AppleIDProvider)):    true,
		strings.ToLower(string(model.GoogleIDProvider)):   false, //TODO: add later
		strings.ToLower(string(model.TwitterIDProvider)):  false, //TODO: add later
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if !ar.SupportedLoginWays.Federated {
			ar.Error(w, ErrorAPIAppFederatedLoginNotSupported, http.StatusBadRequest, "Application does not support federated login", "FederatedLogin.supportedLoginWays")
			return
		}

		d := FederatedLoginData{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		if !federatedProviders[strings.ToLower(d.FederatedIDProvider)] {
			ar.logger.Println("Federated provider is not supported:", d.FederatedIDProvider)
			ar.Error(w, ErrorAPIAppFederatedProviderNotSupported, http.StatusBadRequest, fmt.Sprintf("UnsupportedProvider: %v", d.FederatedIDProvider), "FederatedLogin.federatedProviders[]")
			return
		}

		app := middleware.AppFromContext(r.Context())
		if app == nil {
			ar.logger.Println("Error getting App")
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App id is not specified.", "FederatedLogin.AppFromContext")
			return
		}

		var federatedID string
		var err error

		fidName := strings.ToUpper(d.FederatedIDProvider)
		fidInt32 := proto.FederatedIdentityProvider_value[fidName]
		fid := proto.FederatedIdentityProvider(fidInt32)

		switch fid {
		case proto.FederatedIdentityProvider_FACEBOOK:
			federatedID, err = ar.FacebookUserID(d.AccessToken)
		case proto.FederatedIdentityProvider_APPLE:
			if app.AppleInfo() == nil {
				ar.logger.Println("Empty apple info")
				ar.Error(w, ErrorAPIAppFederatedProviderEmptyAppleInfo, http.StatusBadRequest, "App does not have Apple info.", "FederatedLogin.switch_providers_apple")
				return
			}
			federatedID, err = ar.AppleUserID(d.AuthorizationCode, app.AppleInfo())
		default:
			ar.Error(w, ErrorAPIAppFederatedProviderNotSupported, http.StatusBadRequest, fmt.Sprintf("UnsupportedProvider: %v", fid), "FederatedLogin.switch_providers_default")
			return
		}

		if err != nil {
			ar.logger.Println("Error getting federated user ID:", err)
			ar.Error(w, ErrorAPIAppFederatedProviderEmptyUserID, http.StatusBadRequest, err.Error(), "FederatedLogin.switch_providers.err")
			return
		}

		user, err := ar.userStorage.UserByFederatedID(fid, federatedID)
		// Check error not found, create new user.
		if err == shared.ErrUserNotFound && d.RegisterIfNew {
			user, err = ar.userStorage.AddUserWithFederatedID(fid, federatedID, app.NewUserDefaultRole())
			if err != nil {
				ar.Error(w, ErrorAPIUserUnableToCreate, http.StatusInternalServerError, err.Error(), "FederatedLogin.UserByFederatedID.RegisterNew")
				return
			}
		} else if err == shared.ErrUserNotFound && !d.RegisterIfNew {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusNotFound, err.Error(), "FederatedLogin.UserByFederatedID.NotRegisterNew")
			return
		} else if err != nil {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusInternalServerError, err.Error(), "FederatedLogin.UserByFederatedID")
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
			ar.Error(w, ErrorAPIAppAccessDenied, http.StatusForbidden, err.Error(), "FederatedLogin.Authorizer")
			return
		}

		// Request permissions for the user.
		scopes, err := ar.userStorage.RequestScopes(user.Id, d.Scopes)
		if err != nil {
			ar.Error(w, ErrorAPIRequestScopesForbidden, http.StatusBadRequest, err.Error(), "FederatedLogin.RequestScopes")
			return
		}

		// Generate access token.
		token, err := ar.tokenService.NewAccessToken(user, scopes, app, false)
		if err != nil {
			ar.Error(w, ErrorAPIAppAccessTokenNotCreated, http.StatusUnauthorized, err.Error(), "FederatedLogin.tokenService_NewToken")
			return
		}
		tokenString, err := ar.tokenService.String(token)
		if err != nil {
			ar.Error(w, ErrorAPIAppAccessTokenNotCreated, http.StatusInternalServerError, err.Error(), "FederatedLogin.tokenService_String")
			return
		}

		refreshString := ""
		//requesting offline access ?
		if contains(scopes, jwtService.OfflineScope) {
			refresh, err := ar.tokenService.NewRefreshToken(user, scopes, app)
			if err != nil {
				ar.Error(w, ErrorAPIAppRefreshTokenNotCreated, http.StatusInternalServerError, err.Error(), "FederatedLogin.tokenService_NewRefreshToken")
				return
			}
			refreshString, err = ar.tokenService.String(refresh)
			if err != nil {
				ar.Error(w, ErrorAPIAppRefreshTokenNotCreated, http.StatusInternalServerError, err.Error(), "FederatedLogin.tokenService_String")
				return
			}
		}

		extensions.SanitizeUser(user)
		result := AuthResponse{
			AccessToken:  tokenString,
			RefreshToken: refreshString,
			User:         user,
		}

		ar.userStorage.UpdateLoginMetadata(user.Id)
		ar.ServeJSON(w, http.StatusOK, result)
	}

}
