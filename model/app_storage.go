package model

// AppStorage is an abstract representation of applications data storage.
type AppStorage interface {
	AppByID(id string) (AppData, error)
	ActiveAppByID(appID string) (AppData, error)
	CreateApp(app AppData) (AppData, error)
	DisableApp(app AppData) error
	UpdateApp(appID string, newApp AppData) (AppData, error)
	FetchApps(filter string) ([]AppData, error)
	DeleteApp(id string) error
	ImportJSON(data []byte, cleanOldData bool) error
	TestDatabaseConnection() error
	Close()
}
