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
	Host               string                `yaml:"host,omitempty" json:"host,omitempty"`
	PEMFolderPath      string                `yaml:"pemFolderPath,omitempty" json:"pem_folder_path,omitempty"`
	PrivateKey         string                `yaml:"privateKey,omitempty" json:"private_key,omitempty"`
	PublicKey          string                `yaml:"publicKey,omitempty" json:"public_key,omitempty"`
	Issuer             string                `yaml:"issuer,omitempty" json:"issuer,omitempty"`
	Algorithm          TokenServiceAlgorithm `yaml:"algorithm,omitempty" json:"algorithm,omitempty"`
	MailService        MailServiceType       `yaml:"mailService,omitempty" json:"mail_service,omitempty"`
	SessionStorage     SessionStorageType    `yaml:"sessionStorage,omitempty" json:"session_storage,omitempty"`
	SessionDuration    SessionDuration       `yaml:"sessionDuration,omitempty" json:"session_duration,omitempty"`
	StaticFolderPath   string                `yaml:"staticFolderPath,omitempty" json:"static_folder_path,omitempty"`
	EmailTemplatesPath string                `yaml:"emailTemplatesPath,omitempty" json:"email_templates_path,omitempty"`
	EmailTemplateNames EmailTemplateNames    `yaml:"emailTemplateNames,omitempty" json:"email_template_names,omitempty"`
	AccountConfigPath  string                `yaml:"accountConfigPath,omitempty" json:"account_config_path,omitempty"`
	ServerConfigPath   string                `yaml:"serverConfigPath,omitempty" json:"server_config_path,omitempty"`
	DBSettings         `yaml:"-,inline" json:"db_settings,omitempty"`
}

// DBSettings holds together all possible database-related settings.
type DBSettings struct {
	DBType     string `yaml:"dbType,omitempty" json:"type,omitempty"`
	DBName     string `yaml:"dbName,omitempty" json:"name,omitempty"`
	DBEndpoint string `yaml:"dbEndpoint,omitempty" json:"endpoint,omitempty"`
	DBRegion   string `yaml:"dbRegion,omitempty" json:"region,omitempty"`
	DBPath     string `yaml:"dbPath,omitempty" json:"path,omitempty"`
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

// SessionStorageType - where to store sessions.
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
