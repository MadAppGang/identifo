package api

import (
	"crypto/rand"
	"encoding/base32"
	"io"
	"net/http"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/middleware"
)

// SetTFAOption enables or disables two-factor authentication for the user.
func (ar *Router) SetTFAOption() http.HandlerFunc {
	type tfaSecret struct {
		TFASecret string `json:"tfa_secret"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		d := struct {
			TFAEnabled bool `json:"tfa_enabled"`
		}{}
		if err := ar.MustParseJSON(w, r, &d); err != nil {
			ar.Error(w, ErrorAPIRequestBodyInvalid, http.StatusBadRequest, err.Error(), "SetTFAAvailability.MustParseJSON")
			return
		}

		app := middleware.AppFromContext(r.Context())
		if app == nil {
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App is not in context.", "SetTFAAvailability.AppFromContext")
			return
		}

		if d.TFAEnabled && !app.TFAEnabled() {
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

		tfaSecretStr := ""
		if d.TFAEnabled {
			// Generate 2FA secret.
			tfaSecretStr, err = ar.generate2FASecret(w)
			if err != nil {
				return
			}
		}
		user.SetTFAInfo(d.TFAEnabled, tfaSecretStr)

		if _, err := ar.userStorage.UpdateUser(userID, user); err != nil {
			ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "SetTFAAvailability.UpdateUser")
			return
		}

		ar.ServeJSON(w, http.StatusOK, tfaSecret{TFASecret: tfaSecretStr})
	}
}

func (ar *Router) generate2FASecret(w http.ResponseWriter) (string, error) {
	secret := make([]byte, 10)
	if _, err := io.ReadFull(rand.Reader, secret); err != nil {
		ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "generate2FASecret")
		return "", err
	}
	return base32.StdEncoding.EncodeToString(secret), nil
}

