package admin

import (
	"net/http"

	"github.com/madappgang/identifo/v2/l"
)

// Logout logs admin out.
func (ar *Router) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

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
				ar.Logger.Println("No cookie")
				ar.ServeJSON(w, locale, http.StatusOK, nil)
			default:
				ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAdminPanelMissingCookie, err.Error())
			}
			return
		}

		sessionID, err := decode(cookie.Value)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAdminPanelInvalidCookie, err.Error())
			return
		}
		if err := ar.server.Storages().Session.DeleteSession(sessionID); err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.APIInternalServerError)
			return
		}

		ar.ServeJSON(w, locale, http.StatusOK, nil)
	}
}
