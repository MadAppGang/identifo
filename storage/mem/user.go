package mem

import (
	"github.com/madappgang/identifo/model"
	"github.com/pallinder/go-randomdata"
)

// User data implementation.
type userData struct {
	ID         string                 `json:"id,omitempty"`
	Username   string                 `json:"username,omitempty"`
	Email      string                 `json:"email,omitempty"`
	Pswd       string                 `json:"pswd,omitempty"`
	Profile    map[string]interface{} `json:"profile,omitempty"`
	Active     bool                   `json:"active,omitempty"`
	TFAInfo    model.TFAInfo          `json:"tfa_info"`
	AccessRole string                 `json:"access_role,omitempty"`
}

type user struct {
	userData
}

func (u *user) Sanitize() {
	u.userData.Pswd = ""
	u.userData.Active = false
	u.userData.TFAInfo.Secret = ""
}

// ID implements model.User interface.
func (u *user) ID() string { return u.userData.ID }

// Username implements model.User interface.
func (u *user) Username() string { return u.userData.Username }

// SetUsername implements model.User interface.
func (u *user) SetUsername(username string) { u.userData.Username = username }

// Email implements model.User interface.
func (u *user) Email() string { return u.userData.Email }

// SetEmail implements model.User interface.
func (u *user) SetEmail(email string) { u.userData.Email = email }

// TFAInfo implements model.User interface.
func (u *user) TFAInfo() model.TFAInfo { return u.userData.TFAInfo }

// SetTFAInfo implements model.User interface.
func (u *user) SetTFAInfo(tfaInfo model.TFAInfo) { u.userData.TFAInfo = tfaInfo }

// PasswordHash implements model.User interface.
func (u *user) PasswordHash() string { return u.userData.Pswd }

// Profile implements model.User interface.
func (u *user) Profile() map[string]interface{} { return u.userData.Profile }

// Active implements model.User interface.
func (u *user) Active() bool { return u.userData.Active }

// AccessRole implements model.User interface.
func (u *user) AccessRole() string { return u.userData.AccessRole }

func randUser() *user {
	profile := map[string]interface{}{
		"username": randomdata.SillyName(),
		"id":       randomdata.StringNumber(2, "-"),
		"address":  randomdata.Address(),
	}
	return &user{
		userData: userData{
			ID:       randomdata.StringNumber(2, "-"),
			Username: randomdata.SillyName(),
			Email:    randomdata.Email(),
			Pswd:     randomdata.StringNumber(2, "-"),
			Profile:  profile,
			Active:   randomdata.Boolean(),
		},
	}
}
