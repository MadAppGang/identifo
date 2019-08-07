package admin

import (
	"fmt"
	"net/http"
	"os"

	"github.com/madappgang/identifo/model"
)

// FetchServerSettings returns server settings.
func (ar *Router) FetchServerSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ar.ServeJSON(w, http.StatusOK, ar.ServerSettings)
	}
}

// FetchAccountSettings returns admin account settings.
func (ar *Router) FetchAccountSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conf := new(adminLoginData)
		if ar.getAdminAccountSettings(w, conf) != nil {
			return
		}
		ar.ServeJSON(w, http.StatusOK, conf)
	}
}

// UpdateGeneralSettings changes server's general settings.
func (ar *Router) UpdateGeneralSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var generalSettingsUpdate struct {
			Data model.GeneralServerSettings `json:"general"`
		}

		if ar.mustParseJSON(w, r, &generalSettingsUpdate) != nil {
			return
		}

		newServerSettings := new(model.ServerSettings)
		if err := ar.getServerSettings(w, newServerSettings); err != nil {
			return
		}

		newServerSettings.General = generalSettingsUpdate.Data
		if ar.updateServerSettings(w, newServerSettings) != nil {
			return
		}

		ar.ServeJSON(w, http.StatusOK, newServerSettings.General)
	}
}

// UpdateStorageSettings changes storage connection settings.
func (ar *Router) UpdateStorageSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var storageSettingsUpdate struct {
			Data model.StorageSettings `json:"storage"`
		}

		if ar.mustParseJSON(w, r, &storageSettingsUpdate) != nil {
			return
		}

		newServerSettings := new(model.ServerSettings)
		if err := ar.getServerSettings(w, newServerSettings); err != nil {
			return
		}

		newServerSettings.Storage = storageSettingsUpdate.Data
		if ar.updateServerSettings(w, newServerSettings) != nil {
			return
		}

		ar.ServeJSON(w, http.StatusOK, newServerSettings.Storage)
	}
}

// UpdateSessionStorageSettings changes admin session storage connection settings.
func (ar *Router) UpdateSessionStorageSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sessionStorageSettingsUpdate struct {
			Data model.SessionStorageSettings `json:"session_storage"`
		}

		if ar.mustParseJSON(w, r, &sessionStorageSettingsUpdate) != nil {
			return
		}

		newServerSettings := new(model.ServerSettings)
		if err := ar.getServerSettings(w, newServerSettings); err != nil {
			return
		}

		newServerSettings.SessionStorage = sessionStorageSettingsUpdate.Data
		if ar.updateServerSettings(w, newServerSettings) != nil {
			return
		}

		ar.ServeJSON(w, http.StatusOK, newServerSettings.Storage)
	}
}

// UpdateConfigurationStorageSettings changes storage connection settings.
func (ar *Router) UpdateConfigurationStorageSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var configurationStorageSettingsUpdate struct {
			Data model.ConfigurationStorageSettings `json:"configuration_storage"`
		}

		if ar.mustParseJSON(w, r, &configurationStorageSettingsUpdate) != nil {
			return
		}

		newServerSettings := new(model.ServerSettings)
		if err := ar.getServerSettings(w, newServerSettings); err != nil {
			return
		}

		newServerSettings.ConfigurationStorage = configurationStorageSettingsUpdate.Data
		if ar.updateServerSettings(w, newServerSettings) != nil {
			return
		}

		ar.ServeJSON(w, http.StatusOK, newServerSettings.Storage)
	}
}

// UpdateStaticFilesSettings changes static files settings.
func (ar *Router) UpdateStaticFilesSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var staticFilesSettingsUpdate struct {
			Data model.StaticFilesSettings `json:"static_files"`
		}

		if ar.mustParseJSON(w, r, &staticFilesSettingsUpdate) != nil {
			return
		}

		newServerSettings := new(model.ServerSettings)
		if err := ar.getServerSettings(w, newServerSettings); err != nil {
			return
		}

		newServerSettings.StaticFiles = staticFilesSettingsUpdate.Data
		if ar.updateServerSettings(w, newServerSettings) != nil {
			return
		}

		ar.ServeJSON(w, http.StatusOK, newServerSettings.StaticFiles)
	}
}

// UpdateLoginSettings changes app's login settings.
func (ar *Router) UpdateLoginSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginSettingsUpdate struct {
			Data model.LoginSettings `json:"login"`
		}

		if ar.mustParseJSON(w, r, &loginSettingsUpdate) != nil {
			return
		}

		newServerSettings := new(model.ServerSettings)
		if err := ar.getServerSettings(w, newServerSettings); err != nil {
			return
		}

		newServerSettings.Login = loginSettingsUpdate.Data
		if ar.updateServerSettings(w, newServerSettings) != nil {
			return
		}

		ar.ServeJSON(w, http.StatusOK, newServerSettings.Login)
	}
}

