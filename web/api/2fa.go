package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	ijwt "github.com/madappgang/identifo/v2/jwt"
	l "github.com/madappgang/identifo/v2/localization"
	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/xlzd/gotp"
)

type ResetEmailData struct {
	User  model.User
	Token string
	URL   string
	Host  string
	Data  interface{}
}

// EnableTFA enables two-factor authentication for the user.
func (ar *Router) EnableTFA() http.HandlerFunc {
	type requestBody struct {
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	type tfaSecret struct {
		AccessToken     string `json:"access_token,omitempty"`
		ProvisioningURI string `json:"provisioning_uri,omitempty"`
		ProvisioningQR  string `json:"provisioning_qr,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		d := requestBody{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		locale := r.Header.Get("Accept-Language")

		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequestAPPIDInvalid)
			return
		}

		if tfaStatus := app.TFAStatus; tfaStatus == model.TFAStatusDisabled {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequest2FADisabled)
			return
		}

		accessTokenBytes, ok := r.Context().Value(model.TokenRawContextKey).([]byte)
		if !ok {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequestAPPIDInvalid)
			return
		}

		// Get userID from token and update user with this ID.
		userID, err := ar.getTokenSubject(string(accessTokenBytes))
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorAPIRequestTokenSubError, err)
			return
		}

		user, err := ar.server.Storages().User.UserByID(userID)
		if err != nil {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIUserNotFoundError, err)
			return
		}

		if tfaInfo := user.TFAInfo; tfaInfo.IsEnabled && tfaInfo.Secret != "" {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequest2FAAlreadyEnabled)
			return
		}

		tokenPayload, err := ar.getTokenPayloadForApp(app, user.ID)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorAPIAPPUnableToTokenPayloadForAPPError, app.ID, err)
			return
		}

		accessToken, _, err := ar.loginUser(user, model.AllowedScopesSet{}, app, true, tokenPayload)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorTokenUnableToCreateAccessTokenError, err)
			return
		}

		switch ar.tfaType {
		case model.TFATypeApp:
			// For app we just enable tfa and generate secret
			user.TFAInfo = model.TFAInfo{
				IsEnabled: true,
				Secret:    gotp.RandomSecret(16),
			}

			if _, err := ar.server.Storages().User.UpdateUser(userID, user); err != nil {
				ar.Error(w, locale, http.StatusInternalServerError, l.ErrorStorageUpdateUserError, user.ID, err)
				return
			}

			// Send new provisioning uri for authenticator
			uri := gotp.NewDefaultTOTP(user.TFAInfo.Secret).ProvisioningUri(user.Username, app.Name)

			var png []byte
			png, err := qrcode.Encode(uri, qrcode.Medium, 256)
			if err != nil {
				ar.Error(w, locale, http.StatusInternalServerError, l.ErrorAPIRequest2FAUnableToGenerateQrError, err)
				return
			}
			encoded := base64.StdEncoding.EncodeToString(png)

			ar.ServeJSON(w, locale, http.StatusOK, &tfaSecret{ProvisioningURI: uri, ProvisioningQR: encoded, AccessToken: accessToken})
			return
		case model.TFATypeSMS, model.TFATypeEmail:
			// If 2fa is SMS or Email we set it in TFAInfo and enable it only when it will be verified
			if ar.tfaType == model.TFATypeSMS {
				if d.Phone == "" {
					ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequest2FASetPhone)
					return
				}
				user.TFAInfo = model.TFAInfo{Phone: d.Phone}
			}
			if ar.tfaType == model.TFATypeEmail {
				if d.Email == "" {
					ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequest2FASetEmail)
					return
				}
				user.TFAInfo = model.TFAInfo{Email: d.Email}
			}

			if _, err := ar.server.Storages().User.UpdateUser(userID, user); err != nil {
				ar.Error(w, locale, http.StatusInternalServerError, l.ErrorStorageUpdateUserError, user.ID, err)
				return
			}

			// And send OTP code for 2fa
			if err := ar.sendOTPCode(app, user); err != nil {
				ar.Error(w, locale, http.StatusInternalServerError, l.ErrorAPIRequest2FAUnableToSendOtpError, err)
				return
			}

			ar.ServeJSON(w, locale, http.StatusOK, &tfaSecret{AccessToken: accessToken})
			return
		}
		ar.Error(w, locale, http.StatusInternalServerError, l.ErrorAPIRequest2FAUnknownType, ar.tfaType)
	}
}

func (ar *Router) ResendTFA() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		tfaToken, ok := r.Context().Value(model.TokenRawContextKey).([]byte)
		if !ok {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIContextNoToken)
			return
		}

		token, err := ar.server.Services().Token.Parse(string(tfaToken))
		if err != nil {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPITokenParseError, err)
			return
		}

		now := ijwt.TimeFunc().Unix()

		fromIssued := now - token.IssuedAt().Unix()

		if fromIssued < int64(ar.tfaResendTimeout) {
			ar.Error(w, locale, http.StatusBadRequest, l.Error2FAResendTimeout)
			return
		}

		userID := token.Subject()

		user, err := ar.server.Storages().User.UserByID(userID)
		if err != nil {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorStorageUpdateUserError, user.ID, err)
			return
		}

		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIAPPNoAPPInContext)
			return
		}

		scopes := strings.Split(token.Scopes(), " ")

		authResult, resultScopes, err := ar.loginFlow(AuditOperationLoginWith2FA, app, user, scopes, nil)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.APIInternalServerErrorWithError, err)
			return
		}

		ar.server.Storages().Blocklist.Add(string(tfaToken))

		ar.audit(AuditOperationLoginWith2FA,
			user.ID, app.ID, r.UserAgent(), user.AccessRole, resultScopes.Scopes(),
			authResult.AccessToken, authResult.RefreshToken)

		ar.ServeJSON(w, locale, http.StatusOK, authResult)
	}
}

// FinalizeTFA finalizes two-factor authentication.
func (ar *Router) FinalizeTFA() http.HandlerFunc {
	type requestBody struct {
		TFACode string   `json:"tfa_code"`
		Scopes  []string `json:"scopes"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		d := requestBody{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		if len(d.TFACode) == 0 {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequest2FACodeEmpty)
			return
		}

		oldAccessTokenBytes, ok := r.Context().Value(model.TokenRawContextKey).([]byte)
		if !ok {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIContextNoToken)
			return
		}
		oldAccessTokenString := string(oldAccessTokenBytes)

		userID, err := ar.getTokenSubject(oldAccessTokenString)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorAPIRequestTokenSubError, err)
			return
		}

		user, err := ar.server.Storages().User.UserByID(userID)
		if err != nil {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorStorageUpdateUserError, user.ID, err)
			return
		}

		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIAPPNoAPPInContext)
			return
		}

		otpVerified, err := ar.verifyOTPCode(user, d.TFACode)
		if err != nil {
			ar.Error(w, locale, http.StatusForbidden, l.Error2FAVerifyFailError, err)
			return
		}

		dontNeedVerification := app.DebugTFACode != "" && d.TFACode == app.DebugTFACode

		if !(otpVerified || dontNeedVerification) {
			ar.Error(w, locale, http.StatusUnauthorized, l.ErrorAPIRequest2FACodeInvalid)
			return
		}

		scopes := model.AllowedScopes(d.Scopes, app.Scopes, app.Offline)

		tokenPayload, err := ar.getTokenPayloadForApp(app, user.ID)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorAPIAPPUnableToTokenPayloadForAPPError, app.ID, err)
			return
		}

		accessToken, refreshToken, err := ar.loginUser(user, scopes, app, false, tokenPayload)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorTokenUnableToCreateAccessTokenError, err)
			return
		}

		// Blacklist old access token.
		if err := ar.server.Storages().Blocklist.Add(oldAccessTokenString); err != nil {
			ar.logger.Error("Cannot blacklist old access token",
				logging.FieldError, err)
		}

		user = user.Sanitized()
		result := &AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			User:         user,
		}

		// Enable TFA after verify if it not enabled
		if !user.TFAInfo.IsEnabled {
			user.TFAInfo = model.TFAInfo{
				IsEnabled: true,
				Secret:    gotp.RandomSecret(16),
			}

			if _, err := ar.server.Storages().User.UpdateUser(userID, user); err != nil {
				ar.Error(w, locale, http.StatusInternalServerError, l.ErrorStorageUpdateUserError, user.ID, err)
				return
			}
		}

		ar.server.Storages().User.UpdateLoginMetadata(
			string(AuditOperationLoginWith2FA),
			user.ID,
			app.ID,
			scopes.Scopes(),
			tokenPayload,
		)

		ar.audit(AuditOperationLoginWith2FA,
			user.ID, app.ID, r.UserAgent(), user.AccessRole, scopes.Scopes(),
			result.AccessToken, result.RefreshToken)

		ar.ServeJSON(w, locale, http.StatusOK, result)
	}
}

