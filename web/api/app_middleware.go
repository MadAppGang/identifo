package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/madappgang/identifo/model"
	"github.com/urfave/negroni"
)

const (
	// HeaderKeyAppID is a header key to keep application ID.
	HeaderKeyAppID = "X-Identifo-Clientid"
)

// AppID extracts application ID from the header and writes corresponding app to the context.
func (ar *Router) AppID() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		appID := strings.TrimSpace(r.Header.Get(HeaderKeyAppID))
		app, err := ar.appStorage.ActiveAppByID(appID)
		if err != nil {
			err = fmt.Errorf("Error getting App by ID: %s", err)
			ar.Error(rw, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, err.Error(), "AppID.AppFromContext")
			return
		}
		ctx := context.WithValue(r.Context(), model.AppDataContextKey, app)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	}
}
