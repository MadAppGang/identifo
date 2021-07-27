package mem

import (
	"github.com/madappgang/identifo/model"
)

// NewVerificationCodeStorage creates and inits in-memory verification code storage.
func NewVerificationCodeStorage() (model.VerificationCodeStorage, error) {
	return &VerificationCodeStorage{}, nil
}

// VerificationCodeStorage implements verification code storage interface.
type VerificationCodeStorage struct{}

// IsVerificationCodeFound is always optimistic.
func (vcs *VerificationCodeStorage) IsVerificationCodeFound(phone, code string) (bool, error) {
	return true, nil
}

// CreateVerificationCode is always optimistic.
func (vcs *VerificationCodeStorage) CreateVerificationCode(phone, code string) error {
	return nil
}

// Close does nothing here.
func (vcs *VerificationCodeStorage) Close() {}
