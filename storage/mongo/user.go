package mongo

import (
	"encoding/json"

	"github.com/madappgang/identifo/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is a data structure for MongoDB storage.
type User struct {
	userData
}

// User data implementation.
type userData struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username        string             `bson:"username,omitempty" json:"username,omitempty"`
	Email           string             `bson:"email,omitempty" json:"email,omitempty"`
	Phone           string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Pswd            string             `bson:"pswd,omitempty" json:"pswd,omitempty"`
	Active          bool               `bson:"active,omitempty" json:"active,omitempty"`
	TFAInfo         model.TFAInfo      `bson:"tfa_info" json:"tfa_info"`
	FederatedIDs    []string           `bson:"federated_ids,omitempty" json:"federated_ids,omitempty"`
	NumOfLogins     int                `bson:"num_of_logins" json:"num_of_logins,omitempty"`
	LatestLoginTime int64              `bson:"latest_login_time,omitempty" json:"latest_login_time,omitempty"`
	AccessRole      string             `bson:"access_role,omitempty" json:"access_role,omitempty"`
	Anonymous       bool               `json:"anonymous,omitempty"`
}

// Sanitize removes sensitive data.
func (u *User) Sanitize() {
	u.userData.Pswd = ""
	u.userData.TFAInfo.Secret = ""
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

// Email implements model.User interface.
func (u *User) Email() string { return u.userData.Email }

// SetEmail implements model.User interface.
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

// IsAnonymous implements model.User interface.
func (u *User) IsAnonymous() bool { return u.userData.Anonymous }

// Deanonimize implements model.User interface.
func (u *User) Deanonimize() { u.userData.Anonymous = false }
