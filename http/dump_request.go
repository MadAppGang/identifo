package http

import (
	"net/http"
	"net/http/httputil"

	"github.com/urfave/negroni"
)

func (ar *apiRouter) DumpRequest() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		dump, _ := httputil.DumpRequest(r, true)
		ar.logger.Printf("Request: %s\n", string(dump))
		next(rw, r)
	}
}
