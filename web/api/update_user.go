package api

import (
	"net/http"
	"time"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
)

// UpdateUser allows to change user login and password.
func (ar *Router) UpdateUser() http.HandlerFunc {
	type updateUserRequestData struct {
		ID string `json:"id"`

		Username          *string `json:"username"` // it is a nickname for login purposes
		Email             *string `json:"email"`
		GivenName         *string `json:"given_name"`
		FamilyName        *string `json:"family_name"`
		MiddleName        *string `json:"middle_name"`
		Nickname          *string `json:"nickname"`
		PreferredUsername *string `json:"preferred_username"`
		PhoneNumber       *string `json:"phone_number"`
		// TODO: add update password process to it as part of update process.
		// Password          *string `json:"password"`

		// oidc claims
		// https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims
		Profile  *string                      `json:"profile"`
		Picture  *string                      `json:"picture"`
		Website  *string                      `json:"website"`
		Gender   *string                      `json:"gender"`
		Birthday *time.Time                   `json:"birthday"`
		Timezone *string                      `json:"timezone"`
		Locale   *string                      `json:"locale"`
		Address  map[string]model.UserAddress `json:"address"` // addresses for home, work, etc
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		d := model.User{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		userID := tokenFromContext(r.Context()).UserID()
		// only user himself could update his data. We force userID field to requester one.
		// if Admin needs update someones data - use management data for that.
		// or maybe we can handle that by roles in the future?
		d.ID = userID

		u := model.User{}
		err := model.CopyDstFields(d, &u)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAPIRequestBodyInvalidError, err)
			return
		}

		fields := model.Filled(d)
		u, err = ar.server.Storages().UMC.UpdateUser(r.Context(), u, fields)
		if err != nil {
			ar.Error(w, l.ErrorWithLocale(err, locale))
			return
		}

		// if err := d.validate(user); err != nil {
		// 	ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyInvalidError, err)
		// 	return
		// }
		// // Check that new username is not taken.
		// if d.updateUsername {
		// 	if _, err := ar.server.Storages().UC.UserByID(d.NewUsername); err == nil {
		// 		ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPIUsernameTaken)
		// 		return
		// 	}
		// }

		// // Check that email is not taken.
		// if d.updateEmail {
		// 	if _, err := ar.server.Storages().User.UserByEmail(d.NewEmail); err == nil {
		// 		ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPIEmailTaken)
		// 		return
		// 	}
		// }

		// // Check that phone is not taken.
		// if d.updatePhone {
		// 	if _, err := ar.server.Storages().User.UserByPhone(d.NewPhone); err == nil {
		// 		ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPIPhoneTaken)
		// 		return
		// 	}
		// }

		// // Update password.
		// if d.updatePassword {
		// 	// Check old password.
		// 	if err := ar.server.Storages().User.CheckPassword(user.ID, d.OldPassword); err != nil {
		// 		ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyOldpasswordInvalid)
		// 		return
		// 	}

		// 	// Save new password.
		// 	err = ar.server.Storages().User.ResetPassword(user.ID, d.NewPassword)
		// 	if err != nil {
		// 		ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorStorageResetPasswordUserError, user.ID, err)
		// 		return
		// 	}

		// 	// Refetch user with new password hash.
		// 	if user, err = ar.server.Storages().User.UserByUsername(user.Username); err != nil {
		// 		ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorStorageFindUserEmailPhoneUsernameError, err)
		// 		return
		// 	}
		// }

		// // Change username if user specified new one.
		// if d.updateUsername {
		// 	user.Username = d.NewUsername
		// 	user = user.Deanonimized()
		// }

		// if d.updateEmail {
		// 	user.Email = d.NewEmail
		// }

		// if d.updatePhone {
		// 	user.Phone = d.NewPhone
		// }

		// if d.updateUsername || d.updateEmail || d.updatePhone {
		// 	if _, err = ar.server.Storages().User.UpdateUser(userID, user); err != nil {
		// 		ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorStorageUpdateUserError, userID, err)
		// 		return
		// 	}
		// }

		// // Prepare response.
		// updatedFields := []string{}
		// if d.updateUsername {
		// 	updatedFields = append(updatedFields, "username")
		// }
		// if d.updateEmail {
		// 	updatedFields = append(updatedFields, "email")
		// }
		// if d.updatePhone {
		// 	updatedFields = append(updatedFields, "phone")
		// }
		// if d.updatePassword {
		// 	updatedFields = append(updatedFields, "password")
		// }

		// msg := "Nothing changed."
		// if len(updatedFields) > 0 {
		// 	updatedFields[0] = strings.Title(updatedFields[0])
		// 	msg = strings.Join(updatedFields, ", ") + " changed. "
		// }
		// response := updateResponse{
		// 	Message: msg,
		// }

		ar.ServeJSON(w, locale, http.StatusOK, u)
	}
}
