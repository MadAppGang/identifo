package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/hashicorp/go-plugin"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/plugin/shared"
	"github.com/madappgang/identifo/proto"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

// Here is a real implementation of KV that writes to a local file with
// the key name and the contents are the value of the key.
type UserStorage struct {
	db *bolt.DB
}

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

func main() {
	filepath := os.Getenv("DB_FILE_PATH")

	if filepath == "" {
		panic("Empty DB_FILE_PATH")
	}

	db, err := bolt.Open(filepath, 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	us, err := NewUserStorage(db)
	if err != nil {
		panic(err)
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"user_storage": &shared.UserStorageGRPCPlugin{
				Impl: us,
			},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

// NewUserStorage creates and inits an embedded user storage.
func NewUserStorage(db *bolt.DB) (shared.UserStorage, error) {
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

// UserByID returns user by ID.
func (us *UserStorage) UserByID(id string) (*proto.User, error) {
	res := new(proto.User)
	err := us.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UserBucket))
		u := b.Get([]byte(id))
		if u == nil {
			return shared.ErrUserNotFound
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
func (us *UserStorage) UserByEmail(email string) (*proto.User, error) {
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
func (us *UserStorage) UserByFederatedID(provider proto.FederatedIdentityProvider, id string) (*proto.User, error) {
	res := new(proto.User)
	sid := provider.String() + ":" + id

	err := us.db.View(func(tx *bolt.Tx) error {
		usib := tx.Bucket([]byte(UserBySocialIDBucket))
		// get userID from index.
		userID := usib.Get([]byte(sid))
		if userID == nil {
			return shared.ErrUserNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))
		// get user by userID.
		u := ub.Get(userID)
		if u == nil {
			return shared.ErrUserNotFound
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
			return shared.ErrUserNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))
		if u := ub.Get([]byte(userID)); u == nil {
			return shared.ErrUserNotFound
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
func (us *UserStorage) UserByPhone(phone string) (*proto.User, error) {
	res := new(proto.User)
	err := us.db.View(func(tx *bolt.Tx) error {
		upnb := tx.Bucket([]byte(UserByPhoneNumberBucket))
		// We use phone number as a key.
		// Get user ID.
		userID := upnb.Get([]byte(phone))
		if userID == nil {
			return shared.ErrUserNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))
		// Get user by userID.
		if u := ub.Get(userID); u == nil {
			return shared.ErrUserNotFound
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UserByNamePassword returns user by name and password.
func (us *UserStorage) UserByNamePassword(name, password string) (*proto.User, error) {
	res := new(proto.User)
	err := us.db.View(func(tx *bolt.Tx) error {
		unpb := tx.Bucket([]byte(UserByNameAndPassword))
		// we use username and password hash as a key
		key := name
		// get user ID from index
		userID := unpb.Get([]byte(key))
		if userID == nil {
			return shared.ErrUserNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))
		// get user by userID
		u := ub.Get(userID)
		if u == nil {
			return shared.ErrUserNotFound
		}

		var err error
		res, err = UserFromJSON(u)
		if err != nil {
			return err
		}
		if err = bcrypt.CompareHashAndPassword([]byte(res.PasswordHash), []byte(password)); err != nil {
			// return this error to hide the existence of the user.
			return shared.ErrUserNotFound
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// AddNewUser adds new user to the storage.
func (us *UserStorage) AddNewUser(u *proto.User, password string) (*proto.User, error) {
	u.PasswordHash = PasswordHash(password)
	u.NumOfLogins = 0

	err := us.db.Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(u)
		if err != nil {
			return err
		}

		ub := tx.Bucket([]byte(UserBucket))
		if err := ub.Put([]byte(u.Id), data); err != nil {
			return err
		}

		// we use username and password hash as a key
		key := u.Username
		unpb := tx.Bucket([]byte(UserByNameAndPassword))
		return unpb.Put([]byte(key), []byte(u.Id))
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}

// AddUserByPhone registers new user with phone number.
func (us *UserStorage) AddUserByPhone(phone, role string) (*proto.User, error) {
	u := &proto.User{
		Id:          xid.New().String(),
		Username:    phone,
		IsActive:    true,
		Phone:       phone,
		AccessRole:  role,
		NumOfLogins: 0,
	}

	err := us.db.Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(u)
		if err != nil {
			return err
		}

		ub := tx.Bucket([]byte(UserBucket))
		if err := ub.Put([]byte(u.Id), data); err != nil {
			return err
		}

		// We use phone number as a key.
		upnb := tx.Bucket([]byte(UserByPhoneNumberBucket))
		return upnb.Put([]byte(phone), []byte(u.Id))
	})
	if err != nil {
		return nil, err
	}

	return u, err
}

// AddUserWithFederatedID adds new user with social ID.
func (us *UserStorage) AddUserWithFederatedID(provider proto.FederatedIdentityProvider, federatedID, role string) (*proto.User, error) {
	sid := provider.String() + ":" + federatedID
	// Using user name as a key. If there is no error, it means user already exists.
	if _, err := us.UserByFederatedID(provider, federatedID); err == nil {
		return nil, model.ErrorUserExists
	}

	user := &proto.User{IsActive: true, Username: sid, AccessRole: role, NumOfLogins: 0}
	user.Id = sid // not sure it's a good idea

	err := us.db.Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(user)
		if err != nil {
			return err
		}

		ub := tx.Bucket([]byte(UserBucket))
		if err := ub.Put([]byte(user.Id), data); err != nil {
			return err
		}

		usib := tx.Bucket([]byte(UserBySocialIDBucket))
		return usib.Put([]byte(sid), []byte(user.Id))
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

// AddUserByNameAndPassword creates new user and saves it in the database.
func (us *UserStorage) AddUserByNameAndPassword(username, password, role string, isAnonymous bool) (*proto.User, error) {
	if us.UserExists(username) {
		return nil, model.ErrorUserExists
	}

	user := &proto.User{
		Id:          xid.New().String(),
		IsActive:    true,
		Username:    username,
		AccessRole:  role,
		IsAnonymous: isAnonymous,
	}

	if shared.EmailRegexp.MatchString(username) {
		user.Email = username
	}
	if shared.PhoneRegexp.MatchString(username) {
		user.Phone = username
	}

	return us.AddNewUser(user, password)
}

// UpdateUser updates user in BoltDB storage.
func (us *UserStorage) UpdateUser(userID string, newUser *proto.User) (*proto.User, error) {
	// use ID from the request if it's not set
	if len(newUser.Id) == 0 {
		newUser.Id = userID
	}

	err := us.db.Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(newUser)
		if err != nil {
			return err
		}

		ub := tx.Bucket([]byte(UserBucket))
		if err := ub.Delete([]byte(userID)); err != nil {
			return err
		}

		return ub.Put([]byte(newUser.Id), data)
	})
	if err != nil {
		return nil, err
	}

	updatedUser, err := us.UserByID(newUser.Id)
	return updatedUser, err
}

// ResetPassword sets new user password.
func (us *UserStorage) ResetPassword(id, password string) error {
	return us.db.Update(func(tx *bolt.Tx) error {
		ub := tx.Bucket([]byte(UserBucket))
		u := ub.Get([]byte(id))
		if u == nil {
			return shared.ErrUserNotFound
		}

		user, err := UserFromJSON(u)
		if err != nil {
			return err
		}

		user.PasswordHash = PasswordHash(password)

		u, err = json.Marshal(user)
		if err != nil {
			return err
		}
		return ub.Put([]byte(user.Id), u)
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
			return shared.ErrUserNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))
		u := ub.Get([]byte(userID))
		if u == nil {
			return shared.ErrUserNotFound
		}

		user, err := UserFromJSON(u)
		if err != nil {
			return err
		}

		if !user.IsActive {
			return errors.New("User is inactive")
		}

		id = user.Id
		return nil
	})

	if err != nil {
		return "", err
	}
	return id, nil
}

// FetchUsers fetches users which name satisfies provided filterString.
// Supports pagination.
func (us *UserStorage) FetchUsers(filterString string, skip, limit int) ([]*proto.User, int, error) {
	users := []*proto.User{}
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
		return []*proto.User{}, 0, err
	}
	return users, total, nil
}

// ImportJSON imports data from JSON.
func (us *UserStorage) ImportJSON(data []byte) error {
	ud := []*proto.User{}
	if err := json.Unmarshal(data, &ud); err != nil {
		return err
	}
	for _, u := range ud {
		pswd := u.PasswordHash
		u.PasswordHash = ""
		if _, err := us.AddNewUser(u, pswd); err != nil {
			return err
		}
	}
	return nil
}

// UpdateLoginMetadata updates user's login metadata.
func (us *UserStorage) UpdateLoginMetadata(userID string) {
	user, err := us.UserByID(userID)
	if err != nil {
		log.Printf("Cannot get user by ID %s: %s\n", userID, err)
	}

	user.NumOfLogins++
	user.LatestLoginTime = time.Now().Unix()

	if _, err := us.UpdateUser(UserBucket, user); err != nil {
		log.Println("Cannot update user login info: ", err)
	}
}

// Close closes underlying database.
func (us *UserStorage) Close() {
	if err := us.db.Close(); err != nil {
		log.Printf("Error closing user storage: %s\n", err)
	}
}

// PasswordHash creates hash with salt for password.
func PasswordHash(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash)
}

// UserFromJSON deserializes user data from JSON.
func UserFromJSON(d []byte) (*proto.User, error) {
	user := new(proto.User)
	if err := json.Unmarshal(d, user); err != nil {
		return nil, err
	}
	return user, nil
}
