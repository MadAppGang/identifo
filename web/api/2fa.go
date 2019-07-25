package api

import (
	"fmt"
	"net/http"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/middleware"
	"github.com/xlzd/gotp"
)

// EnableTFA enables two-factor authentication for the user.
func (ar *Router) EnableTFA() http.HandlerFunc {
	type tfaSecret struct {
		TFASecret string `json:"tfa_secret"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		app := middleware.AppFromContext(r.Context())
		if app == nil {
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App is not in context.", "EnableTFA.AppFromContext")
			return
		}

		if tfaStatus := app.TFAStatus(); tfaStatus != model.TFAStatusOptional {
			ar.Error(w, ErrorAPIRequestBodyParamsInvalid, http.StatusBadRequest, fmt.Sprintf("App TFA status is '%s', not 'optional'", tfaStatus), "EnableTFA.TFAStatus")
			return
		}

		accessTokenBytes, ok := r.Context().Value(model.TokenRawContextKey).([]byte)
		if !ok {
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "Token bytes are not in context.", "EnableTFA.TokenBytesFromContext")
			return
		}

		// Get userID from token and update user with this ID.
		userID, err := ar.getTokenSubject(string(accessTokenBytes))
		if err != nil {
			ar.Error(w, ErrorAPIAppCannotExtractTokenSubject, http.StatusInternalServerError, err.Error(), "EnableTFA.getTokenSubject")
			return
		}

		user, err := ar.userStorage.UserByID(userID)
		if err != nil {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusBadRequest, err.Error(), "EnableTFA.UserByID")
			return
		}

		if user.TFAInfo().IsEnabled {
			ar.Error(w, ErrorAPIRequestTFAAlreadyEnabled, http.StatusBadRequest, "TFA already enabled for this user", "EnableTFA.alreadyEnabled")
			return
		}

		tfa := model.TFAInfo{
			IsEnabled: true,
			Secret:    gotp.RandomSecret(16),
		}
		user.SetTFAInfo(tfa)

		if _, err := ar.userStorage.UpdateUser(userID, user); err != nil {
			ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "EnableTFA.UpdateUser")
			return
		}

		switch ar.tfaType {
		case model.TFATypeApp:
			ar.ServeJSON(w, http.StatusOK, &tfaSecret{TFASecret: tfa.Secret})
			return
		case model.TFATypeSMS:
			ar.sendTFASecretInSMS(w, tfa.Secret)
			return
		case model.TFATypeEmail:
			ar.sendTFASecretInEmail(w, tfa.Secret)
			return
		}
		ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, fmt.Sprintf("Unknown tfa type '%s'", ar.tfaType), "switch.tfaType")
	}
}

func (ar *Router) sendTFASecretInSMS(w http.ResponseWriter, tfaSecret string) {
	ar.Error(w, ErrorAPIInternalServerError, http.StatusBadRequest, "Not yet implemented", "sendTFASecretInSMS")
}

func (ar *Router) sendTFASecretInEmail(w http.ResponseWriter, tfaSecret string) {
	ar.Error(w, ErrorAPIInternalServerError, http.StatusBadRequest, "Not yet implemented", "sendTFASecretInEmail")
}
