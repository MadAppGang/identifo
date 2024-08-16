package api

import (
	"fmt"
	"net/http"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/madappgang/identifo/v2/model"
)

// Logout logs user out and deactivates their tokens.
func (ar *Router) Logout() http.HandlerFunc {
	type logoutData struct {
		RefreshToken string `json:"refresh_token,omitempty"`
		DeviceToken  string `json:"device_token,omitempty"`
	}

	result := map[string]string{"result": "ok"}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		locale := r.Header.Get("Accept-Language")

		accessToken := tokenFromContext(ctx)

		accessTokenBytes, ok := ctx.Value(model.TokenRawContextKey).([]byte)
		if !ok {
			ar.logger.Println("Cannot fetch access token bytes from context")
			ar.ServeJSON(w, locale, http.StatusNoContent, nil)
			return
		}
		accessTokenString := string(accessTokenBytes)

		// Blacklist current access token.
		if err := ar.server.Storages().Blocklist.Add(accessTokenString); err != nil {
			ar.logger.Printf("Cannot blacklist access token: %s\n", err)
		}

		if r.Body == http.NoBody {
			ar.ServeJSON(w, locale, http.StatusOK, result)
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
			if err := ar.server.Storages().User.DetachDeviceToken(d.DeviceToken); err != nil {
				ar.logger.Println("Cannot detach device token")
			}
		}

		journal(accessToken.Subject(), accessToken.Audience(), "logout", nil)

		ar.ServeJSON(w, locale, http.StatusOK, result)
	}
}

func (ar *Router) getTokenSubject(tokenString string) (string, error) {
	claims := jwt.RegisteredClaims{}

	if _, err := jwt.ParseWithClaims(tokenString, &claims, nil); err == nil {
		return "", fmt.Errorf("cannot parse token: %s", err)
	}

	return claims.Subject, nil
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
