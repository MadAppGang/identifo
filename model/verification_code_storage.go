package model

// VerificationCodeStorage stores verification codes linked to the phone number.
type VerificationCodeStorage interface {
	IsVerificationCodeFound(phone, code string) (bool, error)
	CreateVerificationCode(phone, code string) error
	Close()
}
