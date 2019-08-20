package server

import (
	"github.com/madappgang/identifo/model"
)

// DatabaseComposer inits database stack.
type DatabaseComposer interface {
	Compose() (
		model.AppStorage,
		model.UserStorage,
		model.TokenStorage,
		model.TokenBlacklist,
		model.VerificationCodeStorage,
		error,
	)
}

// PartialDatabaseComposer can init services backed with different databases.
type PartialDatabaseComposer interface {
	AppStorageComposer() func() (model.AppStorage, error)
	UserStorageComposer() func() (model.UserStorage, error)
	TokenStorageComposer() func() (model.TokenStorage, error)
	TokenBlacklistComposer() func() (model.TokenBlacklist, error)
	VerificationCodeStorageComposer() func() (model.VerificationCodeStorage, error)
}

// Composer is a service composer which is agnostic to particular database implementations.
type Composer struct {
	settings                   model.ServerSettings
	newAppStorage              func() (model.AppStorage, error)
	newUserStorage             func() (model.UserStorage, error)
	newTokenStorage            func() (model.TokenStorage, error)
	newTokenBlacklist          func() (model.TokenBlacklist, error)
	newVerificationCodeStorage func() (model.VerificationCodeStorage, error)
}

// Compose composes all services.
func (c *Composer) Compose() (
	model.AppStorage,
	model.UserStorage,
	model.TokenStorage,
	model.TokenBlacklist,
	model.VerificationCodeStorage,
	error,
) {
	appStorage, err := c.newAppStorage()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	userStorage, err := c.newUserStorage()
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

	return appStorage, userStorage, tokenStorage, tokenBlacklist, verificationCodeStorage, nil
}

// NewComposer returns new database composer based on passed server settings.
func NewComposer(settings model.ServerSettings, partialComposers []PartialDatabaseComposer, options ...func(*Composer) error) (*Composer, error) {
	c := &Composer{settings: settings}

	for _, pc := range partialComposers {
		if pc.AppStorageComposer() != nil {
			c.newAppStorage = pc.AppStorageComposer()
		}
		if pc.UserStorageComposer() != nil {
			c.newUserStorage = pc.UserStorageComposer()
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
