package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo/v2/l"
	"gopkg.in/go-playground/validator.v9"
)

const (
	// errInvalidSkip occurs when 'skip' URL variable cannot be converted to integer.
	errInvalidSkip = "skip value %s cannot be converted to integer"
	// errInvalidLimit occurs when 'limit' URL variable cannot be converted to integer.
	errInvalidLimit = "limit value %s cannot be converted to integer"
)

// mustParseJSON parses request body json data to the `out` interface and then validates it.
// Writes error to ResponseWriter if error happens.
func (ar *Router) mustParseJSON(w http.ResponseWriter, r *http.Request, out interface{}) error {
	locale := r.Header.Get("Accept-Language")
	if err := json.NewDecoder(r.Body).Decode(out); err != nil {
		ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPIJsonParseError, err)
		return err
	}

	validate := validator.New()
	if err := validate.Struct(out); err != nil {
		ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPIJsonParseError, err)
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
			return 0, 0, fmt.Errorf(errInvalidLimit, limitStr)
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

// parseWithArchivedParam parses withValid param.
// If it's empty returns false.
func (ar *Router) parseWithArchivedParam(r *http.Request) (bool, error) {
	withArchivedStr := r.URL.Query().Get("withArchived")
	if withArchivedStr == "" {
		return false, nil
	}

	return strconv.ParseBool(withArchivedStr)
}
