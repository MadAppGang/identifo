package ses

import (
	"bytes"
	"html/template"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/madappgang/identifo/model"
)

const (
	// SESRegionKey is a region setting for SES.
	SESRegionKey = "SES_REGION"
	// SESSenderKey is a sender key for SES.
	SESSenderKey = "SES_SENDER"
)

// NewEmailService creates new email service.
func NewEmailService(sender, region string, templater *model.EmailTemplater) (model.EmailService, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return nil, err
	}

	return &EmailService{Sender: sender, service: ses.New(sess), tmpltr: templater}, nil
}

// NewEmailServiceFromEnv creates new email service getting settings from env variables.
func NewEmailServiceFromEnv(templater *model.EmailTemplater) (model.EmailService, error) {
	region := os.Getenv(SESRegionKey)
	sender := os.Getenv(SESSenderKey)
	return NewEmailService(sender, region, templater)
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
func (es *EmailService) SendResetEmail(subject, recipient string, data interface{}) error {

	return es.SendTemplateEmail(subject, recipient, es.tmpltr.ResetPasswordTemplate, data)
}

// SendWelcomeEmail sends welcoming emails.
func (es *EmailService) SendWelcomeEmail(subject, recipient string, data interface{}) error {
	return es.SendTemplateEmail(subject, recipient, es.tmpltr.WelcomeTemplate, data)
}

// SendVerifyEmail sends email address verification emails.
func (es *EmailService) SendVerifyEmail(subject, recipient string, data interface{}) error {
	return es.SendTemplateEmail(subject, recipient, es.tmpltr.VerifyEmailTemplate, data)
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
