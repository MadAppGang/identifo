package html

import (
	"net/http"
	"net/http/httputil"

	"github.com/urfave/negroni"
)

// DumpRequest dumps request to logger.
func (ar *Router) DumpRequest() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			ar.Logger.Println("Error dumping request:", err)
		}
		ar.Logger.Printf("Request: %s\n", string(dump))
		next(rw, r)
	}
}
