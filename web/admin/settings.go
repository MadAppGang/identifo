package admin

import (
	"fmt"
	"net/http"

	"github.com/madappgang/identifo/model"
)

type ServerSettingsAPI struct {
	General        *model.GeneralServerSettings  `json:"general,omitempty"`
	AdminAccount   *model.AdminAccountSettings   `json:"admin_account,omitempty"`
	Storage        *StorageSettingsAPI           `json:"storage,omitempty"`
	SessionStorage *model.SessionStorageSettings `json:"session_storage,omitempty"`
	Services       *model.ServicesSettings       `json:"external_services,omitempty"`
	Login          *model.LoginSettings          `json:"login,omitempty"`
	KeyStorage     *model.KeyStorageSettings     `json:"keyStorage,omitempty"`
	Config         *model.ConfigStorageSettings  `json:"config,omitempty"`
	Logger         *model.LoggerSettings         `json:"logger,omitempty"`
	AdminPanel     *model.AdminPanelSettings     `json:"admin_panel"`
	LoginWebApp    *model.FileStorageSettings    `json:"login_web_app"`
	EmailTemplates *model.FileStorageSettings    `json:"email_templaits"`
}

type StorageSettingsAPI struct {
	AppStorage              *model.DatabaseSettings `json:"app_storage,omitempty"`
	UserStorage             *model.DatabaseSettings `json:"user_storage,omitempty"`
	TokenStorage            *model.DatabaseSettings `json:"token_storage,omitempty"`
	TokenBlacklist          *model.DatabaseSettings `json:"token_blacklist,omitempty"`
	VerificationCodeStorage *model.DatabaseSettings `json:"verification_code_storage,omitempty"`
	InviteStorage           *model.DatabaseSettings `json:"invite_storage,omitempty"`
}

// FetchSettings returns server settings.
func (ar *Router) FetchSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := ar.server.Settings()
		ar.ServeJSON(w, http.StatusOK, s)
	}
}

// UpdateSettings handles the request to update server settings.
func (ar *Router) UpdateSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := ar.server.Settings()

		us := ServerSettingsAPI{}
		if err := ar.mustParseJSON(w, r, &us); err != nil {
			ar.Error(w, fmt.Errorf("error parsing api settings: %v", err), http.StatusBadRequest, "")
			return
		}

		merged, changed := mergeSettings(s, us)
		if changed == false {
			ar.Error(w, fmt.Errorf("no settings has been changed, skipping the update"), http.StatusBadRequest, "")
			return
		}

		if err := merged.Validate(); err != nil {
			ar.Error(w, fmt.Errorf("settings validation failed with error: %v", err), http.StatusBadRequest, "")
			return
		}

		if err := ar.server.Storages().Config.WriteConfig(merged); err != nil {
			ar.logger.Println("Cannot insert new settings into configuration storage:", err)
			ar.Error(w, fmt.Errorf("error saving new config: %v", err), http.StatusInternalServerError, "")
			return
		}

		// if the config storage is not supporting instant reloading - let's force restart it
		if ar.forceRestart != nil && ar.server.Storages().Config.ForceReloadOnWriteConfig() {
			go func() {
				ar.logger.Println("sending server restart")
				ar.forceRestart <- true
			}()
		}

		ar.ServeJSON(w, http.StatusOK, merged)
	}
}

// mergeSettings merges updatedSettings with settings and produces the new setttings
func mergeSettings(settings model.ServerSettings, updatedSettings ServerSettingsAPI) (model.ServerSettings, bool) {
	changed := false
	if updatedSettings.General != nil {
		settings.General = *updatedSettings.General
		changed = true
	}

	if updatedSettings.AdminAccount != nil {
		settings.AdminAccount = *updatedSettings.AdminAccount
		changed = true
	}

	if updatedSettings.Storage != nil {
		if updatedSettings.Storage.AppStorage != nil {
			settings.Storage.AppStorage = *updatedSettings.Storage.AppStorage
			changed = true
		}
		if updatedSettings.Storage.UserStorage != nil {
			settings.Storage.UserStorage = *updatedSettings.Storage.UserStorage
			changed = true
		}
		if updatedSettings.Storage.TokenStorage != nil {
			settings.Storage.TokenStorage = *updatedSettings.Storage.TokenStorage
			changed = true
		}
		if updatedSettings.Storage.TokenBlacklist != nil {
			settings.Storage.TokenBlacklist = *updatedSettings.Storage.TokenBlacklist
			changed = true
		}
		if updatedSettings.Storage.VerificationCodeStorage != nil {
			settings.Storage.VerificationCodeStorage = *updatedSettings.Storage.VerificationCodeStorage
			changed = true
		}
		if updatedSettings.Storage.InviteStorage != nil {
			settings.Storage.InviteStorage = *updatedSettings.Storage.InviteStorage
			changed = true
		}
	}

	if updatedSettings.SessionStorage != nil {
		settings.SessionStorage = *updatedSettings.SessionStorage
		changed = true
	}

	if updatedSettings.LoginWebApp != nil {
		settings.LoginWebApp = *updatedSettings.LoginWebApp
		changed = true
	}

	if updatedSettings.AdminPanel != nil {
		settings.AdminPanel = *updatedSettings.AdminPanel
		changed = true
	}

	if updatedSettings.EmailTemplates != nil {
		settings.EmailTemplates = *updatedSettings.EmailTemplates
		changed = true
	}

	if updatedSettings.Services != nil {
		settings.Services = *updatedSettings.Services
		changed = true
	}

	if updatedSettings.Login != nil {
		settings.Login = *updatedSettings.Login
		changed = true
	}

	if updatedSettings.KeyStorage != nil {
		settings.KeyStorage = *updatedSettings.KeyStorage
		changed = true
	}

	if updatedSettings.Config != nil {
		settings.Config = *updatedSettings.Config
		changed = true
	}

	if updatedSettings.Logger != nil {
		settings.Logger = *updatedSettings.Logger
		changed = true
	}

	// we need to go section by section and check nee settings
	return settings, changed
}

// TestDatabaseConnection tests database connection.
func (ar *Router) TestDatabaseConnection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := ar.server.Storages().App.TestDatabaseConnection(); err != nil {
			ar.ServeJSON(w, http.StatusInternalServerError, nil)
		} else {
			ar.ServeJSON(w, http.StatusOK, nil)
		}
	}
}
