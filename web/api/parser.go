package api

import (
	"encoding/json"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
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
