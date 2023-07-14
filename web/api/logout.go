package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/madappgang/identifo/v2/jwt"
)

// Logout logs user out and deactivates their tokens.
// add access token and refresh token to block list
// we need to detach device from user
func (ar *Router) Logout() http.HandlerFunc {
	type logoutData struct {
		RefreshToken string `json:"refresh_token,omitempty"`
		DeviceToken  string `json:"device,omitempty"`
	}

	result := map[string]string{"result": "ok"}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		access := tokenFromContext(r.Context())

		if r.Body != http.NoBody {
			ld := logoutData{}
			json.NewDecoder(r.Body).Decode(&ld)
			if ld.RefreshToken != "" {
				refresh, err := jwt.ParseTokenString(ld.RefreshToken)
				if err == nil && refresh != nil {
				}
			}
		}

		ar.server.Storages().Token.SaveToken(r.Context())

		// TODO: Detach device from user
		ar.ServeJSON(w, locale, http.StatusOK, result)
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

	if err := ar.server.Storages().Token.DeleteToken(refreshTokenString); err != nil {
		return fmt.Errorf("Cannot delete refresh token: %s", err)
	}

	if err := ar.server.Storages().Blocklist.Add(refreshTokenString); err != nil {
		return fmt.Errorf("Cannot blacklist refresh token: %s", err)
	}
	return nil
}
