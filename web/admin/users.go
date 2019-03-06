package admin

import (
	"net/http"
	"strings"

	"github.com/madappgang/identifo/model"
)

const (
	defaultUserSkip  = 0
	defaultUserLimit = 20
)

// GetUser fetches user by ID from the database.
func (ar *Router) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getRouteVar("id", r)

		user, err := ar.userStorage.UserByID(userID)
		if err != nil {
			if err == model.ErrorNotFound {
				ar.Error(w, err, http.StatusNotFound, "")
			} else {
				ar.Error(w, err, http.StatusInternalServerError, "")
			}
			return
		}

		ar.ServeJSON(w, http.StatusOK, user)
		return
	}
}

// FetchUsers fetches users from the database.
func (ar *Router) FetchUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filterStr := strings.TrimSpace(r.URL.Query().Get("search"))

		limit, skip, err := ar.parseSkipAndLimit(r, defaultUserSkip, defaultUserLimit, 0)
		if err != nil {
			ar.Error(w, ErrorWrongInput, http.StatusBadRequest, "")
			return
		}

		users, err := ar.userStorage.FetchUsers(filterStr, skip, limit)
		if err != nil {
			ar.Error(w, ErrorInternalError, http.StatusInternalServerError, "")
			return
		}

		ar.ServeJSON(w, http.StatusOK, users)
		return
	}
}

// CreateUser registers new user.
func (ar *Router) CreateUser() http.HandlerFunc {
	type registrationData struct {
		Username string                 `json:"username,omitempty" validate:"required,gte=6,lte=50"`
		Password string                 `json:"password,omitempty" validate:"required,gte=7,lte=50"`
		Profile  map[string]interface{} `json:"user_profile,omitempty"`
		Scope    []string               `json:"scope,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		d := registrationData{}
		if ar.mustParseJSON(w, r, &d) != nil {
			return
		}

		if err := model.StrongPswd(d.Password); err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		user, err := ar.userStorage.AddUserByNameAndPassword(d.Username, d.Password, d.Profile)
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, "")
			return
		}

		user.Sanitize()

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
		return
	}
}
