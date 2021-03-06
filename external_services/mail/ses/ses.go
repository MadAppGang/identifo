package ses

import (
	"bytes"
	"html/template"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/madappgang/identifo/model"
)

// NewEmailService creates new email service.
func NewEmailService(ess model.EmailServiceSettings, templater *model.EmailTemplater) (model.EmailService, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(ess.Region),
	},
	)
	if err != nil {
		return nil, err
	}

	return &EmailService{Sender: ess.Sender, service: ses.New(sess), tmpltr: templater}, nil
}

// EmailService sends email with Amazon Simple Email Service.
type EmailService struct {
	Sender  string
	service *ses.SES
	tmpltr  *model.EmailTemplater
}

// SendMessage sends email with plain text.
func (es *EmailService) SendMessage(subject, body, recipient string) error {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(es.Sender),
	}
	_, err := es.service.SendEmail(input)
	logAWSError(err)
	return err
}

// SendHTML sends email with html.
func (es *EmailService) SendHTML(subject, html, recipient string) error {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(html),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(es.Sender),
	}
	_, err := es.service.SendEmail(input)
	logAWSError(err)
	return err
}

// Templater returns email service templater.
func (es *EmailService) Templater() *model.EmailTemplater {
	return es.tmpltr
}

// SendTemplateEmail applies html template to the specified data and sends it in an email.
func (es *EmailService) SendTemplateEmail(subject, recipient string, template *template.Template, data interface{}) error {
	var tpl bytes.Buffer
	if err := template.Execute(&tpl, data); err != nil {
		return err
	}
	return es.SendHTML(subject, tpl.String(), recipient)
}

// SendResetEmail sends reset password emails.
func (es *EmailService) SendResetEmail(subject, recipient string, data model.ResetEmailData) error {
	return es.SendTemplateEmail(subject, recipient, es.tmpltr.ResetPasswordTemplate, data)
}

// SendInviteEmail sends invite email to the recipient.
func (es *EmailService) SendInviteEmail(subject, recipient string, data model.InviteEmailData) error {
	return es.SendTemplateEmail(subject, recipient, es.tmpltr.InviteTemplate, data)
}

// SendWelcomeEmail sends welcoming emails.
func (es *EmailService) SendWelcomeEmail(subject, recipient string, data model.WelcomeEmailData) error {
	return es.SendTemplateEmail(subject, recipient, es.tmpltr.WelcomeTemplate, data)
}

// SendVerifyEmail sends email address verification emails.
func (es *EmailService) SendVerifyEmail(subject, recipient string, data model.VerifyEmailData) error {
	return es.SendTemplateEmail(subject, recipient, es.tmpltr.VerifyTemplate, data)
}

// SendTFAEmail sends emails with one-time password.
func (es *EmailService) SendTFAEmail(subject, recipient string, data model.SendTFAEmailData) error {
	return es.SendTemplateEmail(subject, recipient, es.tmpltr.TFATemplate, data)
}

func logAWSError(err error) {
	if err == nil {
		return
	}

	aerr, ok := err.(awserr.Error)
	if !ok {
		log.Println("Could not cast the error to AWS error:", err)
		return
	}

	switch aerr.Code() {
	case ses.ErrCodeMessageRejected:
		log.Println(ses.ErrCodeMessageRejected, aerr.Error())
	case ses.ErrCodeMailFromDomainNotVerifiedException:
		log.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
	case ses.ErrCodeConfigurationSetDoesNotExistException:
		log.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
	default:
		log.Println(aerr.Error())
	}
}
