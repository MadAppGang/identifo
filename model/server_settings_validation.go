package model

import (
	"fmt"
	"net/url"
	"os"
)

const (
	identifoLoginWebAppBucket    = "IDENTIFO_LOGIN_APP_BUCKET"
	identifoJWTKeysBucketEnvName = "IDENTIFO_JWT_KEYS_BUCKET"
	// RouteMobileRegionUAE is a regional UAE RouteMobileR platform.
	RouteMobileRegionUAE        = "uae"
	identifoConfigBucketEnvName = "IDENTIFO_CONFIG_BUCKET"
)

// Validate makes sure that all crucial fields are set.
func (ss *ServerSettings) Validate() error {
	if len(ss.AdminAccount.LoginEnvName) == 0 {
		return fmt.Errorf("Admin login env variable name not specified")
	}
	if len(ss.AdminAccount.PasswordEnvName) == 0 {
		return fmt.Errorf("Admin password env variable name not specified")
	}

	if err := ss.General.Validate(); err != nil {
		return err
	}
	if err := ss.Storage.Validate(); err != nil {
		return err
	}
	if err := ss.SessionStorage.Validate(); err != nil {
		return err
	}
	if err := ss.LoginWebApp.Validate(); err != nil {
		return err
	}
	if err := ss.AdminPanel.Validate(); err != nil {
		return err
	}
	if err := ss.EmailTemplates.Validate(); err != nil {
		return err
	}

	if err := ss.Services.Validate(); err != nil {
		return err
	}
	return nil
}

// Validate validates general services settings.
func (gss *GeneralServerSettings) Validate() error {
	subject := "GeneralServerSettings"
	if gss == nil {
		return fmt.Errorf("Nil %s", subject)
	}

	if _, err := url.ParseRequestURI(gss.Host); err != nil {
		return fmt.Errorf("%s. Host is invalid. %s", subject, err)
	}
	if len(gss.Issuer) == 0 {
		return fmt.Errorf("%s. Issuer is not set", subject)
	}
	return nil
}

// Validate validates storage settings.
func (ss *StorageSettings) Validate() error {
	if ss == nil {
		return fmt.Errorf("Nil StorageSettings")
	}
	if err := ss.AppStorage.Validate(); err != nil {
		return fmt.Errorf("AppStorage: %s", err)
	}
	if err := ss.UserStorage.Validate(); err != nil {
		return fmt.Errorf("UserStorage: %s", err)
	}
	if err := ss.TokenStorage.Validate(); err != nil {
		return fmt.Errorf("TokenStorage: %s", err)
	}
	if err := ss.TokenBlacklist.Validate(); err != nil {
		return fmt.Errorf("TokenBlacklist: %s", err)
	}
	if err := ss.VerificationCodeStorage.Validate(); err != nil {
		return fmt.Errorf("VerificationCodeStorage: %s", err)
	}
	return nil
}

// Validate validates database settings.
func (dbs *DatabaseSettings) Validate() error {
	if len(dbs.Type) == 0 {
		return fmt.Errorf("empty database type")
	}

	switch dbs.Type {
	case DBTypeFake:
		return nil
	case DBTypeBoltDB:
		if len(dbs.BoltDB.Path) == 0 {
			return fmt.Errorf("empty database path")
		}
	case DBTypeDynamoDB:
		if len(dbs.Dynamo.Region) == 0 {
			return fmt.Errorf("empty Dynamo region")
		}
		if len(dbs.Dynamo.Endpoint) == 0 {
			return fmt.Errorf("empty Dynamo endpoint")
		}
	case DBTypeMongoDB:
		if _, err := url.ParseRequestURI(dbs.Mongo.ConnectionString); err != nil {
			return fmt.Errorf("invalid mongo connection string. %s", err)
		}
		if len(dbs.Mongo.DatabaseName) == 0 {
			return fmt.Errorf("empty mongo database name")
		}
	default:
		return fmt.Errorf("unsupported database type %s", dbs.Type)
	}
	return nil
}

