package admin

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/madappgang/identifo/v2/model"
)

const (
	defaultUserSkip  = 0
	defaultUserLimit = 20
)

type registrationData struct {
	Username   string `json:"username,omitempty"`
	Email      string `json:"email,omitempty"`
	FullName   string `json:"full_name,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Password   string `json:"pswd,omitempty"`
	AccessRole string `json:"access_role,omitempty"`
}

type passwordResetData struct {
	UserID       string `json:"user_id,omitempty"`
	AppID        string `json:"app_id,omitempty"`
	ResetPageURL string `json:"reset_page_url,omitempty"`
}

type resetEmailData struct {
	User  model.User
	Token string
	URL   string
	Host  string
}

func (rd *registrationData) validate() error {
	if usernameLen := len(rd.Username); usernameLen < 6 || usernameLen > 50 {
		return fmt.Errorf("Incorrect username length %d, expected a number between 6 and 50", usernameLen)
	}
	if pswdLen := len(rd.Password); pswdLen < 6 || pswdLen > 50 {
		return fmt.Errorf("Incorrect password length %d, expected a number between 6 and 50", pswdLen)
	}
	return nil
}

// GetUser fetches user by ID from the database.
func (ar *Router) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getRouteVar("id", r)

		ar.ServeJSON(w, http.StatusOK, user)
	}
}

// FetchUsers fetches users from the database.
func (ar *Router) FetchUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filterStr := strings.TrimSpace(r.URL.Query().Get("search"))

		skip, limit, err := ar.parseSkipAndLimit(r, defaultUserSkip, defaultUserLimit, 0)
		if err != nil {
			ar.Error(w, ErrorWrongInput, http.StatusBadRequest, "")
			return
		}

		users, total, err := ar.server.Storages().User.FetchUsers(filterStr, skip, limit)
		if err != nil {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, "")
			return
		}
		for i, user := range users {
			users[i] = user.Sanitized()
		}

		searchResponse := struct {
			Users []model.User `json:"users"`
			Total int          `json:"total"`
		}{
			Users: users,
			Total: total,
		}

		ar.ServeJSON(w, http.StatusOK, &searchResponse)
	}
}

// CreateUser registers new user.
func (ar *Router) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rd := registrationData{}
		if ar.mustParseJSON(w, r, &rd) != nil {
			return
		}

		if err := rd.validate(); err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		if err := model.StrongPswd(rd.Password); err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		// Create new user.
		um := model.User{
			Username: rd.Username,
			FullName: rd.FullName,
			Email:    rd.Email,
			Phone:    rd.Phone,
			TFAInfo:  rd.TFAInfo,
			Scopes:   rd.Scopes, // we are creating user from admin panel - we can set any scopes we want
		}

		user, err := ar.server.Storages().User.AddUserWithPassword(um, rd.Password, rd.AccessRole, false)
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
