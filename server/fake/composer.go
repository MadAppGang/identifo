package fake

import (
	"fmt"

	"github.com/madappgang/identifo/jwt"
	jwtService "github.com/madappgang/identifo/jwt/service"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/storage/mem"
)

// NewComposer creates new database composer with in-memory storage support.
func NewComposer(settings model.ServerSettings) (*DatabaseComposer, error) {
	c := DatabaseComposer{
		settings:                   settings,
		newAppStorage:              mem.NewAppStorage,
		newUserStorage:             mem.NewUserStorage,
		newTokenStorage:            mem.NewTokenStorage,
		newTokenBlacklist:          mem.NewTokenBlacklist,
		newVerificationCodeStorage: mem.NewVerificationCodeStorage,
	}
	return &c, nil
}

// DatabaseComposer composes in-memory services.
type DatabaseComposer struct {
	settings                   model.ServerSettings
	newAppStorage              func() (model.AppStorage, error)
	newUserStorage             func() (model.UserStorage, error)
	newTokenStorage            func() (model.TokenStorage, error)
	newTokenBlacklist          func() (model.TokenBlacklist, error)
	newVerificationCodeStorage func() (model.VerificationCodeStorage, error)
}

// Compose composes all services with in-memory storage support.
func (dc *DatabaseComposer) Compose() (
	model.AppStorage,
	model.UserStorage,
	model.TokenStorage,
	model.TokenBlacklist,
	model.VerificationCodeStorage,
	jwtService.TokenService,
	error,
) {
	appStorage, err := dc.newAppStorage()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	userStorage, err := dc.newUserStorage()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	tokenStorage, err := dc.newTokenStorage()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	tokenBlacklist, err := dc.newTokenBlacklist()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	verificationCodeStorage, err := dc.newVerificationCodeStorage()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	tokenServiceAlg, ok := jwt.StrToTokenSignAlg[dc.settings.General.Algorithm]
	if !ok {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("Unknown token service algorithm %s", dc.settings.General.Algorithm)
	}

	tokenService, err := jwtService.NewJWTokenService(
		dc.settings.General.PrivateKeyPath,
		dc.settings.General.PublicKeyPath,
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

// NewPartialComposer returns new partial composer with in-memory storage support.
func NewPartialComposer(settings model.StorageSettings, options ...func(*PartialDatabaseComposer) error) (*PartialDatabaseComposer, error) {
	pc := &PartialDatabaseComposer{}

	if settings.AppStorage.Type == model.DBTypeFake {
		pc.newAppStorage = mem.NewAppStorage
	}

	if settings.UserStorage.Type == model.DBTypeFake {
		pc.newUserStorage = mem.NewUserStorage
	}

	if settings.TokenStorage.Type == model.DBTypeFake {
		pc.newTokenStorage = mem.NewTokenStorage
	}

	if settings.TokenBlacklist.Type == model.DBTypeFake {
		pc.newTokenBlacklist = mem.NewTokenBlacklist
	}

	if settings.VerificationCodeStorage.Type == model.DBTypeFake {
		pc.newVerificationCodeStorage = mem.NewVerificationCodeStorage
	}

	for _, option := range options {
		if err := option(pc); err != nil {
			return nil, err
		}
	}
	return pc, nil
}

// PartialDatabaseComposer composes only those services that support in-memory storage.
type PartialDatabaseComposer struct {
	newAppStorage              func() (model.AppStorage, error)
	newUserStorage             func() (model.UserStorage, error)
	newTokenStorage            func() (model.TokenStorage, error)
	newTokenBlacklist          func() (model.TokenBlacklist, error)
	newVerificationCodeStorage func() (model.VerificationCodeStorage, error)
}

// AppStorageComposer returns app storage composer.
func (pc *PartialDatabaseComposer) AppStorageComposer() func() (model.AppStorage, error) {
	return func() (model.AppStorage, error) {
		return pc.newAppStorage()
	}
}

// UserStorageComposer returns user storage composer.
func (pc *PartialDatabaseComposer) UserStorageComposer() func() (model.UserStorage, error) {
	return func() (model.UserStorage, error) {
		return pc.newUserStorage()
	}
}

// TokenStorageComposer returns token storage composer.
func (pc *PartialDatabaseComposer) TokenStorageComposer() func() (model.TokenStorage, error) {
	if pc.newTokenStorage != nil {
		return func() (model.TokenStorage, error) {
			return pc.newTokenStorage()
		}
	}
	return nil
}

// TokenBlacklistComposer returns token blacklist composer.
func (pc *PartialDatabaseComposer) TokenBlacklistComposer() func() (model.TokenBlacklist, error) {
	if pc.newTokenBlacklist != nil {
		return func() (model.TokenBlacklist, error) {
			return pc.newTokenBlacklist()
		}
	}
	return nil
}

// VerificationCodeStorageComposer returns verification code storage composer.
func (pc *PartialDatabaseComposer) VerificationCodeStorageComposer() func() (model.VerificationCodeStorage, error) {
	return func() (model.VerificationCodeStorage, error) {
		return pc.newVerificationCodeStorage()
	}
}
