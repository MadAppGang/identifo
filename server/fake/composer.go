package fake

import (
	"path"

	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/mem"
	"github.com/madappgang/identifo/model"
)

// NewComposer creates new database composer.
func NewComposer(settings Settings) (*DatabaseComposer, error) {
	c := DatabaseComposer{
		settings: settings,
	}
	return &c, nil
}

// DatabaseComposer composes in-memory services.
type DatabaseComposer struct {
	settings Settings
}

// Compose composes all services with in-memory storage support.
func (dc *DatabaseComposer) Compose() (
	model.AppStorage,
	model.UserStorage,
	model.TokenStorage,
	model.TokenService,
	error,
) {

	appStorage := mem.NewAppStorage()
	userStorage := mem.NewUserStorage()
	tokenStorage := mem.NewTokenStorage()

	tokenService, err := jwt.NewTokenService(
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
