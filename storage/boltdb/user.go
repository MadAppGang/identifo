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
	ID              string        `json:"id,omitempty"`
	Username        string        `json:"username,omitempty"`
	Email           string        `json:"email,omitempty"`
	Phone           string        `json:"phone,omitempty"`
	Pswd            string        `json:"pswd,omitempty"`
	Active          bool          `json:"active,omitempty"`
	TFAInfo         model.TFAInfo `json:"tfa_info"`
	NumOfLogins     int           `json:"num_of_logins,omitempty"`
	LatestLoginTime int64         `json:"latest_login_time,omitempty"`
	AccessRole      string        `json:"access_role,omitempty"`
}

// Marshal serializes data to byte array.
func (u User) Marshal() ([]byte, error) {
	return json.Marshal(u.userData)
}

// Sanitize removes all sensitive data.
func (u *User) Sanitize() {
	u.userData.Pswd = ""
	u.userData.TFAInfo.Secret = ""
}

// ID implements model.User interface.
func (u *User) ID() string { return u.userData.ID }

// Username implements model.User interface.
func (u *User) Username() string { return u.userData.Username }

// SetUsername implements model.User interface.
func (u *User) SetUsername(username string) { u.userData.Username = username }

// Email implements model.User interface.
func (u *User) Email() string { return u.userData.Email }

// SetEmail implements model.Email interface.
func (u *User) SetEmail(email string) { u.userData.Email = email }

// Phone implements model.User interface.
func (u *User) Phone() string { return u.userData.Phone }

// TFAInfo implements model.User interface.
func (u *User) TFAInfo() model.TFAInfo { return u.userData.TFAInfo }

// SetTFAInfo implements model.User interface.
func (u *User) SetTFAInfo(tfaInfo model.TFAInfo) { u.userData.TFAInfo = tfaInfo }

// PasswordHash implements model.User interface.
func (u *User) PasswordHash() string { return u.userData.Pswd }

// Active implements model.User interface.
func (u *User) Active() bool { return u.userData.Active }

// AccessRole implements model.User interface.
func (u *User) AccessRole() string { return u.userData.AccessRole }

// UserFromJSON deserializes user data from JSON.
func UserFromJSON(d []byte) (*User, error) {
	user := userData{}
	if err := json.Unmarshal(d, &user); err != nil {
		return &User{}, err
	}
	return &User{userData: user}, nil
}
