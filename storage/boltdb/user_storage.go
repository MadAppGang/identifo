package boltdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/madappgang/identifo/model"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserBucket              = "Users"             // UserBucket is a name for bucket with users.
	UserBySocialIDBucket    = "UserBySocialID"    // UserBySocialIDBucket is a name for bucket with social IDs as keys.
	UserByUsername          = "UserByUsername"    // UserByUsername  is a name for bucket with user names as keys.
	UserByPhoneNumberBucket = "UserByPhoneNumber" // UserByPhoneNumberBucket is a name for bucket with phone numbers as keys.
	UserByEmailBucket       = "UserByEmail"       // UserByEmailBucket is a name for bucket with email as keys.
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
		if _, err := tx.CreateBucketIfNotExists([]byte(UserByUsername)); err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(UserByPhoneNumberBucket)); err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(UserByEmailBucket)); err != nil {
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

// UserByID returns user by ID.
func (us *UserStorage) UserByID(id string) (model.User, error) {
	var res model.User
	err := us.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(UserBucket))
		u := b.Get([]byte(id))
		if u == nil {
			return model.ErrUserNotFound
		}

		var err error
		res, err = model.UserFromJSON(u)
		return err
	})
	if err != nil {
		return model.User{}, err
	}
	return res, nil
}

// UserByEmail returns user by its email.
func (us *UserStorage) UserByEmail(email string) (model.User, error) {
	var res model.User
	err := us.db.View(func(tx *bolt.Tx) error {
		ueb := tx.Bucket([]byte(UserByEmailBucket))
		// We use email as a key.
		// Get user ID.
		userID := ueb.Get([]byte(email))
		if userID == nil {
			return model.ErrUserNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))
		// Get user by userID.
		u := ub.Get(userID)
		if u == nil {
			return model.ErrUserNotFound
		}

		var err error
		res, err = model.UserFromJSON(u)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return model.User{}, err
	}
	// clear password hash
	res.Pswd = ""
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
		unpb := tx.Bucket([]byte(UserByUsername))
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

	if err := us.db.Update(func(tx *bolt.Tx) error {
		upnb := tx.Bucket([]byte(UserByPhoneNumberBucket))
		return upnb.Delete([]byte(id))
	}); err != nil {
		return err
	}

	return us.db.Update(func(tx *bolt.Tx) error {
		ueb := tx.Bucket([]byte(UserByEmailBucket))
		return ueb.Delete([]byte(id))
	})
}

// UserByFederatedID returns user by federated ID.
func (us *UserStorage) UserByFederatedID(provider model.FederatedIdentityProvider, id string) (model.User, error) {
	var res model.User
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
		res, err = model.UserFromJSON(u)
		return err
	})
	if err != nil {
		return model.User{}, err
	}
	// clear password hash
	res.Pswd = ""
	return res, nil
}

// UserExists checks if user with provided name exists.
func (us *UserStorage) UserExists(name string) bool {
	err := us.db.View(func(tx *bolt.Tx) error {
		unpb := tx.Bucket([]byte(UserByUsername))
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
	var res model.User
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
		u := ub.Get(userID)
		if u == nil {
			return model.ErrUserNotFound
		}

		var err error
		res, err = model.UserFromJSON(u)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return model.User{}, err
	}
	// clear password hash
	res.Pswd = ""
	return res, nil
}

// UserByUsername returns user by name
func (us *UserStorage) UserByUsername(username string) (model.User, error) {
	var res model.User
	err := us.db.View(func(tx *bolt.Tx) error {
		unpb := tx.Bucket([]byte(UserByUsername))
		// we use username and password hash as a key
		key := username
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
		res, err = model.UserFromJSON(u)
		if err != nil {
			return err
		}
		return err
	})
	if err != nil {
		return model.User{}, err
	}
	// clear password hash
	res.Pswd = ""
	return res, nil
}

// AddNewUser adds new user to the storage.
func (us *UserStorage) AddNewUser(user model.User, password string) (model.User, error) {
	user.Pswd = model.PasswordHash(password)

	err := us.db.Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(user)
		if err != nil {
			return err
		}

		ub := tx.Bucket([]byte(UserBucket))
		if err := ub.Put([]byte(user.ID), data); err != nil {
			return err
		}

		err = us.UpdateUserBuckets(tx, user)
		return err
	})
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