func (ar *Router) verifyOTPCode(user model.User, otp string) (bool, error) {
	result := false
	if ar.tfaType == model.TFATypeApp {
		totp := gotp.NewDefaultTOTP(user.TFAInfo.Secret)
		result = totp.Verify(otp, time.Now().Unix())
	} else {
		if user.TFAInfo.HOTPExpiredAt.Before(time.Now()) {
			return false, errors.New(ar.ls.SD(l.ErrorOtpExpired))
		}
		hotp := gotp.NewDefaultHOTP(user.TFAInfo.Secret)
		result = hotp.Verify(otp, user.TFAInfo.HOTPCounter)
	}
	return result, nil
}

// RequestDisabledTFA requests link for disabling TFA.
func (ar *Router) RequestDisabledTFA() http.HandlerFunc {
	type requestBody struct {
		Email string `json:"email,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		d := requestBody{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		if !model.EmailRegexp.MatchString(d.Email) {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyEmailInvalid)
			return
		}

		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIAPPNoAPPInContext)
			return
		}

		if app.TFAStatus == model.TFAStatusMandatory {
			ar.Error(w, locale, http.StatusForbidden, l.ErrorAPIRequest2FAMandatory)
			return
		}

		user, err := ar.server.Storages().User.UserByEmail(d.Email)
		if err != nil {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorStorageFindUserEmailError, d.Email, err)
			return
		}

		resetToken, err := ar.server.Services().Token.NewResetToken(user.ID)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorTokenUnableToCreateResetTokenError, err)
			return
		}

		resetTokenString, err := ar.server.Services().Token.String(resetToken)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorTokenUnableToCreateResetTokenError, err)
			return
		}

		query := fmt.Sprintf("token=%s", resetTokenString)

		u := &url.URL{
			Scheme:   ar.Host.Scheme,
			Host:     ar.Host.Host,
			Path:     model.DefaultLoginWebAppSettings.TFADisableURL,
			RawQuery: query,
		}

		// rewrite path for app, if app has specific web app login settings
		if app.LoginAppSettings != nil && len(app.LoginAppSettings.TFADisableURL) > 0 {
			appSpecificURL, err := url.Parse(app.LoginAppSettings.TFADisableURL)
			if err != nil {
				ar.Error(w, locale, http.StatusInternalServerError, l.ErrorTokenUnableToCreateResetTokenError, err)
				return
			}

			// app settings could rewrite host or just path, if path is absolute - it rewrites host as well
			if appSpecificURL.IsAbs() {
				u.Scheme = appSpecificURL.Scheme
				u.Host = appSpecificURL.Host
			}

			u.Path = appSpecificURL.Path
		}

		uu := &url.URL{Scheme: u.Scheme, Host: u.Host, Path: u.Path}

		resetEmailData := ResetEmailData{
			User:  user,
			Token: resetTokenString,
			Host:  uu.String(),
			URL:   u.String(),
		}

		if err = ar.server.Services().Email.SendTemplateEmail(
			model.EmailTemplateTypeResetPassword,
			app.GetCustomEmailTemplatePath(),
			"Disable Two-Factor Authentication",
			d.Email,
			model.EmailData{
				User: user,
				Data: resetEmailData,
			},
		); err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorServiceEmailSendError, err)
			return
		}

		result := map[string]string{"result": "ok"}
		ar.ServeJSON(w, locale, http.StatusOK, result)
	}
}

// RequestTFAReset requests link for resetting TFA: deleting old shared secret and establishing the new one.
func (ar *Router) RequestTFAReset() http.HandlerFunc {
	type requestBody struct {
		Email string `json:"email,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		d := requestBody{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		if !model.EmailRegexp.MatchString(d.Email) {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyEmailInvalid)
			return
		}

		user, err := ar.server.Storages().User.UserByEmail(d.Email)
		if err != nil {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorStorageFindUserEmailError, d.Email, err)
			return
		}

		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIAPPNoAPPInContext)
			return
		}

		if app.TFAStatus == model.TFAStatusDisabled {
			ar.Error(w, locale, http.StatusForbidden, l.ErrorAPIRequest2FADisabled)
			return
		}

		resetToken, err := ar.server.Services().Token.NewResetToken(user.ID)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorTokenUnableToCreateResetTokenError, err)
			return
		}

		resetTokenString, err := ar.server.Services().Token.String(resetToken)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorTokenUnableToCreateResetTokenError, err)
			return
		}

		query := fmt.Sprintf("token=%s", resetTokenString)

		u := &url.URL{
			Scheme:   ar.Host.Scheme,
			Host:     ar.Host.Host,
			Path:     model.DefaultLoginWebAppSettings.TFAResetURL,
			RawQuery: query,
		}

		// rewrite path for app, if app has specific web app login settings
		if app.LoginAppSettings != nil && len(app.LoginAppSettings.TFAResetURL) > 0 {
			appResetURL, err := url.Parse(app.LoginAppSettings.TFAResetURL)
			if err != nil {
				ar.Error(w, locale, http.StatusInternalServerError, l.ErrorTokenUnableToCreateResetTokenError, err)
				return
			}

			// app settings could rewrite host or just path, if path is absolute - it rewrites host as well
			if appResetURL.IsAbs() {
				u.Scheme = appResetURL.Scheme
				u.Host = appResetURL.Host
			}

			u.Path = appResetURL.Path
		}

		uu := &url.URL{Scheme: u.Scheme, Host: u.Host, Path: u.Path}

		resetEmailData := ResetEmailData{
			URL:   u.String(),
			User:  user,
			Token: resetTokenString,
			Host:  uu.String(),
		}

		if err = ar.server.Services().Email.SendTemplateEmail(
			model.EmailTemplateTypeResetPassword,
			app.GetCustomEmailTemplatePath(),
			"Disable Two-Factor Authentication",
			d.Email,
			model.EmailData{
				User: user,
				Data: resetEmailData,
			},
		); err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorServiceEmailSendError, err)
			return
		}

		result := map[string]string{"result": "ok"}
		ar.ServeJSON(w, locale, http.StatusOK, result)
	}
}

