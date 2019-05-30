package model

// SMSService is an SMS sending service.
type SMSService interface {
	SendSMS(recipient, message string) error
}
