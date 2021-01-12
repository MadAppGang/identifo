package admin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/plugin/shared"
	"github.com/madappgang/identifo/proto"
	"github.com/madappgang/identifo/proto/extensions"
)

const (
	defaultUserSkip  = 0
	defaultUserLimit = 20
)

type registrationData struct {
	Username   string   `json:"username,omitempty"`
	Password   string   `json:"password,omitempty"`
	AccessRole string   `json:"access_role,omitempty"`
	Scope      []string `json:"scope,omitempty"`
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

		user, err := ar.userStorage.UserByID(userID)
		if err != nil {
			if err == shared.ErrUserNotFound {
				ar.Error(w, err, http.StatusNotFound, "")
			} else {
				ar.Error(w, err, http.StatusInternalServerError, "")
			}
			return
		}

		extensions.SanitizeUser(user)
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

		users, total, err := ar.userStorage.FetchUsers(filterStr, skip, limit)
		if err != nil {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, err.Error())
			return
		}
		for _, user := range users {
			extensions.SanitizeUser(user)
		}

		searchResponse := struct {
			Users []*proto.User `json:"users"`
			Total int           `json:"total"`
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

		user, err := ar.userStorage.AddUserByNameAndPassword(rd.Username, rd.Password, rd.AccessRole, false)
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		extensions.SanitizeUser(user)
		ar.ServeJSON(w, http.StatusOK, user)
	}
}

// UpdateUser updates user in the database.
func (ar *Router) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getRouteVar("id", r)

		u := new(proto.User)
		if ar.mustParseJSON(w, r, u) != nil {
			return
		}

		user, err := ar.userStorage.UpdateUser(userID, u)
		if err != nil {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, "")
			return
		}

		ar.logger.Printf("User %s updated", userID)

		extensions.SanitizeUser(user)
		ar.ServeJSON(w, http.StatusOK, user)
	}
}

// DeleteUser deletes user from the database.
func (ar *Router) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getRouteVar("id", r)
		if err := ar.userStorage.DeleteUser(userID); err != nil {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, "")
			return
		}

		ar.logger.Printf("User %s deleted", userID)
		ar.ServeJSON(w, http.StatusOK, nil)
	}
}
