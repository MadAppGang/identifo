package admin

import (
	"net/http"
	"time"

	"github.com/urfave/negroni"
)

// Session is a middleware to check if admin is logged in with valid cookie.
// If not, forces to login.
func (ar *Router) Session() negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if ar.isLoggedIn(w, r) {
			next(w, r)
		}
	}
}

// IsLoggedIn checks if admin is logged in.
func (ar *Router) IsLoggedIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ar.isLoggedIn(w, r)
	}
}

func (ar *Router) isLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		ar.Error(w, ErrorNotAuthorized, http.StatusUnauthorized, "")
		return false
	}

	sessionID, err := decode(cookie.Value)
	if err != nil {
		ar.Error(w, err, http.StatusInternalServerError, err.Error())
		return false
	}

	session, err := ar.sessionStorage.GetSession(sessionID)
	if err != nil {
		ar.Error(w, err, http.StatusUnauthorized, err.Error())
		return false
	}

	if session.ExpirationDate.Before(time.Now()) {
		ar.Error(w, ErrorNotAuthorized, http.StatusUnauthorized, "")
		return false
	}

	return true
}
