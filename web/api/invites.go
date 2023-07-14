package api

import (
	"net/http"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
)

// RequestInviteLink requests invite link. Invite link will be returned in response even if email or access_role is not specified.
func (ar *Router) RequestInviteLink() http.HandlerFunc {
	type requestData struct {
		Email       string         `json:"email"`
		Tenant      string         `json:"tenant"`
		Group       string         `json:"group"`
		Role        string         `json:"role"`
		CallbackURL string         `json:"callback"`
		SendToEmail bool           `json:"send_to_email"`
		Data        map[string]any `json:"data"`
	}

	type responseData struct {
		Invitation model.Invite `json:"invitation"`
		Link       string       `json:"link"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		d := requestData{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		token := tokenFromContext(r.Context())
		app := middleware.AppFromContext(r.Context())

		invitation, link, err := ar.server.Storages().UMC.CreateInvitation(r.Context(), token, d.Tenant, d.Group, d.Role, d.Email, &app)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.LocalizedString(err.Error()))
		}

		// create link

		if d.SendToEmail {
			err = ar.server.Storages().UMC.SendInvitationEmail(r.Context(), invitation, link, &app)
			if err != nil {
				ar.LocalizedError(w, locale, http.StatusInternalServerError, l.LocalizedString(err.Error()))
				return
			}
		}

		// get requester data
		ar.ServeJSON(w, locale, http.StatusOK, responseData{Invitation: invitation, Link: link.String()})
	}
}
