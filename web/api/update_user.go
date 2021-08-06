package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/madappgang/identifo/model"
)

// UpdateUser allows to change user login and password.
func (ar *Router) UpdateUser() http.HandlerFunc {
	type updateResponse struct {
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		d := updateData{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}
		userID := tokenFromContext(r.Context()).UserID()
		user, err := ar.server.Storages().User.UserByID(userID)
		if err != nil {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusUnauthorized, err.Error(), "UpdateUser.UserByID")
			return
		}

		if err := d.validate(user); err != nil {
			ar.Error(w, ErrorAPIRequestBodyParamsInvalid, http.StatusBadRequest, err.Error(), "UpdateUser.validate")
			return
		}
		// Check that new username is not taken.
		if d.updateUsername && ar.server.Storages().User.UserExists(d.NewUsername) {
			ar.Error(w, ErrorAPIUsernameTaken, http.StatusBadRequest, "", "UpdateUser.updateUsername && userStorage.UserExists")
			return
		}

		// Check that email is not taken.
		if d.updateEmail {
			if _, err := ar.server.Storages().User.UserByEmail(d.NewEmail); err == nil {
				ar.Error(w, ErrorAPIEmailTaken, http.StatusBadRequest, "", "UpdateUser.updateEmail && UserByEmail")
				return
			}
		}

		// Check that phone is not taken.
		if d.updatePhone {
			if _, err := ar.server.Storages().User.UserByPhone(d.NewPhone); err == nil {
				ar.Error(w, ErrorAPIEmailTaken, http.StatusBadRequest, "", "UpdateUser.updatePhone && UserByPhone")
				return
			}
		}

		// Update password.
		if d.updatePassword {
			// Check old password.
			if err := ar.server.Storages().User.CheckPassword(user.ID, d.OldPassword); err != nil {
				ar.Error(w, ErrorAPIRequestBodyOldPasswordInvalid, http.StatusBadRequest, err.Error(), "UpdateUser.updatePassword && CheckPassword")
				return
			}

			// Save new password.
			err = ar.server.Storages().User.ResetPassword(user.ID, d.NewPassword)
			if err != nil {
				ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, "Reset password. Error: "+err.Error(), "UpdateUser.ResetPassword")
				return
			}

			// Refetch user with new password hash.
			if user, err = ar.server.Storages().User.UserByUsername(user.Username); err != nil {
				ar.Error(w, ErrorAPIRequestBodyOldPasswordInvalid, http.StatusBadRequest, err.Error(), "UpdateUser.RefetchUser")
				return
			}
		}

		// Change username if user specified new one.
		if d.updateUsername {
			user.Username = d.NewUsername
			user = user.Deanonimized()
		}

		if d.updateEmail {
			user.Email = d.NewEmail
		}

		if d.updatePhone {
			user.Phone = d.NewPhone
		}

		if d.updateUsername || d.updateEmail || d.updatePhone {
			if _, err = ar.server.Storages().User.UpdateUser(userID, user); err != nil {
				ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, "unable to update username or email. Error: "+err.Error(), " UpdateUser.UpdateUser ")
				return
			}
		}

		// Prepare response.
		updatedFields := []string{}
		if d.updateUsername {
			updatedFields = append(updatedFields, "username")
		}
		if d.updateEmail {
			updatedFields = append(updatedFields, "email")
		}
		if d.updatePhone {
			updatedFields = append(updatedFields, "phone")
		}
		if d.updatePassword {
			updatedFields = append(updatedFields, "password")
		}

		msg := "Nothing changed."
		if len(updatedFields) > 0 {
			updatedFields[0] = strings.Title(updatedFields[0])
			msg = strings.Join(updatedFields, ", ") + " changed. "
		}
		response := updateResponse{
			Message: msg,
		}
		ar.ServeJSON(w, http.StatusOK, response)
	}
}

type updateData struct {
	NewEmail       string `json:"new_email"`
	NewPhone       string `json:"new_phone,omitempty"`
	NewUsername    string `json:"new_username,omitempty"`
	NewPassword    string `json:"new_password,omitempty"`
	OldPassword    string `json:"old_password,omitempty"`
	updatePassword bool
	updateEmail    bool
	updatePhone    bool
	updateUsername bool
}

func (d *updateData) validate(user model.User) error {
	if d.NewUsername != "" && user.Username != d.NewUsername {
		d.updateUsername = true
	}
	if d.NewEmail != "" && user.Email != d.NewEmail {
		d.updateEmail = true
	}
	if d.NewPhone != "" && user.Phone != d.NewPhone {
		d.updatePhone = true
	}
	if d.NewPassword != "" && d.NewPassword != d.OldPassword {
		d.updatePassword = true
	}

	if d.updatePassword {
		if d.OldPassword == "" {
			return errors.New("Old password is not specified. ")
		}
		// validate password
		if err := model.StrongPswd(d.NewPassword); err != nil {
			return errors.New("New password is not strong enough. ")
		}
	}

	if d.updateEmail && !model.EmailRegexp.MatchString(d.NewEmail) {
		return errors.New("Email is not valid. ")
	}
	return nil
}
