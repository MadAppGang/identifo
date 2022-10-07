package config

import (
	"io/fs"
	"log"
	"time"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/server"
	"github.com/madappgang/identifo/v2/services/mail"
	"github.com/madappgang/identifo/v2/services/sms"
	"github.com/madappgang/identifo/v2/storage"

	jwt "github.com/madappgang/identifo/v2/jwt/service"
)

var adminPanelFSSettings = model.FileStorageSettings{
	Type: model.FileStorageTypeLocal,
	Local: model.FileStorageLocal{
		Path: "./static/admin_panel",
	},
}

var defaultLoginWebAppFSSettings = model.FileStorageSettings{
	Type: model.FileStorageTypeLocal,
	Local: model.FileStorageLocal{
		Path: "./static/web",
	},
}

var defaultEmailTemplateFSSettings = model.FileStorageSettings{
	Type: model.FileStorageTypeLocal,
	Local: model.FileStorageLocal{
		Path: "./static/email_templates",
	},
}

// NewServer creates new server instance from ServerSettings
func NewServer(config model.ConfigurationStorage, restartChan chan<- bool) (model.Server, error) {
	// read settings, if they empty or use cached version
	settings, err := config.LoadServerSettings(false)
	if err != nil {
		log.Println("Error on Load Server Settings")
		return nil, err
	}

	// Create all storages
	app, err := storage.NewAppStorage(settings.Storage.AppStorage)
	if err != nil {
		log.Println("Error on Create New AppStorage")
		return nil, err
	}

	user, err := storage.NewUserStorage(settings.Storage.UserStorage)
	if err != nil {
		log.Println("Error on Create New UserStorage")
		return nil, err
	}

	token, err := storage.NewTokenStorage(settings.Storage.TokenStorage)
	if err != nil {
		log.Println("Error on Create New TokenStorage")
		return nil, err
	}

	tokenBlacklist, err := storage.NewTokenBlacklistStorage(settings.Storage.TokenBlacklist)
	if err != nil {
		log.Println("Error on Create New TokenBlacklistStorage")
		return nil, err
	}

	verification, err := storage.NewVerificationCodesStorage(settings.Storage.VerificationCodeStorage)
	if err != nil {
		log.Println("Error on Create New VerificationCodesStorage")
		return nil, err
	}

	invite, err := storage.NewInviteStorage(settings.Storage.InviteStorage)
	if err != nil {
		log.Println("Error on Create New InviteStorage")
		return nil, err
	}

	session, err := storage.NewSessionStorage(settings.SessionStorage)
	if err != nil {
		log.Println("Error on Create New SessionStorage")
		return nil, err
	}

	key, err := storage.NewKeyStorage(settings.KeyStorage)
	if err != nil {
		log.Println("Error on Create New KeyStorage")
		return nil, err
	}

	// maybe just not serve login web app if type is none?
	lwas := settings.LoginWebApp
	if settings.LoginWebApp.Type == model.FileStorageTypeNone || settings.LoginWebApp.Type == model.FileStorageTypeDefault {
		// if not set, use default value
		lwas = defaultLoginWebAppFSSettings
	}
	loginFS, err := storage.NewFS(lwas)
	if err != nil {
		log.Println("Error on Create login fs storage")
		return nil, err
	}

	var adminPanelFS fs.FS
	if settings.AdminPanel.Enabled {
		adminPanelFS, err = storage.NewFS(adminPanelFSSettings)
		if err != nil {
			log.Println("Error on Create admin panel fs storage")
			return nil, err
		}
	}

	sc := model.ServerStorageCollection{
		App:          app,
		User:         user,
		Token:        token,
		Blocklist:    tokenBlacklist,
		Invite:       invite,
		Verification: verification,
		Session:      session,
		Config:       config,
		Key:          key,
		LoginAppFS:   loginFS,
		AdminPanelFS: adminPanelFS,
	}

	// create 3rd party services
	sms, err := sms.NewService(settings.Services.SMS)
	if err != nil {
		return nil, err
	}

	// maybe not use email templates if type is None?
	ets := settings.EmailTemplates
	if ets.Type == model.FileStorageTypeNone || ets.Type == model.FileStorageTypeDefault {
		ets = defaultEmailTemplateFSSettings
	}
	emailTemplateFS, err := storage.NewFS(ets)
	if err != nil {
		return nil, err
	}

	emailServiceSettings := settings.Services.Email
	// update templates every five minutes and look templates in a root folder of FS
	email, err := mail.NewService(emailServiceSettings, emailTemplateFS, time.Minute*5, "")
	if err != nil {
		return nil, err
	}

	tokenS, err := NewTokenService(settings.General, sc)
	if err != nil {
		return nil, err
	}

	sessionS := model.NewSessionManager(settings.SessionStorage.SessionDuration, session)

	srvs := model.ServerServices{
		SMS:     sms,
		Email:   email,
		Token:   tokenS,
		Session: sessionS,
	}

	server, err := server.NewServer(sc, srvs, restartChan)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func NewTokenService(settings model.GeneralServerSettings, storages model.ServerStorageCollection) (model.TokenService, error) {
	key, err := storages.Key.LoadPrivateKey()
	if err != nil {
		return nil, err
	}

	tokenService, err := jwt.NewJWTokenService(
		key,
		settings.Issuer,
		storages.Token,
		storages.App,
		storages.User,
	)
	return tokenService, err
}
