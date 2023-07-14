package api

import (
	"net/http"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
)

// LoginWithPassword logs user in with email and password.
func (ar *Router) LoginWithPassword() http.HandlerFunc {
	type loginData struct {
		Email    string   `json:"email"`
		Phone    string   `json:"phone"`
		Username string   `json:"username"`
		Password string   `json:"password"`
		Device   string   `json:"device"`
		Scopes   []string `json:"scopes"`
		OS       string   `json:"os"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")
		agent := r.Header.Get("User-Agent")

		ld := loginData{}
		if ar.MustParseJSON(w, r, &ld) != nil {
			return
		}

		if err := ld.validate(); err != nil {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyInvalidError, err)
			return
		}

		if err := ar.checkSupportedWays(ld.login); err != nil {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.APIAPPUsernameLoginNotSupported)
			return
		}

		var err error
		user := model.User{}

		if len(ld.Email) > 0 {
			user, err = ar.server.Storages().User.UserByEmail(ld.Email)
		} else if len(ld.Phone) > 0 {
			user, err = ar.server.Storages().User.UserByPhone(ld.Phone)
		} else if len(ld.Username) > 0 {
			user, err = ar.server.Storages().User.UserByUsername(ld.Username)
		}

		if err != nil {
			ar.LocalizedError(w, locale, http.StatusUnauthorized, l.ErrorAPIRequestIncorrectLoginOrPassword)
			return
		}

		if err = ar.server.Storages().User.CheckPassword(user.ID, ld.Password); err != nil {
			// return this error to hide the existence of the user.
			ar.LocalizedError(w, locale, http.StatusUnauthorized, l.ErrorAPIRequestIncorrectLoginOrPassword)
			return
		}

		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPIAPPNoAPPInContext)
			return
		}

		authResult, err := ar.loginFlow(app, user, ld.Scopes)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAPILoginError, err)
			return
		}

		ar.ServeJSON(w, locale, http.StatusOK, authResult)
	}
}

// IsLoggedIn is for checking whether user is logged in or not.
// In fact, all needed work is done in Token middleware.
// If we reached this code, user is logged in (presented valid and not blacklisted access token).
func (ar *Router) IsLoggedIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")
		ar.ServeJSON(w, locale, http.StatusOK, nil)
	}
}

// GetUser return current user info with sanitized tfa
func (ar *Router) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		userID := tokenFromContext(r.Context()).UserID()
		user, err := ar.server.Storages().User.UserByID(r.Context(), userID)
		if err != nil {
			ar.Error(w, locale, http.StatusUnauthorized, l.ErrorStorageFindUserIDError, userID, err)
			return
		}
		ar.ServeJSON(w, locale, http.StatusOK, user.SanitizedTFA())
	}
}
