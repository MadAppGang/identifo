package admin

import (
	"net/http"
	"strings"
)

const (
	defaultAppSkip  = 0
	defaultAppLimit = 20
)

// FetchApps fetches apps from the database.
func (ar *Router) FetchApps() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filterStr := strings.TrimSpace(r.URL.Query().Get("search"))

		limit, skip, err := ar.parseSkipAndLimit(r, defaultAppSkip, defaultAppLimit, 0)
		if err != nil {
			ar.Error(w, ErrorWrongInput, http.StatusBadRequest, "")
			return
		}

		apps, err := ar.appStorage.FetchApps(filterStr, skip, limit)
		if err != nil {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, "")
			return
		}

		ar.ServeJSON(w, http.StatusOK, apps)
		return
	}
}
