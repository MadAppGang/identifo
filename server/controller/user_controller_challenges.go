package controller

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"html/template"
	"io"
	"math"
	"net/url"
	"time"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
)

// we need to get app, and check if it has strategies for this type of challenge
// if not - return it
// if yes - find the better (now just first suitable) strategy for it
// create full challenge
// send it to the user
// save it in db
// return it back
func (c *UserStorageController) RequestChallenge(ctx context.Context, challenge model.UserAuthChallenge, userIDValue string) (model.UserAuthChallenge, error) {
	zr := model.UserAuthChallenge{}
	app, err := c.as.AppByID(challenge.AppID)
	if err != nil {
		return zr, err
	}

	appAuthStrategies := app.AuthStrategies
	compatibleStrategies := model.FilterCompatible(challenge.Strategy, appAuthStrategies)
	// the app does not supports that type of challenge
	if len(compatibleStrategies) == 0 {
		return zr, l.LocalizedError{ErrID: l.ErrorRequestChallengeUnsupportedByAPP}
	}

	// selecting the first strategy from the list.
	// if there are more than one strategy we need to choose better one.
	auth := compatibleStrategies[0]

	// using the challenge he requested
	// if no user found, just silently return with no error for security reason
	u, err := c.UserByAuthStrategy(ctx, auth, userIDValue)
	// check if there is not user, and app allows to register user, let's create challenge for non-existent user

	// check if we can register passwordless users in the app, if so, let's send a code
	if err != nil && errors.Is(err, l.ErrorUserNotFound) && !app.RegistrationForbidden && app.PasswordlessRegistrationAllowed {
		u = ephemeralUserForStrategy(challenge.Strategy, userIDValue)
	} else if err != nil {
		return zr, nil
	}

	// TODO: Add log entry about the challenge
	cha := challenge
	cha.UserID = u.ID
	cha.Strategy = auth
	cha.Solved = false
	cha.CreatedAt = time.Now()
	cha.ExpiresAt = cha.CreatedAt.Add(time.Minute * model.ExpireChallengeDuration(auth)) // challenge should know how long ti expire
	cha.ExpiresMins = int(math.Round(model.ExpireChallengeDuration(auth).Minutes()))

	// this challenge type is about random OTP code, so generate it
	if auth.Type() == model.AuthStrategyFirstFactorInternal {
		f, ok := auth.(model.FirstFactorInternalStrategy)
		if ok {
			if f.Challenge == model.AuthChallengeTypeOTP ||
				f.Challenge == model.AuthChallengeTypeMagicLink {
				cha.OTP = randomOTP(6)
			}
		}
	}

	ch, err := c.uas.AddChallenge(ctx, cha)
	if err != nil {
		return zr, err
	}

	err = c.sendChallengeToUser(ctx, cha, u)
	if err != nil {
		return zr, err
	}
	return ch, nil
}

