package model

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

const defaultEtcdKey = "identifo"

// ServerSettings are server settings.
type ServerSettings struct {
	General            GeneralServerSettings      `yaml:"general,omitempty" json:"general,omitempty"`
	AdminAccount       AdminAccountSettings       `yaml:"adminAccount,omitempty" json:"admin_account,omitempty"`
	Storage            StorageSettings            `yaml:"storage,omitempty" json:"storage,omitempty"`
	SessionStorage     SessionStorageSettings     `yaml:"sessionStorage,omitempty" json:"session_storage,omitempty"`
	StaticFilesStorage StaticFilesStorageSettings `yaml:"staticFilesStorage,omitempty" json:"static_files_storage,omitempty"`
	ExternalServices   ExternalServicesSettings   `yaml:"externalServices,omitempty" json:"external_services,omitempty"`
	Login              LoginSettings              `yaml:"login,omitempty" json:"login,omitempty"`
	KeyStorage         KeyStorageSettings         `yaml:"keyStorage,omitempty" json:"keyStorage,omitempty"`
	Config             ConfigStorageSettings      `yaml:"config,omitempty" json:"config,omitempty"`
	Logger             LoggerSettings             `yaml:"logger,omitempty" json:"logger,omitempty"`
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
	InviteStorage           DatabaseSettings `yaml:"inviteStorage,omitempty" json:"invite_storage,omitempty"`
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
	DBTypeBoltDB   DatabaseType = "boltdb"   // DBTypeBoltDB is for BoltDB.
	DBTypeMongoDB  DatabaseType = "mongodb"  // DBTypeMongoDB is for MongoDB.
	DBTypeDynamoDB DatabaseType = "dynamodb" // DBTypeDynamoDB is for DynamoDB.
	DBTypeFake     DatabaseType = "fake"     // DBTypeFake is for in-memory storage.
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
	ServeNewWeb      bool                   `yaml:"serveNewWeb,omitempty" json:"serve_new_web,omitempty"`
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

type ConfigStorageSettings struct {
	Type      ConfigStorageType    `json:"type,omitempty" yaml:"type,omitempty"`
	RawString string               `json:"raw_string,omitempty" yaml:"raw_string,omitempty"`
	S3        *S3StorageSettings   `json:"s3,omitempty" yaml:"s3,omitempty"`
	File      *FileStorageSettings `json:"file,omitempty" yaml:"file,omitempty"`
	Etcd      *EtcdStorageSettings `json:"etcd,omitempty" yaml:"etcd,omitempty"`
}

// ConfigStorageType describes type of configuration storage.
type ConfigStorageType string

const (
	// ConfigStorageTypeEtcd is an etcd storage.
	ConfigStorageTypeEtcd ConfigStorageType = "etcd"
	// ConfigurationStorageTypeS3 is an AWS S3 storage.
	ConfigStorageTypeS3 ConfigStorageType = "s3"
	// ConfigurationStorageTypeFile is a config file.
	ConfigStorageTypeFile ConfigStorageType = "file"
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
	Type KeyStorageType          `yaml:"type,omitempty" json:"type,omitempty"`
	S3   *S3KeyStorageSettings   `yaml:"s3,omitempty" json:"s3,omitempty"`
	File *KeyStorageFileSettings `yaml:"file,omitempty" json:"file,omitempty"`
}

type KeyStorageFileSettings struct {
	PrivateKeyPath string `json:"private_key_path,omitempty" yaml:"private_key_path,omitempty"`
	PublicKeyPath  string `json:"public_key_path,omitempty" yaml:"public_key_path,omitempty"`
}

type S3KeyStorageSettings struct {
	Region        string `yaml:"region,omitempty" json:"region,omitempty" bson:"region,omitempty"`
	Bucket        string `yaml:"bucket,omitempty" json:"bucket,omitempty" bson:"bucket,omitempty"`
	PublicKeyKey  string `yaml:"public_key_key,omitempty" json:"public_key_key,omitempty" bson:"public_key_key,omitempty"`
	PrivateKeyKey string `yaml:"private_key_key,omitempty" json:"private_key_key,omitempty" bson:"private_key_key,omitempty"`
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

	// Twilio related config.
	AccountSid string `yaml:"accountSid,omitempty" json:"account_sid,omitempty"`
	AuthToken  string `yaml:"authToken,omitempty" json:"auth_token,omitempty"`
	ServiceSid string `yaml:"serviceSid,omitempty" json:"service_sid,omitempty"`

	// Nexmo related config.
	APIKey    string `yaml:"apiKey,omitempty" json:"api_key,omitempty"`
	APISecret string `yaml:"apiSecret,omitempty" json:"api_secret,omitempty"`

	// RouteMobile related config.
	Username string `yaml:"username,omitempty" json:"username,omitempty"`
	Password string `yaml:"password,omitempty" json:"password,omitempty"`
	Source   string `yaml:"source,omitempty" json:"source,omitempty"`
	Region   string `yaml:"region,omitempty" json:"region,omitempty"`
}

// SMSServiceType - service for sending sms messages.
type SMSServiceType string

const (
	SMSServiceTwilio      SMSServiceType = "twilio"      // SMSServiceTwilio is a Twilio SMS service.
	SMSServiceNexmo       SMSServiceType = "nexmo"       // SMSServiceNexmo is a Nexmo SMS service.
	SMSServiceRouteMobile SMSServiceType = "routemobile" // SMSServiceRouteMobile is a RouteMobile SMS service.
	SMSServiceMock        SMSServiceType = "mock"        // SMSServiceMock is an SMS service mock.
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
	Email     bool `yaml:"email" json:"email,omitempty"`
	Federated bool `yaml:"federated" json:"federated,omitempty"`
}

// TFAType is a type of two-factor authentication for apps that support it.
type TFAType string

const (
	TFATypeApp   TFAType = "app"   // TFATypeApp is an app (like Google Authenticator).
	TFATypeSMS   TFAType = "sms"   // TFATypeSMS is an SMS.
	TFATypeEmail TFAType = "email" // TFATypeEmail is an email.
)

// GetPort returns port on which host listens to incoming connections.
func (ss ServerSettings) GetPort() string {
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

type LoggerSettings struct {
	DumpRequest bool `yaml:"dumpRequest,omitempty" json:"dumpRequest,omitempty"`
}

type FileStorageSettings struct {
	// just a file name
	FileName string `yaml:"file_name,omitempty" json:"file_name,omitempty" bson:"file_name,omitempty"`
}

type S3StorageSettings struct {
	Region string `yaml:"region,omitempty" json:"region,omitempty" bson:"region,omitempty"`
	Bucket string `yaml:"bucket,omitempty" json:"bucket,omitempty" bson:"bucket,omitempty"`
	Key    string `yaml:"key,omitempty" json:"key,omitempty" bson:"key,omitempty"`
}

type EtcdStorageSettings struct {
	Endpoints []string `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
	Key       string   `json:"key,omitempty" yaml:"key,omitempty"`
	Username  string   `json:"username,omitempty" yaml:"username,omitempty"`
	Password  string   `json:"password,omitempty" yaml:"password,omitempty"`
}

func ConfigStorageSettingsFromString(config string) (ConfigStorageSettings, error) {
	// Parse the URL and ensure there are no errors.
	u, err := url.Parse(config)
	if err != nil {
		return ConfigStorageSettings{}, fmt.Errorf("Unable to parse config string: %s", config)
	}

	switch strings.ToLower(u.Scheme) {
	case "etcd":
		return ConfigStorageSettingsFromStringEtcd(config)
	case "s3":
		return ConfigStorageSettingsFromStringS3(config)
	default:
		return ConfigStorageSettingsFromStringFile(config)
	}
}

func ConfigStorageSettingsFromStringS3(config string) (ConfigStorageSettings, error) {
	components := strings.Split(config[5:], "@")
	var pathComponents []string
	region := ""
	if len(components) == 2 {
		region = components[0]
		pathComponents = strings.Split(components[1], "/")
	} else if len(components) == 1 {
		pathComponents = strings.Split(components[0], "/")
	} else {
		return ConfigStorageSettings{}, fmt.Errorf("could not get s3 file path from config: %s", config)
	}
	if len(pathComponents) < 2 {
		return ConfigStorageSettings{}, fmt.Errorf("could not get s3 file path from config: %s", config)
	}
	bucket := pathComponents[0]
	path := strings.Join(pathComponents[1:], "/")

	return ConfigStorageSettings{
		Type:      ConfigStorageTypeS3,
		RawString: config,
		S3: &S3StorageSettings{
			Region: region,
			Bucket: bucket,
			Key:    path,
		},
	}, nil
}

func ConfigStorageSettingsFromStringFile(config string) (ConfigStorageSettings, error) {
	filename := config
	if strings.HasPrefix(strings.ToUpper(filename), "FILE://") {
		filename = filename[7:]
	}
	return ConfigStorageSettings{
		Type:      ConfigStorageTypeFile,
		RawString: config,
		File: &FileStorageSettings{
			FileName: filename,
		},
	}, nil
}

func ConfigStorageSettingsFromStringEtcd(config string) (ConfigStorageSettings, error) {
	result := ConfigStorageSettings{
		Type:      ConfigStorageTypeEtcd,
		RawString: config,
		Etcd: &EtcdStorageSettings{
			Key: defaultEtcdKey,
		},
	}
	var es string
	components := strings.Split(config[7:], "@")
	if len(components) > 1 {
		es = components[1]
		creds := strings.Split(components[0], ":")
		if len(creds) == 2 {
			result.Etcd.Username = creds[0]
			result.Etcd.Password = creds[1]
		}
	} else if len(components) == 1 {
		es = components[0]
	} else {
		return ConfigStorageSettings{}, fmt.Errorf("could not get etcd endpoints from config: %s", config)
	}

	components = strings.Split(es, "|")
	if len(components) > 1 {
		result.Etcd.Key = components[1]
	}
	result.Etcd.Endpoints = strings.Split(components[0], ",")
	return result, nil
}
