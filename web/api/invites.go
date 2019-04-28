package api

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/madappgang/identifo/web/middleware"
)

// RequestInviteToken - request invite token. Invite link will be returned in response even if email is not specified.
func (ar *Router) RequestInviteLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d := struct {
			Email string `json:"email"`
		}{}
		if ar.MustParseJSON(w, r, &d) != nil {
			ar.Error(w, ErrorWrongInput, http.StatusBadRequest, "Invalid input data. ")
			return
		}
		if d.Email != "" && !emailRegexp.MatchString(d.Email) {
			ar.Error(w, ErrorWrongInput, http.StatusBadRequest, "Invalid email. ")
			return
		}

		t, err := ar.tokenService.NewInviteToken()
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "Unable to create invite token. Try again or contact support team. ")
			return
		}

		token, err := ar.tokenService.String(t)
		if err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "Unable to create invite token. Try again or contact support team. ")
			return
		}

		app := middleware.AppFromContext(r.Context())
		scopes := strings.Replace(fmt.Sprintf("%q", app.Scopes()), " ", ",", -1)
		query := url.PathEscape(fmt.Sprintf("appId=%s&scopes=%s&token=%s", app.ID(), scopes, token))
		host, _ := url.Parse(ar.Host)
		u := &url.URL{
			Scheme:   host.Scheme,
			Host:     host.Host,
			Path:     path.Join(ar.WebRouterPrefix, "register"),
			RawQuery: query,
		}

		// send email only if it's specified.
		if d.Email != "" {
			err = ar.emailService.SendInviteEmail("Invitation", d.Email, u.String())
			if err != nil {
				ar.Error(w, err, http.StatusInternalServerError, "Unable to send email. Try again or contact support team. ")
				return
			}
		}
		result := map[string]string{"link": u.String()}
		ar.ServeJSON(w, http.StatusOK, result)
	}
}
