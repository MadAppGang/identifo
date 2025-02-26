package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/coreos/go-oidc/v3/oidc"
	l "github.com/madappgang/identifo/v2/localization"
	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/authorization"
	"github.com/madappgang/identifo/v2/web/middleware"
	"golang.org/x/oauth2"
)

const SessionNameOIDC = "_federated_oidc_session"

type oidcInfo struct {
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
}

var (
	oidcProviderCache     = map[string]oidcInfo{}
	oidcProviderCacheLock = sync.RWMutex{}
)

func getCachedOIDCProvider(ctx context.Context, app model.AppData) (*oidc.Provider, *oidc.IDTokenVerifier, error) {
	err := app.OIDCSettings.IsValid()
	if err != nil {
		return nil, nil, fmt.Errorf("OIDC not configured for app %s: %w", app.ID, err)
	}

	key := app.ID + ":oidc"

	oidcProviderCacheLock.RLock()
	oi, ok := oidcProviderCache[key]
	oidcProviderCacheLock.RUnlock()

	if ok {
		return oi.provider, oi.verifier, nil
	}

	oidcProviderCacheLock.Lock()
	defer oidcProviderCacheLock.Unlock()

	oi, ok = oidcProviderCache[key]
	if ok {
		return oi.provider, oi.verifier, nil
	}

	if app.OIDCSettings.Issuer != "" {
		ctx = oidc.InsecureIssuerURLContext(ctx, app.OIDCSettings.Issuer)
	}

	provider, err := oidc.NewProvider(ctx, app.OIDCSettings.ProviderURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create OIDC provider %s: %w", app.ID, err)
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: app.OIDCSettings.ClientID,
	})

	oi = oidcInfo{
		provider: provider,
		verifier: verifier,
	}

	oidcProviderCache[key] = oi

	return provider, verifier, nil
}

func (ar *Router) OIDCLogin(stateManagedByClient bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		locale := r.Header.Get("Accept-Language")

		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIAPPNoAPPInContext)
			return
		}

		if !ar.SupportedLoginWays.FederatedOIDC {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorFederatedOidcDisabled)
			return
		}

		redirect := r.URL.Query().Get("redirectUrl")
		if len(redirect) == 0 {
			ar.Error(w, locale, http.StatusBadRequest, l.APIAPPFederatedProviderEmptyRedirect)
			return
		}

		oidcProvider, _, err := getCachedOIDCProvider(ctx, app)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorStorageUserFederatedCreateError, err)
			return
		}

		// Configure an OpenID Connect aware OAuth2 client.
		oauth2Config := oauth2.Config{
			ClientID:     app.OIDCSettings.ClientID,
			ClientSecret: app.OIDCSettings.ClientSecret,
			RedirectURL:  redirect,

			// Discovery returns the OAuth2 endpoints.
			Endpoint: oidcProvider.Endpoint(),

			// "openid" is a required scope for OpenID Connect flows.
			Scopes: []string{oidc.ScopeOpenID},
		}

		oauth2Config.Scopes = append(oauth2Config.Scopes, app.OIDCSettings.Scopes...)
		oauth2Config.Scopes = append(oauth2Config.Scopes, getScopes(r)...)

		state, err := setState(r, stateManagedByClient)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorFederatedStateInternalError, err)
			return
		}

		authUrl := oauth2Config.AuthCodeURL(state)

		if !stateManagedByClient {
			fsess := model.FederatedSession{
				AppId:        app.ID,
				AuthUrl:      authUrl,
				CallbackUrl:  getCallbackUrl(r),
				RedirectUrl:  redirect,
				Scopes:       oauth2Config.Scopes,
				ProviderName: app.OIDCSettings.ProviderName,
			}

			sn := oidcSessionKey(app.ID, app.OIDCSettings.ProviderName)
			err = storeInSession(SessionNameOIDC, sn, fsess.Marshal(), r, w)
			if err != nil {
				ar.Error(w, locale, http.StatusBadRequest, l.APIFederatedCreateAuthUrlError, err)
				return
			}
		}

		http.Redirect(w, r, authUrl, http.StatusFound)
	}
}

func oidcSessionKey(appId, provider string) string {
	return "_oidc:" + appId + ":" + provider
}

