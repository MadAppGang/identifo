package admin

import (
	"net/http"
	"time"

	"github.com/urfave/negroni"
)

// Session is a middleware to check if admin is logged in with valid cookie.
// If all checks succeeded, prolongs existing session.
// If not, forces to login.
func (ar *Router) Session() negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			ar.Error(w, ErrorNotAuthorized, http.StatusUnauthorized, "")
			return
		}

		sessionID, err := decode(cookie.Value)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, err.Error())
			return
		}

		session, err := ar.sessionStorage.GetSession(sessionID)
		if err != nil {
			ar.Error(w, err, http.StatusUnauthorized, err.Error())
			return
		}

		if session.ExpirationDate.Before(time.Now()) {
			ar.Error(w, ErrorNotAuthorized, http.StatusUnauthorized, "")
			return
		}

		ar.prolongSession(w, sessionID)

		next(w, r)
	}
}

func (ar *Router) prolongSession(w http.ResponseWriter, sessionID string) {
	if err := ar.sessionService.ProlongSession(sessionID); err != nil {
		ar.logger.Println("Error prolonging session:", err)
		return
	}
	c := &http.Cookie{
		Name:     cookieName,
		Value:    encode(sessionID),
		MaxAge:   ar.sessionService.SessionDurationSeconds(),
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}
