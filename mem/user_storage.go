package mem

import (
	"github.com/madappgang/identifo/model"
	randomdata "github.com/pallinder/go-randomdata"
)

//NewUserStorage creates and inits memory based user storage
//use it only for test purposes and in DI
//all data is wiped on exit
func NewUserStorage() model.UserStorage {
	return &UserStorage{}
}

//UserStorage implements user storage in memory
type UserStorage struct {
}

//UserByID returns random generated user
func (us *UserStorage) UserByID(id string) (model.User, error) {
	return randUser(), nil
}

//UserBySocialID returns random generated user
func (us *UserStorage) UserBySocialID(id string) (model.User, error) {
	return randUser(), nil
}

//CheckIfUserExistByName always returns true
func (us *UserStorage) CheckIfUserExistByName(name string) bool {
	return true
}

//AttachDeviceToken does nothing here.
func (us *UserStorage) AttachDeviceToken(id, token string) error {
	return nil
}

// DetachDeviceToken does nothing here.
func (us *UserStorage) DetachDeviceToken(token string) error {
	return nil
}

//RequestScopes mem always returns requested scope
func (us *UserStorage) RequestScopes(userID string, scopes []string) ([]string, error) {
	return scopes, nil
}

//Scopes returns supported scopes, could be static data of database
func (us *UserStorage) Scopes() []string {
	//we allow all scopes for embedded database, you could implement your own logic in external service
	return []string{"offline", "user"}
}

//UserByNamePassword returns random generated user
func (us *UserStorage) UserByNamePassword(name, password string) (model.User, error) {
	return randUser(), nil
}

//AddUserByNameAndPassword returns random generated user
func (us *UserStorage) AddUserByNameAndPassword(name, password string, profile map[string]interface{}) (model.User, error) {
	return randUser(), nil
}

//UserByFederatedID returns randomly generated user.
func (us *UserStorage) UserByFederatedID(provider model.FederatedIdentityProvider, id string) (model.User, error) {
	return randUser(), nil
}

// AddUserWithFederatedID returns randomly generated user.
func (us *UserStorage) AddUserWithFederatedID(provider model.FederatedIdentityProvider, id string) (model.User, error) {
	return randUser(), nil
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

//data implementation
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

//model.User interface implementation
func (u *user) ID() string                      { return u.userData.ID }
func (u *user) Name() string                    { return u.userData.Name }
func (u *user) PasswordHash() string            { return u.userData.Pswd }
func (u *user) Profile() map[string]interface{} { return u.userData.Profile }
func (u *user) Active() bool                    { return u.userData.Active }
