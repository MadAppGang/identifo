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
		newTokenBlacklist:          mongo.NewTokenBlacklist,
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
	newTokenBlacklist          func(*mongo.DB) (model.TokenBlacklist, error)
	newVerificationCodeStorage func(*mongo.DB) (model.VerificationCodeStorage, error)
}

// Compose composes all services with MongoDB support.
func (dc *DatabaseComposer) Compose() (
	model.AppStorage,
	model.UserStorage,
	model.TokenStorage,
	model.TokenBlacklist,
	model.VerificationCodeStorage,
	jwtService.TokenService,
	error,
) {
	// We assume that all MongoDB-backed storages share the same database name and connection string, so we can pick any of them.
	db, err := mongo.NewDB(dc.settings.Storage.AppStorage.Endpoint, dc.settings.Storage.AppStorage.Name)
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

	tokenServiceAlg, ok := jwt.StrToTokenSignAlg[dc.settings.Algorithm]
	if !ok {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("Unknown token service algorithm %s ", dc.settings.Algorithm)
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
		return nil, nil, nil, nil, nil, nil, err
	}

	return appStorage, userStorage, tokenStorage, tokenBlacklist, verificationCodeStorage, tokenService, nil
}

// NewPartialComposer returns new partial composer with MongoDB support.
func NewPartialComposer(settings model.StorageSettings, options ...func(*PartialDatabaseComposer) error) (*PartialDatabaseComposer, error) {
	pc := &PartialDatabaseComposer{}
	// We assume that all MongoDB-backed storages share the same database name and connection string, so we can pick any of them.
	var dbEndpoint, dbName string

	if settings.AppStorage.Type == model.DBTypeMongoDB {
		pc.newAppStorage = mongo.NewAppStorage
		dbEndpoint = settings.AppStorage.Endpoint
		dbName = settings.AppStorage.Name
	}

	if settings.UserStorage.Type == model.DBTypeMongoDB {
		pc.newUserStorage = mongo.NewUserStorage
		dbEndpoint = settings.UserStorage.Endpoint
		dbName = settings.UserStorage.Name
	}

	if settings.TokenStorage.Type == model.DBTypeMongoDB {
		pc.newTokenStorage = mongo.NewTokenStorage
		dbEndpoint = settings.TokenStorage.Endpoint
		dbName = settings.TokenStorage.Name
	}

	if settings.TokenBlacklist.Type == model.DBTypeMongoDB {
		pc.newTokenBlacklist = mongo.NewTokenBlacklist
		dbEndpoint = settings.TokenBlacklist.Endpoint
		dbName = settings.TokenBlacklist.Name
	}

	if settings.VerificationCodeStorage.Type == model.DBTypeMongoDB {
		pc.newVerificationCodeStorage = mongo.NewVerificationCodeStorage
		dbEndpoint = settings.VerificationCodeStorage.Endpoint
		dbName = settings.VerificationCodeStorage.Name
	}

	db, err := mongo.NewDB(dbEndpoint, dbName)
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

// PartialDatabaseComposer composes only MongoDB-supporting services.
type PartialDatabaseComposer struct {
	db                         *mongo.DB
	newAppStorage              func(*mongo.DB) (model.AppStorage, error)
	newUserStorage             func(*mongo.DB) (model.UserStorage, error)
	newTokenStorage            func(*mongo.DB) (model.TokenStorage, error)
	newTokenBlacklist          func(*mongo.DB) (model.TokenBlacklist, error)
	newVerificationCodeStorage func(*mongo.DB) (model.VerificationCodeStorage, error)
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
