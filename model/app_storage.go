package model

import "encoding/json"

// AppStorage is an abstract representation of applications data storage.
type AppStorage interface {
	AppByID(id string) (AppData, error)
	ActiveAppByID(appID string) (AppData, error)
	CreateApp(app AppData) (AppData, error)
	DisableApp(app AppData) error
	UpdateApp(appID string, newApp AppData) (AppData, error)
	FetchApps(filterString string, skip, limit int) ([]AppData, int, error)
	DeleteApp(id string) error
	ImportJSON(data []byte) error
	TestDatabaseConnection() error
	Close()
}

// AppData represents Application data information.
type AppData struct {
	ID                           string           `bson:"_id,omitempty" json:"id,omitempty"` // TODO: use string?
	Secret                       string           `bson:"secret,omitempty" json:"secret,omitempty"`
	Active                       bool             `bson:"active" json:"active"`
	Name                         string           `bson:"name,omitempty" json:"name,omitempty"`
	Description                  string           `bson:"description,omitempty" json:"description,omitempty"`
	Scopes                       []string         `bson:"scopes,omitempty" json:"scopes,omitempty"` // Scopes is the list of all allowed scopes. If it's empty, no limitations (opaque scope).
	Offline                      bool             `bson:"offline" json:"offline"`                   // Offline is a boolean value that indicates whether on not the app supports refresh tokens. Do not use refresh tokens with apps that does not have secure storage.
	Type                         AppType          `bson:"type,omitempty" json:"type,omitempty"`
	RedirectURLs                 []string         `bson:"redirect_urls,omitempty" json:"redirect_urls,omitempty"`                   // RedirectURLs is the list of allowed urls where user will be redirected after successfull login. Useful not only for web apps, mobile and desktop apps could use custom scheme for that.
	TokenLifespan                int64            `bson:"refresh_token_lifespan,omitempty" json:"refresh_token_lifespan,omitempty"` // TokenLifespan is a token lifespan in seconds, if 0 - default one is used.
	InviteTokenLifespan          int64            `bson:"invite_token_lifespan,omitempty" json:"invite_token_lifespan,omitempty"`   // InviteTokenLifespan a inviteToken lifespan in seconds, if 0 - default one is used.
	RefreshTokenLifespan         int64            `bson:"token_lifespan,omitempty" json:"token_lifespan,omitempty"`                 // RefreshTokenLifespan is a refreshToken lifespan in seconds, if 0 - default one is used.
	TokenPayload                 []string         `bson:"token_payload,omitempty" json:"token_payload,omitempty"`                   // Payload is a list of fields that are included in token. If it's empty, there are no fields in payload.
	TFAStatus                    TFAStatus        `bson:"tfa_status" json:"tfa_status"`
	DebugTFACode                 string           `bson:"debug_tfa_code,omitempty" json:"debug_tfa_code,omitempty"`
	RegistrationForbidden        bool             `bson:"registration_forbidden" json:"registration_forbidden"`
	AnonymousRegistrationAllowed bool             `bson:"anonymous_registration_allowed" json:"anonymous_registration_allowed"`
	AuthzWay                     AuthorizationWay `bson:"authorization_way,omitempty" json:"authorization_way,omitempty"`
	AuthzModel                   string           `bson:"authorization_model,omitempty" json:"authorization_model,omitempty"`
	AuthzPolicy                  string           `bson:"authorization_policy,omitempty" json:"authorization_policy,omitempty"`
	RolesWhitelist               []string         `bson:"roles_whitelist,omitempty" json:"roles_whitelist,omitempty"`
	RolesBlacklist               []string         `bson:"roles_blacklist,omitempty" json:"roles_blacklist,omitempty"`
	NewUserDefaultRole           string           `bson:"new_user_default_role,omitempty" json:"new_user_default_role,omitempty"`
	AppleInfo                    *AppleInfo       `bson:"apple_info,omitempty" json:"apple_info,omitempty"`
}

// AppType is a type of application.
type AppType string

const (
	Web     AppType = "web"     // Web is a web app.
	Android AppType = "android" // Android is an Android app.
	IOS     AppType = "ios"     // IOS is an iOS app.
	Desktop AppType = "desktop" // Desktop is a desktop app.
)

// AuthorizationWay is a way of authorization supported by the application.
type AuthorizationWay string

const (
	NoAuthz        AuthorizationWay = "no authorization" // NoAuthz is when the app does not require any authorization.
	Internal       AuthorizationWay = "internal"         // Internal is for embedded authorization rules.
	RolesWhitelist AuthorizationWay = "whitelist"        // RolesWhitelist is the list of roles allowed to register and login into the application.
	RolesBlacklist AuthorizationWay = "blacklist"        // RolesBlacklist is the list of roles forbidden to register and login into the application.
	External       AuthorizationWay = "external"         // External is for external authorization service.
)

// TFAStatus is how the app supports two-factor authentication.
type TFAStatus string

const (
	TFAStatusMandatory = "mandatory" // TFAStatusMandatory for mandatory TFA for all users.
	TFAStatusOptional  = "optional"  // TFAStatusOptional for TFA that can be enabled/disabled for particular user.
	TFAStatusDisabled  = "disabled"  // TFAStatusDisabled is when the app does not support TFA.
)

// AppDataFromJSON unmarshal AppData from JSON string
func AppDataFromJSON(d []byte) (AppData, error) {
	var apd AppData
	if err := json.Unmarshal(d, &apd); err != nil {
		return AppData{}, err
	}
	return apd, nil
}
