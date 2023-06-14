package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	l "github.com/madappgang/identifo/v2/localization"
	"github.com/madappgang/identifo/v2/model"
	pphttp "github.com/madappgang/identifo/v2/user_payload_provider/http"
	"github.com/madappgang/identifo/v2/user_payload_provider/plugin"
	"github.com/madappgang/identifo/v2/web/authorization"
	"github.com/madappgang/identifo/v2/web/middleware"
	"github.com/xlzd/gotp"
	"golang.org/x/oauth2"
)

var (
	errPleaseEnableTFA   = fmt.Errorf("please enable two-factor authentication to be able to use this app")
	errPleaseSetPhoneTFA = fmt.Errorf("please set phone for two-factor authentication to be able to use this app")
	errPleaseSetEmailTFA = fmt.Errorf("please set email for two-factor authentication to be able to use this app")
	errPleaseDisableTFA  = fmt.Errorf("please disable two-factor authentication to be able to use this app")
)

type SendTFAEmailData struct {
	User model.User
	OTP  string
	Data interface{}
}

const (
	smsTFACode        = "%v is your one-time password!"
	hotpLifespanHours = 12 // One time code expiration in hours, default value is 30 secs for TOTP and 12 hours for HOTP
)

// AuthResponse is a response with successful auth data.
type AuthResponse struct {
	AccessToken  string        `json:"access_token,omitempty" bson:"access_token,omitempty"`
	RefreshToken string        `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`
	User         model.User    `json:"user,omitempty" bson:"user,omitempty"`
	Require2FA   bool          `json:"require_2fa" bson:"require_2fa"`
	Enabled2FA   bool          `json:"enabled_2fa" bson:"enabled_2fa"`
	CallbackUrl  string        `json:"callback_url,omitempty" bson:"callback_url,omitempty"`
	Scopes       []string      `json:"scopes,omitempty" bson:"scopes,omitempty"`
	ProviderData *oauth2.Token `json:"provider_data,omitempty" bson:"provider_data,omitempty"`
}

type login struct {
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

type loginData struct {
	login
	Password    string   `json:"password,omitempty"`
	DeviceToken string   `json:"device_token,omitempty"`
	Scopes      []string `json:"scopes,omitempty"`
}

func (ld *login) validate() error {
	emailLen := len(ld.Email)
	phoneLen := len(ld.Phone)
	usernameLen := len(ld.Username)
	if emailLen > 0 {
		if phoneLen > 0 || usernameLen > 0 {
			return fmt.Errorf("don't use phone or username when login with email")
		}
		if !model.EmailRegexp.MatchString(ld.Email) {
			return fmt.Errorf("invalid email")
		}
	}
	if phoneLen > 0 {
		if emailLen > 0 || usernameLen > 0 {
			return fmt.Errorf("don't use email or username when login with phone")
		}
		if !model.PhoneRegexp.MatchString(ld.Email) {
			return fmt.Errorf("invalid phone")
		}
	}
	if usernameLen > 0 {
		if phoneLen > 0 || emailLen > 0 {
			return fmt.Errorf("don't use phone or email when login with username")
		}
		if usernameLen < 6 || usernameLen > 130 {
			return fmt.Errorf("incorrect username length %d, expected a number between 6 and 130", usernameLen)
		}
	}
	return nil
}

func (ld *loginData) validate() error {
	if err := ld.login.validate(); err != nil {
		return err
	}
	pswdLen := len(ld.Password)
	if pswdLen < 6 || pswdLen > 50 {
		return fmt.Errorf("incorrect password length %d, expected a number between 6 and 130", pswdLen)
	}
	return nil
}

func (ar *Router) checkSupportedWays(l login) error {
	if !ar.SupportedLoginWays.Email && len(l.Email) > 0 {
		return fmt.Errorf("application does not support login with email")
	}

	if !ar.SupportedLoginWays.Phone && len(l.Phone) > 0 {
		return fmt.Errorf("application does not support login with phone")
	}

	if !ar.SupportedLoginWays.Username && len(l.Username) > 0 {
		return fmt.Errorf("application does not support login with username")
	}
	return nil
}

// LoginWithPassword logs user in with email and password.
func (ar *Router) LoginWithPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		ld := loginData{}
		if ar.MustParseJSON(w, r, &ld) != nil {
			return
		}

		if err := ld.validate(); err != nil {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyInvalidError, err)
			return
		}

		if err := ar.checkSupportedWays(ld.login); err != nil {
			ar.Error(w, locale, http.StatusBadRequest, l.APIAPPUsernameLoginNotSupported)
			return
		}

		var err error
		user := model.User{}

		if len(ld.Email) > 0 {
			user, err = ar.server.Storages().User.UserByEmail(ld.Email)
		} else if len(ld.Phone) > 0 {
			user, err = ar.server.Storages().User.UserByPhone(ld.Phone)
		} else if len(ld.Username) > 0 {
			user, err = ar.server.Storages().User.UserByUsername(ld.Username)
		}

		if err != nil {
			ar.Error(w, locale, http.StatusUnauthorized, l.ErrorAPIRequestIncorrectLoginOrPassword)
			return
		}

		if err = ar.server.Storages().User.CheckPassword(user.ID, ld.Password); err != nil {
			// return this error to hide the existence of the user.
			ar.Error(w, locale, http.StatusUnauthorized, l.ErrorAPIRequestIncorrectLoginOrPassword)
			return
		}

		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIAPPNoAPPInContext)
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
			ar.Error(w, locale, http.StatusForbidden, l.APIAccessDenied)
			return
		}

		authResult, err := ar.loginFlow(app, user, ld.Scopes)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorAPILoginError, err)
			return
		}

		ar.ServeJSON(w, locale, http.StatusOK, authResult)
	}
}