// UpdateExternalServicesSettings changes settings for external services.
func (ar *Router) UpdateExternalServicesSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var servicesSettingsUpdate struct {
			ExternalServices model.ExternalServicesSettings `json:"external_services"`
		}

		if ar.mustParseJSON(w, r, &servicesSettingsUpdate) != nil {
			return
		}

		newServerSettings := new(model.ServerSettings)
		if err := ar.getServerSettings(w, newServerSettings); err != nil {
			return
		}

		newServerSettings.ExternalServices = servicesSettingsUpdate.ExternalServices
		if ar.updateServerSettings(w, newServerSettings) != nil {
			return
		}

		ar.ServeJSON(w, http.StatusOK, newServerSettings.StaticFiles)
	}
}

// UpdateAccountSettings changes admin account settings.
func (ar *Router) UpdateAccountSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		adminDataUpdate := new(adminLoginData)

		if ar.mustParseJSON(w, r, adminDataUpdate) != nil {
			return
		}

		if adminDataUpdate.Password != "" {
			if err := ar.validateAdminPassword(adminDataUpdate.Password, w); err != nil {
				return
			}
		}

		newAdminData := new(adminLoginData)
		if err := ar.getAdminAccountSettings(w, newAdminData); err != nil {
			return
		}

		if newAdminData.Login == adminDataUpdate.Login && newAdminData.Password == adminDataUpdate.Password {
			ar.ServeJSON(w, http.StatusOK, nil)
			return
		}

		if len(adminDataUpdate.Login) > 0 {
			newAdminData.Login = adminDataUpdate.Login
		}
		if len(adminDataUpdate.Password) > 0 {
			newAdminData.Password = adminDataUpdate.Password
		}

		if ar.updateAdminAccountSettings(w, newAdminData) != nil {
			return
		}

		ar.ServeJSON(w, http.StatusOK, nil)
	}
}

// TestDatabaseConnection tests database connection.
func (ar *Router) TestDatabaseConnection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := ar.appStorage.TestDatabaseConnection(); err != nil {
			ar.ServeJSON(w, http.StatusInternalServerError, nil)
		} else {
			ar.ServeJSON(w, http.StatusOK, nil)
		}
	}
}

// getServerSettings reads server configuration file and parses it to provided struct.
func (ar *Router) getServerSettings(w http.ResponseWriter, ss *model.ServerSettings) error {
	key := ar.ServerSettings.ConfigurationStorage.SettingsKey

	ss.ConfigurationStorage = model.ConfigurationStorageSettings{SettingsKey: key}
	if err := ar.configurationStorage.LoadServerSettings(ss); err != nil {
		ar.logger.Println("Cannot read server configuration:", err)
		ar.Error(w, err, http.StatusInternalServerError, "")
		return err
	}
	return nil
}

func (ar *Router) updateServerSettings(w http.ResponseWriter, newSettings *model.ServerSettings) error {
	if err := ar.configurationStorage.Insert(ar.ServerSettings.ConfigurationStorage.SettingsKey, newSettings); err != nil {
		ar.logger.Println("Cannot insert new settings into configuartion storage:", err)
		ar.Error(w, err, http.StatusInternalServerError, "")
		return err
	}
	return nil
}

// getAdminAccountSettings admin account settings and parses them to adminData struct.
func (ar *Router) getAdminAccountSettings(w http.ResponseWriter, ald *adminLoginData) error {
	adminLogin := os.Getenv(ar.ServerSettings.AdminAccount.LoginEnvName)
	if len(adminLogin) == 0 {
		err := fmt.Errorf("Admin login not set")
		ar.Error(w, err, http.StatusInternalServerError, err.Error())
		return err
	}

	adminPassword := os.Getenv(ar.ServerSettings.AdminAccount.PasswordEnvName)
	if len(adminPassword) == 0 {
		err := fmt.Errorf("Admin password not set")
		ar.Error(w, err, http.StatusInternalServerError, err.Error())
		return err
	}

	ald.Login = adminLogin
	ald.Password = adminPassword

	return nil
}

func (ar *Router) updateAdminAccountSettings(w http.ResponseWriter, newAdminData *adminLoginData) error {
	if err := os.Setenv(ar.ServerSettings.AdminAccount.LoginEnvName, newAdminData.Login); err != nil {
		err = fmt.Errorf("Cannot save new admin login: %s", err)
		ar.Error(w, err, http.StatusInternalServerError, err.Error())
		return err
	}

	if err := os.Setenv(ar.ServerSettings.AdminAccount.PasswordEnvName, newAdminData.Password); err != nil {
		err = fmt.Errorf("Cannot save new admin password: %s", err)
		ar.Error(w, err, http.StatusInternalServerError, err.Error())
		return err
	}
	return nil
}

func (ar *Router) validateAdminPassword(pswd string, w http.ResponseWriter) error {
	if pswdLen := len(pswd); pswdLen < 6 || pswdLen > 130 {
		err := fmt.Errorf("Incorrect password length %d, expecting number between 6 and 130", pswdLen)
		ar.Error(w, err, http.StatusBadRequest, err.Error())
		return err
	}
	return nil
}
