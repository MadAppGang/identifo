package api

import (
	"net/http"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/middleware"
	"github.com/xlzd/gotp"
)

// SetTFAOption enables or disables two-factor authentication for the user.
func (ar *Router) SetTFAOption() http.HandlerFunc {
	type tfaSecret struct {
		TFASecret string `json:"tfa_secret"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var tfa model.TFAInfo
		if err := ar.MustParseJSON(w, r, &tfa); err != nil {
			ar.Error(w, ErrorAPIRequestBodyInvalid, http.StatusBadRequest, err.Error(), "SetTFAAvailability.MustParseJSON")
			return
		}

		app := middleware.AppFromContext(r.Context())
		if app == nil {
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App is not in context.", "SetTFAAvailability.AppFromContext")
			return
		}

		if tfa.IsEnabled && !app.TFAEnabled() {
			ar.Error(w, ErrorAPIRequestBodyParamsInvalid, http.StatusBadRequest, "App does not support two-factor authentication", "SetTFAAvailability.TFAEnabled")
			return
		}

		accessTokenBytes, ok := r.Context().Value(model.TokenRawContextKey).([]byte)
		if !ok {
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "Token bytes are not in context.", "SetTFAAvailability.TokenBytesFromContext")
			return
		}

		// Get userID from token and update user with this ID.
		userID, err := ar.getTokenSubject(string(accessTokenBytes))
		if err != nil {
			ar.Error(w, ErrorAPIAppCannotExtractTokenSubject, http.StatusInternalServerError, err.Error(), "SetTFAAvailability.getTokenSubject")
			return
		}

		user, err := ar.userStorage.UserByID(userID)
		if err != nil {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusBadGateway, err.Error(), "SetTFAAvailability.UserByID")
			return
		}

		if tfa.IsEnabled {
			// Generate 2FA secret.
			tfa.Secret = gotp.RandomSecret(16)
		}
		user.SetTFAInfo(tfa)

		if _, err := ar.userStorage.UpdateUser(userID, user); err != nil {
			ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "SetTFAAvailability.UpdateUser")
			return
		}

		ar.ServeJSON(w, http.StatusOK, tfaSecret{TFASecret: tfa.Secret})
	}
}
