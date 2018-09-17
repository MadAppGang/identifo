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

func randUser() user {
	profile := map[string]interface{}{
		"name":    randomdata.SillyName(),
		"id":      randomdata.StringNumber(2, "-"),
		"address": randomdata.Address(),
	}
	return user{
		userData: userData{
			id:      randomdata.StringNumber(2, "-"),
			name:    randomdata.SillyName(),
			pswd:    randomdata.StringNumber(2, "-"),
			profile: profile,
		},
	}
}

//data implementation
type userData struct {
	id      string
	name    string
	pswd    string
	profile map[string]interface{}
}

type user struct {
	userData
}

//model.User interface implementation
func (u user) ID() string                      { return u.id }
func (u user) Name() string                    { return u.name }
func (u user) PasswordHash() string            { return u.pswd }
func (u user) Profile() map[string]interface{} { return u.profile }
