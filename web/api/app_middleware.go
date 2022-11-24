package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
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

// resolveRedirectURI checks redirects valid case and
func (ar *Router) resolveRedirectURI(r *http.Request, requestedURI string) (*url.URL, error) {
	app := middleware.AppFromContext(r.Context())
	rurl, err := url.ParseRequestURI(requestedURI)
	if err != nil {
		return nil, fmt.Errorf("requested URI is invalid: %s", requestedURI)
	}

	for _, rr := range app.RedirectURLs {
		url, err := url.ParseRequestURI(rr)
		if err != nil {
			return nil, fmt.Errorf("app has invalid redirect URL: %s", rr)
		}
		if strings.ToLower(rurl.Host) == strings.ToLower(url.Host) &&
			strings.ToLower(rurl.Scheme) == strings.ToLower(url.Scheme) &&
			strings.ToLower(rurl.Path) == strings.ToLower(url.Path) {
			return url, nil
		}
	}

	return nil, fmt.Errorf("requested URL is not allowed %s", requestedURI)
}
