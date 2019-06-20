package mgo

import (
	"fmt"
	"path"

	"github.com/madappgang/identifo/jwt"
	jwtService "github.com/madappgang/identifo/jwt/service"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/storage/mongo"
)

// NewComposer creates new database composer.
func NewComposer(settings model.ServerSettings, options ...func(*DatabaseComposer) error) (*DatabaseComposer, error) {
	c := DatabaseComposer{
		settings:                   settings,
		newAppStorage:              mongo.NewAppStorage,
		newUserStorage:             mongo.NewUserStorage,
		newTokenStorage:            mongo.NewTokenStorage,
		newVerificationCodeStorage: mongo.NewVerificationCodeStorage,
	}

	for _, option := range options {
		if err := option(&c); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

// DatabaseComposer composes MongoDB services.
type DatabaseComposer struct {
	settings                   model.ServerSettings
	newAppStorage              func(*mongo.DB) (model.AppStorage, error)
	newUserStorage             func(*mongo.DB) (model.UserStorage, error)
	newTokenStorage            func(*mongo.DB) (model.TokenStorage, error)
	newVerificationCodeStorage func(*mongo.DB) (model.VerificationCodeStorage, error)
}

// Compose composes all services with MongoDB support.
func (dc *DatabaseComposer) Compose() (
	model.AppStorage,
	model.UserStorage,
	model.TokenStorage,
	model.VerificationCodeStorage,
	jwtService.TokenService,
	error,
) {

	db, err := mongo.NewDB(dc.settings.Database.DBEndpoint, dc.settings.Database.DBName)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	appStorage, err := dc.newAppStorage(db)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	userStorage, err := dc.newUserStorage(db)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	tokenStorage, err := dc.newTokenStorage(db)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	verificationCodeStorage, err := dc.newVerificationCodeStorage(db)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	tokenServiceAlg, ok := jwt.StrToTokenSignAlg[dc.settings.Algorithm]
	if !ok {
		return nil, nil, nil, nil, nil, fmt.Errorf("Unknown token service algorithm %s ", dc.settings.Algorithm)
	}

	tokenService, err := jwtService.NewJWTokenService(
		path.Join(dc.settings.PEMFolderPath, dc.settings.PrivateKey),
		path.Join(dc.settings.PEMFolderPath, dc.settings.PublicKey),
		dc.settings.Issuer,
		tokenServiceAlg,
		tokenStorage,
		appStorage,
		userStorage,
	)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return appStorage, userStorage, tokenStorage, verificationCodeStorage, tokenService, nil
}
