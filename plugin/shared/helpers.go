package shared

import (
	"errors"
	"regexp"
)

// ErrUserNotFound is when user not found.
var ErrUserNotFound = errors.New("User not found. ")
var ErrorInternalError = errors.New("Internal error")
var (
	// EmailRegexp is a regexp which all valid emails must match.
	EmailRegexp = regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)
	// PhoneRegexp is a regexp which all valid phone numbers must match.
	PhoneRegexp = regexp.MustCompile(`^[\+][0-9]{9,15}$`)
)

type Plugins struct {
	UserStorage UserStorage
}
