package api

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/gorilla/sessions"
	l "github.com/madappgang/identifo/v2/localization"
	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/authorization"
	"github.com/madappgang/identifo/v2/web/middleware"
	"github.com/markbates/goth"
)

// SessionName is the key used to access the session store.
const SessionName = "_federated_session"

// Store can/should be set by applications using gothic. The default is a cookie store.
var (
	Store        sessions.Store
	defaultStore sessions.Store
)

var keySet = false

func init() {
	key := []byte(model.RandomPassword(64))
	keySet = len(key) != 0

	cookieStore := sessions.NewCookieStore([]byte(key))
	cookieStore.Options.HttpOnly = true
	Store = cookieStore
	defaultStore = Store
}

// SetState sets the state string associated with the given request.
// If no state string is associated with the request, one will be generated.
// This state is sent to the provider and can be retrieved during the
// callback.
var setState = func(req *http.Request, stateRequired bool) (string, error) {
	state := req.URL.Query().Get("state")
	if len(state) > 0 {
		return state, nil
	} else if stateRequired {
		return "", errors.New("state is required")
	}

	// If a state query param is not passed in, generate a random
	// base64-encoded nonce so that the state on the auth URL
	// is unguessable, preventing CSRF attacks, as described in
	//
	// https://auth0.com/docs/protocols/oauth2/oauth-state#keep-reading
	nonceBytes := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, nonceBytes)
	if err != nil {
		return "", fmt.Errorf("source of randomness unavailable: %v", err)
		// panic("gothic: source of randomness unavailable: " + err.Error())
	}

	return base64.URLEncoding.EncodeToString(nonceBytes), nil
}

// GetState gets the state returned by the provider during the callback.
// This is used to prevent CSRF attacks, see
// http://tools.ietf.org/html/rfc6749#section-10.12
var getState = func(req *http.Request) string {
	params := req.URL.Query()
	if params.Encode() == "" && req.Method == http.MethodPost {
		return req.FormValue("state")
	}
	return params.Get("state")
}

// OpenID Connect is based on OpenID Connect Auto Discovery URL (https://openid.net/specs/openid-connect-discovery-1_0-17.html)
// because the OpenID Connect provider initialize it self in the New(), it can return an error which should be handled or ignored
// ignore the error for now
// openidConnect, _ := openidConnect.New(os.Getenv("OPENID_CONNECT_KEY"), os.Getenv("OPENID_CONNECT_SECRET"), ar.Host+"/auth/federated/openid-connect", os.Getenv("OPENID_CONNECT_DISCOVERY_URL"))
// if openidConnect != nil {
// 	goth.UseProviders(openidConnect)
// }

func (ar *Router) FederatedLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		initProviders(app, redirect)

		// Clear and recreate providers from app settings

		url, err := ar.GetAuthURL(w, r)
		if err != nil {
			ar.Error(w, locale, http.StatusBadRequest, l.APIFederatedCreateAuthUrlError, err)
			return
		}

		ar.logger.Info("federated auth url", logging.FieldURL, url)

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func (ar *Router) FederatedLoginComplete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIAPPNoAPPInContext)
			return
		}

		providerName := r.URL.Query().Get("provider")
		if providerName == "" {
			ar.Error(w, locale, http.StatusBadRequest, l.APIAPPFederatedProviderEmpty)
			return
		}

		value, err := ar.getFromSession(SessionName, sessionKey(app.ID, providerName), r)
		if err != nil {
			ar.Error(w, locale, http.StatusBadRequest, "Can't get session", "FederatedLogin.AppByID")
			return
		}

		fsess, err := model.UnmarshalFederatedSession(value)
		if err != nil {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorFederatedUnmarshalSessionError, err)
			return
		}

		if fsess.AppId != app.ID {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorFederatedSessionAPPIDMismatch, fsess.AppId, app.ID)
		}

		initProviders(app, fsess.RedirectUrl)

		gothUser, err := ar.CompleteUserAuth(w, r)
		if err != nil {
			ar.Error(w, locale, http.StatusBadRequest, l.APIAPPFederatedProviderCantCompleteError, err)
			return
		}

		user, err := ar.server.Storages().User.UserByFederatedID(providerName, gothUser.UserID)

		if err == model.ErrUserNotFound && gothUser.Email != "" {
			user, err = ar.server.Storages().User.UserByEmail(gothUser.Email)
			if err == nil {
				user.AddFederatedId(providerName, gothUser.UserID)
				ar.server.Storages().User.UpdateUser(user.ID, user)
			}
		}

		if err == model.ErrUserNotFound && !app.RegistrationForbidden {
			scopes := model.MergeScopes(app.Scopes, app.NewUserDefaultScopes, nil)

			user, err = ar.server.Storages().User.AddUserWithFederatedID(model.User{
				Email:    gothUser.Email,
				FullName: gothUser.FirstName + " " + gothUser.LastName,
				Scopes:   scopes,
			}, providerName, gothUser.UserID, app.NewUserDefaultRole)
			if err != nil {
				ar.Error(w, locale, http.StatusInternalServerError, l.ErrorStorageUserFederatedCreateError, err)
				return
			}
		} else if err == model.ErrUserNotFound && app.RegistrationForbidden {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIAPPRegistrationForbidden)
			return
		} else if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorStorageUserFederatedCreateError, err)
			return
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

		authResult, resultScopes, err := ar.loginFlow(app, user, fsess.Scopes, nil)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorFederatedLoginError, err)
			return
		}

		authResult.CallbackUrl = fsess.CallbackUrl
		authResult.Scopes = fsess.Scopes

		ar.journal(JournalOperationFederatedLogin,
			user.ID, app.ID, r.UserAgent(), user.AccessRole, resultScopes.Scopes())

		ar.ServeJSON(w, locale, http.StatusOK, authResult)
	}
}

