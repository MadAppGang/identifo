package api

import (
	"net/http"

	"github.com/madappgang/identifo/v2/l"
)

// RequestInviteLink requests invite link. Invite link will be returned in response even if email or access_role is not specified.
func (ar *Router) RequestInviteLink() http.HandlerFunc {
	type requestData struct {
		Email       string         `json:"email"`
		Tenant      string         `json:"tenant"`
		Group       string         `json:"group"`
		Role        string         `json:"role"`
		CallbackURL string         `json:"callback"`
		Data        map[string]any `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		d := requestData{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		token := tokenFromContext(r.Context())

		invitation, err := ar.server.Storages().UMC.CreateInvitation(r.Context(), token, d.Tenant, d.Group, d.Role, d.Email)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.LocalizedString(err.Error()))
		}

		// create link

		// get requester data
		ar.ServeJSON(w, locale, http.StatusOK, invitation)
	}
}
