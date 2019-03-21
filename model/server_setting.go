package model

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
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
	SessionDuration    time.Duration         `yaml:"sessionDuration,omitempty"`
	Host               string                `yaml:"host,omitempty"`
	AccountConfigPath  string                `yaml:"accountConfigPath,omitempty"`
	ServerConfigPath   string                `yaml:"serverConfigPath,omitempty"`
	AppsImportPath     string                `yaml:"appsImportPath,omitempty"`
	UsersImportPath    string                `yaml:"usersImportPath,omitempty"`
}

// UnmarshalYAML implements yaml.Unmarshaller.
func (ss *ServerSettings) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux struct {
		StaticFolderPath   string             `yaml:"staticFolderPath,omitempty"`
		EmailTemplatesPath string             `yaml:"emailTemplatesPath,omitempty"`
		EmailTemplateNames EmailTemplateNames `yaml:"emailTemplateNames,omitempty"`
		PEMFolderPath      string             `yaml:"pemFolderPath,omitempty"`
		PrivateKey         string             `yaml:"privateKey,omitempty"`
		PublicKey          string             `yaml:"publicKey,omitempty"`
		Algorithm          string             `yaml:"algorithm,omitempty"`
		Issuer             string             `yaml:"issuer,omitempty"`
		MailService        string             `yaml:"mailService,omitempty"`
		SessionStorage     string             `yaml:"sessionStorage,omitempty"`
		SessionDuration    int                `yaml:"sessionDuration,omitempty"`
		Host               string             `yaml:"host,omitempty"`
		AccountConfigPath  string             `yaml:"accountConfigPath,omitempty"`
		ServerConfigPath   string             `yaml:"serverConfigPath,omitempty"`
		AppsImportPath     string             `yaml:"appsImportPath,omitempty"`
		UsersImportPath    string             `yaml:"usersImportPath,omitempty"`
	}

	if err := unmarshal(&aux); err != nil {
		return err
	}

	alg, ok := map[string]TokenServiceAlgorithm{
		"es256": TokenServiceAlgorithmES256,
		"rs256": TokenServiceAlgorithmRS256,
		"auto":  TokenServiceAlgorithmAuto}[aux.Algorithm]
	if !ok {
		return fmt.Errorf("Invalid TokenServiceAlgorithm %v", aux.Algorithm)
	}

	mailService, ok := map[string]MailServiceType{
		"aws ses": MailServiceAWS,
		"mailgun": MailServiceMailgun}[aux.MailService]
	if !ok {
		return fmt.Errorf("Invalid MailServiceType %v", aux.Algorithm)
	}

	sessionStorageType, ok := map[string]SessionStorageType{
		"memory": SessionStorageMem,
		"redis":  SessionStorageRedis}[aux.SessionStorage]
	if !ok {
		return fmt.Errorf("Invalid SessionStorageType %v", aux.Algorithm)
	}

	sessionDuration := time.Second * time.Duration(aux.SessionDuration)

	ss.StaticFolderPath = aux.StaticFolderPath
	ss.EmailTemplatesPath = aux.EmailTemplatesPath
	ss.EmailTemplateNames = aux.EmailTemplateNames
	ss.PEMFolderPath = aux.PEMFolderPath
	ss.PrivateKey = aux.PrivateKey
	ss.PublicKey = aux.PublicKey
	ss.Algorithm = alg
	ss.Issuer = aux.Issuer
	ss.MailService = mailService
	ss.SessionStorage = sessionStorageType
	ss.SessionDuration = sessionDuration
	ss.Host = aux.Host
	ss.AccountConfigPath = aux.AccountConfigPath
	ss.ServerConfigPath = aux.ServerConfigPath
	ss.AppsImportPath = aux.AppsImportPath
	ss.UsersImportPath = aux.UsersImportPath
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

// MailServiceType - how to send email to clients.
type MailServiceType int

const (
	// MailServiceMailgun is a Mailgun service.
	MailServiceMailgun MailServiceType = iota
	// MailServiceAWS is an AWS SES service.
	MailServiceAWS
)

// SessionStorageType - where to store sessions.
type SessionStorageType int

const (
	// SessionStorageMem means to store sessions in memory.
	SessionStorageMem SessionStorageType = iota
	// SessionStorageRedis means to store sessions in Redis.
	SessionStorageRedis
)
