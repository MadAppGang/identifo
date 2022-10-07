package model

// EmailService manages sending emails.
type EmailService interface {
	SendTemplateEmail(emailType EmailTemplateType, subfolder, subject, recipient string, data EmailData) error
	Transport() EmailTransport
	Start()
	Stop()
}

type EmailTransport interface {
	SendMessage(subject, body, recipient string) error
	SendHTML(subject, html, recipient string) error
}

type EmailData struct {
	User User
	Data interface{}
}
