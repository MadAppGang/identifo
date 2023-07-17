package middleware

import (
	"net/http"
	"strconv"
	"time"
)

func NewCacheDisable(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Etag", strconv.Itoa(int(time.Now().Unix())))
		next.ServeHTTP(w, r)
	})
}
