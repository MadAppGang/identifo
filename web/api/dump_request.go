package api

import (
	"net/http"
	"net/http/httputil"
)

// DumpRequest logs the request.
func (ar *Router) DumpRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			ar.logger.Println("Error dumping request:", err)
		}
		ar.logger.Printf("Request: %s\n", string(dump))

		next.ServeHTTP(rw, r)
	})
}
