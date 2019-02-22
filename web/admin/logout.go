package admin

import (
	"net/http"
)

// Logout logs admin out.
func (ar *Router) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			switch err {
			case http.ErrNoCookie:
				ar.logger.Println("No cookie")
			default:
				ar.Error(w, err, http.StatusInternalServerError, "")
			}
			return
		}

		sessionID, err := decode(cookie.Value)
		if err != nil {
			ar.Error(w, ErrorRequestInvalidCookie, http.StatusBadRequest, "")
			return
		}
		if err := ar.sessionStorage.DeleteSession(sessionID); err != nil {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, "")
			return
		}

		c := &http.Cookie{
			Name:     cookieName,
			Value:    "",
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(w, c)
		return
	}
}
