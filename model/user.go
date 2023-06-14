package model

import (
	"encoding/json"
	"log"
	"time"
)

type User struct {
	ID string `json:"id,omitempty" bson:"_id"`

	// user information
	Username          string `json:"username,omitempty" bson:"username"` // it is a nickname for login purposes
	Email             string `json:"email,omitempty" bson:"email"`
	GivenName         string `json:"given_name,omitempty"`
	FamilyName        string `json:"family_name,omitempty"`
	MiddleName        string `json:"middle_name,omitempty"`
	Nickname          string `json:"nickname,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`
	PhoneNumber       string `json:"phone_number,omitempty"`
	AvatarURL         string `json:"avatar_url,omitempty"`

	//  authentication data for user
	PasswordHash          string `json:"password_hash,omitempty"` // TODO: do we need salt and pepper here as well?
	PasswordResetRequired string `json:"password_reset_required,omitempty"`
	PasswordChangeForced  string `json:"password_change_forced,omitempty"`

	Tags []string `json:"tags,omitempty"`

	// oidc claims
	// https://openid.net/specs/openid-connect-core-1_0.html#StandardClaims
	Profile  string                 `json:"profile,omitempty"`
	Picture  string                 `json:"picture,omitempty"`
	Website  string                 `json:"website,omitempty"`
	Gender   string                 `json:"gender,omitempty"`
	Birthday time.Time              `json:"birthday,omitempty"`
	Timezone string                 `json:"timezone,omitempty"`
	Locale   string                 `json:"locale,omitempty"`
	Address  map[string]UserAddress `json:"address,omitempty"` // addresses for home, work, etc

	// login stats
	LastIP              string    `json:"last_ip,omitempty"`
	LastLoginAt         time.Time `json:"last_login_at,omitempty"`
	LastPasswordResetAt time.Time `json:"last_password_reset_at,omitempty"`
	LoginsCount         int       `json:"logins_count,omitempty"`

	// verification data
	PhoneVerificationDetails struct {
		VerifiedAt      time.Time `json:"verified_at,omitempty"`
		VerifiedDetails string    `json:"verified_details,omitempty"`
	} `json:"phone_verification_details,omitempty"`

	EmailVerificationDetails struct {
		VerifiedAt      time.Time `json:"verified_at,omitempty"`
		VerifiedDetails string    `json:"verified_details,omitempty"`
	} `json:"email_verification_details,omitempty"`

	// blocked user
	Blocked        bool `json:"blocked,omitempty"`
	BlockedDetails struct {
		BlockedReason string    `json:"blocked_reason,omitempty"`
		BlockedAt     time.Time `json:"blocked_at,omitempty"`
		BlockedBy     string    `json:"blocked_by,omitempty"`
	} `json:"blocked_details,omitempty"`

	// mapping to external systems
	ExternalID      string            `json:"external_id,omitempty"`      // external system user ID
	ExternalMapping map[string]string `json:"external_mapping,omitempty"` // external systems mapping

	// user record metadata
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// UserData model represents all collective information about the user
type UserData struct {
	UserID           string `json:"user_id,omitempty"`
	TenantMembership *struct {
		TenantID   string            `json:"tenant_id,omitempty"`
		TenantName string            `json:"tenant_name,omitempty"`
		Groups     map[string]string `json:"groups,omitempty"` // map of group names to ids
	} `json:"tenant_membership,omitempty"`
	AuthEnrollments []UserAuthEnrolment `json:"auth_enrollments,omitempty"`
	Identities      []UserIdentity      `json:"identities,omitempty"`
	MFAEnrollments  []UserMFAEnrollment `json:"mfa_enrollments,omitempty"`
	// User push devices
	ActiveDevices []UserDevice `json:"active_devices,omitempty"`
	// Additional data for user
	AppsData []ApplicationUserData `json:"apps_data,omitempty"`
	Data     []AdditionalUserData  `json:"data,omitempty"`
}

type UserDataField string

const (
	UserDataFieldTenantMembership UserDataField = "tenant_membership"
	UserDataFieldAuthEnrollments  UserDataField = "auth_enrollments"
	UserDataFieldIdentities       UserDataField = "identities"
	UserDataFieldMFAEnrollments   UserDataField = "mfa_enrollments"
	UserDataFieldActiveDevices    UserDataField = "active_devices"
	UserDataFieldAppsData         UserDataField = "apps_data"
	UserDataFieldData             UserDataField = "data"
	UserDataFieldAll              UserDataField = "all"
)

// ApplicationUserData is custom data that could be attached by application to the user,
// could be theme settings, preferences or some other information
type ApplicationUserData struct {
	AppID       string         `json:"app_id,omitempty"`
	Data        map[string]any `json:"data,omitempty"`
	LastUpdated time.Time      `json:"last_updated,omitempty"`
}

// AdditionalUserData is custom data attached to the user
// the data is organized in buckets
type AdditionalUserData struct {
	BucketName  string         `json:"bucket_name,omitempty"`
	BucketData  map[string]any `json:"bucket_data,omitempty"`
	LastUpdated time.Time      `json:"last_updated,omitempty"`
}

const (
	HomeAddress  = "home"
	WorkAddress  = "work"
	OtherAddress = "other"
)

// https://openid.net/specs/openid-connect-core-1_0.html#AddressClaim
type UserAddress struct {
	Formatted     string `json:"formatted,omitempty"`
	StreetAddress string `json:"street_address,omitempty"`
	Locality      string `json:"locality,omitempty"` // City or locality component.
	Region        string `json:"region,omitempty"`
	PostalCode    string `json:"postal_code,omitempty"`
	Country       string `json:"country,omitempty"`
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
