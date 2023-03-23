package model

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ErrUserNotFound is when user not found.
var ErrUserNotFound = errors.New("User not found. ")

var (
	// EmailRegexp is a regexp which all valid emails must match.
	EmailRegexp = regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)
	// PhoneRegexp is a regexp which all valid phone numbers must match.
	PhoneRegexp = regexp.MustCompile(`^[\+][0-9]{9,15}$`)
)

// UserStorage is an abstract user storage.
type UserStorage interface {
	UserByPhone(phone string) (User, error)
	AddUserWithPassword(user User, password, role string, isAnonymous bool) (User, error)
	UserByID(id string) (User, error)
	UserByEmail(email string) (User, error)
	UserByUsername(username string) (User, error)
	UserByFederatedID(provider string, id string) (User, error)
	AddUserWithFederatedID(user User, provider string, id, role string) (User, error)
	UpdateUser(userID string, newUser User) (User, error)
	ResetPassword(id, password string) error
	CheckPassword(id, password string) error
	DeleteUser(id string) error
	FetchUsers(search string, skip, limit int) ([]User, int, error)
	UpdateLoginMetadata(userID string)

	// push device tokens
	AttachDeviceToken(userID, token string) error
	DetachDeviceToken(token string) error
	AllDeviceTokens(userID string) ([]string, error)

	// import data
	ImportJSON(data []byte, clearOldData bool) error

	Close()
}

// we have three sets of scopes
// allowed - the list of scopes allowed for app
// def - default list of scopes for the new user
// requested - requested list of scopes for new user
func MergeScopes(allowed, def, requested []string) []string {
	// if we are not requesting any scope, just use default set
	if len(requested) == 0 {
		return def
	}

	// if allowed list is empty we accepting anythings
	if len(allowed) == 0 {
		return requested
	}

	// if we requested something, ensure we can use only allowed scopes for the app
	return SliceIntersect(allowed, requested)
}

// merge two sets of scopes for requested scope
// we have three sets of scopes
// user - the list of scopes user has
// requested - requested list of scopes for key
func ReqestedScopesApply(user, requested []string) []string {
	// if we are requesting nothing, we are gettings nothing
	if len(requested) == 0 {
		return []string{}
	}

	// if we requested something, ensure we can use only allowed scopes for the app
	return SliceIntersect(user, requested)
}
