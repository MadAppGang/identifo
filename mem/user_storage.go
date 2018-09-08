package mem

import (
	"github.com/madappgang/identifo/model"
	randomdata "github.com/pallinder/go-randomdata"
)

//NewUserStorage creates and inits memory based user storage
//use it only for test purposes and in DI
//all data is wiped on exit
func NewUserStorage() model.UserStorage {
	//us := UserStorage{}
	// return &us
	return nil
}

//UserStorage implements user storage in memory
type UserStorage struct {
}

//UserByID returns randon generated user
func (us *UserStorage) UserByID(id string) (model.User, error) {
	return randUser(), nil
}

//UserBySocialID returns randon generated user
func (us *UserStorage) UserBySocialID(id string) (model.User, error) {
	return randUser(), nil
}

//AttachDeviceToken do nothing here
func (us *UserStorage) AttachDeviceToken(id, token string) error {
	return nil
}

func randUser() user {
	return user{
		userData: userData{
			id:      randomdata.StringNumber(2, "-"),
			name:    randomdata.SillyName(),
			pswd:    randomdata.StringNumber(2, "-"),
			profile: nil,
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
