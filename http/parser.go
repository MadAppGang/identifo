package http

import (
	"encoding/json"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

//MustParseJSON parses request body json data to the `out` struct.
//Writes error to ResponseWriter on error
func (ar *apiRouter) MustParseJSON(w http.ResponseWriter, r *http.Request, out interface{}) error {
	//parse structure

	if err := json.NewDecoder(r.Body).Decode(out); err != nil {
		ar.Error(w, err, http.StatusBadRequest, "")
		return err
	}

	//validate structure
	validate := validator.New()
	if err := validate.Struct(out); err != nil {
		ar.Error(w, err, http.StatusBadRequest, "")
		return err
	}

	return nil
}
