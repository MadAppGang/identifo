package boltdb

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/madappgang/identifo/model"
	"golang.org/x/crypto/bcrypt"
)

const (
	//UserBucket bucket name with users
	UserBucket = "Users"
	//UserBySocialIDBucket bucket name for user index
	UserBySocialIDBucket = "UserBySocialID"
	//UserByNameAndPassword bucket name for user index
	UserByNameAndPassword = "UserByNameAndPassword"
)

//NewUserStorage creates and inits embedded user storage
func NewUserStorage(db *bolt.DB) (model.UserStorage, error) {
	us := UserStorage{}
	us.db = db
	//ensure we have app's bucket in the database
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(UserBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte(UserBySocialIDBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte(UserByNameAndPassword))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &us, nil
}

//UserStorage implements user storage in memory
type UserStorage struct {
	db *bolt.DB
}

//UserByID returns random generated user
func (us *UserStorage) UserByID(id string) (model.User, error) {
	var res User
	err := us.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UserBucket))
		v := b.Get([]byte(id))
		if v == nil {
			return ErrorNotFound
		}
		rr, err := UserFromJSON(v)
		res = rr
		return err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

//UserBySocialID returns random generated user
func (us *UserStorage) UserBySocialID(id string) (model.User, error) {
	var res User
	err := us.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UserBucket))
		bi := tx.Bucket([]byte(UserBySocialIDBucket))
		//get user ID from index
		userID := bi.Get([]byte(id))
		if userID == nil {
			return ErrorNotFound
		}
		//get user by userID
		u := b.Get(userID)
		if u == nil {
			return ErrorNotFound
		}
		rr, err := UserFromJSON(u)
		res = rr
		return err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

//AttachDeviceToken do nothing here
func (us *UserStorage) AttachDeviceToken(id, token string) error {
	//we are not supporting devices for users here
	return nil
}

//RequestScopes mem always returns requested scope
func (us *UserStorage) RequestScopes(userID string, scopes []string) ([]string, error) {
	//we allow all scopes for embedded database, you could implement your own logic in external service
	return scopes, nil
}

//UserByNamePassword returns random generated user
func (us *UserStorage) UserByNamePassword(name, password string) (model.User, error) {
	var res User
	err := us.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UserBucket))
		bi := tx.Bucket([]byte(UserByNameAndPassword))
		//we use username and password hash as the key
		key := name
		//get user ID from index
		userID := bi.Get([]byte(key))
		if userID == nil {
			return ErrorNotFound
		}
		//get user by userID
		u := b.Get(userID)
		if u == nil {
			return ErrorNotFound
		}
		rr, err := UserFromJSON(u)
		if err := bcrypt.CompareHashAndPassword([]byte(rr.PasswordHash()), []byte(password)); err != nil {
			//return this error to hide the existence of the user
			return ErrorNotFound
		}
		res = rr
		return err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

//AddNewUser adds new user
func (us *UserStorage) AddNewUser(usr model.User, password string) (model.User, error) {
	u, ok := usr.(User)
	if !ok {
		return nil, ErrorWrongDataFormat
	}
	u.userData.Pswd = PasswordHash(password)
	err := us.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UserBucket))
		bi := tx.Bucket([]byte(UserByNameAndPassword))
		//we use username and password hash as the key
		key := u.Name()
		data, err := u.Marshal()
		if err != nil {
			return err
		}
		if err := b.Put([]byte(u.ID()), data); err != nil {
			return err
		}
		return bi.Put([]byte(key), []byte(u.ID()))
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}

//data implementation
type userData struct {
	ID      string                 `json:"id,omitempty"`
	Name    string                 `json:"name,omitempty"`
	Pswd    string                 `json:"pswd,omitempty"`
	Profile map[string]interface{} `json:"profile,omitempty"`
	Active  bool                   `json:"active,omitempty"`
}

//User user data structure for embedded storage
type User struct {
	userData
}

//UserFromJSON deserializes data
func UserFromJSON(d []byte) (User, error) {
	user := userData{}
	if err := json.Unmarshal(d, &user); err != nil {
		return User{}, err
	}
	return User{user}, nil
}

//Marshal serialize data to byte array
func (u User) Marshal() ([]byte, error) {
	return json.Marshal(u.userData)
}

//model.User interface implementation
func (u User) ID() string                      { return u.userData.ID }
func (u User) Name() string                    { return u.userData.Name }
func (u User) PasswordHash() string            { return u.userData.Pswd }
func (u User) Profile() map[string]interface{} { return u.userData.Profile }
func (u User) Active() bool                    { return u.userData.Active }

//PasswordHash creates hash with salt for password
func PasswordHash(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash)
}
