package model

//ServerSettings server settings
type ServerSettings struct {
	StaticFolderPath   string
	EmailTemplatesPath string
	EmailTemplates     EmailTemplates
	EncryptionKeyPath  string
	PEMFolderPath      string
	PrivateKey         string
	PublicKey          string
	Algorithm          TokenServiceAlgorithm
	Issuer             string
	MailService        MailServiceType
	Host               string
}

//MailServiceType - how to send email to clients
type MailServiceType int

const (
	//MailServiceMailgun Mailgun service
	MailServiceMailgun MailServiceType = iota
	//MailServiceAWS AWS SMS service
	MailServiceAWS
)
