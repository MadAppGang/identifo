package http

import (
	"net/http"
)

// type loginResponse struct {
// 	Answer string    `json:"answer,omitempty"`
// 	Date   time.Time `json:"date,omitempty"`
// }

//LoginWithPassword - login user with email and password
func (ar *apiRouter) LoginWithPassword() http.HandlerFunc {

	type loginData struct {
		Username    string   `json:"username,omitempty" validate:"required,gte=6,lte=130"`
		Password    string   `json:"password,omitempty" validate:"required,gte=6,lte=130"`
		DeviceToken string   `json:"device_token,omitempty"`
		Scopes      []string `json:"scopes,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		d := loginData{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		user, err := ar.userStorage.UserByNamePassword(d.Username, d.Password)
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		scopes, err := ar.userStorage.RequestScopes(user.ID(), d.Scopes)
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		token, err := ar.tokenService.NewToken(user, scopes)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		ar.ServeJSON(w, http.StatusOK, map[string]string{"token": token.String()})
	}
}
