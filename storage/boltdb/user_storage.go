package boltdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/madappgang/identifo/model"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

const (
	// UserBucket is a name for bucket with users.
	UserBucket = "Users"
	// UserBySocialIDBucket is a name for bucket with social IDs as keys.
	UserBySocialIDBucket = "UserBySocialID"
	// UserByNameAndPassword  is a name for bucket with user names as keys.
	UserByNameAndPassword = "UserByNameAndPassword"
	// UserByPhoneNumberBucket is a name for bucket with phone numbers as keys.
	UserByPhoneNumberBucket = "UserByPhoneNumber"
)

// NewUserStorage creates and inits an embedded user storage.
func NewUserStorage(db *bolt.DB) (model.UserStorage, error) {
	us := UserStorage{db: db}

	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(UserBucket)); err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(UserBySocialIDBucket)); err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(UserByNameAndPassword)); err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(UserByPhoneNumberBucket)); err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &us, nil
}

// UserStorage implements user storage interface for BoltDB.
type UserStorage struct {
	db *bolt.DB
}

// NewUser returns pointer to newly created user.
func (us *UserStorage) NewUser() model.User {
	return &User{}
}

// UserByID returns user by ID.
func (us *UserStorage) UserByID(id string) (model.User, error) {
	var res User
	err := us.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UserBucket))
		u := b.Get([]byte(id))
		if u == nil {
			return model.ErrUserNotFound
		}

		var err error
		res, err = UserFromJSON(u)
		return err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UserByEmail returns user by its email.
func (us *UserStorage) UserByEmail(email string) (model.User, error) {
	// TODO: implement boltdb UserByEmail
	return nil, errors.New("Not implemented. ")
}

// DeleteUser deletes user by ID.
func (us *UserStorage) DeleteUser(id string) error {
	if err := us.db.Update(func(tx *bolt.Tx) error {
		ub := tx.Bucket([]byte(UserBucket))
		return ub.Delete([]byte(id))
	}); err != nil {
		return err
	}

	if err := us.db.Update(func(tx *bolt.Tx) error {
		unpb := tx.Bucket([]byte(UserByNameAndPassword))
		return unpb.Delete([]byte(id))
	}); err != nil {
		return err
	}

	if err := us.db.Update(func(tx *bolt.Tx) error {
		usib := tx.Bucket([]byte(UserBySocialIDBucket))
		return usib.Delete([]byte(id))
	}); err != nil {
		return err
	}

	err := us.db.Update(func(tx *bolt.Tx) error {
		upnb := tx.Bucket([]byte(UserByPhoneNumberBucket))
		return upnb.Delete([]byte(id))
	})
	return err
}

