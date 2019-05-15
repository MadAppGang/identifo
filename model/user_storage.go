package model

import "errors"

// ErrUserNotFound is when user not found.
var ErrUserNotFound = errors.New("User not found. ")

// UserStorage is an abstract user storage.
type UserStorage interface {
	UserByID(id string) (User, error)
	UserByEmail(email string) (User, error)
	IDByName(name string) (string, error)
	AttachDeviceToken(id, token string) error
	DetachDeviceToken(token string) error
	UserByNamePassword(name, password string) (User, error)
	AddUserByNameAndPassword(name, password string, profile map[string]interface{}) (User, error)
	UserExists(name string) bool
	UserByFederatedID(provider FederatedIdentityProvider, id string) (User, error)
	AddUserWithFederatedID(provider FederatedIdentityProvider, id string) (User, error)
	UpdateUser(userID string, newUser User) (User, error)
	ResetPassword(id, password string) error
	DeleteUser(id string) error
	FetchUsers(search string, skip, limit int) ([]User, error)
	NewUser() User

	RequestScopes(userID string, scopes []string) ([]string, error)
	Scopes() []string
	ImportJSON(data []byte) error
}

// User is an abstract representation of the user in auth layer.
// Everything can be User, we do not depend on any particular implementation.
type User interface {
	ID() string
	Name() string
	SetName(string)
	Email() string
	SetEmail(string)
	PasswordHash() string
	Profile() map[string]interface{}
	Active() bool
	Sanitize() User
}
