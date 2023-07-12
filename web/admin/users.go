package admin

import (
	"net/http"
	"strings"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/tools/xmaps"
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
			ar.LocalizedError(w, locale, http.StatusNotFound, l.ErrorStorageFindUserIDError, userID, err)
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
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAdminPanelNoSkipLimit, err)
			return
		}

		users, total, err := ar.server.Storages().UC.GetUsers(r.Context(), filterStr, skip, limit)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAdminPanelGetUsers, err)
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
	// registrationData is a request data to create new user form admin panel.
	type registrationData struct {
		Username          string `json:"username"`
		Email             string `json:"email"`
		GivenName         string `json:"given_name"`
		FamilyName        string `json:"family_name"`
		MiddleName        string `json:"middle_name"`
		Nickname          string `json:"nickname"`
		PreferredUsername string `json:"preferred_username"`
		PhoneNumber       string `json:"phone_number"`
		Password          string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		rd := registrationData{}
		if ar.mustParseJSON(w, r, &rd) != nil {
			return
		}

		um := model.User{}
		xmaps.CopyDstFields(rd, um)
		user, err := ar.server.Storages().UMC.CreateUserWithPassword(r.Context(), um, rd.Password)
		if err != nil { // this error is already localized.
			ar.Error(w, err)
			return
		}

		user = xmaps.CopyFields(user, model.UserFieldsetBasic.Fields())
		ar.ServeJSON(w, locale, http.StatusOK, user)
	}
}

// UpdateUser updates user in the database.
func (ar *Router) UpdateUser() http.HandlerFunc {
	type updateUserData struct {
		Username          *string `json:"username"`
		Email             *string `json:"email"`
		GivenName         *string `json:"given_name"`
		FamilyName        *string `json:"family_name"`
		MiddleName        *string `json:"middle_name"`
		Nickname          *string `json:"nickname"`
		PreferredUsername *string `json:"preferred_username"`
		PhoneNumber       *string `json:"phone_number"`
		Password          *string `json:"password"`
		Blocked           *bool   `json:"blocked"`
		BlockReason       *string `json:"block_reason"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")
		userID := getRouteVar("id", r)

		u := updateUserData{}
		if ar.mustParseJSON(w, r, &u) != nil {
			return
		}

		// update password if password is part of update process
		if u.Password != nil && len(*u.Password) > 0 {
			err := ar.server.Storages().UMC.UpdateUserPassword(r.Context(), userID, *u.Password)
			if err != nil {
				ar.Error(w, l.ErrorWithLocale(err, locale))
				return
			}
		}

		fields := xmaps.Filled(u)
		fields = xmaps.ContainsFields(u, fields)
		user := model.User{}
		xmaps.CopyDstFields(u, &user)
		user, err := ar.server.Storages().UMC.UpdateUser(r.Context(), user, fields)
		if err != nil {
			ar.HTTPError(w, err, http.StatusInternalServerError)
			return
		}
		user = xmaps.CopyFields(user, model.UserFieldsetBasic.Fields())
		ar.ServeJSON(w, locale, http.StatusOK, user)
	}
}

// DeleteUser deletes user from the database.
func (ar *Router) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")
		userID := getRouteVar("id", r)

		if err := ar.server.Storages().UMC.DeleteUser(r.Context(), userID); err != nil {
			ar.HTTPError(w, err, http.StatusInternalServerError)
			return
		}

		ar.ServeJSON(w, locale, http.StatusOK, nil)
	}
}

func (ar *Router) GenerateNewResetTokenUser() http.HandlerFunc {
	// password reset data from admin panel.
	type passwordResetData struct {
		UserID string `json:"user_id,omitempty"`
		AppID  string `json:"app_id,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		resetData := passwordResetData{}
		if ar.mustParseJSON(w, r, &resetData) != nil {
			return
		}

		red, err := ar.server.Storages().UMC.SendPasswordResetEmail(r.Context(), resetData.UserID, resetData.AppID)
		if err != nil {
			// TODO: generate proper localized error with details
			ar.HTTPError(w, err, http.StatusInternalServerError)
		}

		ar.ServeJSON(w, locale, http.StatusOK, red)
	}
}
