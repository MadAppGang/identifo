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

func (ar *Router) updateServerConfigFile(w http.ResponseWriter, newSettings *model.ServerSettings) error {
	dir, err := os.Getwd()
	if err != nil {
		ar.logger.Println("Cannot get server configuration file:", err)
		ar.Error(w, err, http.StatusInternalServerError, "")
		return err
	}

	ss, err := yaml.Marshal(newSettings)
	if err != nil {
		ar.logger.Println("Cannot marshall server configuration:", err)
		ar.Error(w, err, http.StatusBadRequest, "")
		return err
	}

	if err = ioutil.WriteFile(filepath.Join(dir, ar.ServerConfigPath), ss, 0644); err != nil {
		ar.logger.Println("Cannot write server configuration file:", err)
		ar.Error(w, err, http.StatusInternalServerError, "")
		return err
	}

	return nil
}
