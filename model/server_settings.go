package model

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	defaultEtcdKey            = "identifo"
	IdentifoConfigPathEnvName = "IDENTIFO_CONFIG"
)

var s3ServerFlagRegexp = regexp.MustCompile(`^s3://(?P<region>[a-zA-Z0-9\-]{5,})@(?P<bucket>[a-z0-9\.\-]{3,63})(?P<key>[^\r\n\t\f\v|]+)\|?(?P<endpoint>\S+)?$`)

// ServerSettings are server settings.
type ServerSettings struct {
	General        GeneralServerSettings  `yaml:"general" json:"general"`
	AdminAccount   AdminAccountSettings   `yaml:"adminAccount" json:"admin_account"`
	Storage        StorageSettings        `yaml:"storage" json:"storage"`
	SessionStorage SessionStorageSettings `yaml:"sessionStorage" json:"session_storage"`
	Services       ServicesSettings       `yaml:"services" json:"external_services"`
	Login          LoginSettings          `yaml:"login" json:"login"`
	KeyStorage     FileStorageSettings    `yaml:"keyStorage" json:"key_storage"`
	Config         FileStorageSettings    `yaml:"-" json:"config"`
	Logger         LoggerSettings         `yaml:"logger" json:"logger"`
	AdminPanel     AdminPanelSettings     `yaml:"adminPanel" json:"admin_panel"`
	LoginWebApp    FileStorageSettings    `yaml:"loginWebApp" json:"login_web_app"`
	EmailTemplates FileStorageSettings    `yaml:"emailTemplates" json:"email_templates"`
	Impersonation  ImpersonationSettings  `yaml:"impersonation" json:"impersonation"`
}

type ImpersonationServiceType string

const (
	ImpersonationServiceTypeNone   ImpersonationServiceType = "none"
	ImpersonationServiceTypeScope  ImpersonationServiceType = "scope"
	ImpersonationServiceTypeRole   ImpersonationServiceType = "role"
	ImpersonationServiceTypePlugin ImpersonationServiceType = "plugin"
)

// ImpersonationSettings are settings for impersonation.
type ImpersonationSettings struct {
	Type   ImpersonationServiceType   `yaml:"type" json:"type"`
	Plugin PluginSettings             `yaml:"plugin" json:"plugin"`
	Scope  ImpersonationScopeSettings `yaml:"scope" json:"scope"`
	Role   ImpersonationRoleSettings  `yaml:"role" json:"role"`
}

type ImpersonationScopeSettings struct {
	AllowedScopes []string `yaml:"allowed_scopes" json:"allowed_scopes"`
}

type ImpersonationRoleSettings struct {
	AllowedRoles []string `yaml:"allowed_roles" json:"allowed_roles"`
}

// GeneralServerSettings are general server settings.
type GeneralServerSettings struct {
	Locale          string   `yaml:"locale" json:"locale"`
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
	DefaultStorage          DatabaseSettings `yaml:"default" json:"default"`
	AppStorage              DatabaseSettings `yaml:"appStorage" json:"app_storage"`
	UserStorage             DatabaseSettings `yaml:"userStorage" json:"user_storage"`
	TokenStorage            DatabaseSettings `yaml:"tokenStorage" json:"token_storage"`
	TokenBlacklist          DatabaseSettings `yaml:"tokenBlacklist" json:"token_blacklist"`
	VerificationCodeStorage DatabaseSettings `yaml:"verificationCodeStorage" json:"verification_code_storage"`
	InviteStorage           DatabaseSettings `yaml:"inviteStorage" json:"invite_storage"`
	ManagementKeysStorage   DatabaseSettings `yaml:"managementKeysStorage" json:"management_keys_storage"`
}

// DatabaseSettings holds together all settings applicable to a particular database.
type DatabaseSettings struct {
	Type   DatabaseType           `yaml:"type" json:"type"`
	BoltDB BoltDBDatabaseSettings `yaml:"boltdb" json:"boltdb"`
	Mongo  MongoDatabaseSettings  `yaml:"mongo" json:"mongo"`
	Dynamo DynamoDatabaseSettings `yaml:"dynamo" json:"dynamo"`
	Plugin PluginSettings         `yaml:"plugin" json:"plugin"`
	GRPC   GRPCSettings           `yaml:"grpc" json:"grpc"`
}

func (ds *DatabaseSettings) UnmarshalJSON(b []byte) error {
	type DSAlias DatabaseSettings
	aux := struct{ *DSAlias }{DSAlias: (*DSAlias)(ds)}

	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}

	// if database type is not specified, we assumed to use default one
	if len(ds.Type) == 0 {
		ds.Type = DBTypeDefault
	}
	return nil
}

type BoltDBDatabaseSettings struct {
	Path string `yaml:"path" json:"path"`
}

