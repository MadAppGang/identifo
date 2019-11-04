package model

import (
	"net"
	"net/url"
	"strings"
)

// ServerSettings are server settings.
type ServerSettings struct {
	General              GeneralServerSettings        `yaml:"general,omitempty" json:"general,omitempty"`
	AdminAccount         AdminAccountSettings         `yaml:"adminAccount,omitempty" json:"admin_account,omitempty"`
	Storage              StorageSettings              `yaml:"storage,omitempty" json:"storage,omitempty"`
	ConfigurationStorage ConfigurationStorageSettings `yaml:"configurationStorage,omitempty" json:"configuration_storage,omitempty"`
	SessionStorage       SessionStorageSettings       `yaml:"sessionStorage,omitempty" json:"session_storage,omitempty"`
	StaticFilesStorage   StaticFilesStorageSettings   `yaml:"staticFilesStorage,omitempty" json:"static_files_storage,omitempty"`
	ExternalServices     ExternalServicesSettings     `yaml:"externalServices,omitempty" json:"external_services,omitempty"`
	Login                LoginSettings                `yaml:"login,omitempty" json:"login,omitempty"`
}

// GeneralServerSettings are general server settings.
type GeneralServerSettings struct {
	Host      string `yaml:"host,omitempty" json:"host,omitempty"`
	Issuer    string `yaml:"issuer,omitempty" json:"issuer,omitempty"`
	Algorithm string `yaml:"algorithm,omitempty" json:"algorithm,omitempty"`
}

// AdminAccountSettings are names of environment variables that store admin credentials.
type AdminAccountSettings struct {
	LoginEnvName    string `yaml:"loginEnvName" json:"login_env_name,omitempty"`
	PasswordEnvName string `yaml:"passwordEnvName" json:"password_env_name,omitempty"`
}

// StorageSettings holds together storage settings for different services.
type StorageSettings struct {
	AppStorage              DatabaseSettings `yaml:"appStorage,omitempty" json:"app_storage,omitempty"`
	UserStorage             DatabaseSettings `yaml:"userStorage,omitempty" json:"user_storage,omitempty"`
	TokenStorage            DatabaseSettings `yaml:"tokenStorage,omitempty" json:"token_storage,omitempty"`
	TokenBlacklist          DatabaseSettings `yaml:"tokenBlacklist,omitempty" json:"token_blacklist,omitempty"`
	VerificationCodeStorage DatabaseSettings `yaml:"verificationCodeStorage,omitempty" json:"verification_code_storage,omitempty"`
}

// DatabaseSettings holds together all settings applicable to a particular database.
type DatabaseSettings struct {
	Type     DatabaseType `yaml:"type,omitempty" json:"type,omitempty"`
	Name     string       `yaml:"name,omitempty" json:"name,omitempty"`
	Endpoint string       `yaml:"endpoint,omitempty" json:"endpoint,omitempty"`
	Region   string       `yaml:"region,omitempty" json:"region,omitempty"`
	Path     string       `yaml:"path,omitempty" json:"path,omitempty"`
}

// DatabaseType is a type of database.
type DatabaseType string

const (
	// DBTypeBoltDB is for BoltDB.
	DBTypeBoltDB DatabaseType = "boltdb"
	// DBTypeMongoDB is for MongoDB.
	DBTypeMongoDB DatabaseType = "mongodb"
	// DBTypeDynamoDB is for DynamoDB.
	DBTypeDynamoDB DatabaseType = "dynamodb"
	// DBTypeFake is for in-memory storage.
	DBTypeFake DatabaseType = "fake"
)

