package html

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/madappgang/identifo/model"
	"github.com/urfave/negroni"
)

const (
	// FormKeyAppID is a form key to keep application ID.
	FormKeyAppID = "appId"
)

// AppID gets app id from the request body.
func (ar *Router) AppID() negroni.HandlerFunc {
	errorPath := path.Join(ar.PathPrefix, "/misconfiguration")
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		appID := ""

		onError := func(message string) {
			ar.Logger.Print(message)
			http.Redirect(w, r, errorPath, http.StatusFound)
		}

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
			onError("Empty appId param")
			return
		}

		app, err := ar.AppStorage.AppByID(appID)
		if err != nil {
			message := fmt.Sprintf("Error getting App by ID: %v", err)
			onError(message)
			return
		}

		if !app.Active() {
			message := fmt.Sprintf("App with ID: %v is inactive", app.ID())
			onError(message)
			return
		}

		ctx := context.WithValue(r.Context(), model.AppDataContextKey, app)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

// appFromContext returns app data from request conntext.
func appFromContext(ctx context.Context) model.AppData {
	value := ctx.Value(model.AppDataContextKey)

	if value == nil {
		return nil
	}

	return value.(model.AppData)
}
