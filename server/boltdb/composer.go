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

// NewComposer creates new database composer with BoltDB support.
func NewComposer(settings model.ServerSettings, options ...func(*DatabaseComposer) error) (*DatabaseComposer, error) {
	c := DatabaseComposer{
		settings:                   settings,
		newAppStorage:              boltdb.NewAppStorage,
		newUserStorage:             boltdb.NewUserStorage,
		newTokenStorage:            boltdb.NewTokenStorage,
		newTokenBlacklist:          boltdb.NewTokenBlacklist,
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
	newTokenBlacklist          func(*bolt.DB) (model.TokenBlacklist, error)
	newVerificationCodeStorage func(*bolt.DB) (model.VerificationCodeStorage, error)
}

// Compose composes all services with BoltDB support.
func (dc *DatabaseComposer) Compose() (
	model.AppStorage,
	model.UserStorage,
	model.TokenStorage,
	model.TokenBlacklist,
	model.VerificationCodeStorage,
	jwtService.TokenService,
	error,
) {
	// We assume that all BoltDB-backed storages share the same filepath, so we can pick any of them.
	db, err := boltdb.InitDB(dc.settings.Storage.AppStorage.Path)
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

// NewPartialComposer returns new partial composer with BoltDB support.
func NewPartialComposer(settings model.StorageSettings, options ...func(*PartialDatabaseComposer) error) (*PartialDatabaseComposer, error) {
	pc := &PartialDatabaseComposer{}
	// We assume that all BoltDB-backed storages share the same filepath, so we can pick any of them.
	var dbPath string

	if settings.AppStorage.Type == model.DBTypeBoltDB {
		pc.newAppStorage = boltdb.NewAppStorage
		dbPath = settings.AppStorage.Path
	}

	if settings.UserStorage.Type == model.DBTypeBoltDB {
		pc.newUserStorage = boltdb.NewUserStorage
		dbPath = settings.UserStorage.Path
	}

	if settings.TokenStorage.Type == model.DBTypeBoltDB {
		pc.newTokenStorage = boltdb.NewTokenStorage
		dbPath = settings.TokenStorage.Path
	}

	if settings.VerificationCodeStorage.Type == model.DBTypeBoltDB {
		pc.newVerificationCodeStorage = boltdb.NewVerificationCodeStorage
		dbPath = settings.VerificationCodeStorage.Path
	}

	db, err := boltdb.InitDB(dbPath)
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

// PartialDatabaseComposer composes only BoltDB-supporting services.
type PartialDatabaseComposer struct {
	db                         *bolt.DB
	newAppStorage              func(*bolt.DB) (model.AppStorage, error)
	newUserStorage             func(*bolt.DB) (model.UserStorage, error)
	newTokenStorage            func(*bolt.DB) (model.TokenStorage, error)
	newTokenBlacklist          func(*bolt.DB) (model.TokenBlacklist, error)
	newVerificationCodeStorage func(*bolt.DB) (model.VerificationCodeStorage, error)
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
