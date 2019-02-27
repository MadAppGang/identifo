package mem

import (
	"github.com/madappgang/identifo/model"
	randomdata "github.com/pallinder/go-randomdata"
)

// NewUserStorage creates and inits in-memory user storage.
// Use it only for test purposes and in CI, all data is wiped on exit.
func NewUserStorage() model.UserStorage {
	return &UserStorage{}
}

// UserStorage implements user storage in memory.
type UserStorage struct {
}

// UserByID returns randomly generated user.
func (us *UserStorage) UserByID(id string) (model.User, error) {
	return randUser(), nil
}

// UserBySocialID returns randomly generated user.
func (us *UserStorage) UserBySocialID(id string) (model.User, error) {
	return randUser(), nil
}

// UserExists always returns true.
func (us *UserStorage) UserExists(name string) bool {
	return true
}

// AttachDeviceToken does nothing here.
func (us *UserStorage) AttachDeviceToken(id, token string) error {
	return nil
}

// DetachDeviceToken does nothing here.
func (us *UserStorage) DetachDeviceToken(token string) error {
	return nil
}

// RequestScopes always returns requested scopes.
func (us *UserStorage) RequestScopes(userID string, scopes []string) ([]string, error) {
	return scopes, nil
}

// Scopes returns supported scopes, could be static data of database.
func (us *UserStorage) Scopes() []string {
	// we allow all scopes for embedded database, you could implement your own logic in external service
	return []string{"offline", "user"}
}

// UserByNamePassword returns randomly generated user.
func (us *UserStorage) UserByNamePassword(name, password string) (model.User, error) {
	return randUser(), nil
}

// AddUserByNameAndPassword returns randomly generated user.
func (us *UserStorage) AddUserByNameAndPassword(name, password string, profile map[string]interface{}) (model.User, error) {
	return randUser(), nil
}

// UserByFederatedID returns randomly generated user.
func (us *UserStorage) UserByFederatedID(provider model.FederatedIdentityProvider, id string) (model.User, error) {
	return randUser(), nil
}

// AddUserWithFederatedID returns randomly generated user.
func (us *UserStorage) AddUserWithFederatedID(provider model.FederatedIdentityProvider, id string) (model.User, error) {
	return randUser(), nil
}

// ResetPassword does nothing here.
func (us *UserStorage) ResetPassword(id, password string) error {
	return nil
}

// IDByName returns random id.
func (us *UserStorage) IDByName(name string) (string, error) {
	return randomdata.StringNumber(2, "-"), nil
}

// FetchUsers returns randomly generated user enclosed in slice.
func (us *UserStorage) FetchUsers(filterString string, skip, limit int) ([]model.User, error) {
	return []model.User{randUser()}, nil
}

// ImportJSON imports data from JSON.
func (us *UserStorage) ImportJSON(data []byte) error {
	return nil
}

func randUser() *user {
	profile := map[string]interface{}{
		"name":    randomdata.SillyName(),
		"id":      randomdata.StringNumber(2, "-"),
		"address": randomdata.Address(),
	}
	return &user{
		userData: userData{
			ID:      randomdata.StringNumber(2, "-"),
			Name:    randomdata.SillyName(),
			Pswd:    randomdata.StringNumber(2, "-"),
			Profile: profile,
			Active:  randomdata.Boolean(),
		},
	}
}

// User data implementation.
type userData struct {
	ID      string                 `json:"id,omitempty"`
	Name    string                 `json:"name,omitempty"`
	Pswd    string                 `json:"pswd,omitempty"`
	Profile map[string]interface{} `json:"profile,omitempty"`
	Active  bool                   `json:"active,omitempty"`
}

type user struct {
	userData
}

func (u *user) Sanitize() {
	u.userData.Pswd = ""
	u.userData.Active = false
}

// ID implements model.User interface.
func (u *user) ID() string { return u.userData.ID }

// Name implements model.User interface.
func (u *user) Name() string { return u.userData.Name }

// PasswordHash implements model.User interface.
func (u *user) PasswordHash() string { return u.userData.Pswd }

// Profile implements model.User interface.
func (u *user) Profile() map[string]interface{} { return u.userData.Profile }

// Active implements model.User interface.
func (u *user) Active() bool { return u.userData.Active }
