package mail

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"path"
	"sync"
	"text/template"
	"time"

	"github.com/madappgang/identifo/v2/l"
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

type templatePair struct {
	body    *template.Template
	subject *template.Template
}

func (es *EmailService) SendUserEmail(emailType model.EmailTemplateType, subfolder string, user model.User, data any) error {
	p := path.Join(es.templatesPath, subfolder, emailType.FileName())
	// trying to load user localized version
	if len(user.Locale) > 0 {
		pl := path.Join(es.templatesPath, subfolder, emailType.FileNameWithLocale(user.Locale))
		fi, err := fs.Stat(es.fs, pl)
		if err == nil && fi.Mode().IsRegular() {
			p = pl
		}
	}

	tp := templatePair{}
	tpll, ok := es.cache.Load(p)
	if ok {
		tp = tpll.(templatePair)
	} else {
		data, err := fs.ReadFile(es.fs, p)
		if err != nil {
			return err
		}
		// extract subject and body from template
		subjectText, bodyText, err := ExtractSubjectAndBody(data)
		if err != nil {
			if errors.Is(err, ErrorNoSubject) {
				return l.LocalizedError{
					ErrID:   l.ErrorServiceEmailTemplateMissingSubject,
					Details: []any{emailType},
				}
			} else if errors.Is(err, ErrorNoBody) {
				return l.LocalizedError{
					ErrID:   l.ErrorServiceEmailTemplateMissingBody,
					Details: []any{emailType},
				}
			} else {
				return err
			}
		}
		tp.body, err = template.New(p).Parse(string(bodyText))
		if err != nil {
			return err
		}
		tp.subject, err = template.New(p).Parse(string(subjectText))
		if err != nil {
			return err
		}

		es.cache.Store(p, tp)
		es.watcher.AppendForWatching(p) // add for watching to invalidate cache if we need
	}

	// read template, parse it and send it with underlying service
	var subject bytes.Buffer
	var body bytes.Buffer
	d := map[string]any{
		"Data": data,
		"User": user,
	}
	if err := tp.subject.Execute(&subject, d); err != nil {
		return err
	}
	if err := tp.body.Execute(&body, d); err != nil {
		return err
	}

	return es.SendHTML(subject.String(), body.String(), user.Email)
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
			for _, f := range files {
				es.cache.Delete(f)
				log.Printf("email template changed, the email template cache has been invalidated: %v", f)

			}
		}
	}
}
