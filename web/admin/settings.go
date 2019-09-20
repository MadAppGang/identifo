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

// UpdateAccountSettings updates admin account settings.
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

		adminData := new(adminLoginData)
		if err := ar.getAdminAccountSettings(w, adminData); err != nil {
			return
		}

		namesDidNotChange := adminDataUpdate.LoginEnvName == adminData.LoginEnvName && adminDataUpdate.PasswordEnvName == adminData.PasswordEnvName
		valuesDidNotChange := adminDataUpdate.Login == adminData.Login && adminDataUpdate.Password == adminData.Password

		if namesDidNotChange && valuesDidNotChange {
			ar.ServeJSON(w, http.StatusOK, nil)
			return
		}

		if len(adminDataUpdate.Login) > 0 {
			adminData.Login = adminDataUpdate.Login
		}
		if len(adminDataUpdate.LoginEnvName) > 0 {
			adminData.LoginEnvName = adminDataUpdate.LoginEnvName
		} else {
			adminData.LoginEnvName = ar.ServerSettings.AdminAccount.LoginEnvName
		}

		if len(adminDataUpdate.Password) > 0 {
			adminData.Password = adminDataUpdate.Password
		}
		if len(adminDataUpdate.PasswordEnvName) > 0 {
			adminData.PasswordEnvName = adminDataUpdate.PasswordEnvName
		} else {
			adminData.PasswordEnvName = ar.ServerSettings.AdminAccount.PasswordEnvName
		}

		if ar.updateAdminAccountSettings(w, adminData) != nil {
			return
		}

		ar.ServeJSON(w, http.StatusOK, adminData)
	}
}

// FetchGeneralSettings fetches server's general settings.
func (ar *Router) FetchGeneralSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ar.ServeJSON(w, http.StatusOK, ar.ServerSettings.General)
	}
}

// UpdateGeneralSettings changes server's general settings.
func (ar *Router) UpdateGeneralSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var generalSettingsUpdate model.GeneralServerSettings

		if ar.mustParseJSON(w, r, &generalSettingsUpdate) != nil {
			return
		}
		if err := generalSettingsUpdate.Validate(); err != nil {
			ar.Error(w, err, http.StatusBadRequest, err.Error())
			return
		}

		ar.newSettings.General = generalSettingsUpdate
		ar.ServeJSON(w, http.StatusOK, ar.newSettings.General)
	}
}

// FetchStorageSettings fetches server's general settings.
func (ar *Router) FetchStorageSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ar.ServeJSON(w, http.StatusOK, ar.ServerSettings.Storage)
	}
}

// UpdateStorageSettings changes storage connection settings.
func (ar *Router) UpdateStorageSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var storageSettingsUpdate model.StorageSettings

		if ar.mustParseJSON(w, r, &storageSettingsUpdate) != nil {
			return
		}
		if err := storageSettingsUpdate.Validate(); err != nil {
			ar.Error(w, err, http.StatusBadRequest, err.Error())
			return
		}

		ar.newSettings.Storage = storageSettingsUpdate
		ar.ServeJSON(w, http.StatusOK, ar.newSettings.SessionStorage)
	}
}

// FetchSessionStorageSettings fetches session storage settings.
func (ar *Router) FetchSessionStorageSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ar.ServeJSON(w, http.StatusOK, ar.ServerSettings.SessionStorage)
	}
}

// UpdateSessionStorageSettings changes admin session storage connection settings.
func (ar *Router) UpdateSessionStorageSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sessionStorageSettingsUpdate model.SessionStorageSettings

		if ar.mustParseJSON(w, r, &sessionStorageSettingsUpdate) != nil {
			return
		}
		if err := sessionStorageSettingsUpdate.Validate(); err != nil {
			ar.Error(w, err, http.StatusBadRequest, err.Error())
			return
		}

		ar.newSettings.SessionStorage = sessionStorageSettingsUpdate
		ar.ServeJSON(w, http.StatusOK, ar.newSettings.SessionStorage)
	}
}

// FetchConfigurationStorageSettings fetches configuration storage settings.
func (ar *Router) FetchConfigurationStorageSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ar.ServeJSON(w, http.StatusOK, ar.ServerSettings.ConfigurationStorage)
	}
}

// RestartServer restarts server with new settings.
func (ar *Router) RestartServer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := ar.configurationStorage.InsertConfig(ar.ServerSettings.ConfigurationStorage.SettingsKey, ar.newSettings); err != nil {
			ar.logger.Println("Cannot insert new settings into configuartion storage:", err)
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}
		ar.ServeJSON(w, http.StatusOK, nil)
	}
}

