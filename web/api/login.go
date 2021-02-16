package api

import (
	"fmt"
	"net/http"

	jwtService "github.com/madappgang/identifo/jwt/service"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/authorization"
	"github.com/madappgang/identifo/web/middleware"
	"github.com/xlzd/gotp"
)

var (
	errPleaseEnableTFA  = fmt.Errorf("Please enable two-factor authentication to be able to use this app")
	errPleaseDisableTFA = fmt.Errorf("Please disable two-factor authentication to be able to use this app")
)

const smsTFACode = "%v is your one-time password!"

// AuthResponse is a response with successful auth data.
type AuthResponse struct {
	AccessToken    string     `json:"access_token,omitempty"`
	RefreshToken   string     `json:"refresh_token,omitempty"`
	User           model.User `json:"user,omitempty"`
	NeedFurtherTFA bool       `json:"need_further_tfa,omitempty"`
}

type loginData struct {
	Username    string   `json:"username,omitempty"`
	Password    string   `json:"password,omitempty"`
	DeviceToken string   `json:"device_token,omitempty"`
	Scopes      []string `json:"scopes,omitempty"`
}

func (ld *loginData) validate() error {
	usernameLen := len(ld.Username)
	if usernameLen < 6 || usernameLen > 130 {
		return fmt.Errorf("Incorrect username length %d, expected a number between 6 and 130", usernameLen)
	}
	pswdLen := len(ld.Password)
	if pswdLen < 6 || pswdLen > 130 {
		return fmt.Errorf("Incorrect password length %d, expected a number between 6 and 130", pswdLen)
	}
	return nil
}

// LoginWithPassword logs user in with username and password.
func (ar *Router) LoginWithPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !ar.SupportedLoginWays.Username {
			ar.Error(w, ErrorAPIAppLoginWithUsernameNotSupported, http.StatusBadRequest, "Application does not support login with username", "LoginWithPassword.supportedLoginWays")
			return
		}

		ld := loginData{}
		if ar.MustParseJSON(w, r, &ld) != nil {
			return
		}

		if err := ld.validate(); err != nil {
			ar.Error(w, ErrorAPIRequestBodyParamsInvalid, http.StatusBadRequest, err.Error(), "LoginWithPassword.validate")
			return
		}

		user, err := ar.userStorage.UserByNamePassword(ld.Username, ld.Password)
		if err != nil {
			ar.Error(w, ErrorAPIRequestIncorrectEmailOrPassword, http.StatusUnauthorized, err.Error(), "LoginWithPassword.UserByNamePassword")
			return
		}

		scopes, err := ar.userStorage.RequestScopes(user.ID(), ld.Scopes)
		if err != nil {
			ar.Error(w, ErrorAPIRequestScopesForbidden, http.StatusForbidden, err.Error(), "LoginWithPassword.RequestScopes")
			return
		}

		app := middleware.AppFromContext(r.Context())
		if app == nil {
			ar.logger.Println("Error getting App")
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App is not in context.", "LoginWithPassword.AppFromContext")
			return
		}

		// Authorize user if the app requires authorization.
		azi := authorization.AuthzInfo{
			App:         app,
			UserRole:    user.AccessRole(),
			ResourceURI: r.RequestURI,
			Method:      r.Method,
		}
		if err := ar.Authorizer.Authorize(azi); err != nil {
			ar.Error(w, ErrorAPIAppAccessDenied, http.StatusForbidden, err.Error(), "LoginWithPassword.Authorizer")
			return
		}

		// Check if we should require user to authenticate with 2FA.
		require2FA, err := ar.check2FA(w, app.TFAStatus(), user.TFAInfo())
		if err != nil {
			return
		}

		offline := contains(scopes, jwtService.OfflineScope)
		accessToken, refreshToken, err := ar.loginUser(user, scopes, app, offline, require2FA)
		if err != nil {
			ar.Error(w, ErrorAPIAppAccessTokenNotCreated, http.StatusInternalServerError, err.Error(), "LoginWithPassword.loginUser")
			return
		}

		user.Sanitize()
		result := AuthResponse{
			AccessToken:    accessToken,
			RefreshToken:   refreshToken,
			User:           user,
			NeedFurtherTFA: require2FA,
		}

		if !require2FA {
			ar.userStorage.UpdateLoginMetadata(user.ID())
			ar.ServeJSON(w, http.StatusOK, result)
			return
		}

		totp := gotp.NewDefaultTOTP(user.TFAInfo().Secret).Now()

		switch ar.tfaType {
		case model.TFATypeSMS:
			ar.sendTFACodeInSMS(w, user.Phone(), totp)
			return
		case model.TFATypeEmail:
			ar.sendTFACodeOnEmail(w, user.Email(), totp)
			return
		default:
			ar.ServeJSON(w, http.StatusOK, result)
		}
	}
}

