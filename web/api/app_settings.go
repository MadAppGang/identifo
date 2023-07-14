package api

import (
	"net/http"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/web/middleware"
)

// GetAppSettings return app settings
func (ar *Router) GetAppSettings() http.HandlerFunc {
	// get short version of that
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")
		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPPNoAPPInContext)
			return
		}

		// TODO: Implement new app settings with field sets like users
		app.Secret = ""

		ar.ServeJSON(w, locale, http.StatusOK, app)
	}
}
