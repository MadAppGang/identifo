package admin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/madappgang/identifo/model"
)

const (
	defaultUserSkip  = 0
	defaultUserLimit = 20
)

type registrationData struct {
	Username   string        `json:"username,omitempty"`
	Email      string        `json:"email,omitempty"`
	Phone      string        `json:"phone,omitempty"`
	Password   string        `json:"password,omitempty"`
	AccessRole string        `json:"access_role,omitempty"`
	Scope      []string      `json:"scope,omitempty"`
	TFAInfo    model.TFAInfo `json:"tfa_info,omitempty"`
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

		user, err := ar.server.Storages().User.UserByID(userID)
		if err != nil {
			if err == model.ErrUserNotFound {
				ar.Error(w, err, http.StatusNotFound, "")
			} else {
				ar.Error(w, err, http.StatusInternalServerError, "")
			}
			return
		}

		user = user.Sanitized()
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

		user, err := ar.server.Storages().User.AddUserWithPassword(model.User{Username: rd.Username, Email: rd.Email, Phone: rd.Phone}, rd.Password, rd.AccessRole, false)
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		user.TFAInfo = rd.TFAInfo

		user, err = ar.server.Storages().User.UpdateUser(user.ID, user)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "Setting TFA data")
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
