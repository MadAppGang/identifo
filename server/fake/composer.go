package fake

import (
	"path"

	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/storage/mem"
)

// NewComposer creates new database composer.
func NewComposer(settings model.ServerSettings, options ...func(*DatabaseComposer) error) (*DatabaseComposer, error) {
	c := DatabaseComposer{
		settings:        settings,
		newAppStorage:   mem.NewAppStorage,
		newUserStorage:  mem.NewUserStorage,
		newTokenStorage: mem.NewTokenStorage,
	}

	for _, option := range options {
		if err := option(&c); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

// DatabaseComposer composes in-memory services.
type DatabaseComposer struct {
	settings        model.ServerSettings
	newAppStorage   func() (model.AppStorage, error)
	newUserStorage  func() (model.UserStorage, error)
	newTokenStorage func() (model.TokenStorage, error)
}

// InitAppStorage returns an argument that sets the appStorage initialization function.
func InitAppStorage(initAS func() (model.AppStorage, error)) func(*DatabaseComposer) error {
	return func(dc *DatabaseComposer) error {
		dc.newAppStorage = initAS
		return nil
	}
}

// InitUserStorage returns an argument that sets the userStorage initialization function.
func InitUserStorage(initUS func() (model.UserStorage, error)) func(*DatabaseComposer) error {
	return func(dc *DatabaseComposer) error {
		dc.newUserStorage = initUS
		return nil
	}
}

// InitTokenStorage returns an argument that sets the tokenStorage initialization function.
func InitTokenStorage(initTS func() (model.TokenStorage, error)) func(*DatabaseComposer) error {
	return func(dc *DatabaseComposer) error {
		dc.newTokenStorage = initTS
		return nil
	}
}

// Compose composes all services with in-memory storage support.
func (dc *DatabaseComposer) Compose() (
	model.AppStorage,
	model.UserStorage,
	model.TokenStorage,
	tokensrvc.TokenService,
	error,
) {

	appStorage, err := dc.newAppStorage()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	userStorage, err := dc.newUserStorage()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tokenStorage, err := dc.newTokenStorage()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tokenService, err := tokensrvc.NewDefaultTokenService(
		path.Join(dc.settings.PEMFolderPath, dc.settings.PrivateKey),
		path.Join(dc.settings.PEMFolderPath, dc.settings.PublicKey),
		dc.settings.Issuer,
		dc.settings.Algorithm,
		tokenStorage,
		appStorage,
		userStorage,
	)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return appStorage, userStorage, tokenStorage, tokenService, nil
}
