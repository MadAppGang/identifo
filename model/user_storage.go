package model

import (
	"context"
	"regexp"
)

// ErrUserNotFound is when user not found.
// var ErrUserNotFound = errors.New("user not found")

var (
	// EmailRegexp is a regexp which all valid emails must match.
	EmailRegexp = regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)
	// PhoneRegexp is a regexp which all valid phone numbers must match.
	PhoneRegexp = regexp.MustCompile(`^[\+][0-9]{9,15}$`)
	// UsernameRegexp is a regexp which all valid usernames must match.
	UsernameRegexp = regexp.MustCompile(`(\w|\d){3,}`)
)

// UserStorage is an abstract user storage.
type UserStorage interface {
	Storage
	ImportableStorage

	// Get user with key parameters for a user.
	UserByID(ctx context.Context, id string) (User, error)
	UserBySecondaryID(ctx context.Context, idt AuthIdentityType, id string) (User, error)
	GetUserByFederatedID(ctx context.Context, idType UserFederatedType, userIdentityTypeOther, externalID string) (User, error)

	// Get user data, we can filter the fields we need to handle from data, as it is a large structure.
	UserData(ctx context.Context, userID string, fields ...UserDataField) (UserData, error)
}

type UserMutableStorage interface {
	// User mutation
	AddUser(ctx context.Context, user User) (User, error)
	UpdateUser(ctx context.Context, user User, fields ...string) (User, error)
	UpdateUserData(ctx context.Context, userID string, data UserData, fields ...UserDataField) (UserData, error)
	DeleteUser(ctx context.Context, userID string) error
}

// UserAdminStorage is a storage to manage users from admin panel and management api.
type UserAdminStorage interface {
	Storage

	FindUsers(ctx context.Context, search string, skip, limit int) ([]User, int, error)
	DeleteUser(ctx context.Context, id string) error
}

// UserDeviceStorage is a storage which keep all user device information.
type UserDeviceStorage interface {
	Storage
	ImportableStorage

	AddDevice(ctx context.Context, userID string, device UserDevice) (UserDevice, error)
	DetachDeviceWithToken(ctx context.Context, userID, token string) error
	AllUserDevices(ctx context.Context, userID string) ([]UserDevice, error)
}
