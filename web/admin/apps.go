package admin

import (
	"net/http"
	"strings"

	"github.com/madappgang/identifo/model"
)

const (
	defaultAppSkip  = 0
	defaultAppLimit = 20
)

// GetApp fetches app by ID from the database.
func (ar *Router) GetApp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appID := getRouteVar("id", r)

		app, err := ar.appStorage.AppByID(appID)
		if err != nil {
			if err == model.ErrorNotFound {
				ar.Error(w, err, http.StatusNotFound, "")
			} else {
				ar.Error(w, err, http.StatusInternalServerError, "")
			}
			return
		}

		app = app.Sanitize()
		ar.ServeJSON(w, http.StatusOK, app)
	}
}

// FetchApps fetches apps from the database.
func (ar *Router) FetchApps() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filterStr := strings.TrimSpace(r.URL.Query().Get("search"))

		skip, limit, err := ar.parseSkipAndLimit(r, defaultAppSkip, defaultAppLimit, 0)
		if err != nil {
			ar.Error(w, ErrorWrongInput, http.StatusBadRequest, "")
			return
		}

		apps, total, err := ar.appStorage.FetchApps(filterStr, skip, limit)
		if err != nil {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, "")
			return
		}
		for i, app := range apps {
			apps[i] = app.Sanitize()
		}

		searchResponse := struct {
			Apps  []model.AppData `json:"apps"`
			Total int             `json:"total"`
		}{
			Apps:  apps,
			Total: total,
		}
		ar.ServeJSON(w, http.StatusOK, &searchResponse)
	}
}

// CreateApp adds new app to the database.
func (ar *Router) CreateApp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ad := ar.appStorage.NewAppData()
		if ar.mustParseJSON(w, r, ad) != nil {
			return
		}

		app, err := ar.appStorage.CreateApp(ad)
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		app = app.Sanitize()
		ar.ServeJSON(w, http.StatusOK, app)
	}
}

// UpdateApp updates app in the database.
func (ar *Router) UpdateApp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appID := getRouteVar("id", r)

		ad := ar.appStorage.NewAppData()
		if ar.mustParseJSON(w, r, ad) != nil {
			return
		}

		app, err := ar.appStorage.UpdateApp(appID, ad)
		if err != nil {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, "")
			return
		}

		ar.logger.Printf("App %s updated", appID)

		app = app.Sanitize()
		ar.ServeJSON(w, http.StatusOK, app)
	}
}

// DeleteApp deletes app from the database by id.
func (ar *Router) DeleteApp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appID := getRouteVar("id", r)
		if err := ar.appStorage.DeleteApp(appID); err != nil {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, "")
			return
		}

		ar.logger.Printf("App %s deleted", appID)

		ar.ServeJSON(w, http.StatusOK, nil)
	}
}
