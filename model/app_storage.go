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
	ID           string    `bson:"_id" json:"id"` // TODO: use string?
	Secret       string    `bson:"secret" json:"secret"`
	Active       bool      `bson:"active" json:"active"`
	Name         string    `bson:"name" json:"name"`
	Description  string    `bson:"description" json:"description"`
	Scopes       []string  `bson:"scopes" json:"scopes"`   // Scopes is the list of all allowed scopes. If it's empty, no limitations (opaque scope).
	Offline      bool      `bson:"offline" json:"offline"` // Offline is a boolean value that indicates whether on not the app supports refresh tokens. Do not use refresh tokens with apps that does not have secure storage.
	Type         AppType   `bson:"type" json:"type"`
	RedirectURLs []string  `bson:"redirect_urls" json:"redirect_urls"` // RedirectURLs is the list of allowed urls where user will be redirected after successfull login. Useful not only for web apps, mobile and desktop apps could use custom scheme for that.
	TFAStatus    TFAStatus `bson:"tfa_status" json:"tfa_status"`
	DebugTFACode string    `bson:"debug_tfa_code" json:"debug_tfa_code"`

	// Authorization
	AuthzWay       AuthorizationWay `bson:"authorization_way" json:"authorization_way"`
	AuthzModel     string           `bson:"authorization_model" json:"authorization_model"`
	AuthzPolicy    string           `bson:"authorization_policy" json:"authorization_policy"`
	RolesWhitelist []string         `bson:"roles_whitelist" json:"roles_whitelist"`
	RolesBlacklist []string         `bson:"roles_blacklist" json:"roles_blacklist"`

	// Token settings
	TokenLifespan                     int64                                `bson:"token_lifespan" json:"token_lifespan"`                 // TokenLifespan is a token lifespan in seconds, if 0 - default one is used.
	InviteTokenLifespan               int64                                `bson:"invite_token_lifespan" json:"invite_token_lifespan"`   // InviteTokenLifespan a inviteToken lifespan in seconds, if 0 - default one is used.
	RefreshTokenLifespan              int64                                `bson:"refresh_token_lifespan" json:"refresh_token_lifespan"` // RefreshTokenLifespan is a refreshToken lifespan in seconds, if 0 - default one is used.
	TokenPayload                      []string                             `bson:"token_payload" json:"token_payload"`                   // Payload is a list of fields that are included in token. If it's empty, there are no fields in payload.
	TokenPayloadService               TokenPayloadServiceType              `json:"token_payload_service" bson:"token_payload_service"`
	TokenPayloadServicePluginSettings TokenPayloadServicePluginSettings    `json:"token_payload_service_plugin_settings" bson:"token_payload_service_plugin_settings"`
	TokenPayloadServiceHttpSettings   TokenPayloadServiceHttpSettings      `json:"token_payload_service_http_settings" bson:"token_payload_service_http_settings"`
	FederatedProviders                map[string]FederatedProviderSettings `json:"federated_login_settings" bson:"federated_login_settings"`

	// registration settings
	RegistrationForbidden        bool     `bson:"registration_forbidden" json:"registration_forbidden"`
	AnonymousRegistrationAllowed bool     `bson:"anonymous_registration_allowed" json:"anonymous_registration_allowed"`
	NewUserDefaultRole           string   `bson:"new_user_default_role" json:"new_user_default_role"`
	NewUserDefaultScopes         []string `bson:"new_user_default_scopes" json:"new_user_default_scopes"`
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

// TokenPayloadServiceType service to allow fetch additional data to include to access token
type TokenPayloadServiceType string

const (
	TokenPayloadServiceNone   = "none"   // TokenPayloadServiceNone no service is used
	TokenPayloadServicePlugin = "plugin" // TokenPayloadServicePlugin user local identifo plugin with specific name to retreive token payload
	TokenPayloadServiceHttp   = "http"   // TokenPayloadServiceHttp use external service to get token paylad
)

// TokenPayloadServicePluginSettings settings for token payload service
type TokenPayloadServicePluginSettings struct {
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}

// TokenPayloadServiceHttpSettings settings for token payload service
type TokenPayloadServiceHttpSettings struct {
	URL    string `json:"url,omitempty" bson:"url,omitempty"`
	Secret string `json:"secret,omitempty" bson:"secret,omitempty"`
}

// AppDataFromJSON unmarshal AppData from JSON string
func AppDataFromJSON(d []byte) (AppData, error) {
	var apd AppData
	if err := json.Unmarshal(d, &apd); err != nil {
		return AppData{}, err
	}
	return apd, nil
}

func (a AppData) Sanitized() AppData {
	a.Secret = ""
	a.AuthzWay = ""
	a.AuthzModel = ""
	a.AuthzPolicy = ""
	a.TokenPayloadServiceHttpSettings = TokenPayloadServiceHttpSettings{}
	a.TokenPayloadServicePluginSettings = TokenPayloadServicePluginSettings{}
	return a
}
