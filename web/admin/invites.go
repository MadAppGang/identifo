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

		invites, total, err := ar.server.Storages().Invite.GetAll(withValid, skip, limit)
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
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		d := struct {
			AppID string                 `json:"app_id"`
			Email string                 `json:"email"`
			Role  string                 `json:"access_role"`
			Data  map[string]interface{} `json:"data"`
		}{}

		if err := ar.mustParseJSON(w, r, &d); err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAPIJsonParseError, err)
			return
		}
		if d.Email != "" && !model.EmailRegexp.MatchString(d.Email) {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAPIRequestBodyEmailInvalid)
			return
		}

		inviteToken, err := ar.server.Services().Token.NewInviteToken(d.Email, d.Role, "identifo", d.Data)
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

		invite, err := ar.server.Storages().Invite.GetByID(id)
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

		if err := ar.server.Storages().Invite.ArchiveByID(id); err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorStorageRequestError, err)
			return
		}

		ar.ServeJSON(w, locale, http.StatusOK, nil)
	}
}
