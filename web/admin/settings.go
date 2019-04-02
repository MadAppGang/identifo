package admin

import (
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
		return
	}
}

// FetchAccountSettings provides info about used database engine.
func (ar *Router) FetchAccountSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conf := new(adminData)
		if ar.getAccountConf(w, conf) != nil {
			return
		}
		ar.ServeJSON(w, http.StatusOK, conf)
		return
	}
}

// AlterServerSettings changes server settings.
func (ar *Router) AlterServerSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newset := new(model.ServerSettings)

		if ar.mustParseJSON(w, r, newset) != nil {
			return
		}

		if ar.updateServerConfigFile(w, newset) != nil {
			return
		}

		ar.ServeJSON(w, http.StatusOK, newset)
		return
	}
}

// AlterAccountSettings changes admin account settings.
func (ar *Router) AlterAccountSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newset := new(adminData)

		if ar.mustParseJSON(w, r, newset) != nil {
			return
		}

		if ar.updateAccountConfigFile(w, newset) != nil {
			return
		}

		ar.ServeJSON(w, http.StatusOK, newset)
		return
	}
}

func (ar *Router) updateServerConfigFile(w http.ResponseWriter, newSettings *model.ServerSettings) error {
	dir, err := os.Getwd()
	if err != nil {
		ar.logger.Println("Cannot get server configuration file:", err)
		ar.Error(w, err, http.StatusInternalServerError, "")
		return err
	}

	return ar.updateConfigFile(w, newSettings, filepath.Join(dir, ar.ServerConfigPath))
}

func (ar *Router) updateAccountConfigFile(w http.ResponseWriter, newSettings *adminData) error {
	dir, err := os.Getwd()
	if err != nil {
		ar.logger.Println("Cannot get account configuration file:", err)
		ar.Error(w, err, http.StatusInternalServerError, "")
		return err
	}

	return ar.updateConfigFile(w, newSettings, filepath.Join(dir, ar.AccountConfigPath))
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
