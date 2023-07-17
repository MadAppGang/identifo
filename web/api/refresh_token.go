package api

import (
	"encoding/json"
	"net/http"

	"github.com/madappgang/identifo/v2/web/middleware"
)

// RefreshTokens issues new access and, if requested, refresh token for provided refresh token.
// After new tokens are issued, the old refresh token and access token gets invalidated (added to blocklist).
// We validate refresh token
// if its valid - issue new tokens.
// ! Be careful, old access token still could be accepted by some systems, if it is not yet expired and those systems are not checking blocklist (usually the should not in distributed systems).
func (ar *Router) RefreshTokens() http.HandlerFunc {
	type requestData struct {
		Scopes []string `json:"scopes,omitempty"`
		Access string   `json:"access,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		rd := requestData{}
		if err := json.NewDecoder(r.Body).Decode(&rd); err != nil {
			// Assume we have not requested any scopes,  if there is no valid data in the body
			rd = requestData{Scopes: []string{}}
		}

		app := middleware.AppFromContext(r.Context())
		// should not be empty, the middleware should have validated it
		token := middleware.TokenFromContext(r.Context())
		result, err := ar.server.Storages().UC.RefreshJWTToken(r.Context(), token, rd.Access, app, rd.Scopes)
		if err != nil {
			ar.HTTPError(w, err, http.StatusInternalServerError)
			return
		}
		ar.ServeJSON(w, locale, http.StatusOK, result)
	}
}
