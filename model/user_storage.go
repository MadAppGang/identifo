package model

import "errors"

// ErrUserNotFound is when user not found.
var ErrUserNotFound = errors.New("User not found. ")

// UserStorage is an abstract user storage.
type UserStorage interface {
	UserByPhone(phone string) (User, error)
	AddUserByPhone(phone, role string) (User, error)
	UserByID(id string) (User, error)
	UserByEmail(email string) (User, error)
	IDByName(name string) (string, error)
	AttachDeviceToken(id, token string) error
	DetachDeviceToken(token string) error
	UserByNamePassword(name, password string) (User, error)
	AddUserByNameAndPassword(name, password, role string) (User, error)
	UserExists(name string) bool
	UserByFederatedID(provider FederatedIdentityProvider, id string) (User, error)
	AddUserWithFederatedID(provider FederatedIdentityProvider, id, role string) (User, error)
	UpdateUser(userID string, newUser User) (User, error)
	ResetPassword(id, password string) error
	DeleteUser(id string) error
	FetchUsers(search string, skip, limit int) ([]User, int, error)
	NewUser() User

	RequestScopes(userID string, scopes []string) ([]string, error)
	Scopes() []string
	ImportJSON(data []byte) error
	UpdateLoginMetadata(userID string)
	Close()
}

// User is an abstract representation of the user in auth layer.
// Everything can be User, we do not depend on any particular implementation.
type User interface {
	ID() string
	Username() string
	SetUsername(string)
	Email() string
	SetEmail(string)
	Phone() string
	SetPhone(string)
	TFAInfo() TFAInfo
	SetTFAInfo(TFAInfo)
	PasswordHash() string
	Active() bool
	AccessRole() string
	Sanitize()
}

// TFAInfo encapsulates two-factor authentication user info.
type TFAInfo struct {
	IsEnabled bool   `bson:"is_enabled" json:"is_enabled"`
	Secret    string `bson:"secret" json:"-"`
}
