package model

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// ServerSettings are server settings.
type ServerSettings struct {
	StaticFolderPath   string                `yaml:"staticFolderPath,omitempty" json:"static_folder_path,omitempty"`
	EmailTemplatesPath string                `yaml:"emailTemplatesPath,omitempty" json:"email_templates_path,omitempty"`
	EmailTemplateNames EmailTemplateNames    `yaml:"emailTemplateNames,omitempty" json:"email_template_names,omitempty"`
	PEMFolderPath      string                `yaml:"pemFolderPath,omitempty" json:"pem_folder_path,omitempty"`
	PrivateKey         string                `yaml:"privateKey,omitempty" json:"private_key,omitempty"`
	PublicKey          string                `yaml:"publicKey,omitempty" json:"public_key,omitempty"`
	Algorithm          TokenServiceAlgorithm `yaml:"algorithm,omitempty" json:"algorithm,omitempty"`
	Issuer             string                `yaml:"issuer,omitempty" json:"issuer,omitempty"`
	MailService        MailServiceType       `yaml:"mailService,omitempty" json:"mail_service,omitempty"`
	SessionStorage     SessionStorageType    `yaml:"sessionStorage,omitempty" json:"session_storage,omitempty"`
	SessionDuration    SessionDuration       `yaml:"sessionDuration,omitempty" json:"session_duration,omitempty"`
	Host               string                `yaml:"host,omitempty" json:"host,omitempty"`
	AccountConfigPath  string                `yaml:"accountConfigPath,omitempty" json:"account_config_path,omitempty"`
	DBSettings         `yaml:"-,inline" json:"dbDettings,omitempty"`
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
	MailServiceMailgun MailServiceType = iota
	// MailServiceAWS is an AWS SES service.
	MailServiceAWS
)

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
	SessionStorageMem SessionStorageType = iota
	// SessionStorageRedis means to store sessions in Redis.
	SessionStorageRedis
)

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
