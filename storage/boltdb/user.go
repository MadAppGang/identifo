package boltdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/rs/xid"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserBucket              = "Users"             // UserBucket is a name for bucket with users.
	UserDataBucket          = "UserData"          // UserDataBucket is a name for bucket with user data.
	UserBySocialIDBucket    = "UserBySocialID"    // UserBySocialIDBucket is a name for bucket with social IDs as keys.
	UserByUsername          = "UserByUsername"    // UserByUsername  is a name for bucket with user names as keys.
	UserByPhoneNumberBucket = "UserByPhoneNumber" // UserByPhoneNumberBucket is a name for bucket with phone numbers as keys.
	UserByEmailBucket       = "UserByEmail"       // UserByEmailBucket is a name for bucket with email as keys.
)

// NewUserStorage creates and inits an embedded user storage.
func NewUserStorage(settings model.BoltDBDatabaseSettings) (*UserStorage, error) {
	if len(settings.Path) == 0 {
		return nil, ErrorEmptyDatabasePath
	}

	// init database
	db, err := InitDB(settings.Path)
	if err != nil {
		return nil, err
	}

	us := UserStorage{db: db}

	if err := us.createBuckets(); err != nil {
		return nil, err
	}

	return &us, nil
}

// UserStorage implements user storage interface for BoltDB.
type UserStorage struct {
	db *bolt.DB
}

func (us *UserStorage) createBuckets() error {
	return us.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(UserBucket)); err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(UserDataBucket)); err != nil {
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
	})
}

// ================================================================
// UserStorage implementations
// ================================================================

