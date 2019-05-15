package api

import (
	"net/http"

	"github.com/madappgang/identifo/web/middleware"
)

// RefreshToken - refresh access token
func (ar *Router) RefreshToken() http.HandlerFunc {
	type responseData struct {
		AccessToken string `json:"access_token"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		app := middleware.AppFromContext(r.Context())
		if app == nil {
			ar.logger.Println("Error getting App")
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "app id is absent in header params", "RefreshToken.AppFromContext")
			return
		}

		token := tokenFromContext(r.Context())

		accessToken, err := ar.tokenService.RefreshToken(token)
		if err != nil {
			ar.Error(w, ErrorAPIAppAccessTokenNotCreated, http.StatusInternalServerError, err.Error(), "RefreshToken.RefreshToken")
			return
		}

		tokenStr, err := ar.tokenService.String(accessToken)
		if err != nil {
			ar.Error(w, ErrorAPIAppAccessTokenNotCreated, http.StatusInternalServerError, err.Error(), "RefreshToken.tokenService_string")
			return
		}
		ar.ServeJSON(w, http.StatusOK, responseData{tokenStr})
	}
}
