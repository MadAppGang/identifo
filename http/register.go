package http

import (
	"net/http"
	"unicode"

	"github.com/madappgang/identifo/model"
)

/*
 * Password rules:
 * at least 7 letters
 * at least 1 number
 * at least 1 upper case
 * at least 1 special character
 */

//RegisterWithPassword register new user with password
func (ar *apiRouter) RegisterWithPassword() http.HandlerFunc {

	type registrationData struct {
		Username string                 `json:"username,omitempty" validate:"required,gte=6,lte=50"`
		Password string                 `json:"password,omitempty" validate:"required,gte=7,lte=50"`
		Profile  map[string]interface{} `json:"user_profile,omitempty"`
		Scope    []string               `json:"scope,omitempty"`
	}

	type registrationResponse struct {
		AccessToken  string     `json:"access_token,omitempty"`
		RefreshToken string     `json:"refresh_token,omitempty"`
		User         model.User `json:"user,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		//parse data
		d := registrationData{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		//validate password
		if err := strongPswd(d.Password); err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		//create new user
		user, err := ar.userStorage.AddUserByNameAndPassword(d.Username, d.Password, d.Profile)
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		//do login flow
		scopes, err := ar.userStorage.RequestScopes(user.ID(), d.Scope)
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		app := appFromContext(r.Context())
		if app == nil {
			ar.logger.Println("Error getting App")
			ar.Error(w, ErrorRequestInvalidAppID, http.StatusBadRequest, "")
			return
		}

		token, err := ar.tokenService.NewToken(user, scopes, app)
		if err != nil {
			ar.Error(w, err, http.StatusUnauthorized, "")
			return
		}

		tokenString, err := ar.tokenService.String(token)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		refreshString := ""
		//requesting offline access ?
		if contains(scopes, model.OfflineScope) {
			refresh, err := ar.tokenService.NewRefreshToken(user, scopes, app)
			if err != nil {
				ar.Error(w, err, http.StatusInternalServerError, "")
				return
			}
			refreshString, err = ar.tokenService.String(refresh)
			if err != nil {
				ar.Error(w, err, http.StatusInternalServerError, "")
				return
			}
		}

		user.Sanitize()

		result := registrationResponse{
			AccessToken:  tokenString,
			RefreshToken: refreshString,
			User:         user,
		}

		ar.ServeJSON(w, http.StatusOK, result)
	}
}

func strongPswd(pswd string) error {
	seven, number, uppper, _, invalid := verifyPassword(pswd)
	if invalid {
		return ErrorPasswordWrongSymbols
	} else if !seven {
		return ErrorPasswordShouldHave7Letter
	} else if !number {
		return ErrorPasswordNoNumbers
	} else if !uppper {
		return ErrorPasswordNoUppercase
	}
	return nil
}

func verifyPassword(s string) (sevenOrMore, number, upper, special, invalid bool) {
	letters := 0
	for _, s := range s {
		switch {
		case unicode.IsNumber(s):
			number = true
		case unicode.IsUpper(s):
			upper = true
			letters++
		case unicode.IsPunct(s) || unicode.IsSymbol(s):
			special = true
		case unicode.IsLetter(s) || s == ' ':
			letters++
		default:
			return false, false, false, false, true
		}
	}
	invalid = false
	sevenOrMore = letters >= 7
	return
}