// Validate validates admin session storage settings.
func (sss *SessionStorageSettings) Validate() error {
	subject := "SessionStorageSettings"
	if sss == nil {
		return fmt.Errorf("Nil %s", subject)
	}

	if len(sss.Type) == 0 {
		return fmt.Errorf("Empty session storage type")
	}
	if sss.SessionDuration.Duration == 0 {
		return fmt.Errorf("%s. Session duration is 0 seconds", subject)
	}

	switch sss.Type {
	case SessionStorageMem:
		return nil
	case SessionStorageRedis:
		if _, err := url.ParseRequestURI(sss.Redis.Address); err != nil {
			return fmt.Errorf("%s. Invalid address. %s", subject, err)
		}
	case SessionStorageDynamoDB:
		if len(sss.Dynamo.Region) == 0 {
			return fmt.Errorf("%s. Empty AWS region", subject)
		}
		if len(sss.Dynamo.Endpoint) == 0 {
			return fmt.Errorf("%s. Empty AWS Dynamo endpoint", subject)
		}

	default:
		return fmt.Errorf("%s. Unknown type", subject)
	}
	return nil
}

// Validate validates configuration storage settings.
func (css *ConfigStorageSettings) Validate() error {
	subject := "ConfigurationStorageSettings"
	if css == nil {
		return fmt.Errorf("Nil %s", subject)
	}

	switch css.Type {
	case ConfigStorageTypeFile:
		if css.File == nil {
			return fmt.Errorf("%s. empty file key", subject)
		}
		if len(css.File.FileName) == 0 {
			return fmt.Errorf("%s. empty file key ", subject)
		}
		break
	case ConfigStorageTypeS3:
		if css.S3 == nil {
			return fmt.Errorf("%s. empty s3 settings key", subject)
		}
		if len(css.S3.Region) == 0 {
			return fmt.Errorf("%s. Empty AWS region", subject)
		}
		if bucket := os.Getenv(identifoConfigBucketEnvName); len(bucket) != 0 {
			css.S3.Bucket = bucket
		}
		if len(css.S3.Bucket) == 0 {
			return fmt.Errorf("%s. S3 Bucket for config is not set", subject)
		}
		if len(css.S3.Key) == 0 {
			return fmt.Errorf("%s. S3 Key for config is not set", subject)
		}
	default:
		return fmt.Errorf("%s. Unknown type", subject)
	}

	return nil
}

// Validate validates login web app settings
func (sfs *LoginWebAppSettings) Validate() error {
	subject := "LoginWebAppSettings"
	if sfs == nil {
		return fmt.Errorf("Nil %s", subject)
	}

	if len(sfs.Type) == 0 {
		return fmt.Errorf("%s. Empty login web app type", subject)
	}

	switch sfs.Type {
	case LoginWebAppTypeNone:
		return nil
	case LoginWebAppTypeLocal:
		if len(sfs.Local.FolderPath) > 0 {
			return fmt.Errorf("%s. empty folder", subject)
		}
		return nil
	case LoginWebAppTypeS3:
		if len(sfs.S3.Region) == 0 {
			return fmt.Errorf("%s. Empty AWS region", subject)
		}
		if bucket := os.Getenv(identifoLoginWebAppBucket); len(bucket) != 0 {
			sfs.S3.Bucket = bucket
		}
		if len(sfs.S3.Bucket) == 0 {
			return fmt.Errorf("%s. empty s3 bucket for login web app", subject)
		}
	default:
		return fmt.Errorf("%s. Unknown type", subject)
	}
	return nil
}

