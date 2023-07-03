package storage

import (
	"context"

	l "github.com/madappgang/identifo/v2/localization"

	"github.com/madappgang/identifo/v2/jwt"
	"github.com/madappgang/identifo/v2/localization"
	"github.com/madappgang/identifo/v2/model"
)

// just compile-time check interface compliance.
// please don't use it in runtime.
var _uc model.UserController = NewUserStorageController(nil, model.SecurityServerSettings{})

// UserStorageController performs common user operations using a set of storages.
// For example when user logins, we find the user, match the password, and log the login attempt and save it to log storage.
// All user business logic is implemented in controller and all storage things are dedicated to storage implementations.
type UserStorageController struct {
	s  model.SecurityServerSettings
	u  model.UserStorage
	ua model.UserAdminStorage
	LP *l.Printer // localized string
}

// NewUserStorageController composes new storage controller, no validation and connections happens here.
// The functions expects all the storages already initialized and connected.
func NewUserStorageController(u model.UserStorage, s model.SecurityServerSettings) *UserStorageController {
	return &UserStorageController{
		u: u,
		s: s,
	}
}

// UserByID returns User with ID with basic fields
func (c *UserStorageController) UserByID(ctx context.Context, userID string) (model.User, error) {
	return c.UserByIDWithFields(ctx, userID, UserFieldsetBasic)
}

// UserByIDWithFields returns user with specific fieldset
func (c *UserStorageController) UserByIDWithFields(ctx context.Context, userID string, fields UserFieldset) (model.User, error) {
	user, err := c.u.UserByID(ctx, userID)
	if err != nil {
		return model.User{}, err
	}
	// strip user fields
	result := model.CopyFields(user, UserFieldsetMap[fields])
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
		users[i] = model.CopyFields(user, UserFieldsetMap[UserFieldsetBasic])
	}
	return users, total, nil
}

// CreateUserWithPassword validates password policy, creates password hash and creates new user
// it also responsible to call pre-create and post-create callbacks
func (c *UserStorageController) CreateUserWithPassword(ctx context.Context, u model.User, password, locale string) (model.User, error) {
	// TODO: implement check for isCompromised property
	pr := c.LP.PrinterForLocale(locale)
	valid, vr := c.s.PasswordPolicy.Validate(password, false, pr)
	if !valid {
		// find firs violated rule
		es := ""
		for _, r := range vr {
			if !r.Valid {
				es = r.ValidationRule
				break
			}
		}
		return model.User{}, pr.E(localization.ErrorAPIRequestPasswordWeak, es)
	}

	hash, err := jwt.PasswordHash(password, c.s.PasswordHash, []byte(c.s.PasswordHash.Pepper))
	if err != nil {
		return model.User{}, err
	}
	u.PasswordHash = hash

	// TODO: Call pre-create callbacks

	nu, err := c.u.AddUser(ctx, u)
	if err != nil {
		return model.User{}, pr.EL(err)
	}

	// TODO: Call post-create callbacks

	return nu, nil
}

// GenerateNewResetTokenUser
// DeleteUser(userID);
