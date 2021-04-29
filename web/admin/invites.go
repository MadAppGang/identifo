package admin

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo/model"
)

const (
	defaultInviteSkip  = 0
	defaultInviteLimit = 20
)

// FetchInvites returns all invites, active by default. If the withValid param provided and it's true,
// the method returns all invites including expired and invalid.
func (ar *Router) FetchInvites() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		withValid, err := ar.parseWithArchivedParam(r)
		if err != nil {
			ar.Error(w, ErrorAPIRequestBodyParamsInvalid, http.StatusBadRequest, "")
			return
		}

		skip, limit, err := ar.parseSkipAndLimit(r, defaultInviteSkip, defaultInviteLimit, 0)
		if err != nil {
			ar.Error(w, ErrorAPIRequestBodyParamsInvalid, http.StatusBadRequest, "")
			return
		}

		invites, total, err := ar.inviteStorage.GetAll(withValid, skip, limit)
		if err != nil {
			ar.Error(w, ErrorAPIRequestBodyParamsInvalid, http.StatusBadRequest, "")
			return
		}

		searchResponse := struct {
			Invites []model.Invite `json:"invites"`
			Total   int            `json:"total"`
		}{
			Invites: invites,
			Total:   total,
		}

		ar.ServeJSON(w, http.StatusOK, searchResponse)
	}
}

// GetInviteByID returns specific invite by its id.
func (ar *Router) GetInviteByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		invite, err := ar.inviteStorage.GetByID(id)
		if err != nil {
			ar.Error(w, ErrorAPIInviteNotFound, http.StatusInternalServerError, "")
			return
		}

		ar.ServeJSON(w, http.StatusOK, invite)
	}
}

// ArchiveInviteByID sets the 'valid' field of the model.Invite to false.
func (ar *Router) ArchiveInviteByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		if err := ar.inviteStorage.ArchiveByID(id); err != nil {
			ar.Error(w, ErrorAPIInviteUnableToInvalidate, http.StatusInternalServerError, "")
			return
		}

		ar.ServeJSON(w, http.StatusOK, nil)
	}
}
