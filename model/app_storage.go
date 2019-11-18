package model

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
	NewAppData() AppData
	TestDatabaseConnection() error
	Close()
}

// AppData represents Application data information.
type AppData interface {
	ID() string
	Secret() string
	Active() bool
	Name() string
	Description() string
	// Scopes is the list of all allowed scopes. If it's empty, no limitations (opaque scope).
	Scopes() []string
	// Offline is a boolean value that indicates whether on not the app supports refresh tokens.
	// Do not use refresh tokens with apps that does not have secure storage.
	Offline() bool
	Type() AppType
	// RedirectURLs is the list of allowed urls where user will be redirected after successfull login.
	// Useful not only for web apps, mobile and desktop apps could use custom scheme for that.
	RedirectURLs() []string
	// TokenLifespan is a token lifespan in seconds, if 0 - default one is used.
	TokenLifespan() int64
	// InviteTokenLifespan a inviteToken lifespan in seconds, if 0 - default one is used.
	InviteTokenLifespan() int64
	// RefreshTokenLifespan is a refreshToken lifespan in seconds, if 0 - default one is used.
	RefreshTokenLifespan() int64
	// Payload is a list of fields that are included in token. If it's empty, there are no fields in payload.
	TokenPayload() []string
	Sanitize()
	TFAStatus() TFAStatus
	DebugTFACode() string
	RegistrationForbidden() bool
	AnonymousRegistrationAllowed() bool
	AuthzWay() AuthorizationWay
	AuthzModel() string
	AuthzPolicy() string
	RolesWhitelist() []string
	RolesBlacklist() []string
	NewUserDefaultRole() string
	AppleInfo() *AppleInfo
	SetSecret(secret string)
}

// AppType is a type of application.
type AppType string

const (
	// Web is a web app.
	Web AppType = "web"
	// Android is an Android app.
	Android AppType = "android"
	// IOS is an iOS app.
	IOS AppType = "ios"
	// Desktop is a desktop app.
	Desktop AppType = "desktop"
)

// AuthorizationWay is a way of authorization supported by the application.
type AuthorizationWay string

const (
	// NoAuthz is when the app does not require any authorization.
	NoAuthz AuthorizationWay = "no authorization"
	// Internal is for embedded authorization rules.
	Internal AuthorizationWay = "internal"
	// RolesWhitelist is the list of roles allowed to register and login into the application.
	RolesWhitelist AuthorizationWay = "whitelist"
	// RolesBlacklist is the list of roles forbidden to register and login into the application.
	RolesBlacklist AuthorizationWay = "blacklist"
	// External is for external authorization service.
	External AuthorizationWay = "external"
)

// TFAStatus is how the app supports two-factor authentication.
type TFAStatus string

const (
	// TFAStatusMandatory for mandatory TFA for all users.
	TFAStatusMandatory = "mandatory"
	// TFAStatusOptional for TFA that can be enabled/disabled for particular user.
	TFAStatusOptional = "optional"
	// TFAStatusDisabled is when the app does not support TFA.
	TFAStatusDisabled = "disabled"
)