// Validate validates email template storage settings
func (sfs *EmailTemplatesSettings) Validate() error {
	subject := "EmailTemplatesSettings"
	if sfs == nil {
		return fmt.Errorf("Nil %s", subject)
	}

	if len(sfs.Type) == 0 {
		return fmt.Errorf("%s. Empty email templates sotrage type", subject)
	}

	switch sfs.Type {
	case EmailTemplatesStorageTypeNone:
		return nil
	case EmailTemplatesStorageTypeLocal:
		if len(sfs.Local.FolderPath) > 0 {
			return fmt.Errorf("%s. empty folder", subject)
		}
		return nil
	case EmailTemplatesStorageTypeS3:
		if len(sfs.S3.Region) == 0 {
			return fmt.Errorf("%s. Empty AWS region", subject)
		}
		if bucket := os.Getenv(identifoLoginWebAppBucket); len(bucket) != 0 {
			sfs.S3.Bucket = bucket
		}
		if len(sfs.S3.Bucket) == 0 {
			return fmt.Errorf("%s. empty s3 bucket for email templates", subject)
		}
	default:
		return fmt.Errorf("%s. Unknown type", subject)
	}
	return nil
}

func (kss *AdminPanelSettings) Validate() error {
	// not it has just one bool field wich always valid
	return nil
}

// Validate validates key storage settings.
func (kss *KeyStorageSettings) Validate() error {
	subject := "KeyStorageSettings"
	if len(kss.Type) == 0 {
		return fmt.Errorf("%s. Empty key storage type", subject)
	}

	switch kss.Type {
	case KeyStorageTypeLocal:
		if len(kss.File.PrivateKeyPath) == 0 {
			return fmt.Errorf("%s. empty File settings key", subject)
		}
		break
	case KeyStorageTypeS3:
		if len(kss.S3.Region) == 0 {
			return fmt.Errorf("%s. Empty AWS region", subject)
		}
		if bucket := os.Getenv(identifoJWTKeysBucketEnvName); len(bucket) != 0 {
			kss.S3.Bucket = bucket
		}
		if len(kss.S3.Bucket) == 0 {
			return fmt.Errorf("%s. Bucket for keys is not set", subject)
		}
		if len(kss.S3.PrivateKeyKey) == 0 {
			return fmt.Errorf("%s. Private key  key is not set", subject)
		}
	default:
		return fmt.Errorf("%s. Unknown type '%s'", subject, kss.Type)
	}
	return nil
}

// Validate validates external services settings.
func (ess *ServicesSettings) Validate() error {
	subject := "ExternalServicesSettings"
	if ess == nil {
		return fmt.Errorf("Nil %s", subject)
	}

	if err := ess.Email.Validate(); err != nil {
		return fmt.Errorf("%s. %s", subject, err)
	}
	if err := ess.SMS.Validate(); err != nil {
		return fmt.Errorf("%s. %s", subject, err)
	}
	return nil
}

// Validate validates SMS service settings.
func (sss *SMSServiceSettings) Validate() error {
	subject := "SMSServiceSettings"
	if sss == nil {
		return fmt.Errorf("Nil %s", subject)
	}
	if len(sss.Type) == 0 {
		return fmt.Errorf("Empty SMS service type")
	}

	switch sss.Type {
	case SMSServiceMock:
		return nil
	case SMSServiceNexmo:
		if len(sss.Nexmo.APIKey) == 0 || len(sss.Nexmo.APISecret) == 0 {
			return fmt.Errorf("%s. Error creating Nexmo SMS service, missing at least one of the parameters:"+
				"\n apiKey : %v\n apiSecret : %v\n", subject, sss.Nexmo.APIKey, sss.Nexmo.APISecret)
		}
	case SMSServiceTwilio:
		if len(sss.Twilio.AccountSid) == 0 || len(sss.Twilio.AuthToken) == 0 || len(sss.Twilio.ServiceSid) == 0 {
			return fmt.Errorf("%s. Error creating Twilio SMS service, missing at least one of the parameters:"+
				"\n sidKey : %v\n tokenKey : %v\n ServiceSidKey : %v\n", subject, sss.Twilio.AccountSid, sss.Twilio.AuthToken, sss.Twilio.ServiceSid)
		}
	case SMSServiceRouteMobile:
		if len(sss.Routemobile.Username) == 0 || len(sss.Routemobile.Password) == 0 || len(sss.Routemobile.Source) == 0 {
			return fmt.Errorf("%s. Error creating RouteMobile SMS service, missing at least one of the parameters:"+
				"\n username : %v\n password : %v\n", subject, sss.Routemobile.Username, sss.Routemobile.Password)
		}
		if sss.Routemobile.Region != RouteMobileRegionUAE {
			return fmt.Errorf("%s. Error creating RouteMobile SMS service, region %s is not supported", subject, sss.Routemobile.Region)
		}
	default:
		return fmt.Errorf("%s. Unknown type", subject)
	}
	return nil
}

