package mail

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"path"
	"sync"
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

	watcher := storage.NewFSWatcher(fs, []string{}, updIntrv)

	return &EmailService{
		cache:         sync.Map{},
		transport:     t,
		fs:            fs,
		watcher:       watcher,
		templatesPath: templatesPath,
	}, nil
}

type EmailService struct {
	transport     model.EmailTransport
	fs            fs.FS
	cache         sync.Map
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

func (es *EmailService) SendTemplateEmail(emailType model.EmailTemplateType, subfolder, subject string, recipient string, data model.EmailData) error {
	p := path.Join(es.templatesPath, subfolder, emailType.FileName())
	// check template in cache
	var tt template.Template
	tpll, ok := es.cache.Load(p)
	if ok {
		tt = tpll.(template.Template)
	} else {
		data, err := fs.ReadFile(es.fs, p)
		if err != nil {
			return err
		}
		tmpl, err := template.New(p).Parse(string(data))
		if err != nil {
			return err
		}
		tt = *tmpl
		es.cache.Store(p, tt)
		es.watcher.AppendForWatching(p) // add for watching to invalidate cache if we need
	}

	// read template, parse it and send it with underlying service
	var tpl bytes.Buffer
	if err := tt.Execute(&tpl, data); err != nil {
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
	for files := range es.watcher.WatchChan() {
		for _, f := range files {
			es.cache.Delete(f)
			log.Printf("email template changed, the email template cache has been invalidated: %v", f)
		}
	}
}