// AddUserWithFederatedID adds new user with social ID.
func (us *UserStorage) AddUserWithFederatedID(provider model.FederatedIdentityProvider, federatedID, role string) (model.User, error) {
	sid := string(provider) + ":" + federatedID
	// Using user name as a key. If there is no error, it means user already exists.
	if _, err := us.UserByFederatedID(provider, federatedID); err == nil {
		return model.User{}, model.ErrorUserExists
	}

	user := model.User{
		ID:          sid, // not sure it's a good idea
		Active:      true,
		Username:    sid,
		AccessRole:  role,
		NumOfLogins: 0,
	}

	err := us.db.Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(user)
		if err != nil {
			return err
		}

		ub := tx.Bucket([]byte(UserBucket))
		if err := ub.Put([]byte(user.ID), data); err != nil {
			return err
		}

		usib := tx.Bucket([]byte(UserBySocialIDBucket))
		return usib.Put([]byte(sid), []byte(user.ID))
	})
	if err != nil {
		return model.User{}, err
	}
	return user, nil
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

	u := model.User{
		ID:         xid.New().String(),
		Active:     true,
		Username:   user.Username,
		Phone:      user.Phone,
		Email:      user.Email,
		AccessRole: role,
		Anonymous:  isAnonymous,
	}

	return us.AddNewUser(u, password)
}

// UpdateUser updates user in BoltDB storage.
func (us *UserStorage) UpdateUser(userID string, user model.User) (model.User, error) {
	// use ID from the request if it's not set
	if len(user.ID) == 0 {
		user.ID = userID
	}

	err := us.db.Update(func(tx *bolt.Tx) error {
		ub := tx.Bucket([]byte(UserBucket))
		oldBytes := ub.Get([]byte(userID))

		if len(oldBytes) != 0 {
			oldUser, err := model.UserFromJSON(oldBytes)
			if err != nil {
				return err
			}
			if user.Pswd == "" {
				user.Pswd = oldUser.Pswd
			}
			if user.TFAInfo.Secret == "" {
				user.TFAInfo.Secret = oldUser.TFAInfo.Secret
			}
		}

		data, err := json.Marshal(user)
		if err != nil {
			return err
		}

		if err = ub.Put([]byte(user.ID), data); err != nil {
			return err
		}

		return us.UpdateUserBuckets(tx, user)
	})
	if err != nil {
		return model.User{}, err
	}

	updatedUser, err := us.UserByID(user.ID)
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

		user, err := model.UserFromJSON(u)
		if err != nil {
			return err
		}

		user.Pswd = model.PasswordHash(password)

		u, err = json.Marshal(user)
		if err != nil {
			return err
		}
		return ub.Put([]byte(user.ID), u)
	})
}

// CheckPassword check that password is valid for user id.
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

// ResetUsername sets user username.
func (us *UserStorage) ResetUsername(id, username string) error {
	// TODO: implement
	return errors.New("ResetUsername is not implemented. ")
}

// IDByName returns userID by name.
func (us *UserStorage) IDByName(name string) (string, error) {
	var id string
	err := us.db.View(func(tx *bolt.Tx) error {
		unpb := tx.Bucket([]byte(UserByUsername))
		userID := unpb.Get([]byte(name))
		if userID == nil {
			return model.ErrUserNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))
		u := ub.Get([]byte(userID))
		if u == nil {
			return model.ErrUserNotFound
		}

		user, err := model.UserFromJSON(u)
		if err != nil {
			return err
		}

		if !user.Active {
			return ErrorInactiveUser
		}

		id = user.ID
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
		ub := tx.Bucket([]byte(UserBucket))

		if iterErr := ub.ForEach(func(k, u []byte) error {
			user, err := model.UserFromJSON(u)
			if err != nil {
				return err
			}

			if filterString == "" {
				users = append(users, user)
			} else if strings.Contains(strings.ToLower(string(user.Email)), strings.ToLower(filterString)) {
				users = append(users, user)
			} else if strings.Contains(strings.ToLower(string(user.Phone)), strings.ToLower(filterString)) {
				users = append(users, user)
			} else if strings.Contains(strings.ToLower(string(user.Username)), strings.ToLower(filterString)) {
				users = append(users, user)
			}
			return nil
		}); iterErr != nil {
			return iterErr
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

func (us *UserStorage) UpdateUserBuckets(tx *bolt.Tx, user model.User) error {
	if user.Username != "" {
		unpb := tx.Bucket([]byte(UserByUsername))
		if err := unpb.Put([]byte(user.Username), []byte(user.ID)); err != nil {
			return err
		}
	}

	if user.Email != "" {
		ueb := tx.Bucket([]byte(UserByEmailBucket))
		if err := ueb.Put([]byte(user.Email), []byte(user.ID)); err != nil {
			return err
		}
	}
	if user.Phone != "" {
		upnb := tx.Bucket([]byte(UserByPhoneNumberBucket))
		if err := upnb.Put([]byte(user.Phone), []byte(user.ID)); err != nil {
			return err
		}
	}
	return nil
}
