package admin

import "net/http"

type ServerInitErrors struct {
	Errors []error `json:"errors"`
}

// return all server errors
func (ar *Router) GetServerErrors() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		response := ServerInitErrors{
			Errors: ar.server.Errors(),
		}
		ar.ServeJSON(w, locale, http.StatusOK, response)
	}
}
