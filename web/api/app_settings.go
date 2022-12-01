package api

import (
	"net/http"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
)

type appSettings struct {
	AnonymousRegistrationAllowed bool            `json:"anonymousRegistrationAllowed"`
	Active                       bool            `json:"active"`
	Description                  string          `json:"description"`
	ID                           string          `json:"id"`
	NewUserDefaultRole           string          `json:"newUserDefaultRole"`
	Offline                      bool            `json:"offline"`
	RegistrationForbidden        bool            `json:"registrationForbidden"`
	TfaType                      string          `json:"tfaType"`
	TfaStatus                    string          `json:"tfaStatus"`
	TfaResendTimeout             int             `json:"tfaResendTimeout"`
	LoginWith                    model.LoginWith `json:"loginWith"`
	FederatedProviders           []string        `json:"federatedProviders"`
	CustomEmailTemplates         bool            `json:"customEmailTemplates"`
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
			AnonymousRegistrationAllowed: app.AnonymousRegistrationAllowed,
			Active:                       app.Active,
			Description:                  app.Description,
			ID:                           app.ID,
			NewUserDefaultRole:           app.NewUserDefaultRole,
			Offline:                      app.Offline,
			RegistrationForbidden:        app.RegistrationForbidden,
			TfaType:                      string(ar.tfaType),
			TfaStatus:                    string(app.TFAStatus),
			TfaResendTimeout:             ar.tfaResendTimeout,
			LoginWith:                    ar.SupportedLoginWays,
			FederatedProviders:           make([]string, 0, len(app.FederatedProviders)),
			CustomEmailTemplates:         app.CustomEmailTemplates,
		}

		for k := range app.FederatedProviders {
			result.FederatedProviders = append(result.FederatedProviders, k)
		}

		ar.ServeJSON(w, http.StatusOK, result)
	}
}