// IsLoggedIn is for checking whether user is logged in or not.
// In fact, all needed work is done in Token middleware.
// If we reached this code, user is logged in (presented valid and not blacklisted access token).
func (ar *Router) IsLoggedIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ar.ServeJSON(w, http.StatusOK, nil)
	}
}

// loginUser creates and returns access token for a user.
// createRefreshToken boolean param tells if we should issue refresh token as well.
func (ar *Router) loginUser(user model.User, scopes []string, app model.AppData, createRefreshToken, require2FA bool) (accessTokenString, refreshTokenString string, err error) {
	token, err := ar.tokenService.NewAccessToken(user, scopes, app, require2FA)
	if err != nil {
		return
	}

	accessTokenString, err = ar.tokenService.String(token)
	if err != nil {
		return
	}
	if !createRefreshToken || require2FA {
		return
	}

	refresh, err := ar.tokenService.NewRefreshToken(user, scopes, app)
	if err != nil {
		return
	}
	refreshTokenString, err = ar.tokenService.String(refresh)
	if err != nil {
		return
	}
	return
}

// check2FA checks correspondence between app's TFAstatus and user's TFAInfo,
// and decides if we require two-factor authentication after all checks are successfully passed.
func (ar *Router) check2FA(w http.ResponseWriter, appTFAStatus model.TFAStatus, userTFAInfo model.TFAInfo) (bool, error) {
	if appTFAStatus == model.TFAStatusMandatory && !userTFAInfo.IsEnabled {
		ar.Error(w, ErrorAPIRequestPleaseEnableTFA, http.StatusBadRequest, errPleaseEnableTFA.Error(), "check2FA.mandatory")
		return false, errPleaseEnableTFA
	}

	if appTFAStatus == model.TFAStatusDisabled && userTFAInfo.IsEnabled {
		ar.Error(w, ErrorAPIRequestPleaseDisableTFA, http.StatusBadRequest, errPleaseDisableTFA.Error(), "check2FA.appDisabled_userEnabled")
		return false, errPleaseDisableTFA
	}

	// Request two-factor auth if user enabled it and app supports it.
	if userTFAInfo.IsEnabled && appTFAStatus != model.TFAStatusDisabled {
		if userTFAInfo.Secret == "" {
			// Then admin must have enabled TFA for this user manually.
			// User must obtain TFA secret, i.e send EnableTFA request.
			ar.Error(w, ErrorAPIRequestPleaseEnableTFA, http.StatusConflict, errPleaseEnableTFA.Error(), "check2FA.pleaseEnable")
			return false, errPleaseEnableTFA
		}
		return true, nil
	}
	return false, nil
}

func (ar *Router) sendTFACodeInSMS(w http.ResponseWriter, phone, totp string) {
	if phone == "" {
		ar.Error(w, ErrorAPIRequestPleaseSetPhoneForTFA, http.StatusBadRequest, "", "tfaInSMS.empty_phone")
		return
	}

	if err := ar.smsService.SendSMS(phone, fmt.Sprintf(smsTFACode, totp)); err != nil {
		err = fmt.Errorf("Unable to send sms. %s", err)
		ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "tfaInSMS.SendSMS")
		return
	}
	ar.ServeJSON(w, http.StatusOK, nil)
}

func (ar *Router) sendTFACodeOnEmail(w http.ResponseWriter, email, totp string) {
	if email == "" {
		ar.Error(w, ErrorAPIRequestPleaseSetEmailForTFA, http.StatusBadRequest, "", "tfaInSMS.empty_email")
		return
	}

	if err := ar.emailService.SendTFAEmail("One-time password", email, totp); err != nil {
		err = fmt.Errorf("Unable to send email. %s", err)
		ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "tfaInSMS.SendTFAEmail")
		return
	}
	ar.ServeJSON(w, http.StatusOK, nil)
}
