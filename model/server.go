package model

import (
	"io/fs"
	"net/http"
)

// Server holds together all dependencies.
type Server interface {
	Router() Router
	UpdateCORS()
	Storages() ServerStorageCollection
	Services() ServerServices
	Settings() ServerSettings
	Errors() []error
	Close()
}

// Router handles all incoming http requests.
type Router interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// ServerStorageCollection holds the full collections of server storage components
type ServerStorageCollection struct {
	App           AppStorage
	User          UserStorage
	UC            UserController
	UMC           UserMutationController
	Token         TokenStorage
	Blocklist     TokenBlacklist
	Invite        InviteStorage
	Config        ConfigurationStorage
	Session       SessionStorage
	Key           KeyStorage
	ManagementKey ManagementKeysStorage
	LoginAppFS    fs.FS
	AdminPanelFS  fs.FS
}

type ServerServices struct {
	SMS       SMSService
	Email     EmailService
	Token     TokenService
	Session   SessionService
	Challenge ChallengeController
}