func (ar *Router) OIDCLoginComplete(useSession bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		locale := r.Header.Get("Accept-Language")

		app := middleware.AppFromContext(ctx)
		if len(app.ID) == 0 {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIAPPNoAPPInContext)
			return
		}

		if !ar.SupportedLoginWays.FederatedOIDC {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorFederatedOidcDisabled)
			return
		}

		claims, fsess, providerData, providerScopes, err := ar.completeOIDCAuth(r, app, useSession)
		if err != nil {
			ar.ErrorResponse(w, err)
			return
		}

		userField := "sub"
		if app.OIDCSettings.UserIDClaimField != "" {
			userField = app.OIDCSettings.UserIDClaimField
		}

		fedUserID := extractField(claims, userField)
		email := extractField(claims, app.OIDCSettings.EmailClaimField)

		providerName := app.OIDCSettings.ProviderName

		user, err := ar.tryFindFederatedUser(providerName, fedUserID, email)
		if err != nil {
			if !errors.Is(err, model.ErrUserNotFound) {
				ar.Error(w, locale, http.StatusInternalServerError, l.ErrorStorageUserFederatedCreateError, err)
				return
			}

			if app.RegistrationForbidden {
				ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIAPPRegistrationForbidden)
				return
			}

			scopes := model.MergeScopes(app.Scopes, app.NewUserDefaultScopes, nil)

			user, err = ar.server.Storages().User.AddUserWithFederatedID(model.User{
				Email: email,
				// FullName: gothUser.FirstName + " " + gothUser.LastName,
				Scopes: scopes,
			}, providerName, fedUserID, app.NewUserDefaultRole)
			if err != nil {
				ar.Error(w, locale, http.StatusInternalServerError, l.ErrorStorageUserFederatedCreateError, err)
				return
			}
		}

		// Authorize user if the app requires authorization.
		azi := authorization.AuthzInfo{
			App:         app,
			UserRole:    user.AccessRole,
			ResourceURI: r.RequestURI,
			Method:      r.Method,
		}
		if err := ar.Authorizer.Authorize(azi); err != nil {
			ar.Error(w, locale, http.StatusForbidden, l.ErrorFederatedAccessDeniedError, err)
			return
		}

		// requestedScopes will contain OIDC scopes and custom requested scopes
		requestedScopes := []string{}
		if fsess != nil {
			requestedScopes = append(requestedScopes, fsess.Scopes...)
		}

		requestedScopes = append(requestedScopes, parseScopes(providerScopes, " ")...)
		requestedScopes = append(requestedScopes, getScopes(r)...)

		// map OIDC scopes to Identifo scopes
		requestedScopes = mapScopes(app.OIDCSettings.ScopeMapping, requestedScopes)

		authResult, resultScopes, err := ar.loginFlow(AuditOperationOIDCLogin, app, user, requestedScopes, nil)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorFederatedLoginError, err)
			return
		}

		if fsess != nil {
			authResult.CallbackUrl = fsess.CallbackUrl
		}

		authResult.Scopes = resultScopes.Scopes()
		authResult.ProviderData = *providerData

		ar.audit(AuditOperationOIDCLogin,
			user.ID, app.ID, r.UserAgent(), user.AccessRole, resultScopes.Scopes(),
			authResult.AccessToken, authResult.RefreshToken)

		ar.ServeJSON(w, locale, http.StatusOK, authResult)
	}
}

func mapScopes(m map[string]string, scopes []string) []string {
	result := make([]string, 0, len(scopes))

	for _, s := range scopes {
		if mapped, ok := m[s]; ok {
			result = append(result, mapped)
		} else {
			result = append(result, s)
		}
	}

	return result
}

func extractField(data map[string]any, key string) string {
	val := data[key]

	switch v := val.(type) {
	case string:
		return v
	case []string:
		if len(v) > 0 {
			return v[0]
		}
	case []any:
		if len(v) > 0 {
			sv, _ := v[0].(string)
			return sv
		}
	}

	return ""
}

