package admin

import "net/http"

// TestConnection validates different connection types, if server could connect to that.
func (ar *Router) TestConnection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
