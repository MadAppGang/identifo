package api

import (
	"encoding/json"
	"net/http"

	l "github.com/madappgang/identifo/v2/localization"
	"gopkg.in/go-playground/validator.v9"
)

// MustParseJSON parses request body json data to the `out` struct.
// If error happens, writes it to ResponseWriter.
func (ar *Router) MustParseJSON(w http.ResponseWriter, r *http.Request, out interface{}) error {
	locale := r.Header.Get("Accept-Language")

	if err := json.NewDecoder(r.Body).Decode(out); err != nil {
		ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyInvalidError, err)
		return err
	}

	validate := validator.New()
	if err := validate.Struct(out); err != nil {
		ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyInvalidError, err)
		return err
	}

	return nil
}
