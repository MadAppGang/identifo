package api

import (
	"net/http"
	"net/http/httputil"

	"github.com/urfave/negroni"
)

// DumpRequest logs the request.
func (ar *Router) DumpRequest() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			ar.logger.Println("Error dumping request:", err)
		}
		ar.logger.Printf("Request: %s\n", string(dump))
		next(rw, r)
	}
}
