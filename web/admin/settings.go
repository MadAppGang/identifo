package admin

import (
	"net/http"
)

// FetchDatabaseSettings provides info about used database engine.
func (ar *Router) FetchDatabaseSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ar.ServeJSON(w, http.StatusOK, ar.ServerSettings)
		return
	}
}
