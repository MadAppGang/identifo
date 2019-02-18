package admin

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
)

// mustParseJSON parses request body json data to the `out` interface and then validates it.
// Writes error to ResponseWriter if error happens.
func (ar *Router) mustParseJSON(w http.ResponseWriter, r *http.Request, out interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(out); err != nil {
		ar.Error(w, err, http.StatusBadRequest, "")
		return err
	}

	validate := validator.New()
	if err := validate.Struct(out); err != nil {
		ar.Error(w, err, http.StatusBadRequest, "")
		return err
	}

	return nil
}

// getConf reads admin panel configuration file and parses it to adminData struct.
func (ar *Router) getConf(w http.ResponseWriter, ad *adminData) error {
	dir, err := os.Getwd()
	if err != nil {
		ar.logger.Println("Cannot get configuration file:", err)
		ar.Error(w, err, http.StatusInternalServerError, "")
		return err
	}
	yamlFile, err := ioutil.ReadFile(filepath.Join(dir, ar.ConfigPath))
	if err != nil {
		ar.logger.Println("Cannot read configuration file:", err)
		ar.Error(w, err, http.StatusInternalServerError, "")
		return err
	}

	if err = yaml.Unmarshal(yamlFile, ad); err != nil {
		ar.logger.Println("Cannot unmarshal configuration file:", err)
		ar.Error(w, err, http.StatusInternalServerError, "")
		return err
	}

	return nil
}
