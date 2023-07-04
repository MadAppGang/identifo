package storage

import (
	"context"
	"errors"
	"time"

	l "github.com/madappgang/identifo/v2/localization"

	"github.com/madappgang/identifo/v2/jwt"
	"github.com/madappgang/identifo/v2/model"
)

// just compile-time check interface compliance.
// please don't use it in runtime.
var (
	_uc  model.UserController         = NewUserStorageController(nil, model.SecurityServerSettings{})
	_umc model.UserMutationController = NewUserStorageController(nil, model.SecurityServerSettings{})
)

// UserStorageController performs common user operations using a set of storages.
// For example when user logins, we find the user, match the password, and log the login attempt and save it to log storage.
// All user business logic is implemented in controller and all storage things are dedicated to storage implementations.
type UserStorageController struct {
	s   model.SecurityServerSettings
	u   model.UserStorage
	ums model.UserMutableStorage
	ua  model.UserAdminStorage
	LP  *l.Printer // localized string
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

// ====================================
// Data mutation
// ====================================

// CreateUserWithPassword validates password policy, creates password hash and creates new user
// it also responsible to call pre-create and post-create callbacks
func (c *UserStorageController) CreateUserWithPassword(ctx context.Context, u model.User, password string) (model.User, error) {
	// TODO: implement check for isCompromised property
	isCompromised := false
	valid, vr := c.s.PasswordPolicy.Validate(password, isCompromised)
	if !valid {
		ee := []error{}
		for _, r := range vr {
			if !r.Valid {
				ee = append(ee, r.Error())
			}
		}
		// return all violated rules
		return model.User{}, errors.Join(ee...)
	}

	hash, err := jwt.PasswordHash(password, c.s.PasswordHash, []byte(c.s.PasswordHash.Pepper))
	if err != nil {
		return model.User{}, err
	}
	u.PasswordHash = hash

	// TODO: Call pre-create callbacks

	nu, err := c.ums.AddUser(ctx, u)
	if err != nil {
		return model.User{}, err
	}

	// TODO: Call post-create callbacks

	return nu, nil
}

func (c *UserStorageController) UpdateUserPassword(ctx context.Context, userID, password string) error {
	// TODO: implement check for isCompromised property
	isCompromised := false
	valid, vr := c.s.PasswordPolicy.Validate(password, isCompromised)
	if !valid {
		ee := []error{}
		for _, r := range vr {
			if !r.Valid {
				ee = append(ee, r.Error())
			}
		}
		// return all violated rules
		return errors.Join(ee...)
	}

	hash, err := jwt.PasswordHash(password, c.s.PasswordHash, []byte(c.s.PasswordHash.Pepper))
	if err != nil {
		return err
	}

	// TODO: Call pre-change password callback

	user := model.User{
		ID:                  userID,
		PasswordHash:        hash,
		UpdatedAt:           time.Now(),
		LastPasswordResetAt: time.Now(),
	}
	_, err = c.ums.UpdateUser(ctx, user, model.UserFieldsetPassword.UpdateFields()...)
	if err != nil {
		return err
	}

	// TODO: Call  post-change password callback

	return nil
}

func (c *UserStorageController) ChangeBlockStatus(ctx context.Context, userID, reason, whoName, whoID string, blocked bool) error {
	user := model.User{
		ID:        userID,
		UpdatedAt: time.Now(),
		Blocked:   blocked,
	}
	if blocked {
		user.BlockedDetails = &model.UserBlockedDetails{
			Reason:        reason,
			BlockedByName: whoName,
			BlockedById:   whoID,
			BlockedAt:     time.Now(),
		}
	}

	// TODO: Call  pre-block callback

	_, err := c.ums.UpdateUser(ctx, user, model.UserFieldsetBlockStatus.Fields()...)
	if err != nil {
		return err
	}

	// TODO: Call  post-block callback

	return nil
}

func (c *UserStorageController) UpdateUser(ctx context.Context, user model.User, fields []string) (model.User, error) {
	// TODO: Call  pre-update callback

	u, err := c.ums.UpdateUser(ctx, user, fields...)
	if err != nil {
		return model.User{}, err
	}

	// TODO: Call  post-update callback
	return u, nil
}

func (c *UserStorageController) DeleteUser(ctx context.Context, userID string) error {
	// TODO: Call  pre-update callback

	err := c.ums.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}

	// TODO: Call  post-update callback
	return nil
}
