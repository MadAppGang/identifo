package controller

import (
	"context"
	"fmt"
	"net/url"

	"github.com/madappgang/identifo/v2/model"
)

// we need to create challenge for user
// challenge type is phone confirmation
// user should enter the code from SMS
// the challenge has a specific TTL
func (c *UserStorageController) SendPhoneConfirmation(ctx context.Context, userID string) error {
	// TODO: Create auth data storage to keep all challenges and use enrollments there
	// TODO: Add log events there as well
	// challenge, err := c.authDataStorage.CreateChallenge(userID, challengeType, challengeTTL)
	// err = c.smsSender.Send(ctx, phoneNumber, challenge.Code)
	return nil
}

// we need to create challenge for user
// challenge type is phone confirmation
// user should enter the code from EMAIL
// it could be done with
// the challenge has a specific TTL
// email confirmation is a link with a code in it
func (c *UserStorageController) SendEmailConfirmation(ctx context.Context, userID string) error {
	return nil
}

func (c *UserStorageController) SendPasswordResetEmail(ctx context.Context, userID, appID string) (model.ResetEmailData, error) {
	user, err := c.u.UserByID(ctx, userID)
	if err != nil {
		return model.ResetEmailData{}, err
	}

	resetToken, err := c.ts.NewToken(model.TokenTypeReset, user, []string{appID}, nil, nil)
	if err != nil {
		return model.ResetEmailData{}, err
	}

	app, err := c.as.AppByID(appID)
	if err != nil {
		return model.ResetEmailData{}, err
	}

	host := c.h
	path := model.DefaultLoginWebAppSettings.ResetPasswordURL

	if app.LoginAppSettings != nil && app.LoginAppSettings.ResetPasswordURL != "" {
		ah, err := url.ParseRequestURI(app.LoginAppSettings.ResetPasswordURL)
		if err == nil {
			// if custom url is valid, use it
			host = ah
			path = ah.Path
		}
	}

	query := fmt.Sprintf("appId=%s&token=%s", appID, resetToken.Raw)

	u := &url.URL{
		Scheme:   host.Scheme,
		Host:     host.Host,
		Path:     path,
		RawQuery: query,
	}
	// url with no query
	hostUrl := &url.URL{Scheme: u.Scheme, Host: u.Host, Path: u.Path}

	red := model.ResetEmailData{
		Token: resetToken.Raw,
		URL:   u.String(),
		Host:  hostUrl.String(),
	}

	subfolder := ""
	if app.CustomEmailTemplates {
		subfolder = app.ID
	}

	err = c.es.SendUserEmail(
		model.EmailTemplateTypeResetPassword,
		subfolder,
		user,
		red,
	)
	if err != nil {
		return model.ResetEmailData{}, err
	}

	return red, nil
}

func (c *UserStorageController) SendInvitationEmail(ctx context.Context, inv model.Invite, u *url.URL, app *model.AppData) error {
	subfolder := "" // custom app email templates
	if app != nil {
		if app.CustomEmailTemplates {
			subfolder = app.ID
		}
	}
	hostUrl := u
	u.RawQuery = ""

	data := model.InvitationEmailData{
		Invitation: inv,
		URL:        u.String(),
		Host:       hostUrl.String(),
	}

	err := c.es.SendUserEmail(
		model.EmailTemplateTypeInvite,
		subfolder,
		model.User{}, // we have no user here, we have only invite
		data,
	)
	return err
}

func (c *UserStorageController) invitationLink(inv model.Invite, app *model.AppData) (*url.URL, error) {
	host := c.h
	path := model.DefaultLoginWebAppSettings.RegisterURL
	query := url.Values{}
	query.Add("token", inv.Token)

	if app != nil {
		// add appID in query
		query.Add("appId", app.ID)

		// if app has custom URL for RegisterURL - user it.
		if app.LoginAppSettings != nil && app.LoginAppSettings.RegisterURL != "" {
			ah, err := url.ParseRequestURI(app.LoginAppSettings.RegisterURL)
			if err == nil {
				// if custom url is valid, use it
				host = ah
				path = ah.Path
			}
		}
	}

	return &url.URL{Scheme: host.Scheme, Host: host.Host, Path: path, RawQuery: query.Encode()}, nil
}
