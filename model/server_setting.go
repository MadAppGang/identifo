package model

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// ServerSettings are server settings.
type ServerSettings struct {
	StaticFolderPath   string                `yaml:"staticFolderPath,omitempty"`
	EmailTemplatesPath string                `yaml:"emailTemplatesPath,omitempty"`
	EmailTemplateNames EmailTemplateNames    `yaml:"emailTemplateNames,omitempty"`
	PEMFolderPath      string                `yaml:"pemFolderPath,omitempty"`
	PrivateKey         string                `yaml:"privateKey,omitempty"`
	PublicKey          string                `yaml:"publicKey,omitempty"`
	Algorithm          TokenServiceAlgorithm `yaml:"algorithm,omitempty"`
	Issuer             string                `yaml:"issuer,omitempty"`
	MailService        MailServiceType       `yaml:"mailService,omitempty"`
	SessionStorage     SessionStorageType    `yaml:"sessionStorage,omitempty"`
	SessionDuration    SessionDuration       `yaml:"sessionDuration,omitempty"`
	Host               string                `yaml:"host,omitempty"`
	AccountConfigPath  string                `yaml:"accountConfigPath,omitempty"`
	DBSettings
}

// DBSettings holds together all possible database-related settings.
type DBSettings struct {
	DBType     string `yaml:"type,omitempty"`
	DBName     string `yaml:"name,omitempty"`
	DBEndpoint string `yaml:"endpoint,omitempty"`
	DBRegion   string `yaml:"region,omitempty"`
	DBPath     string `yaml:"path,omitempty"`
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
