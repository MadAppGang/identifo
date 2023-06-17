package model

import "time"

// UnitedUserIdentity connect all federated identities for a user in one entity.
type UnitedUserIdentity struct {
	ID     string `json:"id,omitempty"`
	UserID string `json:"user_id,omitempty"`

	// ExternalID is the ID of the user identity in the external system.
	ExternalID string           `json:"external_id,omitempty"`
	Type       UserIdentityType `json:"type,omitempty"`

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

type UserIdentityType string

const (
	UserIdentityTypeApple    UserIdentityType = "apple"
	UserIdentityTypeGoogle   UserIdentityType = "google"
	UserIdentityTypeFacebook UserIdentityType = "facebook"
	UserIdentityTypeTwitter  UserIdentityType = "twitter"
	UserIdentityTypeOIDC     UserIdentityType = "oidc"
	UserIdentityTypeUnknown  UserIdentityType = "unknown"
	UserIdentityTypeOther    UserIdentityType = "other"
)

// TODO: Jack: implement FIM fields for user identity.

type UserIdentityFacebook struct{}

type UserIdentityGoogle struct{}

type UserIdentityTwitter struct{}

type UserIdentityApple struct{}

type UserIdentityOIDC struct{}
