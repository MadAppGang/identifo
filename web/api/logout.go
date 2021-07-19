package api

import (
	"fmt"
	"net/http"

	"github.com/form3tech-oss/jwt-go"
	"github.com/madappgang/identifo/model"
)

// Logout logs user out and deactivates their tokens.
func (ar *Router) Logout() http.HandlerFunc {
	type logoutData struct {
		RefreshToken string `json:"refresh_token,omitempty"`
		DeviceToken  string `json:"device_token,omitempty"`
	}

	result := map[string]string{"result": "ok"}

	return func(w http.ResponseWriter, r *http.Request) {
		accessTokenBytes, ok := r.Context().Value(model.TokenRawContextKey).([]byte)
		if !ok {
			ar.logger.Println("Cannot fetch access token bytes from context")
			ar.ServeJSON(w, http.StatusNoContent, nil)
			return
		}
		accessTokenString := string(accessTokenBytes)

		// Blacklist current access token.
		if err := ar.tokenBlacklist.Add(accessTokenString); err != nil {
			ar.logger.Printf("Cannot blacklist access token: %s\n", err)
		}

		if r.Body == http.NoBody {
			ar.ServeJSON(w, http.StatusOK, result)
			return
		}

		d := logoutData{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		// Revoke refresh token, if present.
		if err := ar.revokeRefreshToken(d.RefreshToken, accessTokenString); err != nil {
			ar.logger.Printf("Cannot revoke refresh token: %s\n", err)
		}

		// Detach device token, if present.
		if len(d.DeviceToken) > 0 {
			// TODO: check for ownership when device tokens are supported.
			if err := ar.userStorage.DetachDeviceToken(d.DeviceToken); err != nil {
				ar.logger.Println("Cannot detach device token")
			}
		}

		ar.ServeJSON(w, http.StatusOK, result)
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

func (ar *Router) revokeRefreshToken(refreshTokenString, accessTokenString string) error {
	if len(refreshTokenString) == 0 {
		return nil
	}

	atSub, err := ar.getTokenSubject(accessTokenString)
	if err != nil {
		return err
	}
	rtSub, err := ar.getTokenSubject(refreshTokenString)
	if err != nil {
		return err
	}

	if atSub != rtSub {
		return fmt.Errorf("%s tried to revoke refresh token that belong to %s", atSub, rtSub)
	}

	if err := ar.tokenStorage.DeleteToken(refreshTokenString); err != nil {
		return fmt.Errorf("Cannot delete refresh token: %s", err)
	}

	if err := ar.tokenBlacklist.Add(refreshTokenString); err != nil {
		return fmt.Errorf("Cannot blacklist refresh token: %s", err)
	}
	return nil
}
