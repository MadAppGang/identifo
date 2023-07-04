package management

import (
	"errors"
	"net/http"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
)

func (ar *Router) getInviteToken(w http.ResponseWriter, r *http.Request) {
	locale := r.Header.Get("Accept-Language")

	var d InvitationTokenRequest
	if err := ar.MustParseJSON(w, r, &d); err != nil {
		ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyInvalidError, err)
		return
	}

	if d.Email != "" && !model.EmailRegexp.MatchString(d.Email) {
		ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyEmailInvalid)
		return
	}

	_, err := ar.server.Storages().Invite.GetByEmail(d.Email)
	if err != nil && !errors.Is(err, model.ErrorNotFound) {
		ar.Error(w, locale, http.StatusInternalServerError, l.ErrorStorageInviteFindEmailError, err)
		return
	}

	if err := ar.server.Storages().Invite.ArchiveAllByEmail(d.Email); err != nil {
		ar.Error(w, locale, http.StatusInternalServerError, l.ErrorStorageInviteArchiveEmailError, err)
		return
	}

	inviteToken, err := ar.server.Services().Token.NewInviteToken(d.Email, d.Role, d.ApplicationID, d.Data)
	if err != nil {
		ar.Error(w, locale, http.StatusInternalServerError, l.ErrorTokenInviteCreateError, err)
		return
	}

	inviteTokenString, err := ar.server.Services().Token.String(inviteToken)
	if err != nil {
		ar.Error(w, locale, http.StatusInternalServerError, l.ErrorTokenInviteCreateError, err)
		return
	}

	result := map[string]string{"result": "ok", "token": inviteTokenString}
	ar.ServeJSON(w, locale, http.StatusOK, result)
}
