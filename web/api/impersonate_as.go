package api

import (
	"context"
	"errors"
	"net/http"

	l "github.com/madappgang/identifo/v2/localization"
	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/middleware"
)

func (ar *Router) ImpersonateAs() http.HandlerFunc {
	type impersonateData struct {
		UserID string `json:"user_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		locale := r.Header.Get("Accept-Language")

		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIAPPNoAPPInContext)
			return
		}

		userID := tokenFromContext(ctx).UserID()
		adminUser, err := ar.server.Storages().User.UserByID(userID)
		if err != nil {
			ar.Error(w, locale, http.StatusUnauthorized, l.ErrorStorageFindUserIDError, userID, err)
			return
		}

		d := impersonateData{}
		if ar.MustParseJSON(w, r, &d) != nil {
			return
		}

		ar.logger.Info("admin for impersonation",
			"adminUserId", adminUser.ID,
			"adminScopes", adminUser.Scopes,
			"adminAccessRole", adminUser.AccessRole,
			"impersonatedId", d.UserID)

		user, err := ar.server.Storages().User.UserByID(d.UserID)
		if err != nil {
			ar.Error(w, locale, http.StatusUnauthorized, l.ErrorStorageFindUserIDError, d.UserID, err)
			return
		}

		ok, err := ar.checkImpersonationPermissions(ctx, app, adminUser, user)
		if err != nil {
			ar.logger.Error("cannot check impersonation",
				logging.FieldError, err)
			ar.Error(w, locale, http.StatusForbidden, l.ErrorAPIImpersonationForbidden)
			return
		}

		if !ok {
			ar.Error(w, locale, http.StatusForbidden, l.ErrorAPIImpersonationForbidden)
			return
		}

		ap := map[string]any{
			"impersonated_by": adminUser.ID,
		}

		authResult, resultScopes, err := ar.loginFlow(app, user, nil, ap)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorAPILoginError, err)
			return
		}

		// do not allow refresh for impersonated user
		authResult.RefreshToken = ""

		ar.journal(JournalOperationImpersonatedAs,
			userID, app.ID, r.UserAgent(), user.AccessRole, resultScopes.Scopes())

		ar.ServeJSON(w, locale, http.StatusOK, authResult)
	}
}

func (ar *Router) checkImpersonationPermissions(
	ctx context.Context,
	app model.AppData,
	adminUser, impUser model.User,
) (bool, error) {
	imps := ar.server.Services().Impersonation
	if imps == nil {
		return false, errors.New("impersonation service is not set")
	}

	ok, err := imps.CanImpersonate(ctx, app.ID, adminUser, impUser)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, nil
	}

	return true, nil
}