// StaticFilesStorageSettings are settings for static files storage.
type StaticFilesStorageSettings struct {
	Type             StaticFilesStorageType `yaml:"type,omitempty" json:"type,omitempty"`
	ServerConfigPath string                 `yaml:"serverConfigPath,omitempty" json:"server_config_path,omitempty"`
	Folder           string                 `yaml:"folder,omitempty" json:"folder,omitempty"`
	Bucket           string                 `yaml:"bucket,omitempty" json:"bucket,omitempty"`
	Region           string                 `yaml:"region,omitempty" json:"region,omitempty"`
	Endpoint         string                 `yaml:"endpoint,omitempty" json:"endpoint,omitempty"`
	ServeAdminPanel  bool                   `yaml:"serveAdminPanel,omitempty" json:"serve_admin_panel,omitempty"`
}

// StaticFilesStorageType is a type of static files storage.
type StaticFilesStorageType string

const (
	// StaticFilesStorageTypeLocal is for storing static files locally.
	StaticFilesStorageTypeLocal = "local"
	// StaticFilesStorageTypeS3 is for storing static files in S3 bucket.
	StaticFilesStorageTypeS3 = "s3"
	// StaticFilesStorageTypeDynamoDB is for storing static files in DynamoDB table.
	StaticFilesStorageTypeDynamoDB = "dynamodb"
)

// ConfigurationStorageSettings holds together configuration storage settings.
type ConfigurationStorageSettings struct {
	Type        ConfigurationStorageType `yaml:"type,omitempty" json:"type,omitempty"`
	SettingsKey string                   `yaml:"settingsKey,omitempty" json:"settings_key,omitempty"`
	Endpoints   []string                 `yaml:"endpoints,omitempty" json:"endpoints,omitempty"`
	Bucket      string                   `yaml:"bucket,omitempty" json:"bucket,omitempty"`
	Region      string                   `yaml:"region,omitempty" json:"region,omitempty"`
	KeyStorage  KeyStorageSettings       `yaml:"keyStorage,omitempty" json:"key_storage,omitempty"`
}

// ConfigurationStorageType describes type of configuration storage.
type ConfigurationStorageType string

const (
	// ConfigurationStorageTypeEtcd is an etcd storage.
	ConfigurationStorageTypeEtcd ConfigurationStorageType = "etcd"
	// ConfigurationStorageTypeS3 is an AWS S3 storage.
	ConfigurationStorageTypeS3 ConfigurationStorageType = "s3"
	// ConfigurationStorageTypeFile is a config file.
	ConfigurationStorageTypeFile ConfigurationStorageType = "file"
)

// SessionStorageSettings holds together session storage settings.
type SessionStorageSettings struct {
	Type            SessionStorageType `yaml:"type,omitempty" json:"type,omitempty"`
	SessionDuration SessionDuration    `yaml:"sessionDuration,omitempty" json:"session_duration,omitempty"`
	Address         string             `yaml:"address,omitempty" json:"address,omitempty"`
	Password        string             `yaml:"password,omitempty" json:"password,omitempty"`
	DB              int                `yaml:"db,omitempty" json:"db,omitempty"`
	Region          string             `yaml:"region,omitempty" json:"region,omitempty"`
	Endpoint        string             `yaml:"endpoint,omitempty" json:"endpoint,omitempty"`
}

// SessionStorageType - where to store admin sessions.
type SessionStorageType string

const (
	// SessionStorageMem means to store sessions in memory.
	SessionStorageMem = "memory"
	// SessionStorageRedis means to store sessions in Redis.
	SessionStorageRedis = "redis"
	// SessionStorageDynamoDB means to store sessions in DynamoDB.
	SessionStorageDynamoDB = "dynamodb"
)

// KeyStorageSettings are settings for the key storage.
type KeyStorageSettings struct {
	Type   KeyStorageType `yaml:"type,omitempty" json:"type,omitempty"`
	Folder string         `yaml:"folder,omitempty" json:"folder,omitempty"`
	Region string         `yaml:"region,omitempty" json:"region,omitempty"`
	Bucket string         `yaml:"bucket,omitempty" json:"bucket,omitempty"`
}

// KeyStorageType is a type of the key storage.
type KeyStorageType string

