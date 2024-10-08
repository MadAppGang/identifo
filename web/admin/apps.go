package admin

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
)

// GetApp fetches app by ID from the database.
func (ar *Router) GetApp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appID := getRouteVar("id", r)

		app, err := ar.server.Storages().App.AppByID(appID)
		if err != nil {
			if err == model.ErrorNotFound {
				ar.Error(w, err, http.StatusNotFound, "")
				return
			}
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}
		ar.ServeJSON(w, http.StatusOK, app)
	}
}

// FetchApps fetches apps from the database.
func (ar *Router) FetchApps() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filterStr := strings.TrimSpace(r.URL.Query().Get("search"))

		apps, err := ar.server.Storages().App.FetchApps(filterStr)
		if err != nil {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, "")
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
		ar.ServeJSON(w, http.StatusOK, &searchResponse)
	}
}

// CreateApp adds new app to the database.
func (ar *Router) CreateApp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ad := model.AppData{}
		if ar.mustParseJSON(w, r, &ad) != nil {
			return
		}

		appSecret, err := ar.generateAppSecret(w)
		if err != nil {
			return
		}
		ad.Secret = appSecret

		app, err := ar.server.Storages().App.CreateApp(ad)
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}
		ar.ServeJSON(w, http.StatusOK, app)
	}
}

// UpdateApp updates app in the database.
func (ar *Router) UpdateApp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appID := getRouteVar("id", r)

		ad := model.AppData{}
		if ar.mustParseJSON(w, r, &ad) != nil {
			return
		}

		if lenSecret := len(ad.Secret); lenSecret < 24 || lenSecret > 48 {
			err := fmt.Errorf("incorrect appsecret string length %d, expecting 24 to 48 symbols inclusively", lenSecret)
			ar.Error(w, err, http.StatusBadRequest, err.Error())
			return
		}
		if !isBase64(ad.Secret) {
			err := fmt.Errorf("expecting app secret to be base64 encoded")
			ar.Error(w, err, http.StatusBadRequest, err.Error())
			return
		}

		app, err := ar.server.Storages().App.UpdateApp(appID, ad)
		if err != nil {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, err.Error())
			return
		}

		if err = ar.updateAllowedOrigins(); err != nil {
			ar.logger.Error("Error occurred during updating allowed origins for App",
				"appId", appID,
				"error", err)
		}

		ar.logger.Info("App updated",
			"appId", appID)

		ar.ServeJSON(w, http.StatusOK, app)
	}
}

// DeleteApp deletes app from the database by id.
func (ar *Router) DeleteApp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appID := getRouteVar("id", r)
		if err := ar.server.Storages().App.DeleteApp(appID); err != nil {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, "")
			return
		}

		ar.logger.Info("App deleted", "appId", appID)

		ar.ServeJSON(w, http.StatusOK, nil)
	}
}

// DeleteAllApps delete all current apps for some reason?
// now we are using it for tests
func (ar *Router) DeleteAllApps() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apps, err := ar.server.Storages().App.FetchApps("")
		if err != nil {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, err.Error())
		}
		var errs []error
		for _, a := range apps {
			err := ar.server.Storages().App.DeleteApp(a.ID)
			if err != nil {
				errs = append(errs, err)
				ar.logger.Error("Error deleting app. Ignoring and moving next.",
					logging.FieldAppID, a.ID,
					logging.FieldError, err)
			}
		}
		if len(errs) > 0 {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, fmt.Sprintf("%v", errs))
		}
		ar.ServeJSON(w, http.StatusOK, nil)
	}
}

func (ar *Router) generateAppSecret(w http.ResponseWriter) (string, error) {
	secret := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, secret); err != nil {
		ar.Error(w, err, http.StatusInternalServerError, "Cannot create app secret")
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
