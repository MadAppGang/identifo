package admin

import (
	"net/http"

	"github.com/madappgang/identifo/v2/model"
)

// GetApp fetches app by ID from the database.
func (ar *Router) FederatedProvidersList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")
		ar.ServeJSON(w, locale, http.StatusOK, model.FederatedProviders)
	}
}
