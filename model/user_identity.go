package model

import "time"

// UserIdentity connect to remote user identity.
type UserIdentity struct {
	ID           string
	UserID       string
	ConnectedAt  time.Time
	UpdatedAt    time.Time
	AccessToken  string
	RefreshToken string
	Facebook     *UserIdentityFacebook
	Google       *UserIdentityGoogle
	Twitter      *UserIdentityTwitter
	Apple        *UserIdentityApple
	OIDC         map[string]UserIdentityOIDC
}

type UserIdentityFacebook struct{}

type UserIdentityGoogle struct{}

type UserIdentityTwitter struct{}

type UserIdentityApple struct{}

type UserIdentityOIDC struct{}
