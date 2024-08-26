package api

import (
	"encoding/json"
	"net/http"
	"strings"

	l "github.com/madappgang/identifo/v2/localization"
	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
)

// RefreshTokens issues new access and, if requested, refresh token for provided refresh token.
// After new tokens are issued, the old refresh token gets invalidated (via blacklisting).
func (ar *Router) RefreshTokens() http.HandlerFunc {
	type requestData struct {
		Scopes []string `json:"scopes,omitempty"`
	}

	type responseData struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		locale := r.Header.Get("Accept-Language")

		rd := requestData{}
		if err := json.NewDecoder(r.Body).Decode(&rd); err != nil {
			// Assume we have not requested any scopes,  if there is no valid data in the body
			rd = requestData{Scopes: []string{}}
		}

		app := middleware.AppFromContext(ctx)
		if len(app.ID) == 0 {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIAPPNoAPPInContext)
			return
		}

		// Get refresh token from context.
		oldRefreshToken := tokenFromContext(ctx)

		if err := oldRefreshToken.Validate(); err != nil {
			ar.Error(w, locale, http.StatusUnauthorized, l.ErrorTokenInvalidError, err)
			return
		}

		tokenPayload, err := ar.getTokenPayloadForApp(app, oldRefreshToken.Subject())
		if err != nil {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIAPPUnableToTokenPayloadForAPPError)
			return
		}

		// Issue new access token and stringify it for response.
		accessToken, err := ar.server.Services().Token.RefreshAccessToken(oldRefreshToken, tokenPayload)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorTokenRefreshAccessToken, err)
			return
		}
		accessTokenString, err := ar.server.Services().Token.String(accessToken)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorTokenUnableToCreateAccessTokenError, err)
			return
		}

		// Stringify old refresh token and issue new one.
		oldRefreshTokenBytes, ok := r.Context().Value(model.TokenRawContextKey).([]byte)
		if !ok || oldRefreshTokenBytes == nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorTokenRefreshEmpty)
			return
		}
		oldRefreshTokenString := string(oldRefreshTokenBytes)

		newRefreshTokenString, err := ar.issueNewRefreshToken(oldRefreshTokenString, rd.Scopes, app)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorTokenUnableToCreateRefreshTokenError, err)
			return
		}

		// Invalidate old refresh token - delete it from token storage and add to blacklist.
		ar.invalidateOldRefreshToken(oldRefreshTokenString)

		result := &responseData{
			AccessToken:  accessTokenString,
			RefreshToken: newRefreshTokenString,
		}

		resultScopes := strings.Split(accessToken.Scopes(), " ")
		ar.journal(JournalOperationRefreshToken,
			oldRefreshToken.Subject(), app.ID, r.UserAgent(), "", resultScopes)

		ar.ServeJSON(w, locale, http.StatusOK, result)
	}
}

func (ar *Router) issueNewRefreshToken(
	oldRefreshTokenString string,
	requestedScopes []string,
	app model.AppData,
) (string, error) {
	if !contains(requestedScopes, model.OfflineScope) { // Don't issue new refresh token if not requested.
		return "", nil
	}

	userID, err := ar.getTokenSubject(oldRefreshTokenString)
	if err != nil {
		return "", err
	}

	user, err := ar.server.Storages().User.UserByID(userID)
	if err != nil {
		return "", err
	}

	scopes := model.AllowedScopes(requestedScopes, user.Scopes, app.Offline)

	refreshToken, err := ar.server.Services().Token.NewRefreshToken(user, scopes, app)
	if err != nil {
		return "", err
	}

	refreshTokenString, err := ar.server.Services().Token.String(refreshToken)
	if err != nil {
		return "", err
	}

	return refreshTokenString, err
}

func (ar *Router) invalidateOldRefreshToken(oldRefreshTokenString string) {
	if err := ar.server.Storages().Token.DeleteToken(oldRefreshTokenString); err != nil {
		ar.logger.Error("Cannot delete old refresh token from token storage",
			logging.FieldError, err)
	}

	if err := ar.server.Storages().Blocklist.Add(oldRefreshTokenString); err != nil {
		ar.logger.Error("Cannot blacklist old refresh token",
			logging.FieldError, err)
	}

	ar.logger.Info("Old refresh token successfully invalidated")
}
