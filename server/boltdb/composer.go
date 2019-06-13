package boltdb

import (
	"fmt"
	"path"

	"github.com/boltdb/bolt"
	"github.com/madappgang/identifo/jwt"
	jwtService "github.com/madappgang/identifo/jwt/service"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/storage/boltdb"
)

// NewComposer creates new database composer.
func NewComposer(settings model.ServerSettings, options ...func(*DatabaseComposer) error) (*DatabaseComposer, error) {
	c := DatabaseComposer{
		settings:                   settings,
		newAppStorage:              boltdb.NewAppStorage,
		newUserStorage:             boltdb.NewUserStorage,
		newTokenStorage:            boltdb.NewTokenStorage,
		newVerificationCodeStorage: boltdb.NewVerificationCodeStorage,
	}

	for _, option := range options {
		if err := option(&c); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

// DatabaseComposer composes BoltDB services.
type DatabaseComposer struct {
	settings                   model.ServerSettings
	newAppStorage              func(*bolt.DB) (model.AppStorage, error)
	newUserStorage             func(*bolt.DB) (model.UserStorage, error)
	newTokenStorage            func(*bolt.DB) (model.TokenStorage, error)
	newVerificationCodeStorage func(*bolt.DB) (model.VerificationCodeStorage, error)
}

// Compose composes all services with BoltDB support.
func (dc *DatabaseComposer) Compose() (
	model.AppStorage,
	model.UserStorage,
	model.TokenStorage,
	model.VerificationCodeStorage,
	jwtService.TokenService,
	error,
) {

	db, err := boltdb.InitDB(dc.settings.DBPath)
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
