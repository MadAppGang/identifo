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
			ar.Error(w, ErrorRequestInvalidAppID, http.StatusBadRequest, "")
			return
		}

		token := tokenFromContext(r.Context())

		accessToken, err := ar.tokenService.RefreshToken(token)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		tokenStr, err := ar.tokenService.String(accessToken)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}
		ar.ServeJSON(w, http.StatusOK, responseData{tokenStr})
	}
}
