package model

import (
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
	LoginVia             LoginVia                     `yaml:"loginVia,omitempty" json:"login_via,omitempty"`
	ServerConfigPath     string                       `yaml:"serverConfigPath,omitempty" json:"server_config_path,omitempty"`
	MailService          MailServiceType              `yaml:"mailService,omitempty" json:"mail_service,omitempty"`
	SMSService           SMSServiceSettings           `yaml:"smsService,omitempty" json:"sms_service,omitempty"`
	StaticFolderPath     string                       `yaml:"staticFolderPath,omitempty" json:"static_folder_path,omitempty"`
	EmailTemplatesPath   string                       `yaml:"emailTemplatesPath,omitempty" json:"email_templates_path,omitempty"`
	EmailTemplateNames   EmailTemplateNames           `yaml:"emailTemplateNames,omitempty" json:"email_template_names,omitempty"`
}

// ConfigurationStorageSettings holds together configuration storage settings.
type ConfigurationStorageSettings struct {
	Type        ConfigurationStorageType `yaml:"type,omitempty" json:"type,omitempty"`
	SettingsKey string                   `yaml:"settingsKey,omitempty" json:"settings_key,omitempty"`
	Endpoints   []string                 `yaml:"endpoints,omitempty" json:"endpoints,omitempty"`
}

// ConfigurationStorageType describes type of configuration storage.
type ConfigurationStorageType string

const (
	// ConfigurationStorageTypeEtcd is an etcd storage.
	ConfigurationStorageTypeEtcd ConfigurationStorageType = "etcd"
	// ConfigurationStorageTypeMock is a mocked storage.
	ConfigurationStorageTypeMock ConfigurationStorageType = "mock"
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
)

// StorageSettings holds together storage settings for different services.
type StorageSettings struct {
	AppStorage              DatabaseSettings `yaml:"appStorage,omitempty" json:"app_storage,omitempty"`
	UserStorage             DatabaseSettings `yaml:"userStorage,omitempty" json:"user_storage,omitempty"`
	TokenStorage            DatabaseSettings `yaml:"tokenStorage,omitempty" json:"token_storage,omitempty"`
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

// LoginVia is a type for configuring supported login ways.
type LoginVia struct {
	Username  bool `yaml:"username" json:"username,omitempty"`
	Phone     bool `yaml:"phone" json:"phone,omitempty"`
	Federated bool `yaml:"federated" json:"federated,omitempty"`
}

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