func (c *UserStorageController) sendChallengeToUser(ctx context.Context, challenge model.UserAuthChallenge, u model.User) error {
	// no we need send the challenge to the user:
	// - sms
	// - email
	// - push: not supported yet
	// - socket: not supported yet
	// sending the challenge itself with the rendered text, each transport has it's own template
	// - magic link - build a link like reset password, but with OTP code to validate it
	// - OTP -- generate random number
	// using user

	if challenge.Strategy.Type() == model.AuthStrategyFirstFactorInternal {
		st, ok := challenge.Strategy.(model.FirstFactorInternalStrategy)
		if ok {
			// send magic link
			if st.Challenge == model.AuthChallengeTypeMagicLink {
				host := c.h
				path := model.DefaultLoginWebAppSettings.OTPConfirmationURL
				app, err := c.as.AppByID(challenge.AppID)
				if err != nil {
					return err
				}
				if app.LoginAppSettings != nil && app.LoginAppSettings.OTPConfirmationURL != "" {
					ah, err := url.ParseRequestURI(app.LoginAppSettings.OTPConfirmationURL)
					if err == nil {
						// if custom url is valid, use it
						host = ah
						path = ah.Path
					}
				}
				query := fmt.Sprintf("appId=%s&otp=%s", challenge.AppID, challenge.OTP)
				ur := &url.URL{Scheme: host.Scheme, Host: host.Host, Path: path, RawQuery: query}
				// url with no query
				hostUrl := &url.URL{Scheme: ur.Scheme, Host: ur.Host, Path: ur.Path}
				ocd := model.OPTConfirmationData{
					OTP:     challenge.OTP,
					URL:     ur.String(),
					Host:    hostUrl.String(),
					Expires: challenge.ExpiresMins,
					User:    u,
				}

				if st.Transport == model.AuthTransportTypeEmail {
					subfolder := ""
					if app.CustomEmailTemplates {
						subfolder = app.ID
					}
					return c.es.SendUserEmail(model.EmailTemplateTypeOTPMagicLink, subfolder, u, ocd)
				} else if st.Transport == model.AuthTransportTypeSMS {
					sms, err := c.localizedSMSForUser(model.SMSMessageTypeOTPMagicLink, u, app, ocd)
					if err != nil {
						return err
					}
					return c.ss.SendSMS(u.PhoneNumber, sms)
				} else {
					return l.LocalizedError{
						ErrID:   l.ErrorOtpLoginTransportNotSupported,
						Details: []any{st.Transport},
					}
				}
			} else if st.Challenge == model.AuthChallengeTypeOTP {
				app, err := c.as.AppByID(challenge.AppID)
				if err != nil {
					return err
				}
				ocd := model.OPTConfirmationData{
					OTP:     challenge.OTP,
					Expires: challenge.ExpiresMins,
				}

				if st.Transport == model.AuthTransportTypeEmail {
					subfolder := ""
					if app.CustomEmailTemplates {
						subfolder = app.ID
					}
					return c.es.SendUserEmail(model.EmailTemplateTypeOTPCode, subfolder, u, ocd)
				} else if st.Transport == model.AuthTransportTypeSMS {
					sms, err := c.localizedSMSForUser(model.SMSMessageTypeOTPCode, u, app, ocd)
					if err != nil {
						return err
					}
					return c.ss.SendSMS(u.PhoneNumber, sms)
				} else {
					return l.LocalizedError{
						ErrID:   l.ErrorOtpLoginTransportNotSupported,
						Details: []any{st.Transport},
					}
				}
			} else {
				return l.LocalizedError{
					ErrID:   l.ErrorOtpLoginChallengeNotSupported,
					Details: []any{st.Challenge},
				}
			}
		} else {
			// should not happened!!!
			return l.LocalizedError{ErrID: l.APIInternalServerError}
		}
	}
	return l.LocalizedError{
		ErrID:   l.ErrorOtpLoginStrategyNotSupported,
		Details: []any{challenge.Strategy.Type()},
	}
}

