package model

import (
	"fmt"
	"net/url"
	"strings"
)

const defaultEtcdKey = "identifo"

// ServerSettings are server settings.
type ServerSettings struct {
	General        GeneralServerSettings  `yaml:"general" json:"general"`
	AdminAccount   AdminAccountSettings   `yaml:"adminAccount" json:"admin_account"`
	Storage        StorageSettings        `yaml:"storage" json:"storage"`
	SessionStorage SessionStorageSettings `yaml:"sessionStorage" json:"session_storage"`
	Services       ServicesSettings       `yaml:"services" json:"external_services"`
	Login          LoginSettings          `yaml:"login" json:"login"`
	KeyStorage     KeyStorageSettings     `yaml:"keyStorage" json:"key_storage"`
	Config         ConfigStorageSettings  `yaml:"-" json:"config"`
	Logger         LoggerSettings         `yaml:"logger" json:"logger"`
	AdminPanel     AdminPanelSettings     `yaml:"adminPanel" json:"admin_panel"`
	LoginWebApp    FileStorageSettings    `yaml:"loginWebApp" json:"login_web_app"`
	EmailTemplates FileStorageSettings    `yaml:"emailTemplaits" json:"email_templaits"`
}

// GeneralServerSettings are general server settings.
type GeneralServerSettings struct {
	Host            string   `yaml:"host" json:"host"`
	Port            string   `yaml:"port" json:"port"`
	Issuer          string   `yaml:"issuer" json:"issuer"`
	SupportedScopes []string `yaml:"supported_scopes" json:"supported_scopes"`
}

// AdminAccountSettings are names of environment variables that store admin credentials.
type AdminAccountSettings struct {
	LoginEnvName    string `yaml:"loginEnvName" json:"login_env_name"`
	PasswordEnvName string `yaml:"passwordEnvName" json:"password_env_name"`
}

// StorageSettings holds together storage settings for different services.
type StorageSettings struct {
	AppStorage              DatabaseSettings `yaml:"appStorage" json:"app_storage"`
	UserStorage             DatabaseSettings `yaml:"userStorage" json:"user_storage"`
	TokenStorage            DatabaseSettings `yaml:"tokenStorage" json:"token_storage"`
	TokenBlacklist          DatabaseSettings `yaml:"tokenBlacklist" json:"token_blacklist"`
	VerificationCodeStorage DatabaseSettings `yaml:"verificationCodeStorage" json:"verification_code_storage"`
	InviteStorage           DatabaseSettings `yaml:"inviteStorage" json:"invite_storage"`
}

// DatabaseSettings holds together all settings applicable to a particular database.
type DatabaseSettings struct {
	Type   DatabaseType           `yaml:"type" json:"type"`
	BoltDB BoltDBDatabaseSettings `yaml:"boltdb" json:"boltdb"`
	Mongo  MongodDatabaseSettings `yaml:"mongo" json:"mongo"`
	Dynamo DynamoDatabaseSettings `yaml:"dynamo" json:"dynamo"`
	Plugin PluginSettings         `yaml:"plugin" json:"plugin"`
	GRPC   GRPCSettings           `yaml:"grpc" json:"grpc"`
}

type BoltDBDatabaseSettings struct {
	Path string `yaml:"path" json:"path"`
}

type MongodDatabaseSettings struct {
	ConnectionString string `yaml:"connection" json:"connection"`
	DatabaseName     string `yaml:"database" json:"database"`
}

type DynamoDatabaseSettings struct {
	Region   string `yaml:"region" json:"region"`
	Endpoint string `yaml:"endpoint" json:"endpoint"`
}

type PluginSettings struct {
	Cmd    string            `yaml:"cmd" json:"cmd"`
	Params map[string]string `yaml:"params" json:"params"`
}

type GRPCSettings struct {
	Address string `yaml:"address" json:"address"`
}

// DatabaseType is a type of database.
type DatabaseType string

const (
	DBTypeBoltDB   DatabaseType = "boltdb" // DBTypeBoltDB is for BoltDB.
	DBTypeMongoDB  DatabaseType = "mongo"  // DBTypeMongoDB is for MongoDB.
	DBTypeDynamoDB DatabaseType = "dynamo" // DBTypeDynamoDB is for DynamoDB.
	DBTypeFake     DatabaseType = "fake"   // DBTypeFake is for in-memory storage.
	DBTypePlugin   DatabaseType = "plugin" // DBTypePlugin is used for hashicorp/go-plugin.
	DBTypeGRPC     DatabaseType = "grpc"   // DBTypeGRPC is used for pure grpc.
)

type FileStorageSettings struct {
	Type  FileStorageType  `yaml:"type" json:"type"`
	Local FileStorageLocal `yaml:"local,omitempty" json:"local,omitempty"`
	S3    FileStorageS3    `yaml:"s3,omitempty" json:"s3,omitempty"`
}

