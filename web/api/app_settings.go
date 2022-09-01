package api

import (
	"net/http"
	"strings"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
)

type appSettings struct {
	AnonymousResitrationAllowed bool            `json:"anonymousResitrationAllowed"`
	Active                      bool            `json:"active"`
	Description                 string          `json:"description"`
	ID                          string          `json:"id"`
	NewUserDefaultRole          string          `json:"newUserDefaultRole"`
	Offline                     bool            `json:"offline"`
	RegistrationForbidden       bool            `json:"registrationForbidden"`
	TfaType                     string          `json:"tfaType"`
	TfaStatus                   string          `json:"tfaStatus"`
	TfaResendTimeout            int             `json:"tfaResendTimeout"`
	LoginWith                   model.LoginWith `json:"loginWith"`
	FederatedProviders          []string        `json:"federatedProviders"`
	CustomEmailTemplates        bool            `json:"customEmailTemplates"`
}

// GetAppSettings return app settings
func (ar *Router) GetAppSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		callbackUrl := strings.TrimSpace(r.URL.Query().Get("callbackUrl"))
		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.logger.Println("Error getting App")
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App is not in context.", "LoginWithPassword.AppFromContext")
			return
		}

		if callbackUrl != "" && !contains(app.RedirectURLs, callbackUrl) {
			ar.logger.Println("Error checking callback URL")
			ar.Error(w, ErrorAPIRequestCallbackUrlInvalid, http.StatusBadRequest, "Callback url is invalid", "LoginWithPassword.CallbackURL")
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
			TfaStatus:                   string(app.TFAStatus),
			TfaResendTimeout:            ar.tfaResendTimeout,
			LoginWith:                   ar.SupportedLoginWays,
			FederatedProviders:          make([]string, 0, len(app.FederatedProviders)),
			CustomEmailTemplates:        app.CustomEmailTemplates,
		}

		for k := range app.FederatedProviders {
			result.FederatedProviders = append(result.FederatedProviders, k)
		}

		ar.ServeJSON(w, http.StatusOK, result)
	}
}