const (
	// mailgunDomainKey is a name of env variable that contains Mailgun domain value.
	mailgunDomainKey = "MAILGUN_DOMAIN"
	// mailgunPrivateKey is a name of env variable that contains Mailgun private key value.
	mailgunPrivateKey = "MAILGUN_PRIVATE_KEY"
	// mailgunPublicKey is a name of env variable that contains Mailgun public key value.
	mailgunPublicKey = "MAILGUN_PUBLIC_KEY"
	// mailgunSenderKey is a name of env variable that contains Mailgun sender key value.
	mailgunSenderKey = "MAILGUN_SENDER"

	// awsSESSenderKey is a name of env variable that contains AWS SWS sender value.
	awsSESSenderKey = "AWS_SES_SENDER"
	// awsSESRegionKey is a name of env variable that contains AWS SWS region value.
	awsSESRegionKey = "AWS_SES_REGION"
)

// Validate validates email service settings.
func (ess *EmailServiceSettings) Validate() error {
	subject := "EmailServiceSettings"
	if ess == nil {
		return fmt.Errorf("Nil %s", subject)
	}
	if len(ess.Type) == 0 {
		return fmt.Errorf("%s. Empty email service type", subject)
	}

	switch ess.Type {
	case EmailServiceMock:
		return nil
	case EmailServiceAWS:
		if region := os.Getenv(awsSESRegionKey); len(region) != 0 {
			ess.SES.Region = region
		}
		if sender := os.Getenv(awsSESSenderKey); len(sender) != 0 {
			ess.SES.Sender = sender
		}
		if len(ess.SES.Sender) == 0 {
			return fmt.Errorf("%s. Empty AWS sender", subject)
		}
		if len(ess.SES.Region) == 0 {
			return fmt.Errorf("%s. Empty AWS region", subject)
		}
	case EmailServiceMailgun:
		if domain := os.Getenv(mailgunDomainKey); len(domain) != 0 {
			ess.Mailgun.Domain = domain
		}
		if publicKey := os.Getenv(mailgunPublicKey); len(publicKey) != 0 {
			ess.Mailgun.PublicKey = publicKey
		}
		if privateKey := os.Getenv(mailgunPrivateKey); len(privateKey) != 0 {
			ess.Mailgun.PrivateKey = privateKey
		}
		if sender := os.Getenv(mailgunSenderKey); len(sender) != 0 {
			ess.Mailgun.Sender = sender
		}

		if len(ess.Mailgun.Domain) == 0 {
			return fmt.Errorf("%s. Empty Mailgun domain", subject)
		}
		if len(ess.Mailgun.PublicKey) == 0 || len(ess.Mailgun.PrivateKey) == 0 {
			return fmt.Errorf("%s. At least one of the keys is empty", subject)
		}
		if len(ess.Mailgun.Sender) == 0 {
			return fmt.Errorf("%s. Empty Mailgun sender", subject)
		}
	default:
		return fmt.Errorf("%s. Unknown type", subject)
	}
	return nil
}
