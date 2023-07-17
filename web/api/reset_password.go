package api

import (
	"errors"
	"net/http"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
)

// RequestResetPassword requests password reset
// now we support reset password only with JWT reset token send to email
// if user does not have email - we could not reset the password
// we need to find user by secondary ID (phone, username, email)
// if user does not exist - we just silently return ok
// if user exists - we send reset password email
// if there is not email for user - we return error, as it's configuration error
func (ar *Router) RequestResetPassword() http.HandlerFunc {
	type resetRequest struct {
		Phone        string `json:"phone"`
		Email        string `json:"email"`
		Username     string `json:"username"`
		ResetPageURL string `json:"reset_page_url,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")
		// agent := r.Header.Get("User-Agent")
		app := middleware.AppFromContext(r.Context())

		d := resetRequest{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		idType := model.AuthIdentityTypeNone
		idValue := ""
		if len(d.Email) > 0 {
			idType = model.AuthIdentityTypeEmail
			idValue = d.Email
		} else if len(d.Phone) > 0 {
			idType = model.AuthIdentityTypePhone
			idValue = d.Phone
		} else if len(d.Username) > 0 {
			idType = model.AuthIdentityTypeUsername
			idValue = d.Username
		}

		// nothing filled
		if idType == model.AuthIdentityTypeNone {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorLoginDataEmpty)
			return
		}

		respok := map[string]string{"result": "ok"}
		user, err := ar.server.Storages().UC.UserBySecondaryID(r.Context(), idType, idValue)
		if err != nil && !errors.Is(err, l.ErrorUserNotFound) {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.LocalizedString(err.Error()))
			return
		} else if !errors.Is(err, l.ErrorUserNotFound) {
			// if not user - just report ok for security reasons
			ar.ServeJSON(w, locale, http.StatusOK, respok)
			return
		}

		// we have user
		_, err = ar.server.Storages().UMC.SendPasswordResetEmail(r.Context(), user.ID, app.ID)
		if err != nil {
			// TODO: generate proper localized error with details
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.LocalizedString(err.Error()))
			return
		}

		ar.ServeJSON(w, locale, http.StatusOK, respok)
	}
}

// ChangePassword handles password reset form submission (POST request).
// this method exchanges reset JWT token to access JWT token
// getting the new password and saving it in the database.
func (ar *Router) ChangePassword() http.HandlerFunc {
	type newPassword struct {
		Password string   `json:"password,omitempty"`
		Scopes   []string `json:"scopes,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		d := newPassword{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		token := tokenFromContext(r.Context())
		if token == nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAPIContextNoToken)
			return
		}

		// Validate just in case, middleware should do it for us.
		if token.Type() != model.TokenTypeReset {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPIRequestTokenInvalid)
			return
		}

		// Let's update the password.
		err := ar.server.Storages().UMC.UpdateUserPassword(r.Context(), "", d.Password)
		if err != nil {
			ar.HTTPError(w, err, http.StatusInternalServerError)
			return
		}

		app := middleware.AppFromContext(r.Context())
		u, err := ar.server.Storages().UC.UserByID(r.Context(), token.UserID())
		if err != nil {
			ar.HTTPError(w, err, http.StatusInternalServerError)
			return
		}

		loginResponse, err := ar.server.Storages().UC.GetJWTTokens(r.Context(), app, u, d.Scopes)
		if err != nil {
			ar.HTTPError(w, err, http.StatusInternalServerError)
			return
		}
		ar.ServeJSON(w, locale, http.StatusOK, loginResponse)
	}
}