func initProviders(app model.AppData, redirect string) {
	goth.ClearProviders()
	for k, p := range app.FederatedProviders {
		if provider, ok := model.FederatedProviders[k]; ok {
			params := []interface{}{
				p.Params,
				// callbackURL
				redirect + "?appId=" + app.ID + "&provider=" + k,
			}
			for _, v := range p.Scopes {
				params = append(params, v)
			}
			f := reflect.ValueOf(provider.New)

			in := make([]reflect.Value, len(params))
			for k, param := range params {
				in[k] = reflect.ValueOf(param)
			}
			result := f.Call(in)

			provider := result[0].Interface().(goth.Provider)
			if n, ok := result[1].Interface().(error); !ok && n == nil {
				goth.UseProviders(provider)
			}
		}
	}
}

/*
GetAuthURL starts the authentication process with the requested provided.
It will return a URL that should be used to send users to.
It expects to be able to get the name of the provider from the query parameters
as either "provider"
*/
func (ar *Router) GetAuthURL(res http.ResponseWriter, req *http.Request) (string, error) {
	if !keySet && defaultStore == Store {
		ar.logger.Info("goth/gothic: no SESSION_SECRET environment variable is set. The default cookie store is not available and any calls will fail. Ignore this warning if you are using a different store.")
	}

	providerName := req.URL.Query().Get("provider")
	if providerName == "" {
		return "", errors.New(ar.ls.SD(l.APIAPPFederatedProviderEmpty))
	}

	app := middleware.AppFromContext(req.Context())
	if len(app.ID) == 0 {
		return "", errors.New(ar.ls.SD(l.ErrorAPIAPPNoAPPInContext))
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return "", err
	}

	redirect := req.URL.Query().Get("redirectUrl")
	if len(redirect) == 0 {
		return "", errors.New(ar.ls.SD(l.APIAPPFederatedProviderEmptyRedirect))
	}

	state, err := setState(req, false)
	if err != nil {
		return "", err
	}

	sess, err := provider.BeginAuth(state)
	if err != nil {
		return "", err
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}

	fsess := model.FederatedSession{
		ProviderSession: sess.Marshal(),
		AppId:           app.ID,
		AuthUrl:         url,
		CallbackUrl:     getCallbackUrl(req),
		RedirectUrl:     redirect,
		Scopes:          getScopes(req),
		ProviderName:    providerName,
	}

	err = storeInSession(SessionName, sessionKey(app.ID, providerName), fsess.Marshal(), req, res)

	if err != nil {
		return "", err
	}

	return url, err
}

