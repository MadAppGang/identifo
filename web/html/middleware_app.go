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
	// FormKeyAppID is a form key to keep application ID.
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

		app, err := ar.AppStorage.ActiveAppByID(appID)
		if err != nil {
			ar.Logger.Printf("Error: getting app by id. %s", err)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), model.AppDataContextKey, app)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}
