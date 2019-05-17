package dynamodb

import (
	"fmt"
	"path"

	"github.com/madappgang/identifo/jwt"
	jwtService "github.com/madappgang/identifo/jwt/service"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/storage/dynamodb"
)

// NewComposer creates new database composer.
func NewComposer(settings model.ServerSettings, options ...func(*DatabaseComposer) error) (*DatabaseComposer, error) {
	c := DatabaseComposer{
		settings:        settings,
		newAppStorage:   dynamodb.NewAppStorage,
		newUserStorage:  dynamodb.NewUserStorage,
		newTokenStorage: dynamodb.NewTokenStorage,
	}

	for _, option := range options {
		if err := option(&c); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

// InitAppStorage returns an argument that sets the appStorage initialization function.
func InitAppStorage(initAS func(*dynamodb.DB) (model.AppStorage, error)) func(*DatabaseComposer) error {
	return func(dc *DatabaseComposer) error {
		dc.newAppStorage = initAS
		return nil
	}
}

// InitUserStorage returns an argument that sets the userStorage initialization function.
func InitUserStorage(initUS func(*dynamodb.DB) (model.UserStorage, error)) func(*DatabaseComposer) error {
	return func(dc *DatabaseComposer) error {
		dc.newUserStorage = initUS
		return nil
	}
}

// InitTokenStorage returns an argument that sets the tokenStorage initialization function.
func InitTokenStorage(initTS func(*dynamodb.DB) (model.TokenStorage, error)) func(*DatabaseComposer) error {
	return func(dc *DatabaseComposer) error {
		dc.newTokenStorage = initTS
		return nil
	}
}

// DatabaseComposer composes DynamoDB services.
type DatabaseComposer struct {
	settings        model.ServerSettings
	newAppStorage   func(*dynamodb.DB) (model.AppStorage, error)
	newUserStorage  func(*dynamodb.DB) (model.UserStorage, error)
	newTokenStorage func(*dynamodb.DB) (model.TokenStorage, error)
}

// Compose composes all services with DynamoDB support.
func (dc *DatabaseComposer) Compose() (
	model.AppStorage,
	model.UserStorage,
	model.TokenStorage,
	jwtService.TokenService,
	error,
) {
	db, err := dynamodb.NewDB(dc.settings.DBEndpoint, dc.settings.DBRegion)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	appStorage, err := dc.newAppStorage(db)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	userStorage, err := dc.newUserStorage(db)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tokenStorage, err := dc.newTokenStorage(db)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tokenServiceAlg, ok := jwt.StrToTokenServiceAlg[dc.settings.Algorithm]
	if !ok {
		return nil, nil, nil, nil, fmt.Errorf("Unknow token service algoritm %s", dc.settings.Algorithm)
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
		return nil, nil, nil, nil, err
	}

	return appStorage, userStorage, tokenStorage, tokenService, nil
}
