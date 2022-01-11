package authorization

import (
	"fmt"

	"github.com/casbin/casbin"
	"github.com/madappgang/identifo/v2/model"
	strCasbin "github.com/qiangmzsx/string-adapter"
)

// NewAuthorizer creates a new Authorizer.
func NewAuthorizer() *Authorizer {
	return &Authorizer{
		internalAuthorizers: make(map[string]*casbin.Enforcer),
	}
}

// Authorizer is an entity that authorizes users to an app.
type Authorizer struct {
	internalAuthorizers map[string]*casbin.Enforcer
}

const anonymousRole = "anonymous"

// AuthzInfo holds all the data to perform authorization.
type AuthzInfo struct {
	App         model.AppData
	UserRole    string
	ResourceURI string
	Method      string
}

// Authorize performs authorization.
func (az *Authorizer) Authorize(azi AuthzInfo) error {
	if az == nil {
		return nil
	}
	switch azi.App.AuthzWay {
	case model.NoAuthz, "":
		return nil
	case model.RolesWhitelist:
		return az.authorizeWhitelist(azi)
	case model.RolesBlacklist:
		return az.authorizeBlacklist(azi)
	case model.Internal:
		return az.authorizeInternal(azi)
	case model.External:
		return model.ErrorNotImplemented
	}
	return nil
}

func (az *Authorizer) authorizeWhitelist(azi AuthzInfo) error {
	whitelist := azi.App.RolesWhitelist
	if whitelist == nil {
		err := fmt.Errorf("Access denied")
		return err
	}

	role := azi.UserRole
	if role == "" {
		role = anonymousRole
	}

	if accessGranted := contains(whitelist, role); !accessGranted {
		err := fmt.Errorf("Access denied")
		return err
	}
	return nil
}

func (az *Authorizer) authorizeBlacklist(azi AuthzInfo) error {
	blacklist := azi.App.RolesBlacklist
	if blacklist == nil {
		return nil
	}

	role := azi.UserRole
	if role == "" {
		role = anonymousRole
	}

	if accessDenied := contains(blacklist, role); accessDenied {
		err := fmt.Errorf("Access denied")
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
func (az *Authorizer) authorizeInternal(azi AuthzInfo) error {
	authorizer, err := az.initInternalAuthorizer(azi.App)
	if err != nil {
		err = fmt.Errorf("Cannot init internal authorizer for app %s: %s", azi.App.ID, err)
		return err
	}

	sub := azi.UserRole
	if sub == "" {
		sub = anonymousRole
	}
	obj := azi.ResourceURI
	act := azi.Method

	if accessGranted := authorizer.Enforce(sub, obj, act); !accessGranted {
		err := fmt.Errorf("Access denied")
		return err
	}
	return nil
}

func (az *Authorizer) initInternalAuthorizer(app model.AppData) (*casbin.Enforcer, error) {
	authorizer, ok := az.internalAuthorizers[app.ID]
	if ok {
		return authorizer, nil
	}

	// If authorizer has not been initialized already, try initializing it.
	modelStr, policyStr := app.AuthzModel, app.AuthzPolicy
	if len(modelStr) == 0 || len(policyStr) == 0 {
		return nil, fmt.Errorf("Either authz model or policy is empty for app %s, or both", app.ID)
	}
	strAdapter := strCasbin.NewAdapter(policyStr)

	authorizer = casbin.NewEnforcer(casbin.NewModel(modelStr), strAdapter)
	az.internalAuthorizers[app.ID] = authorizer
	authorizer.EnableLog(true)

	return authorizer, nil
}
