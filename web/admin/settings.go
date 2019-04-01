package admin

import (
	"fmt"
	"net/http"
)

// FetchDatabaseSettings provides info about used database engine.
func (ar *Router) FetchDatabaseSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ar.ServeJSON(w, http.StatusOK, ar.ServerSettings)
		return
	}
}

// AlterDatabaseSettings changes database settings.
func (ar *Router) AlterDatabaseSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newdbset := new(databaseSettings)
		if ar.mustParseJSON(w, r, newdbset) != nil {
			return
		}

		if newdbset.Type != ar.DBType {
			ar.Error(w, fmt.Errorf("Database type %s does not match the current one", newdbset.Type), http.StatusBadRequest, "")
			return
		}

		dbset := databaseSettings{
			Type:     ar.DBType,
			Region:   ar.DBRegion,
			Name:     ar.DBName,
			Endpoint: ar.DBEndpoint,
		}

		ar.ServeJSON(w, http.StatusOK, dbset)
		return
	}
}