func (c *UserStorageController) localizedSMSForUser(smsType model.SMSMessageType, u model.User, app model.AppData, data any) (string, error) {
	// get locale for user if set, otherwise it will returns default locale printer
	p := c.LP.PrinterForLocale(u.Locale)
	stringID := fmt.Sprintf("sms.%s", smsType)

	// default localized string
	sms := p.SD(l.LocalizedString(stringID))

	// get App specific SMS message, if it has one
	if app.CustomSMSMessages != nil {
		// locale specific message?
		if app.CustomSMSMessages[p.DefaultTag().String()] != nil && app.CustomSMSMessages[p.DefaultTag().String()][smsType] != "" {
			sms = app.CustomSMSMessages[p.DefaultTag().String()][smsType]
		} else if app.CustomSMSMessages["default"] != nil && app.CustomSMSMessages["default"][smsType] != "" {
			// fallback to default
			sms = app.CustomSMSMessages["default"][smsType]
		}
		// if no language specific message, and no default message, continue with app default
	}
	tmpl, err := template.New(stringID + "_" + p.DefaultTag().String()).Parse(sms)
	if err != nil {
		return "", err
	}
	var bb bytes.Buffer
	err = tmpl.Execute(&bb, data)
	if err != nil {
		return "", err
	}

	return bb.String(), nil
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func randomOTP(length int) string {
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

// VerifyChallenge verifies challenge from user
func (c *UserStorageController) VerifyChallenge(ctx context.Context, challenge model.UserAuthChallenge, userIDValue string) (model.User, model.AppData, error) {
	app, err := c.as.AppByID(challenge.AppID)
	if err != nil {
		return model.User{}, model.AppData{}, err
	}

	appAuthStrategies := app.AuthStrategies
	compatibleStrategies := model.FilterCompatible(challenge.Strategy, appAuthStrategies)
	// the app does not supports that type of challenge
	if len(compatibleStrategies) == 0 {
		return model.User{}, model.AppData{}, l.LocalizedError{ErrID: l.ErrorRequestChallengeUnsupportedByAPP}
	}

	// selecting the first strategy from the list.
	// if there are more than one strategy we need to choose better one.
	auth := compatibleStrategies[0]

	// using the challenge he requested
	// if no user found, just silently return with no error for security reason
	u, err := c.UserByAuthStrategy(ctx, auth, userIDValue)

	// check if we can register passwordless users in the app, if so, let's send a code
	if err != nil && errors.Is(err, l.ErrorUserNotFound) && !app.RegistrationForbidden && app.PasswordlessRegistrationAllowed {
		u = ephemeralUserForStrategy(challenge.Strategy, userIDValue)
	} else if err != nil {
		return model.User{}, model.AppData{}, err
	}

	// check if user has debug challenge and app allows to use it and it matches the code in request
	// ? does not works for new users, to register you need to use real code (or not?)
	shouldValidateOTP := true
	if app.DebugOTPCodeAllowed && !model.ID(u.ID).IsNewUserID() {
		ud, err := c.u.UserData(ctx, u.ID, model.UserDataFieldDebugOTPCode)
		if err != nil {
			return model.User{}, model.AppData{}, err
		}
		if len(ud.DebugOTPCode) > 0 && ud.DebugOTPCode == challenge.OTP {
			shouldValidateOTP = false
			// TODO: Add log entry about the login with debug challenge
		}
	}

	// if user has not debug challenge, validate the challenge
	if shouldValidateOTP {
		ch, err := c.uas.GetLatestChallenge(ctx, challenge.Strategy, u.ID)
		if err != nil {
			return model.User{}, model.AppData{}, err
		}
		err = ch.Valid()
		if err != nil {
			// TODO: Login attempt to login with invalid code
			return model.User{}, model.AppData{}, err
		}
		if ch.OTP != challenge.OTP {
			// TODO: Login attempt to login with invalid code
			return model.User{}, model.AppData{}, l.LocalizedError{ErrID: l.ErrorOtpIncorrect}
		}
		// add information about context, when the challenge been solved
		ch.SolvedDeviceID = challenge.DeviceID
		ch.SolvedUserAgent = challenge.UserAgent
		c.uas.MarkChallengeAsSolved(ctx, ch)
	}
	return u, app, nil
}

// Passwordless login or register user with challenge
func (c *UserStorageController) LoginOrRegisterUserWithChallenge(ctx context.Context, challenge model.UserAuthChallenge, userIDValue string) (model.User, error) {
	// guard check
	if challenge.Strategy.Type() != model.AuthStrategyFirstFactorInternal {
		return model.User{}, l.LocalizedError{ErrID: l.ErrorLoginTypeNotSupported}
	}

	u, _, err := c.VerifyChallenge(ctx, challenge, userIDValue)
	if err != nil {
		return model.User{}, err
	}

	// let's register the user, if it is new user
	if model.ID(u.ID).IsNewUserID() {
		var err error
		u.ID = "" // clear ID, so database layer should generate new one
		u, err = c.ums.AddUser(ctx, u)
		if err != nil {
			return model.User{}, err
		}
	}

	// TODO: save successful login attempt to the database
	// TODO: update active devices list

	// c.loginFlow(ctx, app, u, challenge.ScopesRequested) // ?? requested scopers

	// check if we have no user exists and the code is valid and app allows to register new passwordless users
	// then we create one
	return model.User{}, nil
}