/*
CompleteUserAuth does what it says on the tin. It completes the authentication
process and fetches all of the basic information about the user from the provider.
It expects to be able to get the name of the provider from the query parameters
as either "provider" or ":provider".
See https://github.com/markbates/goth/examples/main.go to see this in action.
*/
func (ar *Router) CompleteUserAuth(res http.ResponseWriter, req *http.Request) (goth.User, error) {
	if !keySet && defaultStore == Store {
		ar.logger.Info("goth/gothic: no SESSION_SECRET environment variable is set. The default cookie store is not available and any calls will fail. Ignore this warning if you are using a different store.")
	}

	providerName := req.URL.Query().Get("provider")
	if providerName == "" {
		return goth.User{}, errors.New(ar.ls.SD(l.APIAPPFederatedProviderEmpty))
	}

	app := middleware.AppFromContext(req.Context())
	if len(app.ID) == 0 {
		return goth.User{}, errors.New(ar.ls.SD(l.ErrorAPIAPPNoAPPInContext))
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return goth.User{}, err
	}

	value, err := ar.getFromSession(SessionName, sessionKey(app.ID, providerName), req)
	if err != nil {
		return goth.User{}, err
	}
	defer Logout(res, req)

	fsess, err := model.UnmarshalFederatedSession(value)
	if err != nil {
		return goth.User{}, err
	}

	sess, err := provider.UnmarshalSession(fsess.ProviderSession)
	if err != nil {
		return goth.User{}, err
	}

	err = validateGothState(req, sess)
	if err != nil {
		return goth.User{}, err
	}

	user, err := provider.FetchUser(sess)
	if err == nil {
		// user can be found with existing session data
		return user, err
	}

	params := req.URL.Query()
	if params.Encode() == "" && req.Method == "POST" {
		req.ParseForm()
		params = req.Form
	}

	// get new token and retry fetch
	_, err = sess.Authorize(provider, params)
	if err != nil {
		return goth.User{}, err
	}
	err = storeInSession(SessionName, sessionKey(app.ID, providerName), fsess.Marshal(), req, res)

	if err != nil {
		return goth.User{}, err
	}

	gu, err := provider.FetchUser(sess)
	return gu, err
}

// validateGothState ensures that the state token param from the original
// AuthURL matches the one included in the current (callback) request.
func validateGothState(req *http.Request, sess goth.Session) error {
	rawAuthURL, err := sess.GetAuthURL()
	if err != nil {
		return err
	}

	return validateState(req, rawAuthURL)
}

func validateState(req *http.Request, rawAuthURL string) error {
	authURL, err := url.Parse(rawAuthURL)
	if err != nil {
		return err
	}

	reqState := getState(req)

	originalState := authURL.Query().Get("state")
	if originalState != "" && (originalState != reqState) {
		return errors.New("state token mismatch")
	}
	return nil
}

// Logout invalidates a user session.
func Logout(res http.ResponseWriter, req *http.Request) error {
	session, err := Store.Get(req, SessionName)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	session.Values = make(map[interface{}]interface{})
	err = session.Save(req, res)
	if err != nil {
		return errors.New("could not delete user session ")
	}
	return nil
}

func getCallbackUrl(req *http.Request) string {
	return req.URL.Query().Get("callbackUrl")
}

func getScopes(req *http.Request) []string {
	return parseScopes(req.URL.Query().Get("scopes"), ",")
}

func parseScopes(scopes, sep string) []string {
	rs := strings.Split(scopes, sep)

	result := []string{}
	for _, scope := range rs {
		if scope != "" {
			result = append(result, scope)
		}
	}

	return result
}

func sessionKey(appId, provider string) string {
	return appId + ":" + provider
}

// StoreInSession stores a specified key/value pair in the session.
func storeInSession(sessionName, key string, value string, req *http.Request, res http.ResponseWriter) error {
	session, _ := Store.New(req, sessionName)
	session.Options.SameSite = http.SameSiteNoneMode
	session.Options.Secure = true

	if err := updateSessionValue(session, key, value); err != nil {
		return err
	}

	return session.Save(req, res)
}

// GetFromSession retrieves a previously-stored value from the session.
// If no value has previously been stored at the specified key, it will return an error.
func (ar *Router) getFromSession(sessionName, key string, req *http.Request) (string, error) {
	session, _ := Store.Get(req, sessionName)
	value, err := ar.getSessionValue(session, key)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (ar *Router) getSessionValue(session *sessions.Session, key string) (string, error) {
	value := session.Values[key]
	if value == nil {
		return "", errors.New(ar.ls.SD(l.ErrorAPISessionNotFound, key))
	}

	rdata := strings.NewReader(value.(string))
	r, err := gzip.NewReader(rdata)
	if err != nil {
		return "", err
	}
	s, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(s), nil
}

func updateSessionValue(session *sessions.Session, key, value string) error {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(value)); err != nil {
		return err
	}
	if err := gz.Flush(); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}

	session.Values[key] = b.String()
	return nil
}
