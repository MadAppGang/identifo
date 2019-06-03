package boltdb

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/madappgang/identifo/model"
)

const (
	// VerificationCodesBucket is a bucket with verification codes.
	VerificationCodesBucket = "VerificationCodes"
)

// NewVerificationCodeStorage creates and inits BoltDB verification code storage.
func NewVerificationCodeStorage(db *bolt.DB) (model.VerificationCodeStorage, error) {
	vcs := &VerificationCodeStorage{db: db}

	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(VerificationCodesBucket)); err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return vcs, nil
}

// VerificationCodeStorage implements verification code storage interface.
type VerificationCodeStorage struct {
	db *bolt.DB
}

// IsVerificationCodeFound checks whether verification code can be found.
func (vcs *VerificationCodeStorage) IsVerificationCodeFound(phone, code string) (bool, error) {
	err := vcs.db.View(func(tx *bolt.Tx) error {
		vcb := tx.Bucket([]byte(VerificationCodesBucket))
		code := vcb.Get([]byte(phone))
		if code == nil {
			return model.ErrorNotFound
		}
		return nil
	})
	return err == nil, nil
}

// CreateVerificationCode inserts new verification code to the database.
func (vcs *VerificationCodeStorage) CreateVerificationCode(phone, code string) error {
	err := vcs.db.Update(func(tx *bolt.Tx) error {
		vcb := tx.Bucket([]byte(VerificationCodesBucket))
		if err := vcb.Delete([]byte(phone)); err != nil {
			return err
		}

		return vcb.Put([]byte(phone), []byte(code))
	})
	return err
}
