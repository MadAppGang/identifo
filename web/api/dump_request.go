package api

import (
	"net/http"
	"net/http/httputil"
)

// DumpRequest logs the request.
func (ar *Router) DumpRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			ar.Logger.Println("Error dumping request:", err)
		}
		ar.Logger.Printf("Request: %s\n", string(dump))
		next.ServeHTTP(w, r)
	})
}
