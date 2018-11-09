package http

import (
	"net/http"
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
			ar.tokenStorage.RevokeToken(d.RefreshToken)
		}
		if tokenString, ok := r.Context().Value(TokenRawContextKey).(string); ok {
			ar.tokenStorage.RevokeToken(tokenString)
		}
		if len(d.DeviceToken) > 0 {
			ar.userStorage.DetachDeviceToken(d.DeviceToken)
		}
		response := logoutResponse{
			Message: "Done",
		}
		ar.ServeJSON(w, http.StatusOK, response)
	}

}
