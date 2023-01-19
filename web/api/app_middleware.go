package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	l "github.com/madappgang/identifo/v2/localization"
	"github.com/madappgang/identifo/v2/model"

	"github.com/urfave/negroni"
)

const (
	// HeaderKeyAppID is a header key to keep application ID.
	HeaderKeyAppID = "X-Identifo-Clientid"
	QueryKeyAppID  = "appId"
)

// AppID extracts application ID from the header and writes corresponding app to the context.
func (ar *Router) AppID() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		appID := strings.TrimSpace(r.Header.Get(HeaderKeyAppID))
		if appID == "" {
			appID = r.URL.Query().Get(QueryKeyAppID)
		}

		app, err := ar.server.Storages().App.ActiveAppByID(appID)
		if err != nil {
			err = fmt.Errorf("Error getting App by ID: %s", err)
			ar.Error(rw, http.StatusBadRequest, l.ErrorStorageAPPFindByIDError, appID, err)
			return
		}
		ctx := context.WithValue(r.Context(), model.AppDataContextKey, app)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	}
}

func (ar *Router) RemoveTrailingSlash() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		next.ServeHTTP(rw, r)
	}
}
