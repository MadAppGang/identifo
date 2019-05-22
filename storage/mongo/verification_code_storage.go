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
)

// NewUserStorage creates and inits MongoDB user storage.
func NewVerificationCodeStorage(db *DB) (model.VerificationCodeStorage, error) {
	vcs := &VerificationCodeStorage{db: db}

	s := vcs.db.Session(VerificationCodesCollection)
	defer s.Close()

	if err := s.C.EnsureIndex(mgo.Index{
		Key:         []string{"phone"},
		Unique:      true,
		ExpireAfter: verificationCodesExpirationTime,
	}); err != nil {
		return nil, err
	}

	if err := s.C.EnsureIndex(mgo.Index{
		Key:    []string{"code"},
		Unique: true,
	}); err != nil {
		return nil, err
	}

	return vcs, nil
}

// UserStorage implements user storage interface.
type VerificationCodeStorage struct {
	db *DB
}

func (vcs *VerificationCodeStorage) FindVerificationCode(phone, code string) (bool, error) {
	s := vcs.db.Session(VerificationCodesCollection)
	defer s.Close()

	_, err := s.C.Find(bson.M{"phone": phone, "code": code}).Apply(mgo.Change{Remove: true}, nil)
	if err != nil {
		if err == mgo.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (vcs *VerificationCodeStorage) CreateVerificationCode(phone, code string) error {
	s := vcs.db.Session(VerificationCodesCollection)
	defer s.Close()

	err := s.C.Insert(bson.M{"phone": phone, "code": code})
	return err
}
