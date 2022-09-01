package model

// EmailService manages sending emails.
type EmailService interface {
	SendTemplateEmail(emailType EmailTemplateType, path string, subject string, recipient string, data EmailData) error
}

type EmailTransport interface {
	SendMessage(subject, body, recipient string) error
	SendHTML(subject, html, recipient string) error
}

type EmailData struct {
	User User
	Data interface{}
}
