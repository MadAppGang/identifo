package admin

import (
	"fmt"
	"net/http"
)

type databaseSettings struct {
	Type     string `json:"type"`
	Region   string `json:"region"`
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
}

// FetchDatabaseSettings provides info about used database engine.
func (ar *Router) FetchDatabaseSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbset := databaseSettings{
			Type:     ar.DBType,
			Region:   ar.DBRegion,
			Name:     ar.DBName,
			Endpoint: ar.DBEndpoint,
		}

		switch ar.DBType {
		case "boltdb":
			dbset.Region = ""
			dbset.Name = ""
			dbset.Endpoint = ""
		case "mongodb":
			dbset.Endpoint = ""
		case "dynamodb":
			dbset.Name = ""
		}

		ar.ServeJSON(w, http.StatusOK, dbset)
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
