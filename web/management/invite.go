package management

import (
	"errors"
	"net/http"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"golang.org/x/exp/maps"
)

// InviteRequest is a request for invite.
type InvitationTokenRequest struct {
	Email       string         `json:"email"`
	AppID       string         `json:"app_id"`
	Roles       map[string]any `json:"roles"`
	CallbackURL string         `json:"callback"`
	Data        map[string]any `json:"data"`
}

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
	if err != nil && !errors.Is(err, l.ErrorNotFound) {
		ar.Error(w, locale, http.StatusInternalServerError, l.ErrorStorageInviteFindEmailError, err)
		return
	}

	if err := ar.server.Storages().Invite.ArchiveAllByEmail(d.Email); err != nil {
		ar.Error(w, locale, http.StatusInternalServerError, l.ErrorStorageInviteArchiveEmailError, err)
		return
	}

	u := model.User{
		ID:    model.NewUserID.String(), // token sub is new user ID
		Email: d.Email,
	}
	aud := []string{}
	if len(d.AppID) > 0 {
		aud = append(aud, d.AppID)
	}
	fields := model.UserFieldsetMap[model.UserFieldsetInviteToken]
	maps.Copy(d.Data, d.Roles)
	if len(d.CallbackURL) > 0 {
		d.Data["callback"] = d.CallbackURL
	}
	inviteToken, err := ar.server.Services().Token.NewToken(model.TokenTypeInvite, u, aud, fields, d.Data)
	if err != nil {
		ar.Error(w, locale, http.StatusInternalServerError, l.ErrorTokenInviteCreateError, err)
		return
	}

	inviteTokenString, err := ar.server.Services().Token.SignToken(inviteToken)
	if err != nil {
		ar.Error(w, locale, http.StatusInternalServerError, l.ErrorTokenInviteCreateError, err)
		return
	}

	result := map[string]string{"result": "ok", "token": inviteTokenString}
	ar.ServeJSON(w, locale, http.StatusOK, result)
}
