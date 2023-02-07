package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/coreos/go-oidc/v3/oidc"
	l "github.com/madappgang/identifo/v2/localization"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/authorization"
	"github.com/madappgang/identifo/v2/web/middleware"
	"golang.org/x/oauth2"
)

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

func (ar *Router) OIDCLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	locale := r.Header.Get("Accept-Language")

	app := middleware.AppFromContext(r.Context())
	if len(app.ID) == 0 {
		ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIAPPNoAPPInContext)
		return
	}

	redirect := r.URL.Query().Get("redirectUrl")
	if len(redirect) == 0 {
		ar.Error(w, locale, http.StatusBadRequest, l.APIAPPFederatedProviderEmptyRedirect)
		return
	}

	redirect, err := makeRedirectURL(redirect, app)
	if err != nil {
		ar.ErrorResponse(w, err)
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

	oauth2Config.Scopes = append(oauth2Config.Scopes, getScopes(r)...)

	state := setState(r)

	fsess := model.FederatedSession{
		AppId:        app.ID,
		AuthUrl:      oauth2Config.AuthCodeURL(state),
		CallbackUrl:  getCallbackUrl(r),
		RedirectUrl:  redirect,
		Scopes:       oauth2Config.Scopes,
		ProviderName: app.OIDCSettings.ProviderName,
	}

	sn := oidcSessionName(app.ID, app.OIDCSettings.ProviderName)
	err = storeInSession(sn, fsess.Marshal(), r, w)
	if err != nil {
		ar.Error(w, locale, http.StatusBadRequest, l.APIFederatedCreateAuthUrlError, err)
		return
	}

	http.Redirect(w, r, fsess.AuthUrl, http.StatusFound)
}

func oidcSessionName(appId, provider string) string {
	return "_oidc:" + appId + ":" + provider
}

func (ar *Router) OIDCLoginComplete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	locale := r.Header.Get("Accept-Language")

	app := middleware.AppFromContext(ctx)
	if len(app.ID) == 0 {
		ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIAPPNoAPPInContext)
		return
	}

	claims, fsess, err := ar.completeOIDCAuth(r, app)
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

	authResult, err := ar.loginFlow(app, user, fsess.Scopes)
	if err != nil {
		ar.Error(w, locale, http.StatusInternalServerError, l.ErrorFederatedLoginError, err)
		return
	}

	authResult.CallbackUrl = fsess.CallbackUrl
	authResult.Scopes = fsess.Scopes

	ar.ServeJSON(w, locale, http.StatusOK, authResult)
}

func makeRedirectURL(redirect string, app model.AppData) (string, error) {
	u, err := url.Parse(redirect)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("appId", app.ID)

	return u.String(), nil
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

func (ar *Router) completeOIDCAuth(r *http.Request, app model.AppData) (map[string]interface{}, *model.FederatedSession, error) {
	ctx := r.Context()

	var fsess *model.FederatedSession

	locale := r.Header.Get("Accept-Language")

	oidcProvider, verifier, err := getCachedOIDCProvider(ctx, app)
	if err != nil {
		return nil, fsess, NewLocalizedError(http.StatusInternalServerError, locale, l.ErrorStorageUserFederatedCreateError)
	}

	authCode := r.URL.Query().Get("code")
	if len(authCode) == 0 {
		log.Println("failed ot authorize user with OIDC: no code in response", r.URL.Query())
		return nil, fsess, NewLocalizedError(http.StatusBadRequest, locale, l.ErrorFederatedCodeError)
	}

	sn := oidcSessionName(app.ID, app.OIDCSettings.ProviderName)
	value, err := ar.getFromSession(sn, r)
	if err != nil {
		return nil, fsess, NewLocalizedError(http.StatusBadRequest, locale, l.ErrorFederatedUnmarshalSessionError)
	}

	fsess, err = model.UnmarshalFederatedSession(value)
	if err != nil {
		return nil, fsess, NewLocalizedError(http.StatusBadRequest, locale, l.ErrorFederatedUnmarshalSessionError)
	}

	if fsess.AppId != app.ID {
		return nil, fsess, NewLocalizedError(http.StatusBadRequest, locale, l.ErrorFederatedSessionAPPIDMismatch)
	}

	errv := validateState(r, fsess.AuthUrl)
	if errv != nil {
		return nil, fsess, NewLocalizedError(http.StatusBadRequest, locale, l.ErrorFederatedStateError)
	}

	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config := oauth2.Config{
		ClientID:     app.OIDCSettings.ClientID,
		ClientSecret: app.OIDCSettings.ClientSecret,
		RedirectURL:  fsess.RedirectUrl,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: oidcProvider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID},
	}

	oauth2Token, err := oauth2Config.Exchange(ctx, authCode)
	if err != nil {
		return nil, fsess, NewLocalizedError(http.StatusBadRequest, locale, l.ErrorFederatedExchangeError, err)
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, fsess, NewLocalizedError(http.StatusBadRequest, locale, l.ErrorFederatedIDtokenMissing)
	}

	// Parse and verify ID Token payload.
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fsess, NewLocalizedError(http.StatusBadRequest, locale, l.ErrorFederatedIDtokenInvalid, err)
	}

	// Extract custom claims
	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		return nil, fsess, NewLocalizedError(http.StatusBadRequest, locale, l.ErrorFederatedClaimsError, err)
	}

	return claims, fsess, nil
}

func (ar *Router) tryFindFederatedUser(provider, fedUserID, email string) (model.User, error) {
	us := ar.server.Storages().User

	if fedUserID != "" {
		user, err := us.UserByFederatedID(provider, fedUserID)
		if err == nil {
			return user, nil
		}

		if !errors.Is(err, model.ErrUserNotFound) {
			return model.User{}, fmt.Errorf("can not find user by federated ID: %w", err)
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
		log.Printf("can not update user %s with federated id: %v\n", email, uerr)
	}

	return user, nil
}