func (ar *Router) sendOTPCode(app model.AppData, user model.User) error {
	// we don't need to send any code for FTA Type App, it uses TOTP and generated on client side with the app
	if ar.tfaType != model.TFATypeApp {

		// increment hotp code seed
		otp := gotp.NewDefaultHOTP(user.TFAInfo.Secret).At(user.TFAInfo.HOTPCounter + 1)
		tfa := user.TFAInfo
		tfa.HOTPCounter++
		tfa.HOTPExpiredAt = time.Now().Add(time.Hour * hotpLifespanHours)
		user.TFAInfo = tfa
		if _, err := ar.server.Storages().User.UpdateUser(user.ID, user); err != nil {
			return err
		}
		switch ar.tfaType {
		case model.TFATypeSMS:
			return ar.sendTFACodeInSMS(app, user.TFAInfo.Phone, otp)
		case model.TFATypeEmail:
			return ar.sendTFACodeOnEmail(app, user, otp)
		}

	}

	return nil
}

// IsLoggedIn is for checking whether user is logged in or not.
// In fact, all needed work is done in Token middleware.
// If we reached this code, user is logged in (presented valid and not blacklisted access token).
func (ar *Router) IsLoggedIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")
		ar.ServeJSON(w, locale, http.StatusOK, nil)
	}
}

// GetUser return current user info with sanitized tfa
func (ar *Router) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		userID := tokenFromContext(r.Context()).UserID()
		user, err := ar.server.Storages().User.UserByID(userID)
		if err != nil {
			ar.Error(w, locale, http.StatusUnauthorized, l.ErrorStorageFindUserIDError, userID, err)
			return
		}
		ar.ServeJSON(w, locale, http.StatusOK, user.SanitizedTFA())
	}
}

// getTokenPayloadForApp get additional token payload data
func (ar *Router) getTokenPayloadForApp(app model.AppData, userID string) (map[string]interface{}, error) {
	if app.TokenPayloadService == model.TokenPayloadServiceNone ||
		app.TokenPayloadService == "" {
		return nil, nil
	}

	ps, err := ar.getTokenPayloadService(app)
	if err != nil {
		return nil, err
	}

	return ps.TokenPayloadForApp(app.ID, app.Name, userID)
}