type FileStorageType string

const (
	FileStorageTypeNone    FileStorageType = "none"
	FileStorageTypeDefault FileStorageType = "default"
	FileStorageTypeLocal   FileStorageType = "local"
	FileStorageTypeS3      FileStorageType = "s3"
)

type FileStorageS3 struct {
	Region string `yaml:"region" json:"region"`
	Bucket string `yaml:"bucket" json:"bucket"`
	Folder string `yaml:"folder" json:"folder"`
}

type FileStorageLocal struct {
	FolderPath string `yaml:"folder" json:"folder"`
}

type ConfigStorageSettings struct {
	Type      ConfigStorageType         `json:"type"`
	RawString string                    `json:"raw_string"`
	S3        *S3StorageSettings        `json:"s3"`
	File      *LocalFileStorageSettings `json:"file"`
	Etcd      *EtcdStorageSettings      `json:"etcd"`
}

// ConfigStorageType describes type of configuration storage.
type ConfigStorageType string

const (
	// ConfigStorageTypeEtcd is an etcd storage.
	// TODO: etcd not supported now
	// ConfigStorageTypeEtcd ConfigStorageType = "etcd"
	// ConfigurationStorageTypeS3 is an AWS S3 storage.
	ConfigStorageTypeS3 ConfigStorageType = "s3"
	// ConfigurationStorageTypeFile is a config file.
	ConfigStorageTypeFile ConfigStorageType = "file"
)

// SessionStorageSettings holds together session storage settings.
type SessionStorageSettings struct {
	Type            SessionStorageType     `yaml:"type" json:"type"`
	SessionDuration SessionDuration        `yaml:"sessionDuration" json:"session_duration"`
	Redis           RedisDatabaseSettings  `yaml:"redis" json:"redis"`
	Dynamo          DynamoDatabaseSettings `yaml:"dynamo" json:"dynamo"`
}

// SessionStorageType - where to store admin sessions.
type SessionStorageType string

const (
	// SessionStorageMem means to store sessions in memory.
	SessionStorageMem = "memory"
	// SessionStorageRedis means to store sessions in Redis.
	SessionStorageRedis = "redis"
	// SessionStorageDynamoDB means to store sessions in DynamoDB.
	SessionStorageDynamoDB = "dynamo"
)

// RedisDatabaseSettings redis storage settings
type RedisDatabaseSettings struct {
	// host:port address.
	Address string `yaml:"address" json:"address"`
	// Optional password. Must match the password specified in the
	// requirepass server configuration option.
	Password string `yaml:"password" json:"password"`
	// Database to be selected after connecting to the server.
	DB int `yaml:"db" json:"db"`
	// Cluster - if true will connect to redis cluster, address can be comma separated list of addresses.
	Cluster bool `yaml:"cluster" json:"cluster"`
	// Prefix for redis keys
	Prefix string `yaml:"prefix" json:"prefix"`
}

type DynamoDBSessionStorageSettings struct{}

// KeyStorageSettings are settings for the key storage.
type KeyStorageSettings struct {
	Type KeyStorageType         `yaml:"type" json:"type"`
	S3   S3KeyStorageSettings   `yaml:"s3" json:"s3"`
	File KeyStorageFileSettings `yaml:"file" json:"file"`
}

type KeyStorageFileSettings struct {
	PrivateKeyPath string `json:"private_key_path" yaml:"private_key_path"`
}

type S3KeyStorageSettings struct {
	Region        string `yaml:"region" json:"region,omitempty" bson:"region"`
	Bucket        string `yaml:"bucket" json:"bucket,omitempty" bson:"bucket"`
	PrivateKeyKey string `yaml:"private_key_key" json:"private_key_key" bson:"private_key_key"`
}

// KeyStorageType is a type of the key storage.
type KeyStorageType string

const (
	// KeyStorageTypeLocal is for storing keys locally.
	KeyStorageTypeLocal = "local"
	// KeyStorageTypeS3 is for storing keys in the S3 bucket.
	KeyStorageTypeS3 = "s3"
)

// ServicesSettings are settings for external services.
type ServicesSettings struct {
	Email EmailServiceSettings `yaml:"email" json:"email_service"`
	SMS   SMSServiceSettings   `yaml:"sms" json:"sms_service"`
}

// EmailServiceType - how to send email to clients.
type EmailServiceType string

const (
	// EmailServiceMailgun is a Mailgun service.
	EmailServiceMailgun = "mailgun"
	// EmailServiceAWS is an AWS SES service.
	EmailServiceAWS = "ses"
	// EmailServiceMock is an email service mock.
	EmailServiceMock = "mock"
)

