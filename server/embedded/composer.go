package embedded

import (
	"path"

	"github.com/madappgang/identifo/boltdb"
	"github.com/madappgang/identifo/encryptor"
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

//DatabaseComposer composes BoltDB services
type DatabaseComposer struct {
	settings Settings
}

//Compose composes all services with BoltDB support
func (dc *DatabaseComposer) Compose() (
	model.AppStorage,
	model.UserStorage,
	model.TokenStorage,
	model.TokenService,
	model.Encryptor,
	error,
) {

	db, err := boltdb.InitDB(dc.settings.DBPath)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	appStorage, err := boltdb.NewAppStorage(db)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	userStorage, err := boltdb.NewUserStorage(db)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	tokenStorage, err := boltdb.NewTokenStorage(db)
	if err != nil {
		return nil, nil, nil, nil, nil, err
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
		return nil, nil, nil, nil, nil, err
	}

	encryptor, err := encryptor.NewEncryptor(dc.settings.EncryptionKeyPath)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return appStorage, userStorage, tokenStorage, tokenService, encryptor, nil
}
