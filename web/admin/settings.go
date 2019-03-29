package admin

import (
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