// UserByID returns user from storage by ID.
func (us *UserStorage) UserByID(ctx context.Context, ID string) (model.User, error) {
	var res model.User
	err := us.db.View(func(tx *bolt.Tx) error {
		ub := tx.Bucket([]byte(UserBucket))
		// Get user by userID.
		u := ub.Get([]byte(ID))
		if u == nil {
			return l.ErrorUserNotFound
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
	return res, nil
}

// UserByPhone fetches user by phone number.
func (us *UserStorage) UserByPhone(ctx context.Context, Phone string) (model.User, error) {
	var res model.User
	err := us.db.View(func(tx *bolt.Tx) error {
		upnb := tx.Bucket([]byte(UserByPhoneNumberBucket))
		// We use phone number as a key.
		// Get user ID.
		userID := upnb.Get([]byte(Phone))
		if userID == nil {
			return l.ErrorUserNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))
		// Get user by userID.
		u := ub.Get(userID)
		if u == nil {
			return l.ErrorUserNotFound
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
	return res, nil
}

// UserByEmail returns user by its email.
func (us *UserStorage) UserByEmail(ctx context.Context, email string) (model.User, error) {
	var res model.User
	err := us.db.View(func(tx *bolt.Tx) error {
		ueb := tx.Bucket([]byte(UserByEmailBucket))
		// We use email as a key.
		// Get user ID.
		userID := ueb.Get([]byte(email))
		if userID == nil {
			return l.ErrorUserNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))
		// Get user by userID.
		u := ub.Get(userID)
		if u == nil {
			return l.ErrorUserNotFound
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
	return res, nil
}

// UserByIdentity returns user by federated ID.
func (us *UserStorage) UserByIdentity(ctx context.Context, idType model.UserIdentityType, userIdentityTypeOther, externalID string) (model.User, error) {
	var res model.User
	sid := string(idType) + ":" + userIdentityTypeOther + ":" + externalID

	err := us.db.View(func(tx *bolt.Tx) error {
		usib := tx.Bucket([]byte(UserBySocialIDBucket))
		// get userID from index.
		userID := usib.Get([]byte(sid))
		if userID == nil {
			return l.ErrorUserNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))
		// get user by userID.
		u := ub.Get(userID)
		if u == nil {
			return l.ErrorUserNotFound
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

// UserByUsername returns user by name
func (us *UserStorage) UserByUsername(ctx context.Context, username string) (model.User, error) {
	var res model.User
	err := us.db.View(func(tx *bolt.Tx) error {
		unpb := tx.Bucket([]byte(UserByUsername))
		// we use username and password hash as a key
		key := username
		// get user ID from index
		userID := unpb.Get([]byte(key))
		if userID == nil {
			return l.ErrorUserNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))
		// get user by userID
		u := ub.Get(userID)
		if u == nil {
			return l.ErrorUserNotFound
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
	return res, nil
}

func (us *UserStorage) UserData(ctx context.Context, userID string, fields ...model.UserDataField) (model.UserData, error) {
	var res model.UserData
	err := us.db.View(func(tx *bolt.Tx) error {
		ub := tx.Bucket([]byte(UserBucket))
		// Get user by userID.
		u := ub.Get([]byte(userID))
		if u == nil {
			return l.ErrorUserNotFound
		}

		var err error
		res, err = model.UserDataFromJSON(u)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return model.UserData{}, err
	}
	return model.FilterUserDataFields(res, fields...), nil
}

func (us *UserStorage) AddUser(ctx context.Context, user model.User) (model.User, error) {
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

// UpdateUser(ctx context.Context, user User) (User, error)
// UpdateUserData(ctx context.Context, userID string, data UserData, fields ...UserDataField) (UserData, error)

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

// AddNewUser adds new user to the storage.
func (us *UserStorage) AddNewUser(user model.User, password string) (model.User, error) {
}

// AddUserWithFederatedID adds new user with social ID.
func (us *UserStorage) AddUserWithFederatedID(user model.User, provider string, federatedID, role string) (model.User, error) {
	// Using user name as a key. If there is no error, it means user already exists.
	if _, err := us.UserByFederatedID(provider, federatedID); err == nil {
		return model.User{}, model.ErrorUserExists
	}

	user.ID = xid.New().String()
	user.Active = true
	user.AccessRole = role
	user.AddFederatedId(provider, federatedID)

	return us.AddNewUser(user, "")
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
		FullName:   user.FullName,
		Scopes:     user.Scopes,
		Phone:      user.Phone,
		Email:      user.Email,
		AccessRole: role,
		Anonymous:  isAnonymous,
		TFAInfo:    user.TFAInfo,
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
			return l.ErrorUserNotFound
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
		return l.ErrorUserNotFound
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Pswd), []byte(password)); err != nil {
		// return this error to hide the existence of the user.
		return l.ErrorUserNotFound
	}
	return nil
}

// ResetUsername sets user username.
func (us *UserStorage) ResetUsername(id, username string) error {
	// TODO: implement
	return errors.New("ResetUsername is not implemented. ")
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
			} else if strings.Contains(strings.ToLower(string(user.PhoneNumber)), strings.ToLower(filterString)) {
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
func (us *UserStorage) ImportJSON(data []byte, clearOldData bool) error {
	if clearOldData {
		us.db.Update(func(tx *bolt.Tx) error {
			tx.DeleteBucket([]byte(UserBucket))
			tx.DeleteBucket([]byte(UserBySocialIDBucket))
			tx.DeleteBucket([]byte(UserByUsername))
			tx.DeleteBucket([]byte(UserByPhoneNumberBucket))
			tx.DeleteBucket([]byte(UserByEmailBucket))
			return nil
		})
		if err := us.createBuckets(); err != nil {
			return err
		}
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
	if err := CloseDB(us.db); err != nil {
		log.Printf("Error closing user storage: %s\n", err)
	}
}

// update all user mappings for all available buckets
// username -> userID
// email -> userID
// phone -> userID
// []federatedIDs --->> userID
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
	if user.PhoneNumber != "" {
		upnb := tx.Bucket([]byte(UserByPhoneNumberBucket))
		if err := upnb.Put([]byte(user.PhoneNumber), []byte(user.ID)); err != nil {
			return err
		}
	}

	// for _, fid := range user.FederatedIDs {
	// 	usib := tx.Bucket([]byte(UserBySocialIDBucket))
	// 	if err := usib.Put([]byte(fid), []byte(user.ID)); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
