package config

import (
	"fmt"
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
	var errs []error
	// read settings, if they empty or use cached version
	settings, e := config.LoadServerSettings(true)
	if len(e) > 0 {
		log.Printf("Error on Load Server Settings: %+v\n", e)
		settings = model.DefaultServerSettings
		errs = e
	}

	// helper function get settings from default or override
	dbSettings := func(settings, def model.DatabaseSettings) model.DatabaseSettings {
		if settings.Type == model.DBTypeDefault {
			return def
		}
		return settings
	}

	// Create all storages
	app, err := storage.NewAppStorage(dbSettings(settings.Storage.AppStorage, settings.Storage.DefaultStorage))
	if err != nil {
		log.Printf("Error on Create New AppStorage %v", err)
		errs = append(errs, fmt.Errorf("error creating app storage: %v", err))
	}

	user, err := storage.NewUserStorage(dbSettings(settings.Storage.UserStorage, settings.Storage.DefaultStorage))
	if err != nil {
		log.Printf("Error on Create New user storage %v", err)
		errs = append(errs, fmt.Errorf("error creating user storage: %v", err))
	}

	token, err := storage.NewTokenStorage(dbSettings(settings.Storage.TokenStorage, settings.Storage.DefaultStorage))
	if err != nil {
		log.Printf("Error on Create New token storage %v", err)
		errs = append(errs, fmt.Errorf("error creating token storage: %v", err))
	}

	tokenBlacklist, err := storage.NewTokenBlacklistStorage(dbSettings(settings.Storage.TokenBlacklist, settings.Storage.DefaultStorage))
	if err != nil {
		log.Printf("Error on Create New blacklist storage %v", err)
		errs = append(errs, fmt.Errorf("error creating blacklist storage: %v", err))
	}

	verification, err := storage.NewVerificationCodesStorage(dbSettings(settings.Storage.VerificationCodeStorage, settings.Storage.DefaultStorage))
	if err != nil {
		log.Printf("Error on Create New verification codes storage %v", err)
		errs = append(errs, fmt.Errorf("error creating verification codes storage: %v", err))
	}

	invite, err := storage.NewInviteStorage(dbSettings(settings.Storage.InviteStorage, settings.Storage.DefaultStorage))
	if err != nil {
		log.Printf("Error on Create New invite storage %v", err)
		errs = append(errs, fmt.Errorf("error creating invite storage: %v", err))
	}

	managementKeys, err := storage.NewManagementKeys(dbSettings(settings.Storage.ManagementKeysStorage, settings.Storage.DefaultStorage))
	if err != nil {
		log.Printf("Error on Create New management keys storage %v", err)
		errs = append(errs, fmt.Errorf("error creating management keys storage: %v", err))
	}

	session, err := storage.NewSessionStorage(settings.SessionStorage)
	if err != nil {
		log.Printf("Error on Create New session storage %v", err)
		errs = append(errs, fmt.Errorf("error creating session storage: %v", err))
	}

	key, err := storage.NewKeyStorage(settings.KeyStorage)
	if err != nil {
		log.Printf("Error on Create New key storage %v", err)
		errs = append(errs, fmt.Errorf("error creating key storage: %v", err))
	}

	// maybe just not serve login web app if type is none?
	lwas := settings.LoginWebApp
	if settings.LoginWebApp.Type == model.FileStorageTypeNone || settings.LoginWebApp.Type == model.FileStorageTypeDefault {
		// if not set, use default value
		lwas = defaultLoginWebAppFSSettings
	}
	loginFS, err := storage.NewFS(lwas)
	if err != nil {
		log.Printf("Error creating login fs storage %v", err)
		errs = append(errs, fmt.Errorf("error creating login fs storage: %v", err))
	}

	var adminPanelFS fs.FS
	if settings.AdminPanel.Enabled {
		adminPanelFS, err = storage.NewFS(adminPanelFSSettings)
		if err != nil {
			log.Printf("Error creating admin panel fs storage %v", err)
			errs = append(errs, fmt.Errorf("error creating admin panel fs storage: %v", err))
		}
	}

	sc := model.ServerStorageCollection{
		App:           app,
		User:          user,
		Token:         token,
		Blocklist:     tokenBlacklist,
		Invite:        invite,
		Verification:  verification,
		Session:       session,
		Config:        config,
		Key:           key,
		ManagementKey: managementKeys,
		LoginAppFS:    loginFS,
		AdminPanelFS:  adminPanelFS,
	}

	// create 3rd party services
	sms, err := sms.NewService(settings.Services.SMS)
	if err != nil {
		log.Printf("Error creating SMS service %v", err)
		errs = append(errs, fmt.Errorf("error creating SMS service: %v", err))
	}

	// maybe not use email templates if type is None?
	ets := settings.EmailTemplates
	if ets.Type == model.FileStorageTypeNone || ets.Type == model.FileStorageTypeDefault {
		ets = defaultEmailTemplateFSSettings
	}
	emailTemplateFS, err := storage.NewFS(ets)
	if err != nil {
		log.Printf("Error creating email template filesystem %v", err)
		errs = append(errs, fmt.Errorf("error creating email template filesystem: %v", err))
	}

	emailServiceSettings := settings.Services.Email
	// update templates every five minutes and look templates in a root folder of FS
	email, err := mail.NewService(emailServiceSettings, emailTemplateFS, time.Minute*5, "")
	if err != nil {
		log.Printf("Error creating email service %v", err)
		errs = append(errs, fmt.Errorf("error creating email service: %v", err))
	}

	tokenS, err := NewTokenService(settings.General, sc)
	if err != nil {
		log.Printf("Error creating token service %v", err)
		errs = append(errs, fmt.Errorf("error creating token service: %v", err))
	}

	sessionS := model.NewSessionManager(settings.SessionStorage.SessionDuration, session)

	srvs := model.ServerServices{
		SMS:     sms,
		Email:   email,
		Token:   tokenS,
		Session: sessionS,
	}

	server, err := server.NewServer(sc, srvs, errs, restartChan)
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
