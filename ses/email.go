package ses

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/madappgang/identifo/model"
)

const (
	SESRegionKey = "SES_REGION"
	SESSenderKey = "SES_SENDER"
)

//NewEmailService creates new email service
func NewEmailService(sender, region string, templater *model.EmailTemplater) (model.EmailService, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return nil, err
	}
	es := EmailService{}
	es.Sender = sender
	es.service = ses.New(sess)
	if templater == nil {
		es.tmpltr, err = model.DefaultEmailTemplater()
		if err != nil {
			return nil, err
		}
	} else {
		es.tmpltr = templater
	}
	return &es, nil
}

//NewEmailServiceFromEnv creates new email service
func NewEmailServiceFromEnv(templater *model.EmailTemplater) (model.EmailService, error) {
	region := os.Getenv(SESRegionKey)
	sender := os.Getenv(SESSenderKey)
	return NewEmailService(sender, region, templater)
}

//EmailService sends email with Amazon Simple Email Service
type EmailService struct {
	Sender  string
	service *ses.SES
	tmpltr  *model.EmailTemplater
}

//SendMessage send email with plain text
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

//SendHTML send email with html text
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

//Templater returns default templater
func (es *EmailService) Templater() *model.EmailTemplater {
	return es.tmpltr
}

//SendTemplateEmail render data to html template and send it in email
func (es *EmailService) SendTemplateEmail(subject, recipient string, template *template.Template, data interface{}) error {
	var tpl bytes.Buffer
	if err := template.Execute(&tpl, data); err != nil {
		return err
	}
	return es.SendHTML(subject, tpl.String(), recipient)
}

//SendResetEmail sends reset passwords email
func (es *EmailService) SendResetEmail(subject, recipient string, data interface{}) error {

	return es.SendTemplateEmail(subject, recipient, es.tmpltr.ResetPasswordTemplate, data)
}

//SendWelcomeEmail sends welcome email
func (es *EmailService) SendWelcomeEmail(subject, recipient string, data interface{}) error {
	return es.SendTemplateEmail(subject, recipient, es.tmpltr.WelcomeTemplate, data)
}

//SendVerifyEmail sends verify email address email
func (es *EmailService) SendVerifyEmail(subject, recipient string, data interface{}) error {
	return es.SendTemplateEmail(subject, recipient, es.tmpltr.VerifyEmailTemplate, data)
}

func logAWSError(err error) {
	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}
}
