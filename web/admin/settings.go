package admin

import (
	"net/http"
)

// FetchDatabaseSettings provides info about used database engine.
func (ar *Router) FetchDatabaseSettings() http.HandlerFunc {
	type databaseConnection struct {
		Type     string `json:"type,omitempty"`
		Region   string `json:"region,omitempty"`
		Name     string `json:"name,omitempty"`
		Endpoint string `json:"endpoint,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		dcdata := databaseConnection{
			Type:     "",
			Region:   "",
			Name:     "",
			Endpoint: r.Host,
		}

		ar.ServeJSON(w, http.StatusOK, dcdata)
		return
	}
}
