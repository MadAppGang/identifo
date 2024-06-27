package config

import (
	"fmt"
	"io/fs"
	"log/slog"
	"time"

	imp "github.com/madappgang/identifo/v2/impersonation/local"
	impPlugin "github.com/madappgang/identifo/v2/impersonation/plugin"
	"github.com/madappgang/identifo/v2/logging"
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
	defaultLogger := logging.NewDefaultLogger().With(logging.FieldComponent, "configurator")

	var errs []error
	// read settings, if they empty or use cached version
	settings, e := config.LoadServerSettings(true)
	if len(e) > 0 {
		defaultLogger.Error("Error on Load Server Settings", logging.FieldErrors, logging.LogErrors(e))
		settings = model.DefaultServerSettings
		errs = e
	}

	baseLogger := logging.NewLogger(
		settings.Logger.Format,
		settings.Logger.Common.Level)

	logger := baseLogger.With(logging.FieldComponent, "configurator")

	// helper function get settings from default or override
	dbSettings := func(s model.DatabaseSettings) model.DatabaseSettings {
		if s.Type == model.DBTypeDefault {
			return settings.Storage.DefaultStorage
		}
		return s
	}

	// Create all storages
	app, err := storage.NewAppStorage(baseLogger, dbSettings(settings.Storage.AppStorage))
	if err != nil {
		logger.Error("Error on Create New AppStorage", logging.FieldError, err)
		errs = append(errs, fmt.Errorf("error creating app storage: %v", err))
	}

	user, err := storage.NewUserStorage(baseLogger, dbSettings(settings.Storage.UserStorage))
	if err != nil {
		logger.Error("Error on Create New user storage", logging.FieldError, err)
		errs = append(errs, fmt.Errorf("error creating user storage: %v", err))
	}

	token, err := storage.NewTokenStorage(baseLogger, dbSettings(settings.Storage.TokenStorage))
	if err != nil {
		logger.Error("Error on Create New token storage", logging.FieldError, err)
		errs = append(errs, fmt.Errorf("error creating token storage: %v", err))
	}

	tokenBlacklist, err := storage.NewTokenBlacklistStorage(baseLogger, dbSettings(settings.Storage.TokenBlacklist))
	if err != nil {
		logger.Error("Error on Create New blacklist storage", logging.FieldError, err)
		errs = append(errs, fmt.Errorf("error creating blacklist storage: %v", err))
	}

	verification, err := storage.NewVerificationCodesStorage(baseLogger, dbSettings(settings.Storage.VerificationCodeStorage))
	if err != nil {
		logger.Error("Error on Create New verification codes storage", logging.FieldError, err)
		errs = append(errs, fmt.Errorf("error creating verification codes storage: %v", err))
	}

	invite, err := storage.NewInviteStorage(baseLogger, dbSettings(settings.Storage.InviteStorage))
	if err != nil {
		logger.Error("Error on Create New invite storage", logging.FieldError, err)
		errs = append(errs, fmt.Errorf("error creating invite storage: %v", err))
	}

	managementKeys, err := storage.NewManagementKeys(baseLogger, dbSettings(settings.Storage.ManagementKeysStorage))
	if err != nil {
		logger.Error("Error on Create New management keys storage", logging.FieldError, err)
		errs = append(errs, fmt.Errorf("error creating management keys storage: %v", err))
	}

	session, err := storage.NewSessionStorage(baseLogger, settings.SessionStorage)
	if err != nil {
		logger.Error("Error on Create New session storage", logging.FieldError, err)
		errs = append(errs, fmt.Errorf("error creating session storage: %v", err))
	}

	key, err := storage.NewKeyStorage(baseLogger, settings.KeyStorage)
	if err != nil {
		logger.Error("Error on Create New key storage", logging.FieldError, err)
		errs = append(errs, fmt.Errorf("error creating key storage: %v", err))
	}

	// maybe just not serve login web app if type is none?
	lwas := settings.LoginWebApp
	if settings.LoginWebApp.Type == model.FileStorageTypeNone ||
		settings.LoginWebApp.Type == model.FileStorageTypeDefault {
		// if not set, use default value
		lwas = defaultLoginWebAppFSSettings
	}
	loginFS, err := storage.NewFS(lwas)
	if err != nil {
		logger.Error("Error creating login fs storage", logging.FieldError, err)
		errs = append(errs, fmt.Errorf("error creating login fs storage: %v", err))
	}

	var adminPanelFS fs.FS
	if settings.AdminPanel.Enabled {
		adminPanelFS, err = storage.NewFS(adminPanelFSSettings)
		if err != nil {
			logger.Error("Error creating admin panel fs storage", logging.FieldError, err)
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
	sms, err := sms.NewService(baseLogger, settings.Services.SMS)
	if err != nil {
		logger.Error("Error creating SMS service",
			logging.FieldError, err)

		errs = append(errs, fmt.Errorf("error creating SMS service: %v", err))
	}

	// maybe not use email templates if type is None?
	ets := settings.EmailTemplates
	if ets.Type == model.FileStorageTypeNone || ets.Type == model.FileStorageTypeDefault {
		ets = defaultEmailTemplateFSSettings
	}
	emailTemplateFS, err := storage.NewFS(ets)
	if err != nil {
		logger.Error("Error creating email template filesystem", logging.FieldError, err)
		errs = append(errs, fmt.Errorf("error creating email template filesystem: %v", err))
	}

	emailServiceSettings := settings.Services.Email
	// update templates every five minutes and look templates in a root folder of FS
	email, err := mail.NewService(baseLogger, emailServiceSettings, emailTemplateFS, time.Minute*5, "")
	if err != nil {
		logger.Error("Error creating email service", logging.FieldError, err)
		errs = append(errs, fmt.Errorf("error creating email service: %v", err))
	}

	tokenS, err := NewTokenService(baseLogger, settings.General, sc)
	if err != nil {
		logger.Error("Error creating token service", logging.FieldError, err)
		errs = append(errs, fmt.Errorf("error creating token service: %v", err))
	}

	sessionS := model.NewSessionManager(settings.SessionStorage.SessionDuration, session)

	impS, err := NewImpersonationProvider(settings.Impersonation)
	if err != nil {
		logger.Error("Error creating impersonation provider", logging.FieldError, err)
		errs = append(errs, fmt.Errorf("error creating impersonation provider: %v", err))
	}

	srvs := model.ServerServices{
		SMS:           sms,
		Email:         email,
		Token:         tokenS,
		Session:       sessionS,
		Impersonation: impS,
	}

	server, err := server.NewServer(sc, srvs, errs, restartChan)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func NewTokenService(
	logger *slog.Logger,
	settings model.GeneralServerSettings,
	storages model.ServerStorageCollection,
) (model.TokenService, error) {
	key, err := storages.Key.LoadPrivateKey()
	if err != nil {
		return nil, err
	}

	tokenService, err := jwt.NewJWTokenService(
		logger,
		key,
		settings.Issuer,
		storages.Token,
		storages.App,
		storages.User,
	)
	return tokenService, err
}

func NewImpersonationProvider(settings model.ImpersonationSettings) (model.ImpersonationProvider, error) {
	switch settings.Type {
	case model.ImpersonationServiceTypeNone, "":
		return nil, nil
	case model.ImpersonationServiceTypeRole:
		return imp.NewAccessRoleImpersonator(settings.Role.AllowedRoles), nil
	case model.ImpersonationServiceTypeScope:
		return imp.NewScopeImpersonator(settings.Scope.AllowedScopes), nil
	case model.ImpersonationServiceTypePlugin:
		return impPlugin.NewImpersonationProvider(settings.Plugin, time.Second)
	}

	return nil, nil
}
