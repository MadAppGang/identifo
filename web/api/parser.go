package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gopkg.in/go-playground/validator.v9"
)

const (
	// errInvalidSkip occurs when 'skip' URL variable cannot be converted to integer.
	errInvalidSkip = "skip value %s cannot be converted to integer"
	// errInvalidLimit occurs when 'limit' URL variable cannot be converted to integer.
	errInvalidLimit = "limit value %s cannot be converted to integer"
)

// MustParseJSON parses request body json data to the `out` struct.
// If error happens, writes it to ResponseWriter.
func (ar *Router) MustParseJSON(w http.ResponseWriter, r *http.Request, out interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(out); err != nil {
		ar.Error(w, ErrorAPIRequestBodyInvalid, http.StatusBadRequest, err.Error(), "Router.MustParseJSON.Decode")
		return err
	}

	validate := validator.New()
	if err := validate.Struct(out); err != nil {
		ar.Error(w, ErrorAPIRequestBodyInvalid, http.StatusBadRequest, err.Error(), "Router.MustParseJSON.Validate")
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

// parseWithValidParam parses withValid param.
// If it's empty returns false.
func (ar *Router) parseWithValidParam(r *http.Request) (bool, error) {
	withValidStr := r.URL.Query().Get("withValid")
	if withValidStr == "" {
		return false, nil
	}

	return strconv.ParseBool(withValidStr)
}
