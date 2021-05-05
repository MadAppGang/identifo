package api

import (
	"net/http"

	"github.com/madappgang/identifo/web/middleware"
)

type appSettings struct {
	AnonymousResitrationAllowed bool   `json:"anonymousResitrationAllowed"`
	Active                      bool   `json:"active"`
	Description                 string `json:"description"`
	ID                          string `json:"id"`
	NewUserDefaultRole          string `json:"newUserDefaultRole"`
	Offline                     bool   `json:"offline"`
	RegistrationForbidden       bool   `json:"registrationForbidden"`
	TfaType                     string `json:"tfaType"`
}

// GetAppSettings return app settings
func (ar *Router) GetAppSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.logger.Println("Error getting App")
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App is not in context.", "LoginWithPassword.AppFromContext")
			return
		}

		result := appSettings{
			AnonymousResitrationAllowed: app.AnonymousRegistrationAllowed,
			Active:                      app.Active,
			Description:                 app.Description,
			ID:                          app.ID,
			NewUserDefaultRole:          app.NewUserDefaultRole,
			Offline:                     app.Offline,
			RegistrationForbidden:       app.RegistrationForbidden,
			TfaType:                     string(ar.tfaType),
		}

		ar.ServeJSON(w, http.StatusOK, result)
	}
}
