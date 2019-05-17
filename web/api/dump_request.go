package api

import (
	"net/http"
	"net/http/httputil"

	"github.com/urfave/negroni"
)

// DumpRequest logs the request.
func (ar *Router) DumpRequest() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		dump, _ := httputil.DumpRequest(r, true)
		ar.logger.Printf("Request: %s\n", string(dump))
		next(rw, r)
	}
}
