package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
)

type InviteEmailData struct {
	Requester model.User
	Token     string
	URL       string
	Host      string
	Query     string
	App       string
	Scopes    string
	Callback  string
	Data      interface{}
}

// RequestInviteLink requests invite link. Invite link will be returned in response even if email or access_role is not specified.
func (ar *Router) RequestInviteLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get requester data
		requesterID := tokenFromContext(r.Context()).UserID()
		audience := tokenFromContext(r.Context()).Audience()
		requester, err := ar.server.Storages().User.UserByID(requesterID)
		if err != nil {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusUnauthorized, err.Error(), "RequestInviteLink.UserByID")
			return
		}

		d := struct {
			Email       string                 `json:"email"`
			Role        string                 `json:"access_role"`
			CallbackURL string                 `json:"callback_url"`
			Data        map[string]interface{} `json:"data"`
		}{}
		if err := ar.MustParseJSON(w, r, &d); err != nil {
			ar.Error(w, ErrorAPIRequestBodyInvalid, http.StatusBadRequest, err.Error(), "RequestInviteLink.MustParseJSON")
			return
		}
		if d.Email != "" && !model.EmailRegexp.MatchString(d.Email) {
			ar.Error(w, ErrorAPIRequestBodyInvalid, http.StatusBadRequest, "", "RequestInviteLink.emailRegexp_MatchString")
			return
		}

		_, err = ar.server.Storages().Invite.GetByEmail(d.Email)
		if err != nil && !errors.Is(err, model.ErrorNotFound) {
			ar.Error(w, ErrorAPIInviteUnableToGet, http.StatusInternalServerError, err.Error(), "RequestInviteLink.inviteStorage_GetByEmail")
			return
		}

		if err := ar.server.Storages().Invite.ArchiveAllByEmail(d.Email); err != nil {
			ar.Error(w, ErrorAPIInviteUnableToInvalidate, http.StatusInternalServerError, err.Error(), "RequestInviteLink.inviteStorage_InvalidateAllByEmail")
			return
		}

		inviteToken, err := ar.server.Services().Token.NewInviteToken(d.Email, d.Role, audience, d.Data)
		if err != nil {
			ar.Error(w, ErrorAPIInviteTokenServerError, http.StatusInternalServerError, err.Error(), "RequestInviteLink.NewInviteToken")
			return
		}

		inviteTokenString, err := ar.server.Services().Token.String(inviteToken)
		if err != nil {
			ar.Error(w, ErrorAPIInviteTokenServerError, http.StatusInternalServerError, err.Error(), "RequestInviteLink.tokenService_String")
			return
		}

		app := middleware.AppFromContext(r.Context())
		scopes := strings.Replace(fmt.Sprintf("%q", app.Scopes), " ", ",", -1)
		query := url.PathEscape(fmt.Sprintf("email=%s&appId=%s&scopes=%s&token=%s&callbackUrl=%s", d.Email, app.ID, scopes, inviteTokenString, d.CallbackURL))

		u := &url.URL{
			Scheme:   ar.Host.Scheme,
			Host:     ar.Host.Host,
			Path:     model.DefaultLoginWebAppSettings.RegisterURL,
			RawQuery: query,
		}

		// rewrite path for app, if app has specific web app login settings
		if app.LoginAppSettings != nil && len(app.LoginAppSettings.RegisterURL) > 0 {
			appSpecificURL, err := url.Parse(app.LoginAppSettings.RegisterURL)
			if err != nil {
				ar.Error(w, ErrorAPIAppResetTokenNotCreated, http.StatusInternalServerError, err.Error(), "RequestDisabledTFA.app_register_url_parse_error")
				return
			}

			// app settings could rewrite host or just path, if path is absolute - it rewrites host as well
			if appSpecificURL.IsAbs() {
				u.Scheme = appSpecificURL.Scheme
				u.Host = appSpecificURL.Host
			}

			u.Path = appSpecificURL.Path
		}

		// Send email only if it's specified.
		if d.Email != "" {
			uu := &url.URL{Scheme: u.Scheme, Host: u.Host, Path: u.Path}
			err = ar.server.Storages().Invite.Save(d.Email, inviteTokenString, d.Role, app.ID, requester.ID, inviteToken.ExpiresAt())
			if err != nil {
				ar.Error(w, ErrorAPIInviteUnableToSave, http.StatusInternalServerError, err.Error(), "RequestInviteLink.inviteStorage_Save")
				return
			}
			requestData := InviteEmailData{
				Token:     inviteTokenString,
				App:       app.ID,
				Scopes:    scopes,
				Callback:  d.CallbackURL,
				URL:       u.String(),
				Requester: requester,
				Query:     query,
				Host:      uu.String(),
			}

			if err = ar.server.Services().Email.SendTemplateEmail(
				model.EmailTemplateTypeInvite,
				app.GetCustomEmailTemplatePath(),
				"Invitation",
				d.Email,
				model.EmailData{
					Data: requestData,
				},
			); err != nil {
				ar.Error(
					w,
					ErrorAPIEmailNotSent,
					http.StatusInternalServerError,
					"Email sending error: "+err.Error(), "RequestInviteLink.SendInviteEmail",
				)
				return
			}

		}
		result := map[string]string{"result": "ok", "link": u.String()}
		ar.ServeJSON(w, http.StatusOK, result)
	}
}
