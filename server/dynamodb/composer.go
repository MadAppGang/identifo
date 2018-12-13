package dynamodb

import (
	"path"

	"github.com/madappgang/identifo/dynamodb"
	"github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/model"
)

//NewComposer creates new database composer
func NewComposer(settings Settings) (*DatabaseComposer, error) {
	c := DatabaseComposer{
		settings: settings,
	}
	return &c, nil
}

//DatabaseComposer composes dynamodb services
type DatabaseComposer struct {
	settings Settings
}

//Compose composes all services with dynamodb support
func (dc *DatabaseComposer) Compose() (
	model.AppStorage,
	model.UserStorage,
	model.TokenStorage,
	model.TokenService,
	error,
) {

	db, err := dynamodb.NewDB(dc.settings.DBEndpoint, dc.settings.DBRegion)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	appStorage, err := dynamodb.NewAppStorage(db)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	userStorage, err := dynamodb.NewUserStorage(db)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tokenStorage, err := dynamodb.NewTokenStorage(db)
	if err != nil {
		return nil, nil, nil, nil, err
	}

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
