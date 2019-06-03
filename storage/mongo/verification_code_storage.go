package mongo

import (
	"time"

	"github.com/madappgang/identifo/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// VerificationCodesCollection is a collection name for verification codes.
	VerificationCodesCollection = "VerificationCodes"

	// verificationCodesExpirationTime specifies time before deleting records.
	verificationCodesExpirationTime = 5 * time.Minute

	phoneField     = "phone"
	codeField      = "code"
	createdAtField = "createdAt"
)

// NewVerificationCodeStorage creates and inits MongoDB verification code storage.
func NewVerificationCodeStorage(db *DB) (model.VerificationCodeStorage, error) {
	vcs := &VerificationCodeStorage{db: db}

	s := vcs.db.Session(VerificationCodesCollection)
	defer s.Close()

	if err := s.EnsureIndex(mgo.Index{
		Key:    []string{phoneField},
		Unique: true,
	}); err != nil {
		return nil, err
	}

	if err := s.EnsureIndex(mgo.Index{
		Key:    []string{codeField},
		Unique: true,
	}); err != nil {
		return nil, err
	}

	if err := s.EnsureIndex(mgo.Index{
		Key:         []string{createdAtField},
		ExpireAfter: verificationCodesExpirationTime,
	}); err != nil {
		return nil, err
	}
	return vcs, nil
}

// VerificationCodeStorage implements verification code storage interface.
type VerificationCodeStorage struct {
	db *DB
}

// IsVerificationCodeFound checks whether verification code can be found.
func (vcs *VerificationCodeStorage) IsVerificationCodeFound(phone, code string) (bool, error) {
	s := vcs.db.Session(VerificationCodesCollection)
	defer s.Close()

	_, err := s.C.Find(bson.M{phoneField: phone, codeField: code}).Apply(mgo.Change{Remove: true}, nil)
	if err != nil {
		if err == mgo.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CreateVerificationCode inserts new verification code to the database.
func (vcs *VerificationCodeStorage) CreateVerificationCode(phone, code string) error {
	s := vcs.db.Session(VerificationCodesCollection)
	defer s.Close()

	if _, err := s.C.RemoveAll(bson.M{phoneField: phone}); err != nil {
		return err
	}

	err := s.C.Insert(bson.M{phoneField: phone, codeField: code, createdAtField: time.Now()})
	return err
}
