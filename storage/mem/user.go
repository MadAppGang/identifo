package mem

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/madappgang/identifo/v2/model"
	"github.com/pallinder/go-randomdata"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// NewUserStorage creates and inits in-memory user storage.
// Use it only for test purposes and in CI, all data is wiped on exit.
func NewUserStorage() (model.UserStorage, error) {
	return &UserStorage{
		users:       []model.User{},
		userDevices: make(map[string]string),
	}, nil
}

// UserStorage is an in-memory user storage .
type UserStorage struct {
	users       []model.User
	userDevices map[string]string
}

// UserByID returns randomly generated user.
func (us *UserStorage) UserByID(id string) (model.User, error) {
	for _, u := range us.users {
		if strings.EqualFold(u.ID, id) {
			return u, nil
		}
	}
	return model.User{}, errors.New("not found")
}

// UserByEmail returns randomly generated user.
func (us *UserStorage) UserByEmail(email string) (model.User, error) {
	for _, u := range us.users {
		if strings.EqualFold(u.Email, email) {
			return u, nil
		}
	}
	return model.User{}, errors.New("not found")
}

// UserBySocialID returns randomly generated user.
func (us *UserStorage) UserBySocialID(id string) (model.User, error) {
	for _, u := range us.users {
		for _, fi := range u.FederatedIDs {
			if strings.EqualFold(fi, id) {
				return u, nil
			}
		}
	}
	return model.User{}, errors.New("not found")
}

// UserByPhone returns randomly generated user.
func (us *UserStorage) UserByPhone(phone string) (model.User, error) {
	for _, u := range us.users {
		if strings.EqualFold(u.Phone, phone) {
			return u, nil
		}
	}
	return model.User{}, errors.New("not found")
}

// AttachDeviceToken does nothing here.
func (us *UserStorage) AttachDeviceToken(id, token string) error {
	us.userDevices[id] = token
	return nil
}

// DetachDeviceToken does nothing here.
func (us *UserStorage) DetachDeviceToken(token string) error {
	delete(us.userDevices, token)
	return nil
}

// TODO: implement get all device tokens logic
func (us *UserStorage) AllDeviceTokens(userID string) ([]string, error) {
	devices := []string{}
	for uid, did := range us.userDevices {
		if strings.EqualFold(uid, userID) {
			devices = append(devices, did)
		}
	}
	return devices, nil
}

// UserByUsername returns randomly generated user.
func (us *UserStorage) UserByUsername(username string) (model.User, error) {
	for _, u := range us.users {
		if strings.EqualFold(u.Username, username) {
			return u, nil
		}
	}
	return model.User{}, errors.New("not found")
}

// AddUserWithPassword creates new user and saves it in the database.
func (us *UserStorage) AddUserWithPassword(user model.User, password, role string, isAnonymous bool) (model.User, error) {
	if _, err := us.UserByUsername(user.Username); err == nil {
		return model.User{}, model.ErrorUserExists
	}
	if _, err := us.UserByEmail(user.Email); err == nil {
		return model.User{}, model.ErrorUserExists
	}
	if _, err := us.UserByPhone(user.Phone); err == nil {
		return model.User{}, model.ErrorUserExists
	}

	user.Active = true
	user.AccessRole = role
	user.Anonymous = isAnonymous
	user.ID = randomdata.StringNumber(2, "-")

	return us.AddNewUser(user, password)
}

// AddNewUser adds new user to the database.
func (us *UserStorage) AddNewUser(user model.User, password string) (model.User, error) {
	user.Email = strings.ToLower(user.Email)

	user.ID = primitive.NewObjectID().Hex()
	if len(password) > 0 {
		user.Pswd = model.PasswordHash(password)
	}
	user.NumOfLogins = 0

	us.users = append(us.users, user)
	return user, nil
}

// UserByFederatedID returns randomly generated user.
func (us *UserStorage) UserByFederatedID(provider string, id string) (model.User, error) {
	sid := string(provider) + ":" + id
	for _, u := range us.users {
		for _, fi := range u.FederatedIDs {
			if strings.EqualFold(fi, sid) {
				return u, nil
			}
		}
	}
	return model.User{}, errors.New("not found")
}

// AddUserWithFederatedID returns randomly generated user.
func (us *UserStorage) AddUserWithFederatedID(user model.User, provider string, id, role string) (model.User, error) {
	// If there is no error, it means user already exists.
	if _, err := us.UserByFederatedID(provider, id); err == nil {
		return model.User{}, model.ErrorUserExists
	}

	user.ID = primitive.NewObjectID().Hex()
	user.Active = true
	user.AccessRole = role
	user.AddFederatedId(provider, id)

	return us.AddNewUser(user, "")
}

// UpdateUser returns what it receives.
func (us *UserStorage) UpdateUser(userID string, newUser model.User) (model.User, error) {
	newUser.Email = strings.ToLower(newUser.Email)
	newUser.Username = strings.ToLower(newUser.Username)
	newUser.ID = userID

	for i, u := range us.users {
		if strings.EqualFold(userID, u.ID) {
			us.users[i] = newUser
			break
		}
	}

	return newUser, nil
}

// ResetPassword does nothing here.
func (us *UserStorage) ResetPassword(id, password string) error {
	for i, u := range us.users {
		if strings.EqualFold(id, u.ID) {
			u.Pswd = model.PasswordHash(password)
			us.users[i] = u
			break
		}
	}
	return nil
}

// CheckPassword does nothig here.
func (us *UserStorage) CheckPassword(id, password string) error {
	user, err := us.UserByID(id)
	if err != nil {
		return model.ErrUserNotFound
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Pswd), []byte(password)); err != nil {
		// return this error to hide the existence of the user.
		return model.ErrUserNotFound
	}
	return nil
}

// DeleteUser does nothing here.
func (us *UserStorage) DeleteUser(id string) error {
	for i, u := range us.users {
		if strings.EqualFold(id, u.ID) {
			us.users = append(us.users[:i], us.users[i+1:]...)
			break
		}
	}
	return nil
}

// UpdateLoginMetadata does nothing here.
func (us *UserStorage) UpdateLoginMetadata(userID string) {
	for i, u := range us.users {
		if strings.EqualFold(userID, u.ID) {
			u.NumOfLogins += 1
			u.LatestLoginTime = time.Now().Unix()
			us.users[i] = u
			break
		}
	}
}

// FetchUsers returns randomly generated user enclosed in slice.
func (us *UserStorage) FetchUsers(filterString string, skip, limit int) ([]model.User, int, error) {
	// no skip, no filters
	return us.users, len(us.users), nil
}

// ImportJSON imports data from JSON.
func (us *UserStorage) ImportJSON(data []byte, clearOldData bool) error {
	if clearOldData {
		us.users = []model.User{}
	}
	ud := []model.User{}
	if err := json.Unmarshal(data, &ud); err != nil {
		return err
	}
	for _, u := range ud {
		pswd := u.Pswd
		u.Pswd = ""
		if _, err := us.AddNewUser(u, pswd); err != nil {
			return err
		}
	}
	return nil
}

// Close does nothing here.
func (us *UserStorage) Close() {}
