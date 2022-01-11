package mail

import (
	"bytes"
	"fmt"
	"io/fs"
	"text/template"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/services/mail/mailgun"
	"github.com/madappgang/identifo/v2/services/mail/mock"
	"github.com/madappgang/identifo/v2/services/mail/ses"
)

const (
	DefaultEmailTemplatePath string = "./email_templates"
)

func NewService(ess model.EmailServiceSettings, fs fs.FS) (model.EmailService, error) {
	var t model.EmailTransport

	switch ess.Type {
	case model.EmailServiceMailgun:
		t = mailgun.NewTransport(ess.Mailgun)
	case model.EmailServiceAWS:
		tt, err := ses.NewTransport(ess.SES)
		if err != nil {
			return nil, err
		}
		t = tt
	case model.EmailServiceMock:
		t = mock.NewTransport()
	default:
		return nil, fmt.Errorf("Email service of type '%s' is not supported", ess.Type)
	}

	return &EmailService{
		cache:     make(map[string]template.Template),
		transport: t,
		fs:        fs,
	}, nil
}

type EmailService struct {
	transport model.EmailTransport
	fs        fs.FS
	cache     map[string]template.Template
}

// proxy to underlying service
func (es *EmailService) SendMessage(subject, body, recipient string) error {
	return es.transport.SendMessage(subject, body, recipient)
}

func (es *EmailService) SendHTML(subject, html, recipient string) error {
	return es.transport.SendHTML(subject, html, recipient)
}

func (es *EmailService) SendTemplateEmail(emailType model.EmailTemplateType, subject string, recipient string, data model.EmailData) error {
	var template template.Template

	// check template in cache
	template, ok := es.cache[emailType.String()]

	// if no, let's try to load it and save to cache
	if !ok {
		data, err := fs.ReadFile(es.fs, emailType.FileName())
		if err != nil {
			return err
		}
		tmpl, err := template.New(emailType.String()).Parse(string(data))
		if err != nil {
			return err
		}
		template = *tmpl
		es.cache[emailType.String()] = template
	}

	// read template, parse it and send it with underlying service
	var tpl bytes.Buffer
	if err := template.Execute(&tpl, data); err != nil {
		return err
	}
	return es.SendHTML(subject, tpl.String(), recipient)
}
