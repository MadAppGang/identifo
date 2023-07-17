package api

import (
	"net/http"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/web/middleware"
	"github.com/urfave/negroni"
	"golang.org/x/exp/slices"
)

// AppToken validates that token is valid for app.
func (ar *Router) AppToken() negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		locale := r.Header.Get("Accept-Language")

		// Get refresh token from context.
		token := middleware.TokenFromContext(r.Context())
		if token == nil {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorTokenRefreshEmpty)
			return
		}

		aud, err := token.Claims.GetAudience()
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorValidationTokenInvalidAudience)
			return
		}

		app := middleware.AppFromContext(r.Context())
		if !slices.Contains(aud, app.ID) {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorTokenIsForOtherAPP)
		}

		next.ServeHTTP(w, r)
	}
}
