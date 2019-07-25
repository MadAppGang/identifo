package api

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin"
	"github.com/madappgang/identifo/model"
	strCasbin "github.com/qiangmzsx/string-adapter"
)

const anonymousRole = "anonymous"

type authzInfo struct {
	app         model.AppData
	userRole    string
	resourceURI string
	method      string
}

// authorize checks if user has an access to the requested resource.
// If error happens, writes it to ResponseWriter.
// Also, writes an error on failed authorization.
func (ar *Router) authorize(w http.ResponseWriter, azi authzInfo) error {
	switch azi.app.AuthzWay() {
	case model.NoAuthz, "":
		return nil
	case model.RolesWhitelist:
		return ar.authorizeWhitelist(w, azi)
	case model.RolesBlacklist:
		return ar.authorizeBlacklist(w, azi)
	case model.Internal:
		return ar.authorizeInternal(w, azi)
	case model.External:
		return model.ErrorNotImplemented
	}
	return nil
}

func (ar *Router) authorizeWhitelist(w http.ResponseWriter, azi authzInfo) error {
	whitelist := azi.app.RolesWhitelist()
	if whitelist == nil {
		err := fmt.Errorf("Access denied")
		ar.Error(w, ErrorAPIAppAccessDenied, http.StatusForbidden, err.Error(), "Authorizer.AccessDenied")
		return err
	}

	role := azi.userRole
	if role == "" {
		role = anonymousRole
	}

	if accessGranted := contains(whitelist, role); !accessGranted {
		err := fmt.Errorf("Access denied")
		ar.Error(w, ErrorAPIAppAccessDenied, http.StatusForbidden, err.Error(), "Authorizer.AccessDenied")
		return err
	}
	return nil
}

func (ar *Router) authorizeBlacklist(w http.ResponseWriter, azi authzInfo) error {
	blacklist := azi.app.RolesBlacklist()
	if blacklist == nil {
		return nil
	}

	role := azi.userRole
	if role == "" {
		role = anonymousRole
	}

	if accessDenied := contains(blacklist, role); accessDenied {
		err := fmt.Errorf("Access denied")
		ar.Error(w, ErrorAPIAppAccessDenied, http.StatusForbidden, err.Error(), "Authorizer.AccessDenied")
		return err
	}
	return nil
}

// authorizeInternal performs authorization based on the model and policy rules,
// which are stored in the application entity as strings.

/* Example model:
`[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act`
*/

/* Example policy:
`p, admin, /auth/login, POST
p, manager, /auth/phone_login, POST
p, manager, /auth/federated, POST
p, anonymous, /auth/register, POST`
*/
func (ar *Router) authorizeInternal(w http.ResponseWriter, azi authzInfo) error {
	authorizer, err := ar.initInternalAuthorizer(azi.app)
	if err != nil {
		err = fmt.Errorf("Cannot init internal authorizer for app %s: %s", azi.app.ID(), err)
		ar.Error(w, ErrorAPIAppCannotInitAuthorizer, http.StatusInternalServerError, err.Error(), "Authorizer.authorizeInternal")
		return err
	}

	sub := azi.userRole
	if sub == "" {
		sub = anonymousRole
	}
	obj := azi.resourceURI
	act := azi.method

	if accessGranted := authorizer.Enforce(sub, obj, act); !accessGranted {
		err := fmt.Errorf("Access denied")
		ar.Error(w, ErrorAPIAppAccessDenied, http.StatusForbidden, err.Error(), "Authorizer.AccessDenied")
		return err
	}
	return nil
}

func (ar *Router) initInternalAuthorizer(app model.AppData) (*casbin.Enforcer, error) {
	authorizer, ok := ar.Authorizers[app.ID()]
	if ok {
		return authorizer, nil
	}

	// If authorizer has not been initialized already, try initializing it.
	modelStr, policyStr := app.AuthzModel(), app.AuthzPolicy()
	if len(modelStr) == 0 || len(policyStr) == 0 {
		return nil, fmt.Errorf("Either authz model or policy is empty for app %s, or both", app.ID())
	}
	strAdapter := strCasbin.NewAdapter(policyStr)

	authorizer = casbin.NewEnforcer(casbin.NewModel(modelStr), strAdapter)
	ar.Authorizers[app.ID()] = authorizer
	authorizer.EnableLog(true)

	return authorizer, nil
}
