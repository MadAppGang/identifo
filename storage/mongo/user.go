package mongo

import (
	"encoding/json"

	"github.com/madappgang/identifo/model"
	"gopkg.in/mgo.v2/bson"
)

// User is a data structure for MongoDB storage.
type User struct {
	userData
}

// User data implementation.
type userData struct {
	ID              bson.ObjectId          `bson:"_id,omitempty" json:"id,omitempty"`
	Username        string                 `bson:"username,omitempty" json:"username,omitempty"`
	Email           string                 `bson:"email,omitempty" json:"email,omitempty"`
	Phone           string                 `bson:"phone,omitempty" json:"phone,omitempty"`
	Pswd            string                 `bson:"pswd,omitempty" json:"pswd,omitempty"`
	Profile         map[string]interface{} `bson:"profile,omitempty" json:"profile,omitempty"`
	Active          bool                   `bson:"active,omitempty" json:"active,omitempty"`
	FederatedIDs    []string               `bson:"federated_ids,omitempty" json:"federated_ids,omitempty"`
	NumOfLogins     int                    `bson:"num_of_logins" json:"num_of_logins,omitempty"`
	LatestLoginTime int64                  `bson:"latest_login_time,omitempty" json:"latest_login_time,omitempty"`
}

// Sanitize removes sensitive data.
func (u *User) Sanitize() model.User {
	u.userData.Pswd = ""
	return u
}

// UserFromJSON deserializes user from JSON.
func UserFromJSON(d []byte) (*User, error) {
	user := userData{}
	if err := json.Unmarshal(d, &user); err != nil {
		return &User{}, err
	}
	return &User{userData: user}, nil
}

// ID implements model.User interface.
func (u *User) ID() string { return u.userData.ID.Hex() }

// Username implements model.User interface.
func (u *User) Username() string { return u.userData.Username }

// SetUsername implements model.User interface.
func (u *User) SetUsername(username string) { u.userData.Username = username }

// Email implements model.Email interface.
func (u *User) Email() string { return u.userData.Email }

// SetEmail implements model.Email interface.
func (u *User) SetEmail(email string) { u.userData.Email = email }

// PasswordHash implements model.User interface.
func (u *User) PasswordHash() string { return u.userData.Pswd }

// Profile implements model.User interface.
func (u *User) Profile() map[string]interface{} { return u.userData.Profile }

// Active implements model.User interface.
func (u *User) Active() bool { return u.userData.Active }
