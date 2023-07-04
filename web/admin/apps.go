package admin

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
)

const (
	defaultAppSkip  = 0
	defaultAppLimit = 20
)

// GetApp fetches app by ID from the database.
func (ar *Router) GetApp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")
		appID := getRouteVar("id", r)
		app, err := ar.server.Storages().App.AppByID(appID)
		if err != nil {
			status := http.StatusInternalServerError
			if errors.Is(err, model.ErrUserNotFound) {
				status = http.StatusNotFound
			}
			ar.HTTPError(w, err, status)
			return
		}
		ar.ServeJSON(w, locale, http.StatusOK, app)
	}
}

// FetchApps fetches apps from the database.
func (ar *Router) FetchApps() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")
		filterStr := strings.TrimSpace(r.URL.Query().Get("search"))

		apps, err := ar.server.Storages().App.FetchApps(filterStr)
		if err != nil {
			ar.HTTPError(w, err, http.StatusInternalServerError)
			return
		}
		for i, app := range apps {
			apps[i] = app.Sanitized()
		}

		searchResponse := struct {
			Apps []model.AppData `json:"apps"`
		}{
			Apps: apps,
		}
		ar.ServeJSON(w, locale, http.StatusOK, &searchResponse)
	}
}

// CreateApp adds new app to the database.
func (ar *Router) CreateApp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		ad := model.AppData{}
		if ar.mustParseJSON(w, r, &ad) != nil {
			return
		}

		appSecret, err := ar.generateAppSecret()
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAdminPanelGenerateSecret, err)
			return
		}
		ad.Secret = appSecret

		app, err := ar.server.Storages().App.CreateApp(ad)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorStorageRequestError, err)
			return
		}
		ar.ServeJSON(w, locale, http.StatusOK, app)
	}
}

// UpdateApp updates app in the database.
func (ar *Router) UpdateApp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")
		appID := getRouteVar("id", r)

		ad := model.AppData{}
		if ar.mustParseJSON(w, r, &ad) != nil {
			return
		}

		if lenSecret := len(ad.Secret); lenSecret < 24 || lenSecret > 48 {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAdminPanelAPPSecretLength, lenSecret)
			return
		}
		if !isBase64(ad.Secret) {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAdminPanelAPPSecretNotBase64)
			return
		}

		app, err := ar.server.Storages().App.UpdateApp(appID, ad)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorStorageRequestError, err)
			return
		}

		if err = ar.updateAllowedOrigins(); err != nil {
			ar.Logger.Printf("Error occurred during updating allowed origins for App %s, error: %v", appID, err)
		}

		ar.Logger.Printf("App %s updated", appID)

		ar.ServeJSON(w, locale, http.StatusOK, app)
	}
}

// DeleteApp deletes app from the database by id.
func (ar *Router) DeleteApp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		appID := getRouteVar("id", r)
		if err := ar.server.Storages().App.DeleteApp(appID); err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorStorageRequestError, err)
			return
		}

		ar.Logger.Printf("App %s deleted", appID)

		ar.ServeJSON(w, locale, http.StatusOK, nil)
	}
}

// DeleteAllApps delete all current apps for some reason?
// now we are using it for tests
func (ar *Router) DeleteAllApps() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		apps, err := ar.server.Storages().App.FetchApps("")
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorStorageRequestError, err)
		}
		var errs []error
		for _, a := range apps {
			err := ar.server.Storages().App.DeleteApp(a.ID)
			if err != nil {
				errs = append(errs, err)
				ar.Logger.Printf("Error deleting app: %s, error: %s. Ignoring and moving next.", a.ID, err)
			}
		}
		if len(errs) > 0 {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorStorageRequestError, errors.Join(errs...))
		}
		ar.ServeJSON(w, locale, http.StatusOK, nil)
	}
}

func (ar *Router) generateAppSecret() (string, error) {
	secret := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, secret); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(secret), nil
}

func (ar *Router) updateAllowedOrigins() error {
	if ar.originUpdate == nil {
		return nil
	}
	return ar.originUpdate()
}

func isBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}
