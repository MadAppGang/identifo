package mem

import (
	"github.com/madappgang/identifo/model"
	"github.com/pallinder/go-randomdata"
)

// NewUserStorage creates and inits in-memory user storage.
// Use it only for test purposes and in CI, all data is wiped on exit.
func NewUserStorage() (model.UserStorage, error) {
	return &UserStorage{}, nil
}

// UserStorage is an in-memory user storage .
type UserStorage struct{}

// UserByID returns randomly generated user.
func (us *UserStorage) UserByID(id string) (model.User, error) {
	return randUser(), nil
}

// UserByEmail returns randomly generated user.
func (us *UserStorage) UserByEmail(email string) (model.User, error) {
	return randUser(), nil
}

// UserBySocialID returns randomly generated user.
func (us *UserStorage) UserBySocialID(id string) (model.User, error) {
	return randUser(), nil
}

// UserByPhone returns randomly generated user.
func (us *UserStorage) UserByPhone(phone string) (model.User, error) {
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

// UserByUsername returns randomly generated user.
func (us *UserStorage) UserByUsername(username string) (model.User, error) {
	return randUser(), nil
}

// AddUserWithPassword creates new user and saves it in the database.
func (us *UserStorage) AddUserWithPassword(user model.User, password, role string, isAnonymous bool) (model.User, error) {
	return randUser(), nil
}

// UserByFederatedID returns randomly generated user.
func (us *UserStorage) UserByFederatedID(provider model.FederatedIdentityProvider, id string) (model.User, error) {
	return randUser(), nil
}

// AddUserWithFederatedID returns randomly generated user.
func (us *UserStorage) AddUserWithFederatedID(provider model.FederatedIdentityProvider, id, role string) (model.User, error) {
	return randUser(), nil
}

// UpdateUser returns what it receives.
func (us *UserStorage) UpdateUser(userID string, newUser model.User) (model.User, error) {
	return newUser, nil
}

// ResetPassword does nothing here.
func (us *UserStorage) ResetPassword(id, password string) error {
	return nil
}

// CheckPassword does nothig here.
func (us *UserStorage) CheckPassword(id, password string) error {
	return nil
}

// IDByName returns random id.
func (us *UserStorage) IDByName(name string) (string, error) {
	return randomdata.StringNumber(2, "-"), nil
}

// DeleteUser does nothing here.
func (us *UserStorage) DeleteUser(id string) error {
	return nil
}

// UpdateLoginMetadata does nothing here.
func (us *UserStorage) UpdateLoginMetadata(userID string) {}

// FetchUsers returns randomly generated user enclosed in slice.
func (us *UserStorage) FetchUsers(filterString string, skip, limit int) ([]model.User, int, error) {
	return []model.User{randUser()}, 1, nil
}

// ImportJSON imports data from JSON.
func (us *UserStorage) ImportJSON(data []byte) error {
	return nil
}

// Close does nothing here.
func (us *UserStorage) Close() {}

func randUser() model.User {
	return model.User{
		ID:       randomdata.StringNumber(2, "-"),
		Username: randomdata.SillyName(),
		Email:    randomdata.Email(),
		Pswd:     randomdata.StringNumber(2, "-"),
		Active:   randomdata.Boolean(),
	}
}
