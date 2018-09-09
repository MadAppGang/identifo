package model

//AppStorage is abstract representation of applications data storage
type AppStorage interface {
	//AppByID returns application data by AppID
	AppByID(id string) (AppData, error)
	AddNewApp(app AppData) error
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
}
