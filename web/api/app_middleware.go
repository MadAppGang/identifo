package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/madappgang/identifo/model"

	"github.com/urfave/negroni"
)

const (
	//HeaderKeyAppID header key to keep application ID
	HeaderKeyAppID = "X-Identifo-Clientid"
)

func (ar *Router) AppID() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		appID := strings.TrimSpace(r.Header.Get(HeaderKeyAppID))
		app, err := ar.appStorage.ActiveAppByID(appID)
		if err != nil {
			ar.logger.Printf("Error getting App by ID %v", err)
			ar.Error(rw, ErrorRequestInvalidAppID, http.StatusBadRequest, "")
			return
		}
		ctx := context.WithValue(r.Context(), model.AppDataContextKey, app)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	}
}
