package api

import (
	"net/http"

	jwtService "github.com/madappgang/identifo/jwt/service"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/middleware"
)

// RefreshTokens issues new access and, if requsted, refresh tokens for provided refresh token.
// After the new refresh token is issued, the old one gets invalidated.
func (ar *Router) RefreshTokens() http.HandlerFunc {
	type requestData struct {
		Scopes []string `json:"scopes,omitempty"`
	}

	type responseData struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		rd := requestData{}
		if ar.MustParseJSON(w, r, &rd) != nil {
			return
		}

		app := middleware.AppFromContext(r.Context())
		if app == nil {
			ar.logger.Println("Error getting App")
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App ID is absent in header params", "RefreshTokens.AppFromContext")
			return
		}

		// Get refresh token from context.
		oldRefreshToken := tokenFromContext(r.Context())

		// Issue new access token and stringify it for response.
		accessToken, err := ar.tokenService.RefreshAccessToken(oldRefreshToken)
		if err != nil {
			ar.Error(w, ErrorAPIAppAccessTokenNotCreated, http.StatusInternalServerError, err.Error(), "RefreshTokens.RefreshAccessToken")
			return
		}
		accessTokenString, err := ar.tokenService.String(accessToken)
		if err != nil {
			ar.Error(w, ErrorAPIAppRefreshTokenNotCreated, http.StatusInternalServerError, err.Error(), "RefreshTokens.accessTokenString")
			return
		}

		// Stringify old refresh token and issue new one.
		oldRefreshTokenBytes, ok := r.Context().Value(model.TokenRawContextKey).([]byte)
		if !ok || oldRefreshTokenBytes == nil {
			ar.Error(w, ErrorAPIAppRefreshTokenNotCreated, http.StatusInternalServerError, "Token is empty or invalid.", "RefreshTokens.RawTokenFromContext")
			return
		}
		oldRefreshTokenString := string(oldRefreshTokenBytes)

		newRefreshTokenString, err := ar.issueNewRefreshToken(oldRefreshTokenString, rd.Scopes, app)
		if err != nil {
			ar.Error(w, ErrorAPIAppRefreshTokenNotCreated, http.StatusInternalServerError, err.Error(), "RefreshToken.newRefreshTokenString")
			return
		}

		// Invalidate old refresh token - delete it from token storage and add to blacklist.
		ar.invalidateOldRefreshToken(oldRefreshTokenString)

		result := &responseData{
			AccessToken:  accessTokenString,
			RefreshToken: newRefreshTokenString,
		}

		ar.ServeJSON(w, http.StatusOK, result)
	}
}

func (ar *Router) issueNewRefreshToken(oldRefreshTokenString string, scopes []string, app model.AppData) (string, error) {
	if !contains(scopes, jwtService.OfflineScope) { // Don't issue new refresh token if not requested.
		return "", nil
	}

	userID, err := ar.getTokenSubject(oldRefreshTokenString)
	if err != nil {
		return "", err
	}

	user, err := ar.userStorage.UserByID(userID)
	if err != nil {
		return "", err
	}

	refreshToken, err := ar.tokenService.NewRefreshToken(user, scopes, app)
	if err != nil {
		return "", err
	}

	refreshTokenString, err := ar.tokenService.String(refreshToken)
	if err != nil {
		return "", err
	}

	return refreshTokenString, err
}

func (ar *Router) invalidateOldRefreshToken(oldRefreshTokenString string) {
	if err := ar.tokenStorage.DeleteToken(oldRefreshTokenString); err != nil {
		ar.logger.Println("Cannot delete old refresh token from token storage:", err)
	}
	if err := ar.tokenBlacklist.Add(oldRefreshTokenString); err != nil {
		ar.logger.Println("Cannot blacklist old refresh token:", err)
	}
	ar.logger.Println("Old refresh token successfully invalidated")
}