type MongoDatabaseSettings struct {
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
	DBTypeDefault  DatabaseType = "default" // DBTypeDefault it means the settings should be referenced from default database settings.
	DBTypeBoltDB   DatabaseType = "boltdb"  // DBTypeBoltDB is for BoltDB.
	DBTypeMongoDB  DatabaseType = "mongo"   // DBTypeMongoDB is for MongoDB.
	DBTypeDynamoDB DatabaseType = "dynamo"  // DBTypeDynamoDB is for DynamoDB.
	DBTypeFake     DatabaseType = "fake"    // DBTypeFake is return some predefined const data.
	DBTypeMem      DatabaseType = "mem"     // DBTypeMem is for in-memory storage.
	DBTypePlugin   DatabaseType = "plugin"  // DBTypePlugin is used for hashicorp/go-plugin.
	DBTypeGRPC     DatabaseType = "grpc"    // DBTypeGRPC is used for pure grpc.
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
	Region   string `yaml:"region" json:"region"`
	Bucket   string `yaml:"bucket" json:"bucket"`
	Key      string `yaml:"key" json:"key"`
	Endpoint string `yaml:"endpoint" json:"endpoint"`
}

type FileStorageLocal struct {
	Path string `yaml:"path" json:"path"`
}

// if key or path has folder and filename joined, this function returns filename part only
func (fs FileStorageSettings) FileName() string {
	path := ""
	if fs.Type == FileStorageTypeLocal {
		path = fs.Local.Path
	} else if fs.Type == FileStorageTypeS3 {
		path = fs.S3.Key
	}
	return filepath.Base(path)
}

// if key or path has folder and filename joined, this function returns path part only
func (fs FileStorageSettings) Dir() string {
	path := ""
	if fs.Type == FileStorageTypeLocal {
		path = fs.Local.Path
	} else if fs.Type == FileStorageTypeS3 {
		path = fs.S3.Key
	}
	return filepath.Dir(path)
}

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
	SendFrom   string `yaml:"sendFrom" json:"send_from"`
	Region     string `yaml:"region" json:"region"`
	Edge       string `yaml:"edge" json:"edge"`
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
	LoginWith            LoginWith `yaml:"loginWith" json:"login_with"`
	TFAType              TFAType   `yaml:"tfaType" json:"tfa_type"`
	TFAResendTimeout     int       `yaml:"tfaResendTimeout" json:"tfa_resend_timeout"`
	AllowRegisterMissing bool      `yaml:"allowRegisterMissing" json:"allow_register_missing"`
}

// LoginWith is a type for configuring supported login ways.
type LoginWith struct {
	Username      bool `yaml:"username" json:"username"`
	Phone         bool `yaml:"phone" json:"phone"`
	Email         bool `yaml:"email" json:"email"`
	Federated     bool `yaml:"federated" json:"federated"`
	FederatedOIDC bool `yaml:"federatedOIDC" json:"federated_oidc"`
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
	return ":" + port
}

type LoggerSettings struct {
	DumpRequest bool `yaml:"dumpRequest" json:"dumpRequest"`
}

type AdminPanelSettings struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}

func ConfigStorageSettingsFromString(config string) (FileStorageSettings, error) {
	// Parse the URL and ensure there are no errors.
	u, err := url.Parse(config)
	if err != nil {
		return FileStorageSettings{}, fmt.Errorf("Unable to parse config string: %s", config)
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

// example of s3 config could be:
// s3://ap-southwest-2@my-favorite-bucket/identifo/config/config.yaml
// or with endpoint for local testing or non AWS S3 storage:
// s3://ap-southwest-2@my-favorite-bucket/identifo/config/config.yaml|https://10.10.10.19:1122

func ConfigStorageSettingsFromStringS3(config string) (FileStorageSettings, error) {
	if match := s3ServerFlagRegexp.MatchString(config); !match {
		return FileStorageSettings{}, fmt.Errorf("error parsing S3 config location: %s", config)
	}
	matches := s3ServerFlagRegexp.FindStringSubmatch(config)
	region := matches[s3ServerFlagRegexp.SubexpIndex("region")]
	bucket := matches[s3ServerFlagRegexp.SubexpIndex("bucket")]
	key := matches[s3ServerFlagRegexp.SubexpIndex("key")]
	endpoint := matches[s3ServerFlagRegexp.SubexpIndex("endpoint")]

	return FileStorageSettings{
		Type: FileStorageTypeS3,
		S3: FileStorageS3{
			Region:   region,
			Bucket:   bucket,
			Key:      key,
			Endpoint: endpoint,
		},
	}, nil
}

func ConfigStorageSettingsFromStringFile(config string) (FileStorageSettings, error) {
	filename := config
	if strings.HasPrefix(strings.ToUpper(filename), "FILE://") {
		filename = filename[7:]
	}

	return FileStorageSettings{
		Type: FileStorageTypeLocal,
		Local: FileStorageLocal{
			Path: filename,
		},
	}, nil
}