// UpdateConfigurationStorageSettings changes storage connection settings.
func (ar *Router) UpdateConfigurationStorageSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var configurationStorageSettingsUpdate model.ConfigurationStorageSettings

		if ar.mustParseJSON(w, r, &configurationStorageSettingsUpdate) != nil {
			return
		}
		if err := configurationStorageSettingsUpdate.Validate(); err != nil {
			ar.Error(w, err, http.StatusBadRequest, err.Error())
			return
		}

		ar.newSettings.ConfigurationStorage = configurationStorageSettingsUpdate
		ar.ServeJSON(w, http.StatusOK, ar.newSettings.ConfigurationStorage)
	}
}

// FetchStaticFilesStorageSettings fetches static files settings.
func (ar *Router) FetchStaticFilesStorageSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ar.ServeJSON(w, http.StatusOK, ar.ServerSettings.StaticFilesStorage)
	}
}

// UpdateStaticFilesStorageSettings changes static files settings.
func (ar *Router) UpdateStaticFilesStorageSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var staticFilesStorageSettingsUpdate model.StaticFilesStorageSettings

		if ar.mustParseJSON(w, r, &staticFilesStorageSettingsUpdate) != nil {
			return
		}

		ar.newSettings.StaticFilesStorage = staticFilesStorageSettingsUpdate
		ar.ServeJSON(w, http.StatusOK, ar.newSettings.StaticFilesStorage)
	}
}

// FetchLoginSettings fetches app's login settings.
func (ar *Router) FetchLoginSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ar.ServeJSON(w, http.StatusOK, ar.ServerSettings.Login)
	}
}

// UpdateLoginSettings changes app's login settings.
func (ar *Router) UpdateLoginSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginSettingsUpdate model.LoginSettings

		if ar.mustParseJSON(w, r, &loginSettingsUpdate) != nil {
			return
		}

		ar.newSettings.Login = loginSettingsUpdate
		ar.ServeJSON(w, http.StatusOK, ar.newSettings.Login)
	}
}

// FetchExternalServicesSettings fetches settings for external services.
func (ar *Router) FetchExternalServicesSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ar.ServeJSON(w, http.StatusOK, ar.ServerSettings.ExternalServices)
	}
}

// UpdateExternalServicesSettings changes settings for external services.
func (ar *Router) UpdateExternalServicesSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var servicesSettingsUpdate model.ExternalServicesSettings

		if ar.mustParseJSON(w, r, &servicesSettingsUpdate) != nil {
			return
		}
		if err := servicesSettingsUpdate.Validate(); err != nil {
			ar.Error(w, err, http.StatusBadRequest, err.Error())
			return
		}

		ar.newSettings.ExternalServices = servicesSettingsUpdate
		ar.ServeJSON(w, http.StatusOK, ar.newSettings.ExternalServices)
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
	ald.LoginEnvName = ar.ServerSettings.AdminAccount.LoginEnvName
	ald.Password = adminPassword
	ald.PasswordEnvName = ar.ServerSettings.AdminAccount.PasswordEnvName

	return nil
}

func (ar *Router) updateAdminAccountSettings(w http.ResponseWriter, newAdminData *adminLoginData) error {
	var needChangeConfig bool

	loginEnvName := ar.ServerSettings.AdminAccount.LoginEnvName
	if newAdminData.LoginEnvName != loginEnvName {
		loginEnvName = newAdminData.LoginEnvName
		needChangeConfig = true
	}
	if err := os.Setenv(loginEnvName, newAdminData.Login); err != nil {
		err = fmt.Errorf("Cannot save new admin login: %s", err)
		ar.Error(w, err, http.StatusInternalServerError, err.Error())
		return err
	}

	passwordEnvName := ar.ServerSettings.AdminAccount.PasswordEnvName
	if newAdminData.PasswordEnvName != passwordEnvName {
		passwordEnvName = newAdminData.PasswordEnvName
		needChangeConfig = true
	}
	if err := os.Setenv(passwordEnvName, newAdminData.Password); err != nil {
		err = fmt.Errorf("Cannot save new admin password: %s", err)
		ar.Error(w, err, http.StatusInternalServerError, err.Error())
		return err
	}

	if !needChangeConfig {
		return nil
	}

	newSettings := ar.ServerSettings
	newSettings.AdminAccount.LoginEnvName = loginEnvName
	newSettings.AdminAccount.PasswordEnvName = passwordEnvName

	ar.newSettings = newSettings
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
