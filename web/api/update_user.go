package api

import (
	"net/http"
	"time"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/tools/xmaps"
)

// UpdateUser allows to change user login and password.
func (ar *Router) UpdateUser() http.HandlerFunc {
	// this is large update thing
	// we need to dedicate empty field and absence of field in update request
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
		err := xmaps.CopyDstFields(d, &u)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAPIRequestBodyInvalidError, err)
			return
		}

		fields := xmaps.Filled(d)
		u, err = ar.server.Storages().UMC.UpdateUser(r.Context(), u, fields)
		if err != nil {
			ar.Error(w, l.ErrorWithLocale(err, locale))
			return
		}

		ar.ServeJSON(w, locale, http.StatusOK, u)
	}
}
