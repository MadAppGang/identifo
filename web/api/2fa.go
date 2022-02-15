package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	ijwt "github.com/madappgang/identifo/v2/jwt"
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

		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App is not in context.", "EnableTFA.AppFromContext")
			return
		}

		if tfaStatus := app.TFAStatus; tfaStatus == model.TFAStatusDisabled {
			ar.Error(w, ErrorAPIRequestBodyParamsInvalid, http.StatusBadRequest, "TFA is not supported by this app", "EnableTFA.TFAStatus")
			return
		}

		accessTokenBytes, ok := r.Context().Value(model.TokenRawContextKey).([]byte)
		if !ok {
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "Token bytes are not in context.", "EnableTFA.TokenBytesFromContext")
			return
		}

		// Get userID from token and update user with this ID.
		userID, err := ar.getTokenSubject(string(accessTokenBytes))
		if err != nil {
			ar.Error(w, ErrorAPIAppCannotExtractTokenSubject, http.StatusInternalServerError, err.Error(), "EnableTFA.getTokenSubject")
			return
		}

		user, err := ar.server.Storages().User.UserByID(userID)
		if err != nil {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusBadRequest, err.Error(), "EnableTFA.UserByID")
			return
		}

		if tfaInfo := user.TFAInfo; tfaInfo.IsEnabled && tfaInfo.Secret != "" {
			ar.Error(w, ErrorAPIRequestTFAAlreadyEnabled, http.StatusBadRequest, "TFA already enabled for this user", "EnableTFA.alreadyEnabled")
			return
		}

		tokenPayload, err := ar.getTokenPayloadForApp(app, user)
		if err != nil {
			ar.Error(w, ErrorAPIAppAccessTokenNotCreated, http.StatusInternalServerError, err.Error(), "EnableTFA.accessToken")
			return
		}

		accessToken, _, err := ar.loginUser(user, []string{}, app, false, true, tokenPayload)
		if err != nil {
			ar.Error(w, ErrorAPIAppAccessTokenNotCreated, http.StatusInternalServerError, err.Error(), "EnableTFA.accessToken")
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
				ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "EnableTFA.UpdateUser")
				return
			}

			// Send new provising uri for authenticator
			uri := gotp.NewDefaultTOTP(user.TFAInfo.Secret).ProvisioningUri(user.Username, app.Name)

			var png []byte
			png, err := qrcode.Encode(uri, qrcode.Medium, 256)
			if err != nil {
				ar.Error(w, ErrorAPIAppAccessTokenNotCreated, http.StatusInternalServerError, err.Error(), "EnableTFA.QRgenerate")
				return
			}
			encoded := base64.StdEncoding.EncodeToString(png)

			ar.ServeJSON(w, http.StatusOK, &tfaSecret{ProvisioningURI: uri, ProvisioningQR: encoded, AccessToken: accessToken})
			return
		case model.TFATypeSMS, model.TFATypeEmail:
			// If 2fa is SMS or Email we set it in TFAInfo and enable it only when it will be verified
			if ar.tfaType == model.TFATypeSMS {
				if d.Phone == "" {
					ar.Error(w, ErrorAPIRequestPleaseSetPhoneForTFA, http.StatusBadRequest, "Please specify your phone number to be able to receive one-time passwords", "EnableTFA.setPhone")
					return
				}
				user.TFAInfo = model.TFAInfo{Phone: d.Phone}
			}
			if ar.tfaType == model.TFATypeEmail {
				if d.Email == "" {
					ar.Error(w, ErrorAPIRequestPleaseSetEmailForTFA, http.StatusBadRequest, "Please specify your email address to be able to receive one-time passwords", "EnableTFA.setEmail")
					return
				}
				user.TFAInfo = model.TFAInfo{Email: d.Email}
			}

			if _, err := ar.server.Storages().User.UpdateUser(userID, user); err != nil {
				ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "EnableTFA.UpdateUser")
				return
			}

			// And send OTP code for 2fa
			if err := ar.sendOTPCode(user); err != nil {
				ar.Error(w, ErrorAPIRequestUnableToSendOTP, http.StatusInternalServerError, err.Error(), "EnableTFA.sendOTP")
				return
			}

			ar.ServeJSON(w, http.StatusOK, &tfaSecret{AccessToken: accessToken})
			return
		}
		ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, fmt.Sprintf("Unknown tfa type '%s'", ar.tfaType), "switch.tfaType")
	}
}
func (ar *Router) ResendTFA() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tfaToken, ok := r.Context().Value(model.TokenRawContextKey).([]byte)
		if !ok {
			ar.Error(w, ErrorAPIRequestTokenInvalid, http.StatusBadRequest, "Token bytes are not in context.", "ResendTFA.TokenBytesFromContext")
			return
		}

		token, err := ar.server.Services().Token.Parse(string(tfaToken))
		if err != nil {
			ar.Error(w, ErrorAPIRequestTokenInvalid, http.StatusBadRequest, "Can't parse token.", "ResendTFA.Parse")
			return
		}

		now := ijwt.TimeFunc().Unix()

		fromIssued := now - token.IssuedAt().Unix()

		if fromIssued < int64(ar.tfaResendTimeout) {
			ar.Error(w, ErrorAPIRequestTokenInvalid, http.StatusBadRequest, "Please wait before new code resend.", "ResendTFA.timeout")
			return
		}

		userID := token.Subject()
		if err != nil {
			ar.Error(w, ErrorAPIAppCannotExtractTokenSubject, http.StatusInternalServerError, err.Error(), "ResendTFA.getTokenSubject")
			return
		}

		user, err := ar.server.Storages().User.UserByID(userID)
		if err != nil {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusBadRequest, err.Error(), "ResendTFA.UserByID")
			return
		}

		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App is not in context.", "ResendTFA.AppFromContext")
			return
		}

		authResult, err := ar.loginFlow(app, user, strings.Split(token.Scopes(), " "))
		if err != nil {
			ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "LoginWithPassword.LoginFlowError")
			return
		}

		ar.server.Storages().Blocklist.Add(string(tfaToken))

		ar.ServeJSON(w, http.StatusOK, authResult)
	}
}

