package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/gorilla/mux"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/web/middleware"
)

const (
	defaultInviteSkip  = 0
	defaultInviteLimit = 20
)

// RequestInviteLink requests invite link. Invite link will be returned in response even if email or access_role is not specified.
func (ar *Router) RequestInviteLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d := struct {
			Email string `json:"email"`
			Role  string `json:"access_role"`
		}{}
		if err := ar.MustParseJSON(w, r, &d); err != nil {
			ar.Error(w, ErrorAPIRequestBodyInvalid, http.StatusBadRequest, err.Error(), "RequestInviteLink.MustParseJSON")
			return
		}
		if d.Email != "" && !model.EmailRegexp.MatchString(d.Email) {
			ar.Error(w, ErrorAPIRequestBodyInvalid, http.StatusBadRequest, "", "RequestInviteLink.emailRegexp_MatchString")
			return
		}

		_, err := ar.inviteStorage.GetByEmail(d.Email)
		if err != nil && !errors.Is(err, model.ErrorNotFound) {
			ar.Error(w, ErrorAPIInviteUnableToGet, http.StatusInternalServerError, err.Error(), "RequestInviteLink.inviteStorage_GetByEmail")
			return
		}

		if err := ar.inviteStorage.InvalidateAllByEmail(d.Email); err != nil {
			ar.Error(w, ErrorAPIInviteUnableToInvalidate, http.StatusInternalServerError, err.Error(), "RequestInviteLink.inviteStorage_InvalidateAllByEmail")
			return
		}

		inviteToken, err := ar.tokenService.NewInviteToken(d.Email, d.Role)
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
		scopes := strings.Replace(fmt.Sprintf("%q", app.Scopes), " ", ",", -1)
		query := url.PathEscape(fmt.Sprintf("appId=%s&scopes=%s&token=%s", app.ID, scopes, inviteTokenString))

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
			token := tokenFromContext(r.Context())

			err = ar.inviteStorage.Save(d.Email, inviteTokenString, d.Role, app.ID, token.UserID(), inviteToken.ExpiresAt())
			if err != nil {
				ar.Error(w, ErrorAPIInviteUnableToSave, http.StatusInternalServerError, err.Error(), "RequestInviteLink.inviteStorage_Save")
				return
			}

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

// GetAllInvites returns all invites, active by default. If the withValid param provided and it's true,
// the method returns all invites including expired and invalid.
func (ar *Router) GetAllInvites() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		withValid, err := ar.parseWithValidParam(r)
		if err != nil {
			ar.Error(w, ErrorAPIRequestBodyParamsInvalid, http.StatusBadRequest, err.Error(), "GetAllInvites.parseWithValidParam")
			return
		}

		skip, limit, err := ar.parseSkipAndLimit(r, defaultInviteSkip, defaultInviteLimit, 0)
		if err != nil {
			ar.Error(w, ErrorAPIRequestBodyParamsInvalid, http.StatusBadRequest, err.Error(), "GetAllInvites.parseSkipAndLimit")
			return
		}

		invites, total, err := ar.inviteStorage.GetAll(withValid, skip, limit)
		if err != nil {
			ar.Error(w, ErrorAPIRequestBodyParamsInvalid, http.StatusBadRequest, err.Error(), "GetAllInvites.parseSkipAndLimit")
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
			ar.Error(w, ErrorAPIInviteNotFound, http.StatusInternalServerError, err.Error(), "GetInviteByID.GetByID")
			return
		}

		ar.ServeJSON(w, http.StatusOK, invite)
	}
}

// InvalidateInviteByID sets the 'valid' field of the model.Invite to false.
func (ar *Router) InvalidateInviteByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		if err := ar.inviteStorage.InvalidateByID(id); err != nil {
			ar.Error(w, ErrorAPIInviteUnableToInvalidate, http.StatusInternalServerError, err.Error(), "InvalidateInviteByID.InvalidateByID")
			return
		}

		ar.ServeJSON(w, http.StatusOK, nil)
	}
}
