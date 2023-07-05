package model

// EmailService manages sending emails.
type EmailService interface {
	SendUserEmail(emailType EmailTemplateType, subfolder string, user User, data any) error
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
	Data any
}

type ResetEmailData struct {
	Token string
	URL   string
	Host  string
}
