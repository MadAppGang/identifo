package model

// VerificationCodeStorage is able to store verification codes linked to phone number.
type VerificationCodeStorage interface {
	// FindVerificationCode looking for verification code for specified phone number and removes it.
	// True value in response means that code was found.
	// False value in response means that code was not found.
	// Non nil error means that something went wrong
	FindVerificationCode(phone, code string) (bool, error)
	// CreateVerificationCode creates new verification code in the storage.
	CreateVerificationCode(phone, code string) error
}
