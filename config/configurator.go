package config

import (
	"fmt"

	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/server"
	"github.com/madappgang/identifo/services/mail"
	"github.com/madappgang/identifo/services/sms"
	"github.com/madappgang/identifo/storage"

	jwt "github.com/madappgang/identifo/jwt/service"
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

	key, err := storage.NewKeyStorage(settings.KeyStorage)
	if err != nil {
		return nil, err
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
		Static:       static,
		Key:          key,
	}

	// create 3rd party services
	sms, err := sms.NewService(settings.ExternalServices.SMSService)
	if err != nil {
		return nil, err
	}

	email, err := mail.NewService(settings.ExternalServices.EmailService, static)
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

	server, err := server.NewServer(sc, srvs)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func NewTokenService(settings model.GeneralServerSettings, storages model.ServerStorageCollection) (model.TokenService, error) {
	tokenServiceAlg, ok := model.StrToTokenSignAlg[settings.Algorithm]
	if !ok {
		return nil, fmt.Errorf("Unknown token service algorithm %s", settings.Algorithm)
	}

	keys, err := storages.Key.LoadKeys(tokenServiceAlg)
	if err != nil {
		return nil, err
	}

	tokenService, err := jwt.NewJWTokenService(
		keys,
		settings.Issuer,
		storages.Token,
		storages.App,
		storages.User,
	)
	return tokenService, err
}
