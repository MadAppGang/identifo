package model

import "time"

// Default Server Settings settings
var DefaultServerSettings = ServerSettings{
	General: GeneralServerSettings{
		Host:   "http://localhost:8081",
		Port:   "8081",
		Issuer: "Identifo",
	},
	Storage: StorageSettings{
		DefaultStorage:          DatabaseSettings{Type: DBTypeMem},
		AppStorage:              DatabaseSettings{Type: DBTypeDefault},
		UserStorage:             DatabaseSettings{Type: DBTypeDefault},
		TokenStorage:            DatabaseSettings{Type: DBTypeDefault},
		TokenBlacklist:          DatabaseSettings{Type: DBTypeDefault},
		VerificationCodeStorage: DatabaseSettings{Type: DBTypeDefault},
		InviteStorage:           DatabaseSettings{Type: DBTypeDefault},
		ManagementKeysStorage:   DatabaseSettings{Type: DBTypeDefault},
	},
	SessionStorage: SessionStorageSettings{
		Type:            SessionStorageMem,
		SessionDuration: SessionDuration{Duration: time.Second * 300},
	},
	KeyStorage: FileStorageSettings{
		Type: FileStorageTypeLocal,
		Local: FileStorageLocal{
			Path: "./jwt/test_artifacts/private.pem",
		},
	},
	Login: LoginSettings{
		LoginWith: LoginWith{
			Phone:     true,
			Username:  true,
			Federated: false,
		},
		TFAType:          TFATypeApp,
		TFAResendTimeout: 30,
	},
	Services: ServicesSettings{
		Email: EmailServiceSettings{
			Type: EmailServiceMock,
		},
		SMS: SMSServiceSettings{
			Type: SMSServiceMock,
		},
	},
	AdminPanel:     AdminPanelSettings{Enabled: true},
	LoginWebApp:    FileStorageSettings{Type: FileStorageTypeNone},
	EmailTemplates: FileStorageSettings{Type: FileStorageTypeNone},
	AdminAccount: AdminAccountSettings{
		LoginEnvName:    "IDENTIFO_ADMIN_LOGIN",
		PasswordEnvName: "IDENTIFO_ADMIN_PASSWORD",
	},
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

	if len(ss.AdminAccount.LoginEnvName) == 0 {
		ss.AdminAccount = DefaultServerSettings.AdminAccount
	}

	if len(ss.Storage.AppStorage.Type) == 0 {
		ss.Storage.AppStorage.Type = DBTypeDefault
	}
	if len(ss.Storage.UserStorage.Type) == 0 {
		ss.Storage.UserStorage.Type = DBTypeDefault
	}
	if len(ss.Storage.TokenBlacklist.Type) == 0 {
		ss.Storage.TokenStorage.Type = DBTypeDefault
	}
	if len(ss.Storage.TokenBlacklist.Type) == 0 {
		ss.Storage.TokenStorage.Type = DBTypeDefault
	}
	if len(ss.Storage.VerificationCodeStorage.Type) == 0 {
		ss.Storage.VerificationCodeStorage.Type = DBTypeDefault
	}
	if len(ss.Storage.InviteStorage.Type) == 0 {
		ss.Storage.InviteStorage.Type = DBTypeDefault
	}
	if len(ss.Storage.TokenBlacklist.Type) == 0 {
		ss.Storage.TokenBlacklist.Type = DBTypeDefault
	}
}
