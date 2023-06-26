package admin

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	l "github.com/madappgang/identifo/v2/localization"
	"github.com/madappgang/identifo/v2/model"
)

const (
	defaultUserSkip  = 0
	defaultUserLimit = 20
)

// GetUser fetches user by ID from the database.
func (ar *Router) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")
		userID := getRouteVar("id", r)

		user, err := ar.server.Storages().UC.UserByID(r.Context(), userID)
		if err != nil {
			ar.Error(w, locale, http.StatusNotFound, l.ErrorStorageFindUserIDError, userID, err)
			return
		}
		ar.ServeJSON(w, locale, http.StatusOK, user)
	}
}

// FetchUsers fetches users from the database.
func (ar *Router) FetchUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")
		filterStr := strings.TrimSpace(r.URL.Query().Get("search"))

		skip, limit, err := ar.parseSkipAndLimit(r, defaultUserSkip, defaultUserLimit, 0)
		if err != nil {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAdminPanelNoSkipAndLimitParams, err.Error())
			return
		}

		users, total, err := ar.server.Storages().UC.GetUsers(r.Context(), filterStr, skip, limit)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorAdminPanelErrorGettingUserList, err.Error())
			return
		}

		searchResponse := struct {
			Users []model.User `json:"users"`
			Skip  int          `json:"skip"`
			Limit int          `json:"limit"`
			Total int          `json:"total"`
		}{
			Users: users,
			Total: total,
			Skip:  skip,
			Limit: limit,
		}

		ar.ServeJSON(w, locale, http.StatusOK, &searchResponse)
	}
}

// CreateUser registers new user.
func (ar *Router) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rd := registrationData{}
		if ar.mustParseJSON(w, r, &rd) != nil {
			return
		}

		um := model.User{}
		model.CopyDstFields(rd, um)
		user, err := ar.server.Storages().User.AddUserWithPassword(um, rd.Password, false)
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		user = user.Sanitized()
		ar.ServeJSON(w, http.StatusOK, user)
	}
}

// UpdateUser updates user in the database.
func (ar *Router) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getRouteVar("id", r)

		u := model.User{}
		if ar.mustParseJSON(w, r, &u) != nil {
			return
		}

		existing, err := ar.server.Storages().User.UserByID(userID)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		if u.TFAInfo.IsEnabled == existing.TFAInfo.IsEnabled {
			u.TFAInfo = existing.TFAInfo
		}

		if !u.TFAInfo.IsEnabled {
			u.TFAInfo = model.TFAInfo{
				IsEnabled: false,
			}
		}

		// update password if password is part of update process
		if len(u.Pswd) > 0 {
			if err := model.StrongPswd(u.Pswd); err != nil {
				ar.Error(w, err, http.StatusBadRequest, "")
				return
			}

			err := ar.server.Storages().User.ResetPassword(userID, u.Pswd)
			if err != nil {
				ar.Error(w, err, http.StatusInternalServerError, "")
				return
			}
			u.Pswd = ""
		}

		user, err := ar.server.Storages().User.UpdateUser(userID, u)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		ar.logger.Printf("User %s updated", userID)

		user = user.Sanitized()
		ar.ServeJSON(w, http.StatusOK, user)
	}
}

// DeleteUser deletes user from the database.
func (ar *Router) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getRouteVar("id", r)
		if err := ar.server.Storages().User.DeleteUser(userID); err != nil {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, "")
			return
		}

		ar.logger.Printf("User %s deleted", userID)
		ar.ServeJSON(w, http.StatusOK, nil)
	}
}

func (ar *Router) GenerateNewResetTokenUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resetData := passwordResetData{}
		if ar.mustParseJSON(w, r, &resetData) != nil {
			return
		}

		user, err := ar.server.Storages().User.UserByID(resetData.UserID)
		if err != nil {
			if err == model.ErrUserNotFound {
				ar.Error(w, err, http.StatusNotFound, "")
			} else {
				ar.Error(w, err, http.StatusInternalServerError, "")
			}
			return
		}

		resetToken, err := ar.server.Services().Token.NewResetToken(user.ID)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		resetTokenString, err := ar.server.Services().Token.String(resetToken)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		query := fmt.Sprintf("appId=%s&token=%s", resetData.AppID, resetTokenString)

		u := &url.URL{
			Scheme:   ar.Host.Scheme,
			Host:     ar.Host.Host,
			Path:     model.DefaultLoginWebAppSettings.ResetPasswordURL,
			RawQuery: query,
		}
		uu := &url.URL{Scheme: u.Scheme, Host: u.Host, Path: u.Path}

		resetEmailData := ResetEmailData{
			Token: resetTokenString,
			URL:   u.String(),
			Host:  uu.String(),
		}

		app, err := ar.server.Storages().App.AppByID(resetData.AppID)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}

		if err = ar.server.Services().Email.SendTemplateEmail(
			model.EmailTemplateTypeResetPassword,
			app.GetCustomEmailTemplatePath(),
			"Reset Password",
			user.Email,
			model.EmailData{
				User: user,
				Data: resetEmailData,
			},
		); err != nil {
			ar.Error(
				w,
				err,
				http.StatusInternalServerError,
				"Email sending error: "+err.Error(),
			)
			return
		}

		ar.ServeJSON(w, http.StatusOK, resetEmailData)
	}
}
