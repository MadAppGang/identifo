package boltdb

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	bolt "go.etcd.io/bbolt"
)

const (
	UserBucket              = "Users"
	UserDataBucket          = "UserData"
	UserByFederatedIDBucket = "UserBySocialID"
	UserByUsername          = "UserByUsername"
	UserByPhoneNumberBucket = "UserByPhoneNumber"
	UserByEmailBucket       = "UserByEmail"
)

var _ model.UserStorage = &UserStorage{}

func NewUserStorage(settings model.BoltDBDatabaseSettings) (*UserStorage, error) {
	if len(settings.Path) == 0 {
		return nil, ErrorEmptyDatabasePath
	}

	db, err := InitDB(settings.Path)
	if err != nil {
		return nil, err
	}

	us := &UserStorage{db: db, path: settings.Path}

	if err := us.createBuckets(); err != nil {
		return nil, fmt.Errorf("failed to create buckets: %w", err)
	}

	return us, nil

}

type UserStorage struct {
	db   *bolt.DB
	path string
}

func (us *UserStorage) createBuckets() error {
	return us.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(UserBucket)); err != nil {
			return fmt.Errorf("failed to create user bucket: %w", err)
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(UserDataBucket)); err != nil {
			return fmt.Errorf("failed to create user data bucket: %w", err)
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(UserByFederatedIDBucket)); err != nil {
			return fmt.Errorf("failed to create federated id bucket: %w", err)
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(UserByUsername)); err != nil {
			return fmt.Errorf("failed to create username bucket: %w", err)
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(UserByPhoneNumberBucket)); err != nil {
			return fmt.Errorf("failed to create phone number bucket: %w", err)
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(UserByEmailBucket)); err != nil {
			return fmt.Errorf("failed to create email bucket: %w", err)
		}
		return nil
	})
}

func (us *UserStorage) Close(context.Context) error {
	if err := CloseDB(us.db); err != nil {
		return fmt.Errorf("failed to close user storage: %w")
	}

	return nil
}

func (us *UserStorage) Connect(context.Context) error {
	db, err := InitDB(us.path)
	if err != nil {
		return fmt.Errorf("failed to connect to db")
	}

	us.db = db

	return nil
}

func (us *UserStorage) Ready(context.Context) error {
	if err := Ready(us.db); err != nil {
		return fmt.Errorf("not ready: %w", err)
	}

	return nil
}

func (us *UserStorage) ImportJSON(data []byte, clearOldData bool) error {
	if clearOldData {
		us.db.Update(func(tx *bolt.Tx) error {
			tx.DeleteBucket([]byte(UserBucket))
			tx.DeleteBucket([]byte(UserByFederatedIDBucket))
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

	return nil
}

func (us *UserStorage) UserByID(ctx context.Context, id string) (model.User, error) {
	var res model.User

	err := us.db.View(func(tx *bolt.Tx) error {
		ub := tx.Bucket([]byte(UserBucket))

		u := ub.Get([]byte(id))
		if u == nil {
			return l.ErrorUserNotFound
		}

		var err error

		res, err = model.UserFromJSON(u)
		if err != nil {
			return fmt.Errorf("failed to convert user to json: %w", err)
		}

		return nil
	})
	if err != nil {
		return model.User{}, fmt.Errorf("failed to view into db: %w", err)
	}
	return res, nil
}

func (us *UserStorage) GetUserByFederatedID(ctx context.Context, idType model.UserFederatedType, userIdentityTypeOther, externalID string) (model.User, error) {
	var res model.User
	sid := string(idType) + ":" + userIdentityTypeOther + ":" + externalID

	err := us.db.View(func(tx *bolt.Tx) error {
		usib := tx.Bucket([]byte(UserByFederatedIDBucket))

		userID := usib.Get([]byte(sid))
		if userID == nil {
			return l.ErrorUserNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))

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

func (us *UserStorage) UserData(ctx context.Context, userID string, fields ...model.UserDataField) (model.UserData, error) {
	var res model.UserData

	err := us.db.View(func(tx *bolt.Tx) error {
		ub := tx.Bucket([]byte(UserBucket))

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

func (us *UserStorage) UserBySecondaryID(ctx context.Context, idt model.AuthIdentityType, id string) (model.User, error) {
	switch idt {
	case model.AuthIdentityTypePhone:
		user, err := us.userByPhone(ctx, id)
		if err != nil {
			return model.User{}, fmt.Errorf("failed to find user by phone: %w", err)
		}

		return user, nil
	case model.AuthIdentityTypeEmail:
		user, err := us.userByEmail(ctx, id)
		if err != nil {
			return model.User{}, fmt.Errorf("failed to find user by email: %w", err)
		}

		return user, nil
	case model.AuthIdentityTypeUsername:
		user, err := us.userByUsername(ctx, id)
		if err != nil {
			return model.User{}, fmt.Errorf("failed to find user by name: %w", err)
		}

		return user, nil
	}

	return model.User{}, fmt.Errorf("invalid id type")
}

func (us *UserStorage) userByPhone(ctx context.Context, Phone string) (model.User, error) {
	var res model.User
	err := us.db.View(func(tx *bolt.Tx) error {
		upnb := tx.Bucket([]byte(UserByPhoneNumberBucket))

		userID := upnb.Get([]byte(Phone))
		if userID == nil {
			return l.ErrorUserNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))

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

func (us *UserStorage) userByEmail(ctx context.Context, email string) (model.User, error) {
	var res model.User
	err := us.db.View(func(tx *bolt.Tx) error {
		ueb := tx.Bucket([]byte(UserByEmailBucket))

		userID := ueb.Get([]byte(email))
		if userID == nil {
			return l.ErrorUserNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))

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

func (us *UserStorage) userByUsername(ctx context.Context, username string) (model.User, error) {
	var res model.User
	err := us.db.View(func(tx *bolt.Tx) error {
		unpb := tx.Bucket([]byte(UserByUsername))
		key := username

		userID := unpb.Get([]byte(key))
		if userID == nil {
			return l.ErrorUserNotFound
		}

		ub := tx.Bucket([]byte(UserBucket))

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
