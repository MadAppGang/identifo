package admin

import (
	"net/http"

	"github.com/madappgang/identifo/model"
)

// GetApp fetches app by ID from the database.
func (ar *Router) FederatedProvidersList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ar.ServeJSON(w, http.StatusOK, model.FederatedProviders)
	}
}
