package http

import (
	"net/http"
)

//RefreshToken - refresh access token
func (ar *apiRouter) RefreshToken() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		app := appFromContext(r.Context())
		if app == nil {
			ar.logger.Println("Error getting App")
			ar.Error(w, ErrorRequestInvalidAppID, http.StatusBadRequest, "")
			return
		}

		ar.ServeJSON(w, http.StatusOK, app)
	}
}
