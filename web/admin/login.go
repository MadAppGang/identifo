package admin

import (
	"net/http"
	"os"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
)

const (
	AoothAdminUsername = "AOOTH_ADMIN_USERNAME"
	AoothAdminPassword = "AOOTH_ADMIN_PASSWORD"
)

type adminLoginData struct {
	Login    string `json:"email"`
	Password string `json:"password"`
}

type adminLoginResponse struct {
	Token string `json:"token"`
}

// Login logins admin with admin name and password.
func (ar *Router) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		ld := adminLoginData{}
		if err := ar.mustParseJSON(w, r, &ld); err != nil {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPIJsonParseError, err.Error())
			return
		}

		err := checkCredentials(ld.Login, ld.Password)
		if err != nil {
			ar.HTTPError(w, l.ErrorWithLocale(err, locale), http.StatusUnauthorized)
			return
		}
		// create JWT token
		u := model.User{ID: model.RootUserID.String()}
		token, err := ar.server.Services().Token.NewToken(model.TokenTypeManagement, u, nil, nil, nil)
		if err != nil {
			ar.HTTPError(w, l.ErrorWithLocale(err, locale), http.StatusInternalServerError)
		}

		ar.ServeJSON(w, locale, http.StatusOK, adminLoginResponse{Token: token.Raw})
	}
}

func checkCredentials(login, password string) error {
	u := os.Getenv(AoothAdminUsername)
	p := os.Getenv(AoothAdminUsername)

	if len(u) == 0 || len(p) == 0 {
		return l.ErrorAdminPanelAdminCredentialsNotSet
	}

	if login != u || password != p {
		return l.ErrorAdminPanelAdminCredentialsMismatch
	}

	return nil
}
