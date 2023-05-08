package model

import (
	"encoding/json"
	"log"
	"time"
)

// 2FA connectors
// https://community.auth0.com/t/when-will-api-v2-users-userid-authenticators-be-documented/52722

// Guardian authenticators
// https://auth0.com/docs/secure/multi-factor-authentication/auth0-guardian

// User is a user identitty data
// do not mix it with login strategy
// https://auth0.com/docs/manage-users/user-accounts/user-profiles

// oidc claims
// https://www.iana.org/assignments/jwt/jwt.xhtml#claims
type User struct {
	ID string `json:"id" bson:"_id"`

	TenantMembership *struct {
		TenantID   string            `json:"tenant_id"`
		TenantName string            `json:"tenant_name"`
		Groups     map[string]string // map of group names to ids
	}

	// user information
	Username          string `json:"username" bson:"username"` // it is a nickname for login purposes
	Email             string `json:"email" bson:"email"`
	GivenName         string
	FamilyName        string
	MiddleName        string
	Nickname          string
	PreferredUsername string
	PhoneNumber       string
	AvatarURL         string

	//  authentication data for user
	PasswordHash          string // TODO: do we need salt and pepper here as well?
	PasswordResetRequired string
	PasswordChangeForced  string
	Authenticators        []UserAuthenticator
	Identities            []UserIdentity
	MFAEnrollments        []UserMFAEnrollment

	// login stats
	LastIP              string
	LastLoginAt         time.Time
	LastPasswordResetAt time.Time
	LoginsCount         int

	// verification data
	PhoneVerificationDetails struct {
		VerifiedAt      time.Time
		VerifiedDetails string
	}

	EmailVerificationDetails struct {
		VerifiedAt      time.Time
		VerifiedDetails string
	}

	// blocked user
	Blocked        bool
	BlockedDetails struct {
		BlockedReason string
		BlockedAt     time.Time
		BlockedBy     string
	}

	// mapping to external systems
	ExternalID      string            // external system user ID
	ExternalMapping map[string]string // external systems mapping

	// User push devices
	ActiveDevices []UserDevice

	// Additional data for user
	AppsData []ApplicationUserData
	Data     []AdditionalUserData

	// oidc claims
	// https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims
	Profile  string
	Picture  string
	Website  string
	Gender   string
	Birthday time.Time
	Timezone string
	Locale   string
	Address  map[string]UserAddress // addresses for home, work, etc

	// user record metadata
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ApplicationUserData is custom data that could be attached by application to the user,
// could be theme settings, preferences or some other information
type ApplicationUserData struct {
	AppID       string
	Data        map[string]any
	LastUpdated time.Time
}

// AdditionalUserData is custom data attached to the user
// the data is organized in buckets
type AdditionalUserData struct {
	BucketName string
	BucketData map[string]any
}

const (
	HomeAddress  = "home"
	WorkAddress  = "work"
	OtherAddress = "other"
)

// https://openid.net/specs/openid-connect-core-1_0.html#AddressClaim
type UserAddress struct {
	Formatted     string
	StreetAddress string
	Locality      string // City or locality component.
	Region        string
	PostalCode    string
	Country       string
}

// UserFromJSON deserialize user data from JSON.
func UserFromJSON(d []byte) (User, error) {
	var user User
	if err := json.Unmarshal(d, &user); err != nil {
		log.Println("Error while unmarshal user:", err)
		return User{}, err
	}
	return user, nil
}
