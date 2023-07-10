package storage

import (
	"bytes"
	"context"
	"crypto/rand"
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
func (c *UserStorageController) RequestChallenge(ctx context.Context, challenge model.UserAuthChallenge) (model.UserAuthChallenge, error) {
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
	u, err := c.UserByAuthStrategy(ctx, auth)
	if err != nil {
		return zr, nil
	}

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

	_ = c.sendChallengeToUser(ctx, cha, u)
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
