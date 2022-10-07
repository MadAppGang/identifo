package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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
		if len(appID) == 0 {
			appID = r.URL.Query().Get(QueryKeyAppID)
		}
		if len(appID) == 0 {
			err := fmt.Errorf("no appId provided in header or query")
			ar.Error(rw, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, err.Error(), "AppID.AppFromContext")
			return
		}
		app, err := ar.server.Storages().App.ActiveAppByID(appID)
		if err != nil {
			err = fmt.Errorf("Error getting App by ID(%s): %s", appID, err)
			ar.Error(rw, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, err.Error(), "AppID.AppFromContext")
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
