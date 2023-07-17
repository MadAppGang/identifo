package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
)

const (
	// HeaderKeyAppID is a header key to keep application ID.
	HeaderKeyAppID = "X-Identifo-Clientid"
	QueryKeyAppID  = "appId"
)

// AppID extracts application ID from the header and writes corresponding app to the context.
func (ar *Router) AppID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		appID := strings.TrimSpace(r.Header.Get(HeaderKeyAppID))
		if appID == "" {
			appID = r.URL.Query().Get(QueryKeyAppID)
		}

		if appID == "" {
			appID = mux.Vars(r)[QueryKeyAppID]
		}

		app, err := ar.server.Storages().App.ActiveAppByID(appID)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorStorageAPPFindByIDError, appID)
			return
		}
		ctx := context.WithValue(r.Context(), model.AppDataContextKey, app)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
