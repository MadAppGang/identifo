package api

import (
	"encoding/json"
	"net/http"

	"github.com/madappgang/identifo/v2/jwt"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
)

// Logout logs user out and deactivates their tokens.
// add access token and refresh token to block list
// we need to detach device from user
func (ar *Router) Logout() http.HandlerFunc {
	type logoutData struct {
		RefreshToken string `json:"refresh_token,omitempty"`
		DeviceToken  string `json:"device,omitempty"`
	}

	result := map[string]string{"result": "ok"}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		access := middleware.TokenFromContext(r.Context())
		var refresh *model.JWToken

		if r.Body != http.NoBody {
			ld := logoutData{}
			json.NewDecoder(r.Body).Decode(&ld)
			if ld.RefreshToken != "" {
				refresh, _ = jwt.ParseTokenString(ld.RefreshToken)
			}
		}
		ar.server.Storages().UMC.InvalidateTokens(r.Context(), refresh, access, "user logout api request")
		// TODO: Detach device from user
		ar.ServeJSON(w, locale, http.StatusOK, result)
	}
}
