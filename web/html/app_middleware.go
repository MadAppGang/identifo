package html

import (
	"context"
	"net/http"
	"path"
	"strings"

	"github.com/madappgang/identifo/model"
	"github.com/urfave/negroni"
)

const (
	//FormKeyAppID form key to keep application ID.
	FormKeyAppID = "appId"
)

// AppID gets app id from the request body.
func (ar *Router) AppID() negroni.HandlerFunc {
	errorPath := path.Join(ar.PathPrefix, "/misconfiguration")
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		appID := ""
		switch r.Method {
		case http.MethodGet:
			appID = strings.TrimSpace(r.URL.Query().Get(FormKeyAppID))
		case http.MethodPost:
			if err := r.ParseForm(); err != nil {
				break
			}
			appID = strings.TrimSpace(r.FormValue(FormKeyAppID))
		}

		if appID == "" {
			ar.Logger.Printf("Error: empty appId param")
			http.Redirect(w, r, errorPath, http.StatusMovedPermanently)
			return
		}

		app, err := ar.AppStorage.AppByID(appID)
		if err != nil {
			ar.Logger.Printf("Error getting App by ID %v", err)
			http.Redirect(w, r, errorPath, http.StatusMovedPermanently)
			return
		}

		if !app.Active() {
			ar.Logger.Printf("App with ID: %v is inactive", app.ID())
			http.Redirect(w, r, errorPath, http.StatusMovedPermanently)
			return
		}

		ctx := context.WithValue(r.Context(), model.AppDataContextKey, app)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

// appFromContext returns app data from request conntext.
func appFromContext(ctx context.Context) model.AppData {
	return ctx.Value(model.AppDataContextKey).(model.AppData)
}