// FinalizeTFA finalizes two-factor authentication.
func (ar *Router) FinalizeTFA() http.HandlerFunc {
	type requestBody struct {
		TFACode string   `json:"tfa_code"`
		Scopes  []string `json:"scopes"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		d := requestBody{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		if len(d.TFACode) == 0 {
			ar.Error(w, ErrorAPIRequestTFACodeEmpty, http.StatusBadRequest, "", "FinalizeTFA.empty")
			return
		}

		oldAccessTokenBytes, ok := r.Context().Value(model.TokenRawContextKey).([]byte)
		if !ok {
			ar.Error(w, ErrorAPIRequestTokenInvalid, http.StatusBadRequest, "Token bytes are not in context.", "FinalizeTFA.TokenBytesFromContext")
			return
		}
		oldAccessTokenString := string(oldAccessTokenBytes)

		userID, err := ar.getTokenSubject(oldAccessTokenString)
		if err != nil {
			ar.Error(w, ErrorAPIAppCannotExtractTokenSubject, http.StatusInternalServerError, err.Error(), "FinalizeTFA.getTokenSubject")
			return
		}

		user, err := ar.server.Storages().User.UserByID(userID)
		if err != nil {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusBadRequest, err.Error(), "FinalizeTFA.UserByID")
			return
		}

		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App is not in context.", "FinalizeTFA.AppFromContext")
			return
		}

		otpVerified, err := ar.verifyOTPCode(user, d.TFACode)
		if err != nil {
			ar.Error(w, ErrorAPIRequestScopesForbidden, http.StatusForbidden, err.Error(), "FinalizeTFA.OTP_Invalid")
			return
		}

		dontNeedVerification := app.DebugTFACode != "" && d.TFACode == app.DebugTFACode

		if !(otpVerified || dontNeedVerification) {
			ar.Error(w, ErrorAPIRequestTFACodeInvalid, http.StatusUnauthorized, "", "FinalizeTFA.OTP_Invalid")
			return
		}

		scopes := []string{}
		// if we requested any scope, let's provide all the scopes user has and requested
		if len(d.Scopes) > 0 {
			scopes = model.SliceIntersect(d.Scopes, user.Scopes)
		}
		if model.SliceContains(d.Scopes, "offline") && app.Offline {
			scopes = append(scopes, "offline")
		}

		tokenPayload, err := ar.getTokenPayloadForApp(app, user)
		if err != nil {
			ar.Error(w, ErrorAPIAppAccessTokenNotCreated, http.StatusInternalServerError, err.Error(), "LoginWithPassword.loginUser")
			return
		}

		createRefreshToken := contains(scopes, model.OfflineScope)
		accessToken, refreshToken, err := ar.loginUser(user, scopes, app, createRefreshToken, false, tokenPayload)
		if err != nil {
			ar.Error(w, ErrorAPIAppAccessTokenNotCreated, http.StatusInternalServerError, err.Error(), "FinalizeTFA.loginUser")
			return
		}

		// Blacklist old access token.
		if err := ar.server.Storages().Blocklist.Add(oldAccessTokenString); err != nil {
			ar.logger.Printf("Cannot blacklist old access token: %s\n", err)
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
				ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "EnableTFA.UpdateUser")
				return
			}
		}

		ar.server.Storages().User.UpdateLoginMetadata(user.ID)
		ar.ServeJSON(w, http.StatusOK, result)
	}
}

func (ar *Router) verifyOTPCode(user model.User, otp string) (bool, error) {
	result := false
	if ar.tfaType == model.TFATypeApp {
		totp := gotp.NewDefaultTOTP(user.TFAInfo.Secret)
		result = totp.Verify(otp, int(time.Now().Unix()))
	} else {
		if user.TFAInfo.HOTPExpiredAt.Before(time.Now()) {
			return false, errors.New("OTP token expired, please get the new one and try again")
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
		d := requestBody{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		if !model.EmailRegexp.MatchString(d.Email) {
			ar.Error(w, ErrorAPIRequestBodyInvalid, http.StatusBadRequest, "", "RequestDisabledTFA.emailRegexp_MatchString")
			return
		}

		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App is not in context.", "RequestDisabledTFA.AppFromContext")
			return
		}

		if app.TFAStatus == model.TFAStatusMandatory {
			ar.Error(w, ErrorAPIRequestMandatoryTFA, http.StatusForbidden, "Two-factor authentication is mandatory for this app", "RequestDisabledTFA.TFAStatusMandatory")
			return
		}

		user, err := ar.server.Storages().User.UserByEmail(d.Email)
		if err != nil {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusBadRequest, err.Error(), "RequestDisabledTFA.UserByEmail")
			return
		}

		resetToken, err := ar.server.Services().Token.NewResetToken(user.ID)
		if err != nil {
			ar.Error(w, ErrorAPIAppResetTokenNotCreated, http.StatusInternalServerError, err.Error(), "RequestDisabledTFA.NewResetToken")
			return
		}

		resetTokenString, err := ar.server.Services().Token.String(resetToken)
		if err != nil {
			ar.Error(w, ErrorAPIAppResetTokenNotCreated, http.StatusInternalServerError, err.Error(), "RequestDisabledTFA.tokenService_String")
			return
		}

		host, err := url.Parse(ar.Host)
		if err != nil {
			ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "RequestDisabledTFA.URL_parse")
			return
		}

		query := fmt.Sprintf("token=%s", resetTokenString)

		u := &url.URL{
			Scheme:   host.Scheme,
			Host:     host.Host,
			Path:     path.Join(ar.LoginAppPath, "tfa/disable"),
			RawQuery: query,
		}
		uu := &url.URL{
			Scheme: host.Scheme,
			Host:   host.Host,
			Path:   path.Join(ar.LoginAppPath, "tfa/disable"),
		}
		resetEmailData := ResetEmailData{
			User:  user,
			Token: resetTokenString,
			Host:  uu.String(),
			URL:   u.String(),
		}

		if err = ar.server.Services().Email.SendTemplateEmail(
			model.EmailTemplateTypeResetPassword,
			"Disable Two-Factor Authentication",
			d.Email,
			model.EmailData{
				User: user,
				Data: resetEmailData,
			},
		); err != nil {
			ar.Error(
				w,
				ErrorAPIEmailNotSent,
				http.StatusInternalServerError,
				"Email sending error: "+err.Error(), "RequestDisabledTFA.SendResetEmail",
			)
			return
		}

		result := map[string]string{"result": "ok"}
		ar.ServeJSON(w, http.StatusOK, result)
	}
}

// RequestTFAReset requests link for resetting TFA: deleting old shared secret and establishing the new one.
func (ar *Router) RequestTFAReset() http.HandlerFunc {
	type requestBody struct {
		Email string `json:"email,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		d := requestBody{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		if !model.EmailRegexp.MatchString(d.Email) {
			ar.Error(w, ErrorAPIRequestBodyInvalid, http.StatusBadRequest, "", "RequestTFAReset.emailRegexp_MatchString")
			return
		}

		user, err := ar.server.Storages().User.UserByEmail(d.Email)
		if err != nil {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusBadRequest, err.Error(), "RequestTFAReset.UserByEmail")
			return
		}

		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.Error(w, ErrorAPIRequestAppIDInvalid, http.StatusBadRequest, "App is not in context.", "RequestDisabledTFA.AppFromContext")
			return
		}

		if app.TFAStatus == model.TFAStatusDisabled {
			ar.Error(w, ErrorAPIRequestDisabledTFA, http.StatusForbidden, "Two-factor authentication is disabled for this app", "RequestTFAReset.TFAStatusDisabled")
			return
		}

		resetToken, err := ar.server.Services().Token.NewResetToken(user.ID)
		if err != nil {
			ar.Error(w, ErrorAPIAppResetTokenNotCreated, http.StatusInternalServerError, err.Error(), "RequestTFAReset.NewResetToken")
			return
		}

		resetTokenString, err := ar.server.Services().Token.String(resetToken)
		if err != nil {
			ar.Error(w, ErrorAPIAppResetTokenNotCreated, http.StatusInternalServerError, err.Error(), "RequestTFAReset.tokenService_String")
			return
		}

		host, err := url.Parse(ar.Host)
		if err != nil {
			ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "RequestTFAReset.URL_parse")
			return
		}

		query := fmt.Sprintf("token=%s", resetTokenString)

		u := &url.URL{
			Scheme:   host.Scheme,
			Host:     host.Host,
			Path:     path.Join(ar.LoginAppPath, "tfa/reset"),
			RawQuery: query,
		}
		uu := &url.URL{
			Scheme: host.Scheme,
			Host:   host.Host,
			Path:   path.Join(ar.LoginAppPath, "tfa/reset"),
		}

		resetEmailData := ResetEmailData{
			URL:   u.String(),
			User:  user,
			Token: resetTokenString,
			Host:  uu.String(),
		}

		if err = ar.server.Services().Email.SendTemplateEmail(
			model.EmailTemplateTypeResetPassword,
			"Disable Two-Factor Authentication",
			d.Email,
			model.EmailData{
				User: user,
				Data: resetEmailData,
			},
		); err != nil {
			ar.Error(
				w,
				ErrorAPIEmailNotSent,
				http.StatusInternalServerError,
				"Email sending error: "+err.Error(), "RequestDisabledTFA.SendResetEmail",
			)
			return
		}

		result := map[string]string{"result": "ok"}
		ar.ServeJSON(w, http.StatusOK, result)
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

func (ar *Router) sendTFACodeInSMS(phone, otp string) error {
	if phone == "" {
		return errors.New("unable to send SMS OTP, user has no phone number")
	}

	if err := ar.server.Services().SMS.SendSMS(phone, fmt.Sprintf(smsTFACode, otp)); err != nil {
		return fmt.Errorf("unable to send sms. %s", err)
	}
	return nil
}

func (ar *Router) sendTFACodeOnEmail(user model.User, otp string) error {
	if user.TFAInfo.Email == "" {
		return errors.New("unable to send email OTP, user has no email")
	}

	emailData := SendTFAEmailData{
		User: user,
		OTP:  otp,
	}

	if err := ar.server.Services().Email.SendTemplateEmail(
		model.EmailTemplateTypeTFAWithCode,
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
