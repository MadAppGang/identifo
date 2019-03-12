package model

//AppStorage is abstract representation of applications data storage
type AppStorage interface {
	//AppByID returns application data by AppID
	AppByID(id string) (AppData, error)
	ActiveAppByID(appID string) (AppData, error)
	CreateApp(app AppData) (AppData, error)
	DisableApp(app AppData) error
	UpdateApp(oldAppID string, newApp AppData) error
	FetchApps(filterString string, skip, limit int) ([]AppData, error)
	DeleteApp(id string) error
	ImportJSON(data []byte) error
	NewAppData() AppData
}

//AppData represents Application data information
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
	// RedirectURL is a redirect URL where to redirect the user after successfull login.
	// Useful not only for web apps, mobile and desktop apps could use custom scheme for that.
	RedirectURL() string
	// TokenLifespan is a token lifespan in seconds, if 0 - default one is used.
	TokenLifespan() int64
	// RefreshTokenLifespan is a refreshToken lifespan in seconds, if 0 - default one is used.
	RefreshTokenLifespan() int64
	// Payload is a list of fields that are included in token. If it's empty, there are no fields in payload.
	TokenPayload() []string
	Sanitize() AppData
}
