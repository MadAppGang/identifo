package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/urfave/negroni"
)

func NewCacheDisable() negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Etag", strconv.Itoa(int(time.Now().Unix())))
		next.ServeHTTP(w, r)
	}
}
