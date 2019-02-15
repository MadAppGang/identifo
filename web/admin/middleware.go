package admin

import (
	"net/http"
	"time"

	"github.com/urfave/negroni"
)

const (
	cookieName = "SessionID"
)

// Session is a middleware to check if admin is logged in with valid cookie.
// If not, forces to login.
func (ar *Router) Session() negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			http.Redirect(w, r, ar.RedirectURL, http.StatusMovedPermanently)
			return
		}

		sessionID, err := decode(cookie.Value)
		if err != nil {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, err.Error())
			return
		}

		session, err := ar.sessionStorage.GetSession(sessionID)

		if session.ExpirationDate.Before(time.Now()) {
			ar.Error(w, ErrorNotAuthorized, http.StatusUnauthorized, "")
			return
		}

		next(w, r)
	}
}
