package admin

import (
	"net/http"
)

const (
	defaultUserSkip  = 0
	defaultUserLimit = 20
)

// FetchUsers fetches users from the database.
func (ar *Router) FetchUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, skip, err := ar.parseSkipAndLimit(r, defaultUserSkip, defaultUserLimit, 0)

		users, err := ar.userStorage.FetchUsers()
		if err != nil {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, "")
			return
		}

		if err = ar.sessionStorage.InsertSession(session); err != nil {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, "")
			return
		}

		c := &http.Cookie{
			Name:     cookieName,
			Value:    encode(session.ID),
			MaxAge:   ar.sessionService.SessionDurationSeconds(),
			HttpOnly: true,
		}
		http.SetCookie(w, c)
		return
	}
}
