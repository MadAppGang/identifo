package admin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
)

const (
	// errInvalidSkip occurs when 'skip' URL variable cannot be converted to integer.
	errInvalidSkip = "Skip value %s cannot be converted to integer"
	// errInvalidLimit occurs when 'limit' URL variable cannot be converted to integer.
	errInvalidLimit = "Limit value %s cannot be converted to integer"
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

// parseSkipAndLimit parses pagination parameters from provided request.
// If url query value is empty, sets corresponding parameter to default value.
// If nonzero maxLimit parameter provided, it is used as an upper bound for the limit parameter.
// Returns non-nil error if provided strings cannot be converted to integers.
func (ar *Router) parseSkipAndLimit(r *http.Request, defaultSkip, defaultLimit, maxLimit int) (int, int, error) {
	skipStr := r.URL.Query().Get("skip")
	limitStr := r.URL.Query().Get("limit")

	var skip int
	var err error

	if len(skipStr) != 0 {
		skip, err = strconv.Atoi(skipStr)
		if err != nil {
			return 0, 0, fmt.Errorf(errInvalidSkip, skipStr)
		}
	} else {
		skip = defaultSkip
	}

	var limit int
	if len(limitStr) != 0 {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, fmt.Errorf(errInvalidLimit, skipStr)
		}
	} else {
		limit = defaultLimit
	}

	if maxLimit == 0 {
		return skip, limit, nil
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	return skip, limit, nil
}

// getRouteVar returns the route variable with specified name for the provided request.
func getRouteVar(name string, r *http.Request) string {
	return mux.Vars(r)[name]
}

// getAccountConf reads admin account configuration file and parses it to adminData struct.
func (ar *Router) getAccountConf(w http.ResponseWriter, ad *adminData) error {
	dir, err := os.Getwd()
	if err != nil {
		ar.logger.Println("Cannot get configuration file:", err)
		ar.Error(w, err, http.StatusInternalServerError, "")
		return err
	}

	yamlFile, err := ioutil.ReadFile(filepath.Join(dir, ar.AccountConfigPath))
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
