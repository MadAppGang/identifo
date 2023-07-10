package model

type SMSMessageType string

const (
	SMSMessageTypeOTPMagicLink SMSMessageType = "otp_magic_link"
	SMSMessageTypeOTPCode      SMSMessageType = "otp_code"
)

// SMSService is an SMS sending service.
type SMSService interface {
	SendSMS(recipient, message string) error
}