// check2FA checks correspondence between app's TFAstatus and user's TFAInfo,
// and decides if we require two-factor authentication after all checks are successfully passed.
// require2FA, enabled2FA, err
func (ar *Router) check2FA(appTFAStatus model.TFAStatus, serverTFAType model.TFAType, user model.User) (bool, bool, error) {
	if appTFAStatus == model.TFAStatusMandatory && !user.TFAInfo.IsEnabled {
		return true, false, errPleaseEnableTFA
	}

	// if appTFAStatus == model.TFAStatusDisabled && user.TFAInfo.IsEnabled {
	// 	return false, true, errPleaseDisableTFA
	// }

	// Request two-factor auth if user enabled it and app supports it.
	if user.TFAInfo.IsEnabled && appTFAStatus != model.TFAStatusDisabled {
		if user.TFAInfo.Phone == "" && serverTFAType == model.TFATypeSMS {
			// Server required sms tfa but user phone is empty
			return true, false, errPleaseSetPhoneTFA
		}
		if user.TFAInfo.Email == "" && serverTFAType == model.TFATypeEmail {
			// Server required email tfa but user email is empty
			return true, false, errPleaseSetEmailTFA
		}
		if user.TFAInfo.Secret == "" {
			// Then admin must have enabled TFA for this user manually.
			// User must obtain TFA secret, i.e send EnableTFA request.
			return true, false, errPleaseEnableTFA
		}
		return true, true, nil
	}
	return false, false, nil
}

func (ar *Router) sendTFACodeInSMS(_ model.AppData, phone, otp string) error {
	if phone == "" {
		return errors.New("unable to send SMS OTP, user has no phone number")
	}

	if err := ar.server.Services().SMS.SendSMS(phone, fmt.Sprintf(smsTFACode, otp)); err != nil {
		return fmt.Errorf("unable to send sms. %s", err)
	}
	return nil
}

func (ar *Router) sendTFACodeOnEmail(app model.AppData, user model.User, otp string) error {
	if user.TFAInfo.Email == "" {
		return errors.New("unable to send email OTP, user has no email")
	}

	emailData := SendTFAEmailData{
		User: user,
		OTP:  otp,
	}

	if err := ar.server.Services().Email.SendTemplateEmail(
		model.EmailTemplateTypeTFAWithCode,
		app.GetCustomEmailTemplatePath(),
		"One-time password",
		user.TFAInfo.Email,
		model.EmailData{
			User: user,
			Data: emailData,
		},
	); err != nil {
		return fmt.Errorf("unable to send email with OTP with error: %s", err)
	}

	return nil
}
