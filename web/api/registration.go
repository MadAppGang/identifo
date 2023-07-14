package api

import (
	"net/http"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/tools/xmaps"
	"github.com/madappgang/identifo/v2/web/middleware"
)

// RegisterWithPassword registers new user with password.
// check if app allows registration, if not - fails
// create the user from data and password
// if invite present
// get inviter
func (ar *Router) RegisterWithPassword() http.HandlerFunc {
	type registrationData struct {
		User      UserRequestData `json:"user"`
		Scopes    []string        `json:"scopes"`
		Anonymous bool            `json:"anonymous"`
		Invite    string          `json:"invite"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPPNoAPPInContext)
			return
		}

		if app.RegistrationForbidden {
			ar.LocalizedError(w, locale, http.StatusForbidden, l.ErrorAPPRegistrationForbidden)
			return
		}

		// Parse registration data.
		rd := registrationData{}
		if ar.MustParseJSON(w, r, &rd) != nil {
			return
		}
		rd.User.ID = model.NewUserID.String() // this is a new user, rewrite any data the client could send

		if rd.Anonymous && !app.AnonymousRegistrationAllowed {
			ar.LocalizedError(w, locale, http.StatusForbidden, l.ErrorAPILoginAnonymousForbidden)
			return
		}

		u := model.User{}
		err := xmaps.CopyDstFields(rd, &u) // copy non nil fields from d to u
		if err != nil {
			ar.LocalizedError(w, locale, http.StatusInternalServerError, l.ErrorAPIRequestBodyInvalidError, err)
			return
		}

		// create anonymous user
		if rd.Anonymous == true {
			u.Anonymous = true
			u, err = ar.server.Storages().UMC.CreateUser(r.Context(), u)
			if err != nil {
				ar.HTTPError(w, err, http.StatusInternalServerError)
				return
			}
		} else { // create regular user
			pswd := ""
			if rd.User.Password != nil {
				pswd = *rd.User.Password
			}
			u, err = ar.server.Storages().UMC.CreateUserWithPassword(r.Context(), u, pswd)
			if err != nil {
				ar.HTTPError(w, err, http.StatusInternalServerError)
				return
			}
		}

		// add user to invited tenants
		if len(rd.Invite) > 0 {

			parsedInviteToken, err := ar.server.Services().Token.Parse(rd.Invite)
			if err != nil {
				ar.LocalizedError(w, locale, http.StatusBadRequest, l.ErrorAPIInviteUnableToInvalidateError, err)
				return
			}

			_, err = ar.server.Storages().UMC.AddUserToTenantWithInvitationToken(r.Context(), u, parsedInviteToken)
			if err != nil {
				ar.HTTPError(w, err, http.StatusInternalServerError)
				return
			}
		}

		// return JWT tokens
		auth, err := ar.server.Storages().UC.GetJWTTokens(r.Context(), app, u, rd.Scopes)
		if err != nil {
			ar.HTTPError(w, err, http.StatusInternalServerError)
			return
		}
		ar.ServeJSON(w, locale, http.StatusOK, auth)
	}
}
