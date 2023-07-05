package model

import "time"

// UnitedUserIdentity connect all federated identities for a user in one entity.
type UnitedUserIdentity struct {
	ID     string `json:"id,omitempty"`
	UserID string `json:"user_id,omitempty"`

	// ExternalID is the ID of the user identity in the external system.
	ExternalID string            `json:"external_id,omitempty"`
	Type       UserFederatedType `json:"federated_type,omitempty"`

	// TypeOther, if Type is UserIdentityTypeOther, this field keeps unique name for those other identity type.
	TypeOther   string    `json:"type_other,omitempty"`
	ConnectedAt time.Time `json:"connected_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`

	// Unique field set for each external iDP.
	Facebook *UserIdentityFacebook       `json:"facebook,omitempty"`
	Google   *UserIdentityGoogle         `json:"google,omitempty"`
	Twitter  *UserIdentityTwitter        `json:"twitter,omitempty"`
	Apple    *UserIdentityApple          `json:"apple,omitempty"`
	OIDC     map[string]UserIdentityOIDC `json:"oidc,omitempty"`
}

// UserFederatedType the user federated identity type.
type UserFederatedType string

const (
	UserIdentityTypeApple    UserFederatedType = "apple"
	UserIdentityTypeGoogle   UserFederatedType = "google"
	UserIdentityTypeFacebook UserFederatedType = "facebook"
	UserIdentityTypeTwitter  UserFederatedType = "twitter"
	UserIdentityTypeOIDC     UserFederatedType = "oidc"
	UserIdentityTypeUnknown  UserFederatedType = "unknown"
	UserIdentityTypeOther    UserFederatedType = "other"
)

// TODO: Jack: implement FIM fields for user identity.

// UserIdentityFacebook facebook specific fields.
type UserIdentityFacebook struct{}

// UserIdentityGoogle google specific fields.
type UserIdentityGoogle struct{}

// UserIdentityTwitter twitter specific fields.
type UserIdentityTwitter struct{}

// UserIdentityApple apple specific fields.
type UserIdentityApple struct{}

// UserIdentityOIDC oidc specific fields.
type UserIdentityOIDC struct{}
