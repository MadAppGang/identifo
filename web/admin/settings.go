package admin

import (
	"net/http"
)

type databaseConnection struct {
	Type     string `json:"type,omitempty"`
	Region   string `json:"region,omitempty"`
	Name     string `json:"name,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
}

// FetchDatabaseSettings provides info about used database engine.
func (ar *Router) FetchDatabaseSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dcdata := databaseConnection{
			Type:     ar.DBType,
			Region:   ar.DBRegion,
			Name:     ar.DBName,
			Endpoint: ar.DBEndpoint,
		}

		ar.ServeJSON(w, http.StatusOK, dcdata)
		return
	}
}
