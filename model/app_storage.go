package model

//AppStorage is abstract representation of applications data storage
type AppStorage interface {
	//AppByID returns application data by AppID
	AppByID(id string) (AppData, error)
	AddNewApp(app AppData) (AppData, error)
	DisableApp(app AppData) error
	UpdateApp(oldAppID string, newApp AppData) error
}

//AppData represents Application data information
type AppData interface {
	ID() string
	Secret() string
	Active() bool
	Description() string
	//Scopes is the list of all allowed scopes
	//if it's empty, no limitations (opaque scope)
	Scopes() []string
	//Offline - indicated could the app has refresh token
	//don't use refresh tokens with apps, that don't have secure storage
	Offline() bool
	//RedirectURL - redirect URL, where to redirect the user after seccessfull login
	//useful not only for web apps, mobile and desktop app could use custom scheme for that
	RedirectURL() string
	//TokenLifespan Token lifespan in seconds, if 0 - use default one
	TokenLifespan() int64
	//RefreshTokenLifespan RefreshToken lifespan in seconds, if 0 - use default one
	RefreshTokenLifespan() int64
	//ResetPasswordTokenLifespan ResetPasswordToken lifespan in seconds, if 0 - use default one
	ResetPasswordTokenLifespan() int64
}
