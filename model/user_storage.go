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
	IDByName(name string) (string, error)
	AttachDeviceToken(id, token string) error
	DetachDeviceToken(token string) error
	UserByUsername(username string) (User, error)
	UserExists(name string) bool
	UserByFederatedID(provider string, id string) (User, error)
	AddUserWithFederatedID(user User, provider string, id, role string) (User, error)
	UpdateUser(userID string, newUser User) (User, error)
	ResetPassword(id, password string) error
	CheckPassword(id, password string) error
	DeleteUser(id string) error
	FetchUsers(search string, skip, limit int) ([]User, int, error)

	RequestScopes(userID string, scopes []string) ([]string, error)
	Scopes() []string
	ImportJSON(data []byte) error
	UpdateLoginMetadata(userID string)
	Close()
}

// User is an abstract representation of the user in auth layer.
// Everything can be User, we do not depend on any particular implementation.
type User struct {
	ID              string   `json:"id,omitempty" bson:"_id,omitempty"`
	Username        string   `json:"username" bson:"username,omitempty"`
	Email           string   `json:"email" bson:"email,omitempty"`
	Phone           string   `json:"phone" bson:"phone,omitempty"`
	Pswd            string   `json:"pswd,omitempty" bson:"pswd,omitempty"`
	Active          bool     `json:"active,omitempty" bson:"active,omitempty"`
	TFAInfo         TFAInfo  `json:"tfa_info,omitempty" bson:"tfa_info,omitempty"`
	NumOfLogins     int      `json:"num_of_logins,omitempty" bson:"num_of_logins,omitempty"`
	LatestLoginTime int64    `json:"latest_login_time,omitempty" bson:"latest_login_time,omitempty"`
	AccessRole      string   `json:"access_role,omitempty" bson:"access_role,omitempty"`
	Anonymous       bool     `json:"anonymous,omitempty" bson:"anonymous,omitempty"`
	FederatedIDs    []string `json:"federated_ids,omitempty" bson:"federated_i_ds,omitempty"`
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
	IsEnabled     bool      `json:"is_enabled,omitempty" bson:"is_enabled,omitempty"`
	HOTPCounter   int       `json:"hotp_counter,omitempty" bson:"hotp_counter,omitempty"`
	HOTPExpiredAt time.Time `json:"hotp_expired_at,omitempty" bson:"hotp_expired_at,omitempty"`
	Secret        string    `json:"secret,omitempty" bson:"secret,omitempty"`
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
