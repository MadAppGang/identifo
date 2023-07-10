package storage

import (
	"context"
	"net/url"
	"os"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
)

// just compile-time check interface compliance.
// please don't use it in runtime.
var _uc model.UserController = NewUserStorageController(nil, nil, nil, nil, nil, nil, nil, nil, model.ServerSettings{})

// UserStorageController performs common user operations using a set of storages.
// For example when user logins, we find the user, match the password, and log the login attempt and save it to log storage.
// All user business logic is implemented in controller and all storage things are dedicated to storage implementations.
type UserStorageController struct {
	s model.SecurityServerSettings
	h *url.URL
	// lgs model.LogStorage //TODO: implement logs
	uidfs  []string // unique ID fields
	imudfs []string // immutable ID fields
	u      model.UserStorage
	ums    model.UserMutableStorage
	ua     model.UserAdminStorage
	as     model.AppStorage
	ts     model.TokenService
	es     model.EmailService
	ss     model.SMSService
	uas    model.UserAuthStorage
	LP     *l.Printer // localized string

	// cache
	ffs  []model.FirstFactorInternalStrategy
	idts []model.AuthIdentityType // identity types in use
}

// NewUserStorageController composes new storage controller, no validation and connections happens here.
// The functions expects all the storages already initialized and connected.
func NewUserStorageController(
	u model.UserStorage,
	ums model.UserMutableStorage,
	ua model.UserAdminStorage,
	as model.AppStorage,
	uas model.UserAuthStorage,
	ts model.TokenService,
	es model.EmailService,
	ss model.SMSService,
	s model.ServerSettings,
) *UserStorageController {
	// env variable can rewrite host option
	hostName := os.Getenv("IDENTIFO_HOST_NAME")
	if len(hostName) == 0 {
		hostName = s.General.Host
	}

	host, err := url.ParseRequestURI(hostName)
	if err != nil {
		host, _ = url.ParseRequestURI("http://localhost")
	}

	return &UserStorageController{
		h:      host,
		u:      u,
		ua:     ua,
		ums:    ums,
		as:     as,
		uas:    uas,
		ts:     ts,
		es:     es,
		ss:     ss,
		s:      s.SecuritySettings,
		uidfs:  authTypeToField(s.General.UniqueIDFields),
		imudfs: authTypeToField(s.General.ImmutableIDFields),
	}
}

// InvalidateCache let's clean app cache
func (c *UserStorageController) InvalidateCache() {
	c.ffs = nil
	c.idts = nil
}

// UserByID returns User with ID with basic fields
func (c *UserStorageController) UserByID(ctx context.Context, userID string) (model.User, error) {
	// TODO: use scopes to identifies the fieldset to return
	return c.UserByIDWithFields(ctx, userID, model.UserFieldsetBasic)
}

// UserByIDWithFields returns user with specific fieldset
func (c *UserStorageController) UserByIDWithFields(ctx context.Context, userID string, fields model.UserFieldset) (model.User, error) {
	user, err := c.u.UserByID(ctx, userID)
	if err != nil {
		return model.User{}, err
	}
	// strip user fields
	result := model.CopyFields(user, fields.Fields())
	return result, nil
}

// UserBySecondaryID returns user with basic fieldset by his secondary ID.
func (c *UserStorageController) UserBySecondaryID(ctx context.Context, idt model.AuthIdentityType, id string) (model.User, error) {
	// TODO: use scopes to identifies the fieldset to return
	return c.UserBySecondaryIDWithFields(ctx, idt, id, model.UserFieldsetBasic)
}

// UserBySecondaryIDWithFields returns user with specific fieldset by his secondary ID.
func (c *UserStorageController) UserBySecondaryIDWithFields(ctx context.Context, idt model.AuthIdentityType, id string, fields model.UserFieldset) (model.User, error) {
	user, err := c.u.UserBySecondaryID(ctx, idt, id)
	if err != nil {
		return model.User{}, err
	}
	// strip user fields
	result := model.CopyFields(user, fields.Fields())
	return result, nil
}

// UserByFederatedID returns user profile by federated ID.
func (c *UserStorageController) UserByFederatedID(ctx context.Context, idt model.UserFederatedType, idOther, id string) (model.User, error) {
	// TODO: use scopes to identifies the fieldset to return
	user, err := c.u.GetUserByFederatedID(ctx, idt, idOther, id)
	if err != nil {
		return model.User{}, err
	}
	// strip user fields
	result := model.CopyFields(user, model.UserFieldsetBasic.Fields())
	return result, nil
}

// GetUsers returns users with basic fields in it.
func (c *UserStorageController) GetUsers(ctx context.Context, filter string, skip, limit int) ([]model.User, int, error) {
	users, total, err := c.ua.FindUsers(ctx, filter, skip, limit)
	if err != nil {
		return nil, 0, err
	}

	// strip user fields
	for i, user := range users {
		users[i] = model.CopyFields(user, model.UserFieldsetBasic.Fields())
	}
	return users, total, nil
}

// UserByAuthStrategy extract the authentication id and tries to find the user with this ID.
func (c *UserStorageController) UserByAuthStrategy(ctx context.Context, auth model.AuthStrategy, userIDValue string) (model.User, error) {
	//we need to check each strategy type
	//if it's local:
	//- get identity type id - userbyid
	//- other? user userbysecoindaryID
	//- anonymous - get user by id
	//- fim - find federated identity
	//- second factor - has no identity

	if auth.Type() == model.AuthStrategyAnonymous {
		return c.UserByID(ctx, userIDValue)
	}
	if auth.Type() == model.AuthStrategyFirstFactorInternal {
		a, ok := auth.(model.FirstFactorInternalStrategy)
		if !ok {
			return model.User{}, l.LocalizedError{ErrID: l.ErrorUserNotFound}
		}
		if a.Identity == model.AuthIdentityTypeID {
			return c.UserByID(ctx, userIDValue)
		} else {
			return c.UserBySecondaryID(ctx, a.Identity, userIDValue)
		}
	}
	if auth.Type() == model.AuthStrategyFirstFactorFIM {
		a, ok := auth.(model.FirstFactorFIMStrategy)
		if !ok {
			return model.User{}, l.LocalizedError{ErrID: l.ErrorUserNotFound}
		}
		return c.UserByFederatedID(ctx, model.UserFederatedType(a.FIMType), "", userIDValue)
	}
	if auth.Type() == model.AuthStrategyFirstFactorEnterprise {
		return model.User{}, l.LocalizedError{ErrID: l.ErrorLoginTypeNotSupported}
	}

	if auth.Type() == model.AuthStrategySecondFactor {
		return model.User{}, l.LocalizedError{ErrID: l.ErrorLoginTypeNotSupported}
	}

	return model.User{}, l.LocalizedError{ErrID: l.ErrorLoginTypeNotSupported}
}
