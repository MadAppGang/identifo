package middleware

import (
	"net/http"
	"strings"

	"github.com/urfave/negroni"
)

func RemoveTrailingSlash() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		next.ServeHTTP(rw, r)
	}
}
