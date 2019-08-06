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
func NewComposer(settings model.ServerSettings) (*DatabaseComposer, error) {
	c := DatabaseComposer{
		settings:                   settings,
		newAppStorage:              dynamodb.NewAppStorage,
		newUserStorage:             dynamodb.NewUserStorage,
		newTokenStorage:            dynamodb.NewTokenStorage,
		newTokenBlacklist:          dynamodb.NewTokenBlacklist,
		newVerificationCodeStorage: dynamodb.NewVerificationCodeStorage,
	}
	return &c, nil
}

// DatabaseComposer composes DynamoDB services.
type DatabaseComposer struct {
	settings                   model.ServerSettings
	newAppStorage              func(*dynamodb.DB) (model.AppStorage, error)
	newUserStorage             func(*dynamodb.DB) (model.UserStorage, error)
	newTokenStorage            func(*dynamodb.DB) (model.TokenStorage, error)
	newTokenBlacklist          func(*dynamodb.DB) (model.TokenBlacklist, error)
	newVerificationCodeStorage func(*dynamodb.DB) (model.VerificationCodeStorage, error)
}

// Compose composes all services with DynamoDB support.
func (dc *DatabaseComposer) Compose() (
	model.AppStorage,
	model.UserStorage,
	model.TokenStorage,
	model.TokenBlacklist,
	model.VerificationCodeStorage,
	jwtService.TokenService,
	error,
) {
	// We assume that all DynamoDB-backed storages share the same endpoint and region, so we can pick any of them.
	db, err := dynamodb.NewDB(dc.settings.Storage.AppStorage.Endpoint, dc.settings.Storage.AppStorage.Region)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	appStorage, err := dc.newAppStorage(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	userStorage, err := dc.newUserStorage(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	tokenStorage, err := dc.newTokenStorage(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	tokenBlacklist, err := dc.newTokenBlacklist(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	verificationCodeStorage, err := dc.newVerificationCodeStorage(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	tokenServiceAlg, ok := jwt.StrToTokenSignAlg[dc.settings.General.Algorithm]
	if !ok {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("Unknown token service algorithm %s", dc.settings.General.Algorithm)
	}

	tokenService, err := jwtService.NewJWTokenService(
		path.Join(dc.settings.General.PEMFolderPath, dc.settings.General.PrivateKey),
		path.Join(dc.settings.General.PEMFolderPath, dc.settings.General.PublicKey),
		dc.settings.General.Issuer,
		tokenServiceAlg,
		tokenStorage,
		appStorage,
		userStorage,
	)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	return appStorage, userStorage, tokenStorage, tokenBlacklist, verificationCodeStorage, tokenService, nil
}

// NewPartialComposer returns new partial composer with DynamoDB support.
func NewPartialComposer(settings model.StorageSettings, options ...func(*PartialDatabaseComposer) error) (*PartialDatabaseComposer, error) {
	pc := &PartialDatabaseComposer{}
	// We assume that all DynamoDB-backed storages share the same endpoint and region, so we can pick any of them.
	var dbEndpoint, dbRegion string

	if settings.AppStorage.Type == model.DBTypeDynamoDB {
		pc.newAppStorage = dynamodb.NewAppStorage
		dbEndpoint = settings.AppStorage.Endpoint
		dbRegion = settings.AppStorage.Region
	}

	if settings.UserStorage.Type == model.DBTypeDynamoDB {
		pc.newUserStorage = dynamodb.NewUserStorage
		dbEndpoint = settings.UserStorage.Endpoint
		dbRegion = settings.UserStorage.Region
	}

	if settings.TokenStorage.Type == model.DBTypeDynamoDB {
		pc.newTokenStorage = dynamodb.NewTokenStorage
		dbEndpoint = settings.TokenStorage.Endpoint
		dbRegion = settings.TokenStorage.Region
	}

	if settings.TokenBlacklist.Type == model.DBTypeDynamoDB {
		pc.newTokenBlacklist = dynamodb.NewTokenBlacklist
		dbEndpoint = settings.TokenBlacklist.Endpoint
		dbRegion = settings.TokenBlacklist.Region
	}

	if settings.VerificationCodeStorage.Type == model.DBTypeDynamoDB {
		pc.newVerificationCodeStorage = dynamodb.NewVerificationCodeStorage
		dbEndpoint = settings.VerificationCodeStorage.Endpoint
		dbRegion = settings.VerificationCodeStorage.Region
	}

	db, err := dynamodb.NewDB(dbEndpoint, dbRegion)
	if err != nil {
		return nil, err
	}
	pc.db = db

	for _, option := range options {
		if err := option(pc); err != nil {
			return nil, err
		}
	}
	return pc, nil
}

// PartialDatabaseComposer composes only DynamoDB-supporting services.
type PartialDatabaseComposer struct {
	db                         *dynamodb.DB
	newAppStorage              func(*dynamodb.DB) (model.AppStorage, error)
	newUserStorage             func(*dynamodb.DB) (model.UserStorage, error)
	newTokenStorage            func(*dynamodb.DB) (model.TokenStorage, error)
	newTokenBlacklist          func(*dynamodb.DB) (model.TokenBlacklist, error)
	newVerificationCodeStorage func(*dynamodb.DB) (model.VerificationCodeStorage, error)
}

// AppStorageComposer returns app storage composer.
func (pc *PartialDatabaseComposer) AppStorageComposer() func() (model.AppStorage, error) {
	if pc.newAppStorage != nil {
		return func() (model.AppStorage, error) {
			return pc.newAppStorage(pc.db)
		}
	}
	return nil
}

// UserStorageComposer returns user storage composer.
func (pc *PartialDatabaseComposer) UserStorageComposer() func() (model.UserStorage, error) {
	if pc.newUserStorage != nil {
		return func() (model.UserStorage, error) {
			return pc.newUserStorage(pc.db)
		}
	}
	return nil
}

// TokenStorageComposer returns token storage composer.
func (pc *PartialDatabaseComposer) TokenStorageComposer() func() (model.TokenStorage, error) {
	if pc.newTokenStorage != nil {
		return func() (model.TokenStorage, error) {
			return pc.newTokenStorage(pc.db)
		}
	}
	return nil
}

// TokenBlacklistComposer returns token blacklist composer.
func (pc *PartialDatabaseComposer) TokenBlacklistComposer() func() (model.TokenBlacklist, error) {
	if pc.newTokenBlacklist != nil {
		return func() (model.TokenBlacklist, error) {
			return pc.newTokenBlacklist(pc.db)
		}
	}
	return nil
}

// VerificationCodeStorageComposer returns verification code storage composer.
func (pc *PartialDatabaseComposer) VerificationCodeStorageComposer() func() (model.VerificationCodeStorage, error) {
	if pc.newVerificationCodeStorage != nil {
		return func() (model.VerificationCodeStorage, error) {
			return pc.newVerificationCodeStorage(pc.db)
		}
	}
	return nil
}
