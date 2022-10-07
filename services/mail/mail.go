package mail

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"path"
	"path/filepath"
	"text/template"
	"time"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/services/mail/mailgun"
	"github.com/madappgang/identifo/v2/services/mail/mock"
	"github.com/madappgang/identifo/v2/services/mail/ses"
	"github.com/madappgang/identifo/v2/storage"
)

const (
	DefaultEmailTemplatePath string = "./email_templates"
)

func NewService(ess model.EmailServiceSettings, fs fs.FS, updIntrv time.Duration, templatesPath string) (model.EmailService, error) {
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

	watchFiles := []string{}
	for _, f := range model.AllEmailTemplatesFileNames() {
		watchFiles = append(watchFiles, filepath.Join(templatesPath, f))
	}
	watcher := storage.NewFSWatcher(fs, watchFiles, updIntrv)

	return &EmailService{
		cache:         make(map[string]template.Template),
		transport:     t,
		fs:            fs,
		watcher:       watcher,
		templatesPath: templatesPath,
	}, nil
}

type EmailService struct {
	transport     model.EmailTransport
	fs            fs.FS
	cache         map[string]template.Template
	watcher       *storage.FSWatcher
	templatesPath string
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

	p := path.Join(es.templatesPath, emailType.FileName())
	// check template in cache
	template, ok := es.cache[p]

	// if no, let's try to load it and save to cache
	if !ok {
		data, err := fs.ReadFile(es.fs, p)
		if err != nil {
			return err
		}
		tmpl, err := template.New(p).Parse(string(data))
		if err != nil {
			return err
		}
		template = *tmpl
		es.cache[p] = template
	}

	// read template, parse it and send it with underlying service
	var tpl bytes.Buffer
	if err := template.Execute(&tpl, data); err != nil {
		return err
	}
	return es.SendHTML(subject, tpl.String(), recipient)
}

func (es *EmailService) Start() {
	if !es.watcher.IsWatching() {
		es.watcher.Watch()
		go es.watch()
	}
}

func (es *EmailService) Stop() {
	es.watcher.Stop()
}

func (es *EmailService) Transport() model.EmailTransport {
	return es.transport
}

func (es *EmailService) watch() {
	for {
		select {
		case files, ok := <-es.watcher.WatchChan():
			// the channel is closed
			if ok == false {
				return
			}
			es.cache = make(map[string]template.Template)
			log.Printf("email template changed, the email template cache has been invalidated: %v", files)
		}
	}
}
