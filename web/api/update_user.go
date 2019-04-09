package api

import (
	"errors"
	"net/http"

	"github.com/madappgang/identifo/model"
)

// UpdateUser allows to change user login and password.
func (ar *Router) UpdateUser() http.HandlerFunc {
	type updateData struct {
		NewUsername string `json:"new_username,omitempty"`
		NewPassword string `json:"new_password,omitempty"`
		OldPassword string `json:"old_password,omitempty"`
	}

	type updateResponse struct {
		Message string `json:"message"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		d := updateData{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}
		userID := tokenFromContext(r.Context()).UserID()
		user, err := ar.userStorage.UserByID(userID)
		if err != nil {
			ar.Error(w, err, http.StatusUnauthorized, "Not authorized")
			return
		}

		// check that new username is not taken.
		if d.NewUsername != "" && user.Name() != d.NewUsername {
			if ar.userStorage.UserExists(d.NewUsername) {
				ar.Error(w, errors.New("Username is busy. "), http.StatusBadRequest, "Username is busy. Try to choose another one.")
				return
			}
		}

		usernameChanged := false
		passwordChanged := false

		// change password if user specified new one.
		if d.NewPassword != "" {
			if d.OldPassword == "" {
				ar.Error(w, err, http.StatusBadRequest, "Old password is not specified.")
				return
			}
			// validate password
			if err := model.StrongPswd(d.NewPassword); err != nil {
				ar.Error(w, err, http.StatusBadRequest, "New password is not strong enough.")
				return
			}

			// check old password.
			_, err := ar.userStorage.UserByNamePassword(user.Name(), d.OldPassword)
			if err != nil {
				ar.Error(w, err, http.StatusBadRequest, "Invalid old password.")
				return
			}

			// save new password.
			err = ar.userStorage.ResetPassword(user.ID(), d.NewPassword)
			if err != nil {
				ar.Error(w, err, http.StatusInternalServerError, "Unable to update password. Please try again.")
				return
			}
			passwordChanged = true
		}

		// change username if user specified new one.
		if d.NewUsername != "" && user.Name() != d.NewUsername {
			err = ar.userStorage.ResetUsername(userID, d.NewUsername)
			if err != nil {
				ar.Error(w, err, http.StatusBadRequest, "Username is taken. Try to choose another one.")
				return
			}
			usernameChanged = true
		}

		// prepare response
		msg := "Nothing changed."
		if usernameChanged && passwordChanged {
			msg = "Username and password changed."
		} else if usernameChanged {
			msg = "Username changed."
		} else if passwordChanged {
			msg = "Password changed."
		}

		response := updateResponse{
			Message: msg,
		}
		ar.ServeJSON(w, http.StatusOK, response)
	}

}
