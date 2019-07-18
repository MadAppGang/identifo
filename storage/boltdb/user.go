package boltdb

import (
	"encoding/json"

	"github.com/madappgang/identifo/model"
)

// User is a user data structure for BoltDB storage.
type User struct {
	userData
}

// User data implementation.
type userData struct {
	ID              string                 `json:"id,omitempty"`
	Username        string                 `json:"username,omitempty"`
	Email           string                 `json:"email,omitempty"`
	Phone           string                 `json:"phone,omitempty"`
	Pswd            string                 `json:"pswd,omitempty"`
	Profile         map[string]interface{} `json:"profile,omitempty"`
	Active          bool                   `json:"active,omitempty"`
	TFAInfo         model.TFAInfo          `json:"tfa_info"`
	NumOfLogins     int                    `json:"num_of_logins,omitempty"`
	LatestLoginTime int64                  `json:"latest_login_time,omitempty"`
}

// Marshal serializes data to byte array.
func (u User) Marshal() ([]byte, error) {
	return json.Marshal(u.userData)
}

// Sanitize removes all sensitive data.
func (u User) Sanitize() model.User {
	u.userData.Pswd = ""
	u.userData.Active = false
	u.userData.TFAInfo.Secret = ""
	return u
}

// ID implements model.User interface.
func (u User) ID() string { return u.userData.ID }

// Username implements model.User interface.
func (u User) Username() string { return u.userData.Username }

// SetUsername implements model.User interface.
func (u User) SetUsername(username string) { u.userData.Username = username }

// Email implements model.User interface.
func (u User) Email() string { return u.userData.Email }

// SetEmail implements model.User interface.
func (u User) SetEmail(email string) { u.userData.Email = email }

// TFAInfo implements model.User interface.
func (u User) TFAInfo() model.TFAInfo { return u.userData.TFAInfo }

// SetTFAInfo implements model.User interface.
func (u User) SetTFAInfo(tfaInfo model.TFAInfo) { u.userData.TFAInfo = tfaInfo }

// PasswordHash implements model.User interface.
func (u User) PasswordHash() string { return u.userData.Pswd }

// Profile implements model.User interface.
func (u User) Profile() map[string]interface{} { return u.userData.Profile }

// Active implements model.User interface.
func (u User) Active() bool { return u.userData.Active }

// UserFromJSON deserializes user data from JSON.
func UserFromJSON(d []byte) (User, error) {
	user := userData{}
	if err := json.Unmarshal(d, &user); err != nil {
		return User{}, err
	}
	return User{userData: user}, nil
}
