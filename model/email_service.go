package model

// EmailService manage sending email
type EmailService interface {
	SendMessage(subject, body, recipient string) (string, string, error)
	SendHTML(subject, html, recipient string) (string, string, error)
}
