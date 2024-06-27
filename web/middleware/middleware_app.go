package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"github.com/urfave/negroni"
)

const (
	// FormKeyAppID is a form key to keep application ID.
	FormKeyAppID = "appId"
)

// AppID gets app id from the request body.
func AppID(
	logger *slog.Logger,
	errorPath string,
	appStorage model.AppStorage,
) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if strings.HasSuffix(errorPath, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		var appID string

		switch r.Method {
		case http.MethodGet:
			appID = strings.TrimSpace(r.URL.Query().Get(FormKeyAppID))
		case http.MethodPost:
			if err := r.ParseForm(); err != nil {
				break
			}
			appID = strings.TrimSpace(r.FormValue(FormKeyAppID))
		}

		app, err := appStorage.ActiveAppByID(appID)
		if err != nil {
			logger.Error("Error: getting app by id",
				logging.FieldError, err)
			http.Redirect(w, r, errorPath, http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), model.AppDataContextKey, app)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}
