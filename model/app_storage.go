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
	// Offline is a boolean value that indicates wheter on not the app supports refresh tokens.
	// Do not use refresh tokens with apps that does not have secure storage.
	Offline() bool
	Type() AppType
	// RedirectURL is a redirect URL where to redirect the user after successfull login.
	// Useful not only for web apps, mobile and desktop apps could use custom scheme for that.
	RedirectURL() string
	// TokenLifespan is a token lifespan in seconds, if 0 - default one is used.
	TokenLifespan() int64
	// InviteTokenLifespan a inviteToken lifespan in seconds, if 0 - default one is used.
	InviteTokenLifespan() int64
	// RefreshTokenLifespan is a refreshToken lifespan in seconds, if 0 - default one is used.
	RefreshTokenLifespan() int64
	// Payload is a list of fields that are included in token. If it's empty, there are no fields in payload.
	TokenPayload() []string
	Sanitize() AppData
	RegistrationForbidden() bool
	AppleInfo() *AppleInfo
}

// AppType is a type of application.
type AppType string

const (
	// Web is a web app.
	Web AppType = "web"
	// Android is an Android app.
	Android AppType = "android"
	// IOS in an iOS app.
	IOS AppType = "ios"
)
