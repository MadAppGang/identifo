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
	ImportJSON(data []byte) error

	Close()
}

// User is an abstract representation of the user in auth layer.
// Everything can be User, we do not depend on any particular implementation.
type User struct {
	ID              string   `json:"id" bson:"_id"`
	Username        string   `json:"username" bson:"username"`
	Email           string   `json:"email" bson:"email"`
	FullName        string   `json:"full_name" bson:"full_name"`
	Phone           string   `json:"phone" bson:"phone"`
	Pswd            string   `json:"pswd" bson:"pswd"`
	Active          bool     `json:"active" bson:"active"`
	TFAInfo         TFAInfo  `json:"tfa_info" bson:"tfa_info"`
	NumOfLogins     int      `json:"num_of_logins" bson:"num_of_logins"`
	LatestLoginTime int64    `json:"latest_login_time" bson:"latest_login_time"`
	AccessRole      string   `json:"access_role" bson:"access_role"`
	Anonymous       bool     `json:"anonymous" bson:"anonymous"`
	FederatedIDs    []string `json:"federated_ids" bson:"federated_i_ds"`
	Scopes          []string `json:"scopes" bson:"scopes"`
}

func maskLeft(s string, hideFraction int) string {
	rs := []rune(s)
	for i := 0; i < len(rs)-len(rs)/hideFraction; i++ {
		rs[i] = '*'
	}
	return string(rs)
}

// Sanitized returns data structure without sensitive information
func (u User) Sanitized() User {
	u.Pswd = ""
	u.TFAInfo.Secret = ""
	u.TFAInfo.HOTPCounter = 0
	u.TFAInfo.HOTPExpiredAt = time.Time{}
	return u
}

func (u *User) AddFederatedId(provider string, id string) string {
	fid := provider + ":" + id
	for _, ele := range u.FederatedIDs {
		if ele == fid {
			return fid
		}
	}
	u.FederatedIDs = append(u.FederatedIDs, fid)
	return fid
}

// SanitizedTFA returns data structure with masked sensitive data
func (u User) SanitizedTFA() User {
	u.Sanitized()
	if len(u.Email) > 0 {
		emailParts := strings.Split(u.Email, "@")
		u.Email = maskLeft(emailParts[0], 2) + "@" + maskLeft(emailParts[1], 2)
	}

	if len(u.Phone) > 0 {
		u.Phone = maskLeft(u.Phone, 3)
	}
	return u
}

// Deanonimized returns model with all fields set for deanonimized user
func (u User) Deanonimized() User {
	u.Anonymous = false
	return u
}

// TFAInfo encapsulates two-factor authentication user info.
type TFAInfo struct {
	IsEnabled     bool      `json:"is_enabled" bson:"is_enabled"`
	HOTPCounter   int       `json:"hotp_counter" bson:"hotp_counter"`
	HOTPExpiredAt time.Time `json:"hotp_expired_at" bson:"hotp_expired_at"`
	Email         string    `json:"email" bson:"email"`
	Phone         string    `json:"phone" bson:"phone"`
	Secret        string    `json:"secret" bson:"secret"`
}

// UserFromJSON deserialize user data from JSON.
func UserFromJSON(d []byte) (User, error) {
	var user User
	if err := json.Unmarshal(d, &user); err != nil {
		log.Println("Error while unmarshal user:", err)
		return User{}, err
	}
	return user, nil
}

// PasswordHash creates hash with salt for password.
func PasswordHash(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash)
}

// RandomPassword creates random password
func RandomPassword(length int) string {
	rand.Seed(time.Now().UnixNano())
	return randSeq(length)
}

var rndPassLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890?!@#")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = rndPassLetters[rand.Intn(len(rndPassLetters))]
	}
	return string(b)
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
