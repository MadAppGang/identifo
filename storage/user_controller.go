package storage

import (
	"context"

	"github.com/madappgang/identifo/v2/model"
)

// UserStorageController performs common user operations using a set of storages.
// For example when user logins, we find the user, match the password, and log the login attempt and save it to log storage.
// All user business logic is implemented in controller and all storage things are dedicated to storage implementations.
type UserStorageController struct {
	u model.UserStorage
}

// NewStorageController composes new storage controller, no validation and connections happens here.
// The functions expects all the storages already initialized and connected.
func NewStorageController(u model.UserStorage) *UserStorageController {
	return &UserStorageController{
		u: u,
	}
}

func (c *UserStorageController) UserByID(ctx context.Context, userID string) (model.User, error) {
	user, err := c.u.UserByID(ctx, userID)
	if err != nil {
		return model.User{}, err
	}

	// strip user fields to empty ones
	result := model.CopyFields(user, UserFieldsetMap[UserFieldsetBasic])
	return result, nil
}

// user, err := ar.server.Storages().User.UserByID(userID)
// user = user.Sanitized()
// FetchUsers
// AddUserWithPassword(um, rd.Password, rd.AccessRole, false)
// GenerateNewResetTokenUser
// DeleteUser(userID);

func (c *UserStorageController) TODO(ctx context.Context) {
}
