package model

// SMSService basically service to send SMS
type SMSService interface {
	SendSMS(recipient, message string) error
}
