package model

import (
	"encoding/json"
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
	MailService          MailServiceType              `yaml:"mailService,omitempty" json:"mail_service,omitempty"`
	ConfigurationStorage ConfigurationStorageSettings `yaml:"configurationStorage,omitempty" json:"configuration_storage,omitempty"`
	SessionStorage       SessionStorageSettings       `yaml:"sessionStorage,omitempty" json:"session_storage,omitempty"`
	StaticFolderPath     string                       `yaml:"staticFolderPath,omitempty" json:"static_folder_path,omitempty"`
	EmailTemplatesPath   string                       `yaml:"emailTemplatesPath,omitempty" json:"email_templates_path,omitempty"`
	EmailTemplateNames   EmailTemplateNames           `yaml:"emailTemplateNames,omitempty" json:"email_template_names,omitempty"`
	AdminAccount         AdminAccountSettings         `yaml:"adminAccount,omitempty" json:"admin_account,omitempty"`
	ServerConfigPath     string                       `yaml:"serverConfigPath,omitempty" json:"server_config_path,omitempty"`
	SMSService           SMSServiceSettings           `yaml:"smsService,omitempty" json:"sms_service,omitempty"`
	Database             DBSettings                   `yaml:"database,omitempty" json:"database,omitempty"`
}

// DBSettings holds together all possible database-related settings.
// TODO: support separate database engines for each service
type DBSettings struct {
	DBType     DatabaseType `yaml:"dbType,omitempty" json:"type,omitempty"`
	DBName     string       `yaml:"dbName,omitempty" json:"name,omitempty"`
	DBEndpoint string       `yaml:"dbEndpoint,omitempty" json:"endpoint,omitempty"`
	DBRegion   string       `yaml:"dbRegion,omitempty" json:"region,omitempty"`
	DBPath     string       `yaml:"dbPath,omitempty" json:"path,omitempty"`
}

// DatabaseType is a type of database management engine.
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

// SessionStorageSettings holds together session storage settings.
type SessionStorageSettings struct {
	Type            SessionStorageType `yaml:"type,omitempty" json:"type,omitempty"`
	Address         string             `yaml:"address,omitempty" json:"address,omitempty"`
	Password        string             `yaml:"password,omitempty" json:"password,omitempty"`
	DB              int                `yaml:"db,omitempty" json:"db,omitempty"`
	SessionDuration SessionDuration    `yaml:"sessionDuration,omitempty" json:"session_duration,omitempty"`
}

// AdminAccountSettings are names of environment variables that store admin credentials.
type AdminAccountSettings struct {
	LoginEnvName    string `yaml:"loginEnvName" json:"login_env_name,omitempty"`
	PasswordEnvName string `yaml:"passwordEnvName" json:"password_env_name,omitempty"`
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
type MailServiceType int

const (
	// MailServiceMailgun is a Mailgun service.
	MailServiceMailgun MailServiceType = iota + 1
	// MailServiceAWS is an AWS SES service.
	MailServiceAWS
)

// String implements Stringer.
func (mst MailServiceType) String() string {
	switch mst {
	case MailServiceMailgun:
		return "mailgun"
	case MailServiceAWS:
		return "aws ses"
	default:
		return fmt.Sprintf("MailServiceType(%d)", mst)
	}
}

// MarshalJSON implements json.Marshaller.
func (mst MailServiceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(mst.String())
}

// UnmarshalJSON implements json.Unmarshaller.
func (mst *MailServiceType) UnmarshalJSON(data []byte) error {
	if mst == nil {
		return nil
	}

	var aux string
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	mailServiceType, ok := map[string]MailServiceType{
		"aws ses": MailServiceAWS,
		"mailgun": MailServiceMailgun}[aux]
	if !ok {
		return fmt.Errorf("Invalid MailServiceType %v", aux)
	}

	*mst = mailServiceType
	return nil
}

// MarshalYAML implements yaml.Marshaller.
func (mst MailServiceType) MarshalYAML() (interface{}, error) {
	return mst.String(), nil
}

// UnmarshalYAML implements yaml.Unmarshaller.
func (mst *MailServiceType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if mst == nil {
		return nil
	}

	var aux string
	if err := unmarshal(&aux); err != nil {
		return err
	}

	mailServiceType, ok := map[string]MailServiceType{
		"aws ses": MailServiceAWS,
		"mailgun": MailServiceMailgun}[aux]
	if !ok {
		return fmt.Errorf("Invalid MailServiceType %v", aux)
	}

	*mst = mailServiceType
	return nil
}

// SessionStorageType - where to store admin sessions.
type SessionStorageType int

const (
	// SessionStorageMem means to store sessions in memory.
	SessionStorageMem SessionStorageType = iota + 1
	// SessionStorageRedis means to store sessions in Redis.
	SessionStorageRedis
)

// String implements Stringer.
func (sst SessionStorageType) String() string {
	switch sst {
	case SessionStorageMem:
		return "memory"
	case SessionStorageRedis:
		return "redis"
	default:
		return fmt.Sprintf("SessionStorageType(%d)", sst)
	}
}

// MarshalJSON implements json.Marshaller.
func (sst SessionStorageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(sst.String())
}

// UnmarshalJSON implements json.Unmarshaller.
func (sst *SessionStorageType) UnmarshalJSON(data []byte) error {
	if sst == nil {
		return nil
	}

	var aux string
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	sessionStorageType, ok := map[string]SessionStorageType{
		"memory": SessionStorageMem,
		"redis":  SessionStorageRedis}[aux]
	if !ok {
		return fmt.Errorf("Invalid SessionStorageType %v", aux)
	}

	*sst = sessionStorageType
	return nil
}

// MarshalYAML implements yaml.Marshaller.
func (sst SessionStorageType) MarshalYAML() (interface{}, error) {
	return sst.String(), nil
}

// UnmarshalYAML implements yaml.Unmarshaller.
func (sst *SessionStorageType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if sst == nil {
		return nil
	}

	var aux string
	if err := unmarshal(&aux); err != nil {
		return err
	}

	sessionStorageType, ok := map[string]SessionStorageType{
		"memory": SessionStorageMem,
		"redis":  SessionStorageRedis}[aux]
	if !ok {
		return fmt.Errorf("Invalid SessionStorageType %v", aux)
	}

	*sst = sessionStorageType
	return nil
}

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
