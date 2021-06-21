package config

import (
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
	"github.com/madappgang/identifo/services/mail"
	"github.com/madappgang/identifo/services/sms"
	"github.com/madappgang/identifo/storage"
)

// NewServer creates new server instance from ServerSettings
func NewServer(config model.ConfigurationStorage) (model.Server, error) {
	// read settings, if they empty or use cached version
	settings, err := config.LoadServerSettings(false)
	if err != nil {
		return nil, err
	}

	// Create all storages
	app, err := storage.NewAppStorage(settings.Storage.AppStorage)
	if err != nil {
		return nil, err
	}

	user, err := storage.NewUserStorage(settings.Storage.UserStorage)
	if err != nil {
		return nil, err
	}

	token, err := storage.NewTokenStorage(settings.Storage.TokenStorage)
	if err != nil {
		return nil, err
	}

	tokenBlacklist, err := storage.NewTokenBlacklistStorage(settings.Storage.TokenBlacklist)
	if err != nil {
		return nil, err
	}

	verification, err := storage.NewVerificationCodesStorage(settings.Storage.VerificationCodeStorage)
	if err != nil {
		return nil, err
	}

	invite, err := storage.NewInviteStorage(settings.Storage.InviteStorage)
	if err != nil {
		return nil, err
	}

	session, err := storage.NewSessionStorage(settings.SessionStorage)
	if err != nil {
		return nil, err
	}

	static, err := storage.NewStaticFileStorage(settings.StaticFilesStorage)
	if err != nil {
		return nil, err
	}

	// create 3rd party services

	sc := model.ServerStorageCollection{
		App:          app,
		User:         user,
		Token:        token,
		Blocklist:    tokenBlacklist,
		Invite:       invite,
		Verification: verification,
		Session:      session,
		Config:       config,
		Static:       static,
	}

	sms, err := sms.NewService(settings.ExternalServices.SMSService)
	if err != nil {
		return nil, err
	}

	email, err := mail.NewService(settings.ExternalServices.EmailService, static)
	if err != nil {
		return nil, err
	}

	srvs := model.ServerThirdPartyServices{
		SMS:   sms,
		Email: email,
	}

	server, err := server.NewServer(sc, srvs)
	if err != nil {
		return nil, err
	}

	return server, nil
}
