package admin

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
)

const (
	defaultInviteSkip  = 0
	defaultInviteLimit = 20
)

// FetchInvites returns all invites, active by default. If the withValid param provided and it's true,
// the method returns all invites including expired and invalid.
func (ar *Router) FetchInvites() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		withValid, err := ar.parseWithArchivedParam(r)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAPIJsonParseError, err)
			return
		}

		skip, limit, err := ar.parseSkipAndLimit(r, defaultInviteSkip, defaultInviteLimit, 0)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAPIJsonParseError, err)
			return
		}

		invites, total, err := ar.server.Storages().Invite.GetAll(r.Context(), withValid, skip, limit)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorStorageRequestError, err)
			return
		}

		searchResponse := struct {
			Invites []model.Invite `json:"invites"`
			Total   int            `json:"total"`
		}{
			Invites: invites,
			Total:   total,
		}

		ar.ServeJSON(w, locale, http.StatusOK, searchResponse)
	}
}

func (ar *Router) AddInvite() http.HandlerFunc {
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
		if ar.mustParseJSON(w, r, &d) != nil {
			return
		}

		token := tokenFromContext(r.Context())

		if d.Email != "" && !model.EmailRegexp.MatchString(d.Email) {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAPIRequestBodyEmailInvalid)
			return
		}

		u := model.User{
			ID:    model.NewUserID.String(),
			Email: d.Email,
		}
		aud := []string{}
		if len(d.AppID) > 0 {
			aud = append(aud, d.AppID)
		}
		fields := model.UserFieldsetMap[model.UserFieldsetInviteToken]
		inviteToken, err := ar.server.Services().Token.NewToken(model.TokenTypeInvite, u, aud, fields, d.Data)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAPIRequestBodyEmailInvalid, err)
			return
		}

		inviteTokenString, err := ar.server.Services().Token.SignToken(inviteToken)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorTokenInviteCreateError, err)
			return
		}

		err = ar.server.Storages().Invite.Save(d.Email, inviteTokenString, d.Role, d.AppID, "", inviteToken.ExpiresAt())
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorStorageInviteSaveError, err)
			return
		}

		invite, err := ar.server.Storages().Invite.GetByEmail(d.Email)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorStorageInviteFindEmailError, err)
			return
		}

		ar.ServeJSON(w, locale, http.StatusOK, invite)
	}
}

// GetInviteByID returns specific invite by its id.
func (ar *Router) GetInviteByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		id := mux.Vars(r)["id"]

		invite, err := ar.server.Storages().Invite.GetByID(r.Context(), id)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorStorageInviteFindIDError, err)
			return
		}

		ar.ServeJSON(w, locale, http.StatusOK, invite)
	}
}

// ArchiveInviteByID sets the 'valid' field of the model.Invite to false.
func (ar *Router) ArchiveInviteByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")
		id := mux.Vars(r)["id"]

		inv, err := ar.server.Storages().Invite.GetByID(r.Context(), id)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorStorageRequestError, err)
			return
		}
		inv.Archived = true
		err = ar.server.Storages().Invite.Update(r.Context(), inv)
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorStorageRequestError, err)
			return
		}
		ar.ServeJSON(w, locale, http.StatusOK, nil)
	}
}
