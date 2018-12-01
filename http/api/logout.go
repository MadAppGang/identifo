package api

import (
	"net/http"

	"github.com/madappgang/identifo/model"
)

//Logout logouts user, deactivates his tokens
func (ar *apiRouter) Logout() http.HandlerFunc {
	type logoutData struct {
		RefreshToken string `json:"refresh_token,omitempty"`
		DeviceToken  string `json:"device_token,omitempty"`
	}

	type logoutResponse struct {
		Message string `json:"message"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		d := logoutData{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}
		if len(d.RefreshToken) > 0 {
			if err := ar.tokenStorage.RevokeToken(d.RefreshToken); err != nil {
				ar.logger.Println("Cannot revoke refresh token")
			}
		}
		if tokenString, ok := r.Context().Value(model.TokenRawContextKey).(string); ok {
			if err := ar.tokenStorage.RevokeToken(tokenString); err != nil {
				ar.logger.Println("Cannot revoke token")
			}
		}
		if len(d.DeviceToken) > 0 {
			if err := ar.userStorage.DetachDeviceToken(d.DeviceToken); err != nil {
				ar.logger.Println("Cannot detach device token")
			}
		}
		response := logoutResponse{
			Message: "Done",
		}
		ar.ServeJSON(w, http.StatusOK, response)
	}

}
