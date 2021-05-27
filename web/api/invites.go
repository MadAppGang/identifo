package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/middleware"
)

// RequestInviteLink requests invite link. Invite link will be returned in response even if email or access_role is not specified.
func (ar *Router) RequestInviteLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get requester data
		requesterID := tokenFromContext(r.Context()).UserID()
		requester, err := ar.userStorage.UserByID(requesterID)
		if err != nil {
			ar.Error(w, ErrorAPIUserNotFound, http.StatusUnauthorized, err.Error(), "RequestInviteLink.UserByID")
			return
		}

		d := struct {
			Email       string `json:"email"`
			Role        string `json:"access_role"`
			CallbackURL string `json:"callback_url"`
		}{}
		if err := ar.MustParseJSON(w, r, &d); err != nil {
			ar.Error(w, ErrorAPIRequestBodyInvalid, http.StatusBadRequest, err.Error(), "RequestInviteLink.MustParseJSON")
			return
		}
		if d.Email != "" && !model.EmailRegexp.MatchString(d.Email) {
			ar.Error(w, ErrorAPIRequestBodyInvalid, http.StatusBadRequest, "", "RequestInviteLink.emailRegexp_MatchString")
			return
		}

		_, err = ar.inviteStorage.GetByEmail(d.Email)
		if err != nil && !errors.Is(err, model.ErrorNotFound) {
			ar.Error(w, ErrorAPIInviteUnableToGet, http.StatusInternalServerError, err.Error(), "RequestInviteLink.inviteStorage_GetByEmail")
			return
		}

		if err := ar.inviteStorage.ArchiveAllByEmail(d.Email); err != nil {
			ar.Error(w, ErrorAPIInviteUnableToInvalidate, http.StatusInternalServerError, err.Error(), "RequestInviteLink.inviteStorage_InvalidateAllByEmail")
			return
		}

		inviteToken, err := ar.tokenService.NewInviteToken(d.Email, d.Role)
		if err != nil {
			ar.Error(w, ErrorAPIInviteTokenServerError, http.StatusInternalServerError, err.Error(), "RequestInviteLink.NewInviteToken")
			return
		}

		inviteTokenString, err := ar.tokenService.String(inviteToken)
		if err != nil {
			ar.Error(w, ErrorAPIInviteTokenServerError, http.StatusInternalServerError, err.Error(), "RequestInviteLink.tokenService_String")
			return
		}

		app := middleware.AppFromContext(r.Context())
		scopes := strings.Replace(fmt.Sprintf("%q", app.Scopes), " ", ",", -1)
		query := url.PathEscape(fmt.Sprintf("appId=%s&scopes=%s&token=%s&callbackUrl=%s", app.ID, scopes, inviteTokenString, d.CallbackURL))

		host, err := url.Parse(ar.Host)
		if err != nil {
			ar.Error(w, ErrorAPIInternalServerError, http.StatusInternalServerError, err.Error(), "RequestInviteLink.URL_parse")
			return
		}

		u := &url.URL{
			Scheme:   host.Scheme,
			Host:     host.Host,
			Path:     path.Join(ar.WebRouterPrefix, "register"),
			RawQuery: query,
		}

		// Send email only if it's specified.
		if d.Email != "" {
			uu := &url.URL{Scheme: host.Scheme, Host: host.Host, Path: path.Join(ar.WebRouterPrefix, "register")}
			err = ar.inviteStorage.Save(d.Email, inviteTokenString, d.Role, app.ID, requester.ID, inviteToken.ExpiresAt())
			if err != nil {
				ar.Error(w, ErrorAPIInviteUnableToSave, http.StatusInternalServerError, err.Error(), "RequestInviteLink.inviteStorage_Save")
				return
			}
			requestData := model.InviteEmailData{
				Token:     inviteTokenString,
				App:       app.ID,
				Scopes:    scopes,
				Callback:  d.CallbackURL,
				URL:       u.String(),
				Requester: requester,
				Query:     query,
				Host:      uu.String(),
			}
			err = ar.emailService.SendInviteEmail("Invitation", d.Email, requestData)
			if err != nil {
				ar.Error(w, ErrorAPIEmailNotSent, http.StatusInternalServerError, err.Error(), "RequestInviteLink.SendInviteEmail")
				return
			}
		}
		result := map[string]string{"link": u.String()}
		ar.ServeJSON(w, http.StatusOK, result)
	}
}
