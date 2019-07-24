package model

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// ServerSettings are server settings.
type ServerSettings struct {
	Host                 string                       `yaml:"host,omitempty" json:"host,omitempty"`
	PEMFolderPath        string                       `yaml:"pemFolderPath,omitempty" json:"pem_folder_path,omitempty"`
	PrivateKey           string                       `yaml:"privateKey,omitempty" json:"private_key,omitempty"`
	PublicKey            string                       `yaml:"publicKey,omitempty" json:"public_key,omitempty"`
	Issuer               string                       `yaml:"issuer,omitempty" json:"issuer,omitempty"`
	Algorithm            string                       `yaml:"algorithm,omitempty" json:"algorithm,omitempty"`
	ConfigurationStorage ConfigurationStorageSettings `yaml:"configurationStorage,omitempty" json:"configuration_storage,omitempty"`
	SessionStorage       SessionStorageSettings       `yaml:"sessionStorage,omitempty" json:"session_storage,omitempty"`
	Storage              StorageSettings              `yaml:"storage,omitempty" json:"storage,omitempty"`
	AdminAccount         AdminAccountSettings         `yaml:"adminAccount,omitempty" json:"admin_account,omitempty"`
	LoginWith            LoginWith                    `yaml:"loginWith,omitempty" json:"login_with,omitempty"`
	TFAType              TFAType                      `yaml:"tfaType,omitempty" json:"tfa_type,omitempty"`
	ServerConfigPath     string                       `yaml:"serverConfigPath,omitempty" json:"server_config_path,omitempty"`
	MailService          MailServiceType              `yaml:"mailService,omitempty" json:"mail_service,omitempty"`
	SMSService           SMSServiceSettings           `yaml:"smsService,omitempty" json:"sms_service,omitempty"`
	StaticFolderPath     string                       `yaml:"staticFolderPath,omitempty" json:"static_folder_path,omitempty"`
	EmailTemplatesPath   string                       `yaml:"emailTemplatesPath,omitempty" json:"email_templates_path,omitempty"`
	EmailTemplateNames   EmailTemplateNames           `yaml:"emailTemplateNames,omitempty" json:"email_template_names,omitempty"`
	AdminPanelBuildPath  string                       `yaml:"adminPanelBuildPath,omitempty" json:"admin_panel_build_path,omitempty"`
}

// ConfigurationStorageSettings holds together configuration storage settings.
type ConfigurationStorageSettings struct {
	Type        ConfigurationStorageType `yaml:"type,omitempty" json:"type,omitempty"`
	SettingsKey string                   `yaml:"settingsKey,omitempty" json:"settings_key,omitempty"`
	Endpoints   []string                 `yaml:"endpoints,omitempty" json:"endpoints,omitempty"`
	Region      string                   `yaml:"region,omitempty" json:"region,omitempty"`
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

// AdminAccountSettings are names of environment variables that store admin credentials.
type AdminAccountSettings struct {
	LoginEnvName    string `yaml:"loginEnvName" json:"login_env_name,omitempty"`
	PasswordEnvName string `yaml:"passwordEnvName" json:"password_env_name,omitempty"`
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

// SMSServiceSettings holds together settings for SMS service.
type SMSServiceSettings struct {
	Type       SMSServiceType `yaml:"type,omitempty" json:"type,omitempty"`
	AccountSid string         `yaml:"accountSid,omitempty" json:"account_sid,omitempty"`
	AuthToken  string         `yaml:"authToken,omitempty" json:"auth_token,omitempty"`
	ServiceSid string         `yaml:"serviceSid,omitempty" json:"service_sid,omitempty"`
}

// SMSServiceType - service to to use for sending sms.
type SMSServiceType string

const (
	// SMSServiceTwilio is a Twillo SMS service.
	SMSServiceTwilio SMSServiceType = "twilio"
	// SMSServiceMock is an SMS service mock.
	SMSServiceMock SMSServiceType = "mock"
)

// MailServiceType - how to send email to clients.
type MailServiceType string

const (
	// MailServiceMailgun is a Mailgun service.
	MailServiceMailgun = "mailgun"
	// MailServiceAWS is an AWS SES service.
	MailServiceAWS = "aws ses"
)

// GetPort returns port on which host listens to incoming connections.
func (ss *ServerSettings) GetPort() string {
	u, err := url.Parse(ss.Host)
	if err != nil {
		panic(err)
	}

	_, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		panic(err)
	}
	return strings.Join([]string{":", port}, "")
}

// Validate makes sure that all crucial fields are set.
func (ss *ServerSettings) Validate() error {
	if len(ss.AdminAccount.LoginEnvName) == 0 {
		return fmt.Errorf("Admin login env variable name not specified")
	}
	if len(ss.AdminAccount.PasswordEnvName) == 0 {
		return fmt.Errorf("Admin password env variable name not specified")
	}
	return nil
}
