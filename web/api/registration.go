package api

import (
	"fmt"
	"net/http"

	l "github.com/madappgang/identifo/v2/localization"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/web/authorization"
	"github.com/madappgang/identifo/v2/web/middleware"
)

type registrationData struct {
	Username  string   `json:"username"`
	Phone     string   `json:"phone"`
	Email     string   `json:"email"`
	FullName  string   `json:"full_name"`
	Password  string   `json:"password"`
	Scopes    []string `json:"scopes"`
	Anonymous bool     `json:"anonymous"`
	Invite    string   `json:"invite"`
}

func (rd *registrationData) validate() error {
	emailLen := len(rd.Email)
	phoneLen := len(rd.Phone)
	usernameLen := len(rd.Username)
	pswdLen := len(rd.Password)

	if emailLen > 0 {
		if !model.EmailRegexp.MatchString(rd.Email) {
			return fmt.Errorf("invalid email")
		}
	}
	if phoneLen > 0 {
		if !model.PhoneRegexp.MatchString(rd.Phone) {
			return fmt.Errorf("invalid phone")
		}
	}

	if usernameLen > 0 {
		if usernameLen < 6 || usernameLen > 130 {
			return fmt.Errorf("incorrect username length %d, expected a number between 6 and 130", usernameLen)
		}
	}

	if emailLen == 0 && phoneLen == 0 && usernameLen == 0 {
		return fmt.Errorf("username, phone or/and email are quired for registration")
	}

	if pswdLen < 6 || pswdLen > 50 {
		return fmt.Errorf("incorrect password length %d, expected a number between 6 and 130", pswdLen)
	}
	return nil
}

/*
 * Password rules:
 * at least 6 letters
 * at least 1 upper case
 */

// RegisterWithPassword registers new user with password.
func (ar *Router) RegisterWithPassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := r.Header.Get("Accept-Language")

		app := middleware.AppFromContext(r.Context())
		if len(app.ID) == 0 {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIAPPNoAPPInContext)
			return
		}

		if app.RegistrationForbidden {
			ar.Error(w, locale, http.StatusForbidden, l.ErrorAPIAPPRegistrationForbidden)
			return
		}

		// Check if it makes sense to create new user.
		azi := authorization.AuthzInfo{
			App:         app,
			UserRole:    app.NewUserDefaultRole,
			ResourceURI: r.RequestURI,
			Method:      r.Method,
		}
		if err := ar.Authorizer.Authorize(azi); err != nil {
			ar.Error(w, locale, http.StatusUnauthorized, l.APIAccessDenied)
			return
		}

		// Parse registration data.
		rd := registrationData{}
		if ar.MustParseJSON(w, r, &rd) != nil {
			return
		}

		if rd.Anonymous && !app.AnonymousRegistrationAllowed {
			ar.Error(w, locale, http.StatusForbidden, l.ErrorAPILoginAnonymousForbidden)
			return
		}

		if err := rd.validate(); err != nil {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequestBodyInvalidError, err)
			return
		}

		// Validate password.
		if err := model.StrongPswd(rd.Password); err != nil {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIRequestPasswordWeak, err)
			return
		}

		// merge scopes for user
		scopes := model.MergeScopes(app.Scopes, app.NewUserDefaultScopes, rd.Scopes)

		// Create new user.
		um := model.User{
			Username: rd.Username,
			Email:    rd.Email,
			Phone:    rd.Phone,
			FullName: rd.FullName,
			Scopes:   model.SliceIntersect(app.Scopes, scopes), // add scopes to user which are limited by app only
		}

		userRole := app.NewUserDefaultRole

		if rd.Invite != "" {
			parsedInviteToken, err := ar.server.Services().Token.Parse(rd.Invite)
			if err != nil {
				ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIInviteUnableToInvalidateError, err)
				return
			}

			email, ok := parsedInviteToken.Payload()["email"].(string)
			if !ok || email != rd.Email {
				ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIInviteEmailMismatch)
				return
			}

			role, ok := parsedInviteToken.Payload()["role"].(string)
			if !ok {
				ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIInviteRoleMissing)
				return
			}
			userRole = role
		}

		user, err := ar.server.Storages().User.AddUserWithPassword(um, rd.Password, userRole, rd.Anonymous)
		if err == model.ErrorUserExists {
			ar.Error(w, locale, http.StatusBadRequest, l.ErrorAPIUsernamePhoneEmailTaken)
			return
		}
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorStorageUserCreateError, err)
			return
		}
		// if err = ar.server.Services().Email.SendTemplateEmail(
		// 	model.EmailTemplateTypeResetPassword,
		// 	app.GetCustomEmailTemplatePath(),
		// 	"Reset Password",
		// 	user.Email,
		// 	model.EmailData{
		// 		User: user,
		// 		Data: resetEmailData,
		// 	},
		// ); err != nil {
		// 	ar.Error(
		// 		w,
		// 		err,
		// 		http.StatusInternalServerError,
		// 		"Email sending error: "+err.Error(),
		// 	)
		// 	return

		// Do login flow.
		authResult, resultScopes, err := ar.loginFlow(app, user, rd.Scopes, nil)
		if err != nil {
			ar.Error(w, locale, http.StatusInternalServerError, l.ErrorAPILoginError, err)
			return
		}

		ar.journal(JournalOperationRegistration,
			user.ID, app.ID, r.UserAgent(), user.AccessRole, resultScopes.Scopes())

		ar.ServeJSON(w, locale, http.StatusOK, authResult)
	}
}