func (ar *Router) completeOIDCAuth(
	r *http.Request,
	app model.AppData,
	useSession bool,
) (map[string]interface{}, *model.FederatedSession, *providerData, string, error) {
	ctx := r.Context()

	locale := r.Header.Get("Accept-Language")

	oidcProvider, verifier, err := getCachedOIDCProvider(ctx, app)
	if err != nil {
		return nil, nil, nil, "", NewLocalizedError(http.StatusInternalServerError, locale, l.ErrorFederatedOidcProviderError, err)
	}

	authCode := r.URL.Query().Get("code")
	if len(authCode) == 0 {
		ar.logger.Info("failed to authorize user with OIDC: no code in response",
			logging.FieldURL, r.URL.Query())
		return nil, nil, nil, "", NewLocalizedError(http.StatusBadRequest, locale, l.ErrorFederatedCodeError)
	}

	var fsess *model.FederatedSession

	if useSession {
		sn := oidcSessionKey(app.ID, app.OIDCSettings.ProviderName)
		value, err := ar.getFromSession(SessionNameOIDC, sn, r)
		if err != nil {
			return nil, fsess, nil, "", NewLocalizedError(http.StatusBadRequest, locale, l.ErrorFederatedUnmarshalSessionError, err)
		}

		fsess, err = model.UnmarshalFederatedSession(value)
		if err != nil {
			return nil, fsess, nil, "", NewLocalizedError(http.StatusBadRequest, locale, l.ErrorFederatedUnmarshalSessionError, err)
		}

		if fsess.AppId != app.ID {
			return nil, fsess, nil, "", NewLocalizedError(http.StatusBadRequest, locale, l.ErrorFederatedSessionAPPIDMismatch, fsess.AppId, app.ID)
		}

		errv := validateState(r, fsess.AuthUrl)
		if errv != nil {
			return nil, fsess, nil, "", NewLocalizedError(http.StatusBadRequest, locale, l.ErrorFederatedStateError)
		}
	}

	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config := oauth2.Config{
		ClientID:     app.OIDCSettings.ClientID,
		ClientSecret: app.OIDCSettings.ClientSecret,
		Endpoint:     oidcProvider.Endpoint(),
	}

	if fsess != nil {
		oauth2Config.RedirectURL = fsess.RedirectUrl
	}

	oauth2Token, err := oauth2Config.Exchange(ctx, authCode)
	if err != nil {
		return nil, fsess, nil, "", NewLocalizedError(http.StatusBadRequest, locale, l.ErrorFederatedExchangeError, err)
	}

	providerScopeVal := oauth2Token.Extra("scope")

	providerScope, ok := providerScopeVal.(string)
	if !ok {
		ar.logger.Warn("oidc returned scope is not string",
			"scopeType", fmt.Sprintf("%T", providerScopeVal),
			"scopeValue", fmt.Sprintf("%+v", providerScopeVal))
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, fsess, nil, "", NewLocalizedError(http.StatusBadRequest, locale, l.ErrorFederatedIDtokenMissing)
	}

	// Parse and verify ID Token payload.
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fsess, nil, "", NewLocalizedError(http.StatusBadRequest, locale, l.ErrorFederatedIDtokenInvalid, err)
	}

	// Extract custom claims
	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		return nil, fsess, nil, "", NewLocalizedError(http.StatusBadRequest, locale, l.ErrorFederatedClaimsError, err)
	}

	providerData := &providerData{
		AccessToken:  oauth2Token.AccessToken,
		RefreshToken: oauth2Token.RefreshToken,
		TokenType:    oauth2Token.TokenType,
		Expiry:       oauth2Token.Expiry,
	}

	return claims, fsess, providerData, providerScope, nil
}

func (ar *Router) tryFindFederatedUser(provider, fedUserID, email string) (model.User, error) {
	us := ar.server.Storages().User

	if fedUserID != "" {
		user, err := us.UserByFederatedID(provider, fedUserID)
		if err == nil {
			return user, nil
		}

		if !errors.Is(err, model.ErrUserNotFound) {
			return model.User{}, fmt.Errorf("cannot find user by federated ID: %w", err)
		}
	}

	if email == "" {
		return model.User{}, model.ErrUserNotFound
	}

	user, err := us.UserByEmail(email)
	if err != nil {
		return model.User{}, err
	}

	user.AddFederatedId(provider, fedUserID)

	_, uerr := us.UpdateUser(user.ID, user)
	if uerr != nil {
		ar.logger.Error("cannot update user with federated id",
			logging.FieldEmail, email,
			"fedreatedID", fedUserID,
			logging.FieldError, uerr)
	}

	return user, nil
}