func (ar *Router) getTokenPayloadService(app model.AppData) (model.TokenPayloadProvider, error) {
	ar.tokenPayloadServicesLock.RLock()

	ps, exists := ar.tokenPayloadServices[app.ID]

	ar.tokenPayloadServicesLock.RUnlock()

	if exists {
		return ps, nil
	}

	ar.tokenPayloadServicesLock.Lock()
	defer ar.tokenPayloadServicesLock.Unlock()

	ps, exists = ar.tokenPayloadServices[app.ID]
	if exists {
		return ps, nil
	}

	var err error

	switch app.TokenPayloadService {
	case model.TokenPayloadServiceHttp:
		ps, err = pphttp.NewTokenPayloadProvider(
			app.TokenPayloadServiceHttpSettings.Secret,
			app.TokenPayloadServiceHttpSettings.URL,
		)

	case model.TokenPayloadServicePlugin:
		ps, err = plugin.NewTokenPayloadProvider(
			model.PluginSettings{
				Cmd:    app.TokenPayloadServicePluginSettings.Cmd,
				Params: app.TokenPayloadServicePluginSettings.Params,
			}, app.TokenPayloadServicePluginSettings.ClientTimeout)
	}

	if err != nil {
		return nil, err
	}

	ar.tokenPayloadServices[app.ID] = ps

	return ps, nil
}

// loginUser creates and returns access token for a user.
// createRefreshToken boolean param tells if we should issue refresh token as well.
func (ar *Router) loginUser(user model.User, scopes []string, app model.AppData, createRefreshToken, require2FA bool, tokenPayload map[string]interface{}) (string, string, error) {
	token, err := ar.server.Services().Token.NewAccessToken(user, scopes, app, require2FA, tokenPayload)
	if err != nil {
		return "", "", err
	}

	accessTokenString, err := ar.server.Services().Token.String(token)
	if err != nil {
		return "", "", err
	}
	if !createRefreshToken || require2FA {
		return accessTokenString, "", nil
	}

	refresh, err := ar.server.Services().Token.NewRefreshToken(user, scopes, app)
	if err != nil {
		ar.logger.Println(err)
		return accessTokenString, "", nil
	}
	refreshTokenString, err := ar.server.Services().Token.String(refresh)
	if err != nil {
		return "", "", err
	}
	return accessTokenString, refreshTokenString, nil
}

func (ar *Router) loginFlow(app model.AppData, user model.User, requestedScopes []string) (AuthResponse, error) {
	// check if the user has the scope, that allows to login to the app
	// user has to have at least one scope app expecting
	if len(app.Scopes) > 0 && len(model.SliceIntersect(app.Scopes, user.Scopes)) == 0 {
		return AuthResponse{}, errors.New("user does not have required scope for the app")
	}

	// Do login flow.
	scopes := []string{}
	// if we requested any scope, let's provide all the scopes user has and requested
	if len(requestedScopes) > 0 {
		scopes = model.SliceIntersect(requestedScopes, user.Scopes)
	}
	if model.SliceContains(requestedScopes, "offline") && app.Offline {
		scopes = append(scopes, "offline")
	}

	// Check if we should require user to authenticate with 2FA.
	require2FA, enabled2FA, err := ar.check2FA(app.TFAStatus, ar.tfaType, user)
	if !require2FA && enabled2FA && err != nil {
		return AuthResponse{}, err
	}

	offline := contains(scopes, model.OfflineScope)
	tokenPayload, err := ar.getTokenPayloadForApp(app, user.ID)
	if err != nil {
		return AuthResponse{}, err
	}

	accessToken, refreshToken, err := ar.loginUser(user, scopes, app, offline, require2FA, tokenPayload)
	if err != nil {
		return AuthResponse{}, err
	}

	result := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Require2FA:   require2FA,
		Enabled2FA:   enabled2FA,
	}

	if require2FA && enabled2FA {
		if err := ar.sendOTPCode(app, user); err != nil {
			return AuthResponse{}, err
		}
	} else {
		ar.server.Storages().User.UpdateLoginMetadata(user.ID)
	}

	user = user.Sanitized()
	result.User = user
	return result, nil
}
