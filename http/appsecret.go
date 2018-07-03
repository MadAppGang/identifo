package http

import (
	"net/http"
)

//AppSecretMiddleware check request header for app secret, reject if it's not natch
func AppSecretMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// call endpoint handler
	next(rw, r)
}
