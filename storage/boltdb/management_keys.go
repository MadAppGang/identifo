package boltdb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madappgang/identifo/v2/model"
	bolt "go.etcd.io/bbolt"
)

const (
	// ManagementKeysBucket is a name for bucket with keys.
	ManagementKeysBucket = "ManagementKeys"
)

// NewTokenStorage creates a BoltDB token storage.
func NewManagementKeysStorage(settings model.BoltDBDatabaseSettings) (model.ManagementKeysStorage, error) {
	if len(settings.Path) == 0 {
		return nil, ErrorEmptyDatabasePath
	}

	// init database
	db, err := InitDB(settings.Path)
	if err != nil {
		return nil, err
	}

	ts := &ManagementKeysStorage{db: db}
	// Ensure that we have needed bucket in the database.
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(ManagementKeysBucket)); err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return ts, nil
}

// TokenStorage is a BoltDB token storage.
type ManagementKeysStorage struct {
	db *bolt.DB
}

func (ms *ManagementKeysStorage) GetKey(ctx context.Context, id string) (model.ManagementKey, error) {
	var res model.ManagementKey
	err := ms.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ManagementKeysBucket))
		u := b.Get([]byte(id))
		if u == nil {
			return model.ErrUserNotFound
		}

		return json.Unmarshal(u, &res)
	})
	return res, err
}

func (ms *ManagementKeysStorage) CreateKey(ctx context.Context, name string, scopes []string) (model.ManagementKey, error) {
	key := model.ManagementKey{
		Name:      name,
		Scopes:    scopes,
		ID:        uuid.New().String(),
		Active:    true,
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
	}
	err := ms.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ManagementKeysBucket))

		data, err := json.Marshal(key)
		if err != nil {
			return err
		}

		return b.Put([]byte(key.ID), data)
	})
	return key, err
}

func (ms *ManagementKeysStorage) DisableKey(ctx context.Context, id string) (model.ManagementKey, error) {
	key, err := ms.GetKey(ctx, id)
	if err != nil {
		return key, err
	}

	key.Active = false
	err = ms.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ManagementKeysBucket))

		data, err := json.Marshal(key)
		if err != nil {
			return err
		}

		return b.Put([]byte(key.ID), data)
	})
	return key, err
}

func (ms *ManagementKeysStorage) RenameKey(ctx context.Context, id, name string) (model.ManagementKey, error) {
	key, err := ms.GetKey(ctx, id)
	if err != nil {
		return key, err
	}

	key.Name = name
	err = ms.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ManagementKeysBucket))

		data, err := json.Marshal(key)
		if err != nil {
			return err
		}

		return b.Put([]byte(key.ID), data)
	})
	return key, err
}

func (ms *ManagementKeysStorage) ChangeScopesForKey(ctx context.Context, id string, scopes []string) (model.ManagementKey, error) {
	key, err := ms.GetKey(ctx, id)
	if err != nil {
		return key, err
	}

	key.Scopes = scopes
	err = ms.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ManagementKeysBucket))

		data, err := json.Marshal(key)
		if err != nil {
			return err
		}

		return b.Put([]byte(key.ID), data)
	})
	return key, err
}

func (ms *ManagementKeysStorage) UseKey(ctx context.Context, id string) (model.ManagementKey, error) {
	key, err := ms.GetKey(ctx, id)
	if err != nil {
		return key, err
	}

	key.LastUsed = time.Now()
	err = ms.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ManagementKeysBucket))

		data, err := json.Marshal(key)
		if err != nil {
			return err
		}

		return b.Put([]byte(key.ID), data)
	})
	return key, err
}

func (ms *ManagementKeysStorage) GeyAllKeys(ctx context.Context) ([]model.ManagementKey, error) {
	keys := []model.ManagementKey{}

	err := ms.db.View(func(tx *bolt.Tx) error {
		ub := tx.Bucket([]byte(ManagementKeysBucket))

		if iterErr := ub.ForEach(func(k, d []byte) error {
			var key model.ManagementKey
			err := json.Unmarshal(d, &key)
			if err != nil {
				return err
			}

			keys = append(keys, key)
			return nil
		}); iterErr != nil {
			return iterErr
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return keys, nil
}
