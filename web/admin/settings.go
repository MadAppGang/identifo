package admin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/madappgang/identifo/model"
	"gopkg.in/yaml.v2"
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

// AlterServerSettings changes the whole set of server settings.
func (ar *Router) AlterServerSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newset := new(model.ServerSettings)

		if ar.mustParseJSON(w, r, newset) != nil {
			return
		}

		if ar.updateServerSettings(w, newset) != nil {
			return
		}

		ar.ServeJSON(w, http.StatusOK, newset)
	}
}

// AlterDatabaseSettings changes database connection settings.
func (ar *Router) AlterDatabaseSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dbSettingsUpdate model.DBSettings
		if ar.mustParseJSON(w, r, &dbSettingsUpdate) != nil {
			return
		}

		newServerSettings := new(model.ServerSettings)
		if err := ar.getServerSettings(w, newServerSettings); err != nil {
			return
		}

		newServerSettings.DBSettings = dbSettingsUpdate
		if ar.updateServerSettings(w, newServerSettings) != nil {
			return
		}

		ar.ServeJSON(w, http.StatusOK, newServerSettings.DBSettings)
	}
}

// AlterAccountSettings changes admin account settings.
func (ar *Router) AlterAccountSettings() http.HandlerFunc {
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

func (ar *Router) validateAdminPassword(pswd string, w http.ResponseWriter) error {
	if len(pswd) < 6 || len(pswd) > 130 {
		err := fmt.Errorf("Incorrect password length %d, expecting number between 6 and 130", len(pswd))
		ar.Error(w, err, http.StatusBadRequest, err.Error())
		return err
	}
	return nil
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
func (ar *Router) getServerSettings(w http.ResponseWriter, sc *model.ServerSettings) error {
	dir, err := os.Getwd()
	if err != nil {
		ar.logger.Println("Cannot get server configuration file:", err)
		ar.Error(w, err, http.StatusInternalServerError, "")
		return err
	}

	yamlFile, err := ioutil.ReadFile(filepath.Join(dir, ar.ServerConfigPath))
	if err != nil {
		ar.logger.Println("Cannot read server configuration file:", err)
		ar.Error(w, err, http.StatusInternalServerError, "")
		return err
	}

	if err = yaml.Unmarshal(yamlFile, sc); err != nil {
		ar.logger.Println("Cannot unmarshal server configuration file:", err)
		ar.Error(w, err, http.StatusInternalServerError, "")
		return err
	}

	return nil
}

func (ar *Router) updateServerSettings(w http.ResponseWriter, newSettings *model.ServerSettings) error {
	dir, err := os.Getwd()
	if err != nil {
		ar.logger.Println("Cannot get server configuration file:", err)
		ar.Error(w, err, http.StatusInternalServerError, "")
		return err
	}

	if err = ar.updateConfigFile(w, newSettings, filepath.Join(dir, ar.ServerConfigPath)); err != nil {
		ar.logger.Println("Cannot update server configuration file:", err)
		ar.Error(w, err, http.StatusInternalServerError, "")
		return err
	}

	if err = ar.configurationStorage.Insert(ar.ServerSettings.SettingsKey, newSettings); err != nil {
		ar.logger.Println("Cannot insert new settings into configuartion storage:", err)
		ar.Error(w, err, http.StatusInternalServerError, "")
		return err
	}
	return nil
}

// getAdminAccountSettings admin account settings and parses them to adminData struct.
func (ar *Router) getAdminAccountSettings(w http.ResponseWriter, ald *adminLoginData) error {
	adminLogin := os.Getenv(ar.ServerSettings.AdminAccountSettings.LoginEnvName)
	if len(adminLogin) == 0 {
		err := fmt.Errorf("Admin login not set")
		ar.Error(w, err, http.StatusInternalServerError, err.Error())
		return err
	}

	adminPassword := os.Getenv(ar.ServerSettings.AdminAccountSettings.PasswordEnvName)
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
	if err := os.Setenv(ar.ServerSettings.AdminAccountSettings.LoginEnvName, newAdminData.Login); err != nil {
		err = fmt.Errorf("Cannot save new admin login: %s", err)
		ar.Error(w, err, http.StatusInternalServerError, err.Error())
		return err
	}

	if err := os.Setenv(ar.ServerSettings.AdminAccountSettings.PasswordEnvName, newAdminData.Password); err != nil {
		err = fmt.Errorf("Cannot save new admin password: %s", err)
		ar.Error(w, err, http.StatusInternalServerError, err.Error())
		return err
	}
	return nil
}

func (ar *Router) updateConfigFile(w http.ResponseWriter, in interface{}, dir string) error {
	ss, err := yaml.Marshal(in)
	if err != nil {
		ar.logger.Println("Cannot marshall configuration:", err)
		ar.Error(w, err, http.StatusBadRequest, "")
		return err
	}

	if err = ioutil.WriteFile(dir, ss, 0644); err != nil {
		ar.logger.Println("Cannot write configuration file:", err)
		ar.Error(w, err, http.StatusInternalServerError, "")
		return err
	}
	return nil
}