// EmailServiceSettings holds together settings for the email service.
type EmailServiceSettings struct {
	Type    EmailServiceType            `yaml:"type" json:"type"`
	Mailgun MailgunEmailServiceSettings `yaml:"mailgun" json:"mailgun"`
	SES     SESEmailServiceSettings     `yaml:"ses" json:"ses"`
}

type MailgunEmailServiceSettings struct {
	Domain     string `yaml:"domain" json:"domain"`
	PrivateKey string `yaml:"privateKey" json:"private_key"`
	PublicKey  string `yaml:"publicKey" json:"public_key"`
	Sender     string `yaml:"sender" json:"sender"`
}

type SESEmailServiceSettings struct {
	Region string `yaml:"region" json:"region"`
	Sender string `yaml:"sender" json:"sender"`
}

// SMSServiceSettings holds together settings for SMS service.
type SMSServiceSettings struct {
	Type        SMSServiceType             `yaml:"type" json:"type"`
	Twilio      TwilioServiceSettings      `yaml:"twilio" json:"twilio"`
	Nexmo       NexmoServiceSettings       `yaml:"nexmo" json:"nexmo"`
	Routemobile RouteMobileServiceSettings `yaml:"routemobile" json:"routemobile"`
}

// SMSServiceType - service for sending sms messages.
type SMSServiceType string

const (
	SMSServiceTwilio      SMSServiceType = "twilio"      // SMSServiceTwilio is a Twilio SMS service.
	SMSServiceNexmo       SMSServiceType = "nexmo"       // SMSServiceNexmo is a Nexmo SMS service.
	SMSServiceRouteMobile SMSServiceType = "routemobile" // SMSServiceRouteMobile is a RouteMobile SMS service.
	SMSServiceMock        SMSServiceType = "mock"        // SMSServiceMock is an SMS service mock.
)

type TwilioServiceSettings struct {
	// Twilio related config.
	AccountSid string `yaml:"accountSid" json:"account_sid"`
	AuthToken  string `yaml:"authToken" json:"auth_token"`
	ServiceSid string `yaml:"serviceSid" json:"service_sid"`
}

type NexmoServiceSettings struct {
	// Nexmo related config.
	APIKey    string `yaml:"apiKey" json:"api_key"`
	APISecret string `yaml:"apiSecret" json:"api_secret"`
}

type RouteMobileServiceSettings struct {
	// RouteMobile related config.
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Source   string `yaml:"source" json:"source"`
	Region   string `yaml:"region" json:"region"`
}

// LoginSettings are settings of login.
type LoginSettings struct {
	LoginWith        LoginWith `yaml:"loginWith" json:"login_with"`
	TFAType          TFAType   `yaml:"tfaType" json:"tfa_type"`
	TFAResendTimeout int       `yaml:"tfaResendTimeout" json:"tfa_resend_timeout"`
}

// LoginWith is a type for configuring supported login ways.
type LoginWith struct {
	Username  bool `yaml:"username" json:"username"`
	Phone     bool `yaml:"phone" json:"phone"`
	Email     bool `yaml:"email" json:"email"`
	Federated bool `yaml:"federated" json:"federated"`
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
	port := ss.General.Port
	if port == "" {
		panic("can't start without port")
	}
	return strings.Join([]string{":", port}, "")
}

type LoggerSettings struct {
	DumpRequest bool `yaml:"dumpRequest" json:"dumpRequest"`
}

type LocalFileStorageSettings struct {
	// just a file name
	FileName string `yaml:"file_name" json:"file_name" bson:"file_name"`
}

type S3StorageSettings struct {
	Region   string `yaml:"region" json:"region" bson:"region"`
	Bucket   string `yaml:"bucket" json:"bucket" bson:"bucket"`
	Endpoint string `yaml:"endpoint" json:"endpoint" bson:"endpoint"`
	Key      string `yaml:"key" json:"key" bson:"key"`
}

type EtcdStorageSettings struct {
	Endpoints []string `json:"endpoints" yaml:"endpoints"`
	Key       string   `json:"key" yaml:"key"`
	Username  string   `json:"username" yaml:"username"`
	Password  string   `json:"password" yaml:"password"`
}

type AdminPanelSettings struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}

func ConfigStorageSettingsFromString(config string) (ConfigStorageSettings, error) {
	// Parse the URL and ensure there are no errors.
	u, err := url.Parse(config)
	if err != nil {
		return ConfigStorageSettings{}, fmt.Errorf("Unable to parse config string: %s", config)
	}

	switch strings.ToLower(u.Scheme) {
	// case "etcd":
	// 	return ConfigStorageSettingsFromStringEtcd(config)
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
		File: &LocalFileStorageSettings{
			FileName: filename,
		},
	}, nil
}
