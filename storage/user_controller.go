package storage

import (
	"context"

	"github.com/madappgang/identifo/v2/model"
)

// just compile-time check interface compliance.
// please don't use it in runtime.
var _uc model.UserController = NewUserStorageController(nil)

// UserStorageController performs common user operations using a set of storages.
// For example when user logins, we find the user, match the password, and log the login attempt and save it to log storage.
// All user business logic is implemented in controller and all storage things are dedicated to storage implementations.
type UserStorageController struct {
	u  model.UserStorage
	ua model.UserAdminStorage
}

// NewUserStorageController composes new storage controller, no validation and connections happens here.
// The functions expects all the storages already initialized and connected.
func NewUserStorageController(u model.UserStorage) *UserStorageController {
	return &UserStorageController{
		u: u,
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

func (c *UserStorageController) CreateUserWithPassword(ctx context.Context, u model.User, password string) (model.User, error) {
	// Password policy: set of parameters for password +++
	// https: // www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go
	// Password hasher: salt, pepper, algorithm
}

// AddUserWithPassword(um, rd.Password, rd.AccessRole, false)
// GenerateNewResetTokenUser
// DeleteUser(userID);

func (c *UserStorageController) TODO(ctx context.Context) {
}

// TODO: implement this logic in User Controller.
// and make it properly using OWASP recommendations.

// // PasswordHash creates hash with salt for password.
// func PasswordHash(pwd string) string {
// 	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
// 	return string(hash)
// }

// // RandomPassword creates random password
// func RandomPassword(length int) string {
// 	rand.Seed(time.Now().UnixNano())
// 	return randSeq(length)
// }

// var rndPassLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890?!@#")

// func randSeq(n int) string {
// 	b := make([]rune, n)
// 	for i := range b {
// 		b[i] = rndPassLetters[rand.Intn(len(rndPassLetters))]
// 	}
// 	return string(b)
// }
