package management

import (
	"net/http"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
)

// InviteRequest is a request for invite.
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

func (ar *Router) getInviteToken(w http.ResponseWriter, r *http.Request) {
	locale := r.Header.Get("Accept-Language")

	d := requestData{}
	if ar.MustParseJSON(w, r, &d) != nil {
		return
	}

	invitation, link, err := ar.server.Storages().UMC.CreateInvitation(r.Context(), nil, d.Tenant, d.Group, d.Role, d.Email, nil)
	if err != nil {
		ar.Error(w, locale, http.StatusInternalServerError, l.LocalizedString(err.Error()))
		return
	}

	if d.SendToEmail {
		err = ar.server.Storages().UMC.SendInvitationEmail(r.Context(), invitation, link, nil)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.LocalizedString(err.Error()))
			return
		}
	}

	ar.ServeJSON(w, locale, http.StatusOK, responseData{Invitation: invitation, Link: link.String()})
}
