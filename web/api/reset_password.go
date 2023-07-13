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

// ResetPassword handles password reset form submission (POST request).
// this method exchanges reset JWT token to access JWT token
// getting the new password and saving it in the database.
func (ar *Router) ResetPassword() http.HandlerFunc {
	type newPassword struct {
		Password string `json:"password,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		d := newPassword{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}
		if err := model.StrongPswd(d.Password); err != nil {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequestPasswordWeak, err)
			return
		}

		accessTokenBytes, ok := r.Context().Value(model.TokenRawContextKey).([]byte)
		if !ok {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIContextNoToken)
			return
		}

		// Get userID from token and update user with this ID.
		userID, err := ar.getTokenSubject(string(accessTokenBytes))
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAPIRequestTokenSubError, err)
			return
		}

		user, err := ar.server.Storages().User.UserByID(userID)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusUnauthorized, l.ErrorStorageFindUserIDError, userID, err)
			return
		}

		// Save new password.
		if err := ar.server.Storages().User.ResetPassword(user.ID, d.Password); err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorStorageResetPasswordUserError, user.ID, err)
			return
		}

		result := map[string]string{"result": "ok"}
		ar.ServeJSON(w, locale, http.StatusOK, result)
	}
}
