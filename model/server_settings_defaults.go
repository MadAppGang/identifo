package model

// Default Server Settings settings
var DefaultServerSettings = ServerSettings{
	LoginWebApp:    FileStorageSettings{Type: FileStorageTypeNone},
	EmailTemplates: FileStorageSettings{Type: FileStorageTypeNone},
}

// Check server settings and apply changes if needed
func (ss *ServerSettings) RewriteDefaults() {
	// if login web app empty - set default values
	if len(ss.LoginWebApp.Type) == 0 {
		ss.LoginWebApp = DefaultServerSettings.LoginWebApp
	}

	if len(ss.EmailTemplates.Type) == 0 {
		ss.EmailTemplates = DefaultServerSettings.EmailTemplates
	}
}
