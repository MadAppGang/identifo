package api

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/madappgang/identifo/model"
)

// Logout logs user out and deactivates their tokens.
func (ar *Router) Logout() http.HandlerFunc {
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

		accessTokenBytes, ok := r.Context().Value(model.TokenRawContextKey).([]byte)
		if !ok {
			ar.logger.Println("Cannot fetch access token bytes from context")
			ar.ServeJSON(w, http.StatusNoContent, nil)
			return
		}
		accessTokenString := string(accessTokenBytes)

		// Revoke current access token.
		if err := ar.tokenStorage.RevokeToken(accessTokenString); err != nil {
			ar.logger.Println("Cannot revoke access token")
			ar.ServeJSON(w, http.StatusNoContent, nil)
			return
		}

		// Revoke refresh token, if present.
		if len(d.RefreshToken) > 0 {
			atSub, err := ar.getTokenSubject(accessTokenString)
			if err != nil {
				ar.logger.Println(err)
			}
			rtSub, err := ar.getTokenSubject(d.RefreshToken)
			if err != nil {
				ar.logger.Println(err)
			}

			if atSub != rtSub {
				ar.logger.Printf("%s tried to revoke refresh token that belong to %s\n", atSub, rtSub)
			}
			if err := ar.tokenStorage.RevokeToken(d.RefreshToken); err != nil {
				ar.logger.Println("Cannot revoke refresh token")
			}
		}

		// Detach device token, if present.
		if len(d.DeviceToken) > 0 {
			// TODO: check for ownership when device tokens are supported.
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

func (ar *Router) getTokenSubject(tokenString string) (string, error) {
	claims := jwt.MapClaims{}

	if _, err := jwt.ParseWithClaims(tokenString, claims, nil); err == nil {
		return "", fmt.Errorf("Cannot parse token: %s", err)
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("Cannot obtain token subject")
	}

	return sub, nil
}
