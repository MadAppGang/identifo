package mem

import (
	"github.com/madappgang/identifo/model"
	randomdata "github.com/pallinder/go-randomdata"
)

//NewUserStorage creates and inits memory based user storage
//use it only for test purposes and in DI
//all data is wiped on exit
func NewUserStorage() model.UserStorage {
	us := UserStorage{}
	return &us
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

//AttachDeviceToken do nothing here
func (us *UserStorage) AttachDeviceToken(id, token string) error {
	return nil
}

//RequestScopes mem always returns requested scope
func (us *UserStorage) RequestScopes(userID string, scopes []string) ([]string, error) {
	return scopes, nil
}

//UserByNamePassword returns random generated user
func (us *UserStorage) UserByNamePassword(name, password string) (model.User, error) {
	return randUser(), nil
}

//AddUserByNameAndPassword returns random generated user
func (us *UserStorage) AddUserByNameAndPassword(name, password string, profile map[string]interface{}) (model.User, error) {
	return randUser(), nil
}

//UserByName return random generated user
func (us *UserStorage) UserByName(name string) (model.User, error) {
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
