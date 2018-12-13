package model

//ServerSettings server settings
type ServerSettings struct {
	StaticFolderPath string
	PEMFolderPath    string
	PrivateKey       string
	PublicKey        string
	Algorithm        TokenServiceAlgorithm
	Issuer           string
	MailService      MailServiceType
}

//MailServiceType - how to send email to clients
type MailServiceType int

const (
	//MailServiceMailgun Mailgun service
	MailServiceMailgun MailServiceType = iota
	//MailServiceAWS AWS SMS service
	MailServiceAWS
)
