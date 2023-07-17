package model

import "time"

var SupportedFIMType = []UserFederatedType{}

// TODO: implement creating and updating those identities.
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
}

// UserFederatedType the user federated identity type.
type UserFederatedType string

const (
	UserFIMApple    UserFederatedType = "apple"
	UserFIMGoogle   UserFederatedType = "google"
	UserFIMFacebook UserFederatedType = "facebook"
	UserFIMTwitter  UserFederatedType = "twitter"
	UserFIMOIDC     UserFederatedType = "oidc"
	UserFIMUnknown  UserFederatedType = "unknown"
	UserFIMOther    UserFederatedType = "other"
)
