package server

import (
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/plugin/shared"
)

// DatabaseComposer inits database stack.
type DatabaseComposer interface {
	Compose() (
		model.AppStorage,
		shared.UserStorage,
		model.TokenStorage,
		model.TokenBlacklist,
		model.VerificationCodeStorage,
		error,
	)
}

// PartialDatabaseComposer can init services backed with different databases.
type PartialDatabaseComposer interface {
	AppStorageComposer() func() (model.AppStorage, error)
	TokenStorageComposer() func() (model.TokenStorage, error)
	TokenBlacklistComposer() func() (model.TokenBlacklist, error)
	VerificationCodeStorageComposer() func() (model.VerificationCodeStorage, error)
}

// Composer is a service composer which is agnostic to particular database implementations.
type Composer struct {
	settings                   model.ServerSettings
	newAppStorage              func() (model.AppStorage, error)
	userStorage                shared.UserStorage
	newTokenStorage            func() (model.TokenStorage, error)
	newTokenBlacklist          func() (model.TokenBlacklist, error)
	newVerificationCodeStorage func() (model.VerificationCodeStorage, error)
}

// Compose composes all services.
func (c *Composer) Compose() (
	model.AppStorage,
	shared.UserStorage,
	model.TokenStorage,
	model.TokenBlacklist,
	model.VerificationCodeStorage,
	error,
) {
	appStorage, err := c.newAppStorage()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	tokenStorage, err := c.newTokenStorage()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	tokenBlacklist, err := c.newTokenBlacklist()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	verificationCodeStorage, err := c.newVerificationCodeStorage()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return appStorage, c.userStorage, tokenStorage, tokenBlacklist, verificationCodeStorage, nil
}

// NewComposer returns new database composer based on passed server settings.
func NewComposer(settings model.ServerSettings, partialComposers []PartialDatabaseComposer, plugins shared.Plugins, options ...func(*Composer) error) (*Composer, error) {
	c := &Composer{settings: settings}

	c.userStorage = plugins.UserStorage

	for _, pc := range partialComposers {
		if pc.AppStorageComposer() != nil {
			c.newAppStorage = pc.AppStorageComposer()
		}
		if pc.TokenStorageComposer() != nil {
			c.newTokenStorage = pc.TokenStorageComposer()
		}
		if pc.TokenBlacklistComposer() != nil {
			c.newTokenBlacklist = pc.TokenBlacklistComposer()
		}
		if pc.VerificationCodeStorageComposer() != nil {
			c.newVerificationCodeStorage = pc.VerificationCodeStorageComposer()
		}
	}

	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}
