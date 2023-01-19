package api

import (
	"net/http"

	l "github.com/madappgang/identifo/v2/localization"
	"github.com/urfave/negroni"
)

// Config middleware return error, if server config is invalid
func (ar *Router) ConfigCheck() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		// server has errors while initialized
		// errors could be config file errors
		// or errors could be connection to services and databases errors
		if len(ar.server.Errors()) > 0 {
			ar.Error(rw, http.StatusInternalServerError, l.ErrorNativeLoginConfigErrors, ar.server.Errors())
			return
		}
		next.ServeHTTP(rw, r)
	}
}
