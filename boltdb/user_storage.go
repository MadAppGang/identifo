package boltdb

import (
	"encoding/json"
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
)

// NewUserStorage creates and inits an embedded user storage.
func NewUserStorage(db *bolt.DB) (model.UserStorage, error) {
	us := UserStorage{db: db}
	// ensure we have app's bucket in the database.
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
			return model.ErrorNotFound
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

	err := us.db.Update(func(tx *bolt.Tx) error {
		usib := tx.Bucket([]byte(UserBySocialIDBucket))
		return usib.Delete([]byte(id))
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
			return model.ErrorNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))
		// get user by userID.
		u := ub.Get(userID)
		if u == nil {
			return model.ErrorNotFound
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
			return model.ErrorNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))
		if u := ub.Get([]byte(userID)); u == nil {
			return model.ErrorNotFound
		}
		return nil
	})
	return err == nil
}

//AttachDeviceToken does nothing here.
func (us *UserStorage) AttachDeviceToken(id, token string) error {
	//we are not supporting devices for users here
	return model.ErrorNotImplemented
}

//DetachDeviceToken does nothing here.
func (us *UserStorage) DetachDeviceToken(token string) error {
	//we are not supporting devices for users here
	return model.ErrorNotImplemented
}

//RequestScopes mem always returns requested scope
func (us *UserStorage) RequestScopes(userID string, scopes []string) ([]string, error) {
	//we allow all scopes for embedded database, you could implement your own logic in external service
	return scopes, nil
}

//Scopes returns supported scopes, could be static data of database
func (us *UserStorage) Scopes() []string {
	//we allow all scopes for embedded database, you could implement your own logic in external service
	return []string{"offline", "user"}
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
			return model.ErrorNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))
		// get user by userID
		u := ub.Get(userID)
		if u == nil {
			return model.ErrorNotFound
		}

		var err error
		res, err = UserFromJSON(u)
		if err != nil {
			return err
		}
		if err = bcrypt.CompareHashAndPassword([]byte(res.PasswordHash()), []byte(password)); err != nil {
			// return this error to hide the existence of the user.
			return model.ErrorNotFound
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
		key := u.Name()
		unpb := tx.Bucket([]byte(UserByNameAndPassword))
		return unpb.Put([]byte(key), []byte(u.ID()))
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}

// AddUserWithFederatedID adds new user with social ID.
func (us *UserStorage) AddUserWithFederatedID(provider model.FederatedIdentityProvider, federatedID string) (model.User, error) {
	sid := string(provider) + ":" + federatedID
	// Using user name as a key. If there is no error, it means user already exists.
	if _, err := us.UserByFederatedID(provider, federatedID); err == nil {
		return nil, model.ErrorUserExists
	}

	u := userData{Active: true, Name: sid}
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
	// Using user name as a key. If there is no error, it means user already exists.
	if _, err := us.UserByID(name); err == nil {
		return nil, model.ErrorUserExists
	}
	u := userData{Active: true, Name: name, Profile: profile, ID: name}
	return us.AddNewUser(User{userData: u}, password)
}

// UpdateUser updates user in BoltDB storage.
func (us *UserStorage) UpdateUser(userID string, newUser model.User) (model.User, error) {
	res, ok := newUser.(*User)
	if !ok || res == nil {
		return nil, ErrorWrongDataFormat
	}

	// generate new ID if it's not set
	if len(newUser.ID()) == 0 {
		res.userData.ID = xid.New().String()
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

	updatedUser, err := us.UserByID(userID)
	return updatedUser, err
}

// ResetPassword sets new user password.
func (us *UserStorage) ResetPassword(id, password string) error {
	return us.db.Update(func(tx *bolt.Tx) error {
		ub := tx.Bucket([]byte(UserBucket))
		u := ub.Get([]byte(id))
		if u == nil {
			return model.ErrorNotFound
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

// IDByName returns userID by name.
func (us *UserStorage) IDByName(name string) (string, error) {
	var id string
	err := us.db.View(func(tx *bolt.Tx) error {
		unpb := tx.Bucket([]byte(UserByNameAndPassword))
		userID := unpb.Get([]byte(name))
		if userID == nil {
			return model.ErrorNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))
		u := ub.Get([]byte(userID))
		if u == nil {
			return model.ErrorNotFound
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
func (us *UserStorage) FetchUsers(filterString string, skip, limit int) ([]model.User, error) {
	var users []model.User

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
		for _, uid := range userIDs {
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
		return nil, err
	}
	return users, nil
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
	ID      string                 `json:"id,omitempty"`
	Name    string                 `json:"name,omitempty"`
	Pswd    string                 `json:"pswd,omitempty"`
	Profile map[string]interface{} `json:"profile,omitempty"`
	Active  bool                   `json:"active,omitempty"`
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

// Name implements model.User interface.
func (u User) Name() string { return u.userData.Name }

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