const (
	// KeyStorageTypeLocal is for storing keys locally.
	KeyStorageTypeLocal = "local"
	// KeyStorageTypeS3 is for storing keys in the S3 bucket.
	KeyStorageTypeS3 = "s3"
)

// ExternalServicesSettings are settings for external services.
type ExternalServicesSettings struct {
	EmailService EmailServiceSettings `yaml:"emailService,omitempty" json:"email_service,omitempty"`
	SMSService   SMSServiceSettings   `yaml:"smsService,omitempty" json:"sms_service,omitempty"`
}

// EmailServiceType - how to send email to clients.
type EmailServiceType string

const (
	// EmailServiceMailgun is a Mailgun service.
	EmailServiceMailgun = "mailgun"
	// EmailServiceAWS is an AWS SES service.
	EmailServiceAWS = "aws ses"
	// EmailServiceMock is an email service mock.
	EmailServiceMock = "mock"
)

// EmailServiceSettings holds together settings for the email service.
type EmailServiceSettings struct {
	Type       EmailServiceType `yaml:"type,omitempty" json:"type,omitempty"`
	Domain     string           `yaml:"domain,omitempty" json:"domain,omitempty"`
	PublicKey  string           `yaml:"publicKey,omitempty" json:"public_key,omitempty"`
	PrivateKey string           `yaml:"privateKey,omitempty" json:"private_key,omitempty"`
	Sender     string           `yaml:"sender,omitempty" json:"sender,omitempty"`
	Region     string           `yaml:"region,omitempty" json:"region,omitempty"`
}

// SMSServiceSettings holds together settings for SMS service.
type SMSServiceSettings struct {
	Type SMSServiceType `yaml:"type,omitempty" json:"type,omitempty"`
	// Twilio related config
	AccountSid string `yaml:"accountSid,omitempty" json:"account_sid,omitempty"`
	AuthToken  string `yaml:"authToken,omitempty" json:"auth_token,omitempty"`
	ServiceSid string `yaml:"serviceSid,omitempty" json:"service_sid,omitempty"`
	// Nexmo related config
	APIKey    string `yaml:"apiKey,omitempty" json:"api_key,omitempty"`
	APISecret string `yaml:"apiSecret,omitempty" json:"api_secret,omitempty"`
}

// SMSServiceType - service for sending sms messages.
type SMSServiceType string

const (
	// SMSServiceTwilio is a Twillo SMS service.
	SMSServiceTwilio SMSServiceType = "twilio"
	// SMSServiceNexmo is a Nexmo SMS service.
	SMSServiceNexmo SMSServiceType = "nexmo"
	// SMSServiceMock is an SMS service mock.
	SMSServiceMock SMSServiceType = "mock"
)

// LoginSettings are settings of login.
type LoginSettings struct {
	LoginWith LoginWith `yaml:"loginWith,omitempty" json:"login_with,omitempty"`
	TFAType   TFAType   `yaml:"tfaType,omitempty" json:"tfa_type,omitempty"`
}

// LoginWith is a type for configuring supported login ways.
type LoginWith struct {
	Username  bool `yaml:"username" json:"username,omitempty"`
	Phone     bool `yaml:"phone" json:"phone,omitempty"`
	Federated bool `yaml:"federated" json:"federated,omitempty"`
}

// TFAType is a type of two-factor authentication for apps that support it.
type TFAType string

const (
	// TFATypeApp is an app (like Google Authenticator).
	TFATypeApp TFAType = "app"
	// TFATypeSMS is an SMS.
	TFATypeSMS TFAType = "sms"
	// TFATypeEmail is an email.
	TFATypeEmail TFAType = "email"
)

// GetPort returns port on which host listens to incoming connections.
func (ss *ServerSettings) GetPort() string {
	u, err := url.Parse(ss.General.Host)
	if err != nil {
		panic(err)
	}

	_, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		panic(err)
	}
	return strings.Join([]string{":", port}, "")
}
