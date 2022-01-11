package ses

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/madappgang/identifo/v2/model"
)

// NewTransport creates AWS SES email service.
func NewTransport(ess model.SESEmailServiceSettings) (model.EmailTransport, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(ess.Region),
	},
	)
	if err != nil {
		return nil, err
	}

	return &transport{Sender: ess.Sender, service: ses.New(sess)}, nil
}

// EmailService sends email with Amazon Simple Email Service.
type transport struct {
	Sender  string
	service *ses.SES
}

// SendMessage sends email with plain text.
func (es *transport) SendMessage(subject, body, recipient string) error {
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
func (es *transport) SendHTML(subject, html, recipient string) error {
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
