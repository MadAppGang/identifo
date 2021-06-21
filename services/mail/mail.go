package mail

import (
	"errors"
	"fmt"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/services/mail/mailgun"
	"github.com/madappgang/identifo/services/mail/mock"
	"github.com/madappgang/identifo/services/mail/ses"
)

func NewService(ess model.EmailServiceSettings, sfs model.StaticFilesStorage) (model.EmailService, error) {
	templater, err := model.NewEmailTemplater(sfs)
	if err != nil {
		return nil, err
	}
	if templater == nil {
		return nil, errors.New("email templater is nil")
	}

	switch ess.Type {
	case model.EmailServiceMailgun:
		return mailgun.NewEmailService(ess, templater), nil
	case model.EmailServiceAWS:
		return ses.NewEmailService(ess, templater)
	case model.EmailServiceMock:
		return mock.NewEmailService(), nil
	}
	return nil, fmt.Errorf("Email service of type '%s' is not supported", ess.Type)
}
