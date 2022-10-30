package api

import (
	"fmt"
	"net/http"

	"github.com/urfave/negroni"
)

// Config middleware return error, if server config is invalid
func (ar *Router) ConfigCheck() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		// server has errors while initialized
		// errors could be config file errors
		// or errors could be connection to services and databases errors
		if len(ar.server.Errors()) > 0 {
			errs := fmt.Errorf("identifo initialized with errors: %+v", ar.server.Errors())
			ar.Error(rw, ErrorAPIServerInitializedWithErrors, http.StatusInternalServerError, errs.Error(), "API.ConfigCheck")
			return
		}
		next.ServeHTTP(rw, r)
	}
}
