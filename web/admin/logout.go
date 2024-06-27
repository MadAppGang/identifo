package admin

import (
	"net/http"
)

// Logout logs admin out.
func (ar *Router) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := &http.Cookie{
			Name:     cookieName,
			Value:    "",
			MaxAge:   0,
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(w, c)

		cookie, err := r.Cookie(cookieName)
		if err != nil {
			switch err {
			case http.ErrNoCookie:
				ar.logger.Warn("No cookie during logout")
				ar.ServeJSON(w, http.StatusOK, nil)
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
		if err := ar.server.Storages().Session.DeleteSession(sessionID); err != nil {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, "")
			return
		}

		ar.ServeJSON(w, http.StatusOK, nil)
	}
}
