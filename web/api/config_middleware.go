package api

import (
	"net/http"

	"github.com/madappgang/identifo/v2/l"
)

func (ar *Router) ConfigCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		// server has errors while initialized
		// errors could be config file errors
		// or errors could be connection to services and databases errors
		if len(ar.server.Errors()) > 0 {
			err := l.LocalizedError{
				ErrID:   l.ErrorNativeLoginConfigErrors,
				Details: []any{ar.server.Errors()},
				Locale:  locale,
			}
			ar.HTTPError(w, err, http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
