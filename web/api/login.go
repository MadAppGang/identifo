package api

import (
	"net/http"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/tools/xmaps"
	"github.com/madappgang/identifo/v2/web/middleware"
)

// LoginWithPassword logs user in with email and password.
func (ar *Router) LoginWithPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")
		// agent := r.Header.Get("User-Agent")
		// ip := middleware.IPFromContext(r.Context())

		ld := model.LoginRequest{}
		if ar.MustParseJSON(w, r, &ld) != nil {
			return
		}

		strategy, id := ld.Strategy()
		uc := ar.server.Storages().UC

		user, err := uc.UserBySecondaryID(r.Context(), strategy.Identity, id)
		if err != nil {
			ar.HTTPError(w, l.ErrorWithLocale(err, locale), http.StatusUnauthorized)
			return
		}

		err = uc.VerifyPassword(r.Context(), user, ld.Password)
		if err != nil {
			ar.HTTPError(w, l.ErrorWithLocale(err, locale), http.StatusUnauthorized)
			return
		}

		app := middleware.AppFromContext(r.Context())
		response, err := ar.server.Storages().UC.GetJWTTokens(r.Context(), app, user, ld.Scopes)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusUnauthorized, l.ErrorAPILoginError, err)
			return
		}

		ar.ServeJSON(w, locale, http.StatusOK, response)
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

// GetUser return current user info with sanitized data
func (ar *Router) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")
		scopes := r.URL.Query()["scopes"]

		// add basic profile information to scope
		scopes = append(scopes, model.ProfileScope)
		fields := model.FieldsetForScopes(scopes)

		userID := tokenFromContext(r.Context()).UserID()
		user, err := ar.server.Storages().User.UserByID(r.Context(), userID)
		if err != nil {
			ar.HTTPError(w, l.ErrorWithLocale(err, locale), http.StatusUnauthorized)
			return
		}
		u := xmaps.CopyFields(user, fields)
		ar.ServeJSON(w, locale, http.StatusOK, u)
	}
}
