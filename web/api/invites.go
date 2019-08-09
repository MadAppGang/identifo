package api

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/middleware"
)

// RequestInviteLink requests invite link. Invite link will be returned in response even if email is not specified.
func (ar *Router) RequestInviteLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d := struct {
			Email string `json:"email"`
		}{}
		if err := ar.MustParseJSON(w, r, &d); err != nil {
			ar.Error(w, ErrorAPIRequestBodyInvalid, http.StatusBadRequest, err.Error(), "RequestInviteLink.MustParseJSON")
			return
		}
		if d.Email != "" && !model.EmailRegexp.MatchString(d.Email) {
			ar.Error(w, ErrorAPIRequestBodyInvalid, http.StatusBadRequest, "", "RequestInviteLink.emailRegexp_MatchString")
			return
		}

		inviteToken, err := ar.tokenService.NewInviteToken()
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
		scopes := strings.Replace(fmt.Sprintf("%q", app.Scopes()), " ", ",", -1)
		query := url.PathEscape(fmt.Sprintf("appId=%s&scopes=%s&token=%s", app.ID(), scopes, inviteTokenString))

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
			err = ar.emailService.SendInviteEmail("Invitation", d.Email, u.String())
			if err != nil {
				ar.Error(w, ErrorAPIEmailNotSent, http.StatusInternalServerError, err.Error(), "RequestInviteLink.SendInviteEmail")
				return
			}
		}
		result := map[string]string{"link": u.String()}
		ar.ServeJSON(w, http.StatusOK, result)
	}
}