// UserByFederatedID returns user by federated ID.
func (us *UserStorage) UserByFederatedID(provider model.FederatedIdentityProvider, id string) (model.User, error) {
	var res User
	sid := string(provider) + ":" + id

	err := us.db.View(func(tx *bolt.Tx) error {
		usib := tx.Bucket([]byte(UserBySocialIDBucket))
		// get userID from index.
		userID := usib.Get([]byte(sid))
		if userID == nil {
			return model.ErrUserNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))
		// get user by userID.
		u := ub.Get(userID)
		if u == nil {
			return model.ErrUserNotFound
		}

		var err error
		res, err = UserFromJSON(u)
		return err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UserExists checks if user with provided name exists.
func (us *UserStorage) UserExists(name string) bool {
	err := us.db.View(func(tx *bolt.Tx) error {
		unpb := tx.Bucket([]byte(UserByNameAndPassword))
		userID := unpb.Get([]byte(name))

		if userID == nil {
			return model.ErrUserNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))
		if u := ub.Get([]byte(userID)); u == nil {
			return model.ErrUserNotFound
		}
		return nil
	})
	return err == nil
}

// AttachDeviceToken does nothing here.
func (us *UserStorage) AttachDeviceToken(id, token string) error {
	// BoltDB-backed implementation does not support user devices.
	return model.ErrorNotImplemented
}

// DetachDeviceToken does nothing here.
func (us *UserStorage) DetachDeviceToken(token string) error {
	// BoltDB-backed implementation does not support user devices.
	return model.ErrorNotImplemented
}

// RequestScopes returns requested scopes.
func (us *UserStorage) RequestScopes(userID string, scopes []string) ([]string, error) {
	// We allow all scopes for embedded database, you can implement your own logic in the external service.
	return scopes, nil
}

// Scopes returns supported scopes.
func (us *UserStorage) Scopes() []string {
	// We allow all scopes for embedded database, you can implement your own logic in the external service.
	return []string{"offline", "user"}
}

// UserByPhone fetches user by phone number.
func (us *UserStorage) UserByPhone(phone string) (model.User, error) {
	var res User
	err := us.db.View(func(tx *bolt.Tx) error {
		upnb := tx.Bucket([]byte(UserByPhoneNumberBucket))
		// We use phone number as a key.
		// Get user ID.
		userID := upnb.Get([]byte(phone))
		if userID == nil {
			return model.ErrUserNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))
		// Get user by userID.
		if u := ub.Get(userID); u == nil {
			return model.ErrUserNotFound
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UserByNamePassword returns user by name and password.
func (us *UserStorage) UserByNamePassword(name, password string) (model.User, error) {
	var res User
	err := us.db.View(func(tx *bolt.Tx) error {
		unpb := tx.Bucket([]byte(UserByNameAndPassword))
		// we use username and password hash as a key
		key := name
		// get user ID from index
		userID := unpb.Get([]byte(key))
		if userID == nil {
			return model.ErrUserNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))
		// get user by userID
		u := ub.Get(userID)
		if u == nil {
			return model.ErrUserNotFound
		}

		var err error
		res, err = UserFromJSON(u)
		if err != nil {
			return err
		}
		if err = bcrypt.CompareHashAndPassword([]byte(res.PasswordHash()), []byte(password)); err != nil {
			// return this error to hide the existence of the user.
			return model.ErrUserNotFound
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// AddNewUser adds new user to the storage.
func (us *UserStorage) AddNewUser(usr model.User, password string) (model.User, error) {
	u, ok := usr.(User)
	if !ok {
		return nil, ErrorWrongDataFormat
	}
	u.userData.Pswd = PasswordHash(password)

	err := us.db.Update(func(tx *bolt.Tx) error {
		data, err := u.Marshal()
		if err != nil {
			return err
		}

		ub := tx.Bucket([]byte(UserBucket))
		if err := ub.Put([]byte(u.ID()), data); err != nil {
			return err
		}

		// we use username and password hash as a key
		key := u.Username()
		unpb := tx.Bucket([]byte(UserByNameAndPassword))
		return unpb.Put([]byte(key), []byte(u.ID()))
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}

// AddUserByPhone registers new user with phone number.
func (us *UserStorage) AddUserByPhone(phone string) (model.User, error) {
	u := &User{userData: userData{Active: true, Phone: phone, ID: xid.New().String()}}

	err := us.db.Update(func(tx *bolt.Tx) error {
		data, err := u.Marshal()
		if err != nil {
			return err
		}

		ub := tx.Bucket([]byte(UserBucket))
		if err := ub.Put([]byte(u.ID()), data); err != nil {
			return err
		}

		// We use phone number as a key.
		upnb := tx.Bucket([]byte(UserByPhoneNumberBucket))
		return upnb.Put([]byte(phone), []byte(u.userData.ID))
	})
	if err != nil {
		return nil, err
	}

	return u, err
}

// AddUserWithFederatedID adds new user with social ID.
func (us *UserStorage) AddUserWithFederatedID(provider model.FederatedIdentityProvider, federatedID string) (model.User, error) {
	sid := string(provider) + ":" + federatedID
	// Using user name as a key. If there is no error, it means user already exists.
	if _, err := us.UserByFederatedID(provider, federatedID); err == nil {
		return nil, model.ErrorUserExists
	}

	u := userData{Active: true, Username: sid}
	u.ID = sid // not sure it's a good idea
	user := User{userData: u}

	err := us.db.Update(func(tx *bolt.Tx) error {
		data, err := user.Marshal()
		if err != nil {
			return err
		}

		ub := tx.Bucket([]byte(UserBucket))
		if err := ub.Put([]byte(user.ID()), data); err != nil {
			return err
		}

		usib := tx.Bucket([]byte(UserBySocialIDBucket))
		return usib.Put([]byte(sid), []byte(user.ID()))
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

// AddUserByNameAndPassword creates new user and saves it in the database.
func (us *UserStorage) AddUserByNameAndPassword(name, password string, profile map[string]interface{}) (model.User, error) {
	if us.UserExists(name) {
		return nil, model.ErrorUserExists
	}

	u := userData{Active: true, Username: name, Profile: profile, ID: xid.New().String()}
	return us.AddNewUser(User{userData: u}, password)
}

// UpdateUser updates user in BoltDB storage.
func (us *UserStorage) UpdateUser(userID string, newUser model.User) (model.User, error) {
	res, ok := newUser.(*User)
	if !ok || res == nil {
		return nil, ErrorWrongDataFormat
	}

	// use ID from the request if it's not set
	if len(res.ID()) == 0 {
		res.userData.ID = userID
	}

	err := us.db.Update(func(tx *bolt.Tx) error {
		data, err := res.Marshal()
		if err != nil {
			return err
		}

		ub := tx.Bucket([]byte(UserBucket))
		if err := ub.Delete([]byte(userID)); err != nil {
			return err
		}

		return ub.Put([]byte(res.ID()), data)
	})
	if err != nil {
		return nil, err
	}

	updatedUser, err := us.UserByID(res.ID())
	return updatedUser, err
}

// ResetPassword sets new user password.
func (us *UserStorage) ResetPassword(id, password string) error {
	return us.db.Update(func(tx *bolt.Tx) error {
		ub := tx.Bucket([]byte(UserBucket))
		u := ub.Get([]byte(id))
		if u == nil {
			return model.ErrUserNotFound
		}

		user, err := UserFromJSON(u)
		if err != nil {
			return err
		}

		user.userData.Pswd = PasswordHash(password)

		u, err = user.Marshal()
		if err != nil {
			return err
		}
		return ub.Put([]byte(user.ID()), u)
	})
}

// ResetUsername sets user username.
func (us *UserStorage) ResetUsername(id, username string) error {
	// TODO: implement
	return errors.New("ResetUsername is not implemented. ")
}

// IDByName returns userID by name.
func (us *UserStorage) IDByName(name string) (string, error) {
	var id string
	err := us.db.View(func(tx *bolt.Tx) error {
		unpb := tx.Bucket([]byte(UserByNameAndPassword))
		userID := unpb.Get([]byte(name))
		if userID == nil {
			return model.ErrUserNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))
		u := ub.Get([]byte(userID))
		if u == nil {
			return model.ErrUserNotFound
		}

		user, err := UserFromJSON(u)
		if err != nil {
			return err
		}

		if !user.Active() {
			return ErrorInactiveUser
		}

		id = user.ID()
		return nil
	})

	if err != nil {
		return "", err
	}
	return id, nil
}

// FetchUsers fetches users which name satisfies provided filterString.
// Supports pagination.
func (us *UserStorage) FetchUsers(filterString string, skip, limit int) ([]model.User, int, error) {
	users := []model.User{}
	var total int

	err := us.db.View(func(tx *bolt.Tx) error {
		ubnp := tx.Bucket([]byte(UserByNameAndPassword))
		var userIDs [][]byte

		if iterErr := ubnp.ForEach(func(k, v []byte) error {
			if strings.Contains(strings.ToLower(string(k)), strings.ToLower(filterString)) {
				userIDs = append(userIDs, v)
			}
			return nil
		}); iterErr != nil {
			return iterErr
		}

		ub := tx.Bucket([]byte(UserBucket))
		total = len(userIDs)

		for i, uid := range userIDs {
			if i < skip {
				continue
			}
			if limit != 0 && len(users) == limit {
				break
			}

			u := ub.Get(uid)
			if u == nil {
				log.Printf("User %s does not exist in %s, but does exist in %s", uid, UserBucket, UserByNameAndPassword)
				continue
			}

			user, err := UserFromJSON(u)
			if err != nil {
				return err
			}
			users = append(users, user)
		}
		return nil
	})

	if err != nil {
		return []model.User{}, 0, err
	}
	return users, total, nil
}

// ImportJSON imports data from JSON.
func (us *UserStorage) ImportJSON(data []byte) error {
	ud := []userData{}
	if err := json.Unmarshal(data, &ud); err != nil {
		return err
	}
	for _, u := range ud {
		pswd := u.Pswd
		u.Pswd = ""
		if _, err := us.AddNewUser(&User{userData: u}, pswd); err != nil {
			return err
		}
	}
	return nil
}

// User data implementation.
type userData struct {
	ID       string                 `json:"id,omitempty"`
	Username string                 `json:"username,omitempty"`
	Email    string                 `json:"email,omitempty"`
	Phone    string                 `bson:"phone,omitempty" json:"phone,omitempty"`
	Pswd     string                 `json:"pswd,omitempty"`
	Profile  map[string]interface{} `json:"profile,omitempty"`
	Active   bool                   `json:"active,omitempty"`
}

// User is a user data structure for embedded storage.
type User struct {
	userData
}

// UserFromJSON deserializes user data from JSON.
func UserFromJSON(d []byte) (User, error) {
	user := userData{}
	if err := json.Unmarshal(d, &user); err != nil {
		return User{}, err
	}
	return User{userData: user}, nil
}

// Marshal serializes data to byte array.
func (u User) Marshal() ([]byte, error) {
	return json.Marshal(u.userData)
}

// Sanitize removes all sensitive data.
func (u User) Sanitize() model.User {
	u.userData.Pswd = ""
	u.userData.Active = false
	return u
}

// ID implements model.User interface.
func (u User) ID() string { return u.userData.ID }

// Username implements model.User interface.
func (u User) Username() string { return u.userData.Username }

// SetUsername implements model.User interface.
func (u User) SetUsername(username string) { u.userData.Username = username }

// Email implements model.User interface.
func (u User) Email() string { return u.userData.Email }

// SetEmail implements model.Email interface.
func (u User) SetEmail(email string) { u.userData.Email = email }

// PasswordHash implements model.User interface.
func (u User) PasswordHash() string { return u.userData.Pswd }

// Profile implements model.User interface.
func (u User) Profile() map[string]interface{} { return u.userData.Profile }

// Active implements model.User interface.
func (u User) Active() bool { return u.userData.Active }

// PasswordHash creates hash with salt for password.
func PasswordHash(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash)
}
