package model

import (
	"time"
)

//ServerSettings server settings
type ServerSettings struct {
	StaticFolderPath   string
	EmailTemplatesPath string
	EmailTemplates     EmailTemplates
	PEMFolderPath      string
	PrivateKey         string
	PublicKey          string
	Algorithm          TokenServiceAlgorithm
	Issuer             string
	MailService        MailServiceType
	SessionStorage     SessionStorageType
	SessionDuration    time.Duration
	Host               string
	ConfigPath         string
}

// MailServiceType - how to send email to clients.
type MailServiceType int

const (
	//MailServiceMailgun Mailgun service
	MailServiceMailgun MailServiceType = iota
	//MailServiceAWS AWS SMS service
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
