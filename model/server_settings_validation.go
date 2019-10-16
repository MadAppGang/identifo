package model

import (
	"fmt"
	"net/url"
	"os"
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
	if err := ss.ConfigurationStorage.Validate(); err != nil {
		return err
	}
	if err := ss.StaticFilesStorage.Validate(); err != nil {
		return err
	}
	if err := ss.ExternalServices.Validate(); err != nil {
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
	if len(gss.Algorithm) == 0 {
		return fmt.Errorf("%s. Algorithm is not set", subject)
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
	subject := "DatabaseSettings"
	if dbs == nil {
		return fmt.Errorf("Nil %s", subject)
	}
	if len(dbs.Type) == 0 {
		return fmt.Errorf("Empty database type")
	}

	switch dbs.Type {
	case DBTypeFake:
		return nil
	case DBTypeBoltDB:
		if len(dbs.Path) == 0 {
			return fmt.Errorf("Empty database path")
		}
	case DBTypeDynamoDB:
		if len(dbs.Region) == 0 {
			return fmt.Errorf("Empty AWS region")
		}
	case DBTypeMongoDB:
		if _, err := url.ParseRequestURI(dbs.Endpoint); err != nil {
			return fmt.Errorf("Invalid endpoint. %s", err)
		}
		if len(dbs.Name) == 0 {
			return fmt.Errorf("Empty database name")
		}
	default:
		return fmt.Errorf("%s. Unknown type", subject)
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
		if _, err := url.ParseRequestURI(sss.Address); err != nil {
			return fmt.Errorf("%s. Invalid address. %s", subject, err)
		}
	case SessionStorageDynamoDB:
		if len(sss.Region) == 0 {
			return fmt.Errorf("%s. Empty AWS region", subject)
		}
	default:
		return fmt.Errorf("%s. Unknown type", subject)
	}
	return nil
}

const identifoConfigBucketEnvName = "IDENTIFO_CONFIG_BUCKET"

// Validate validates configuration storage settings.
func (css *ConfigurationStorageSettings) Validate() error {
	subject := "ConfigurationStorageSettings"
	if css == nil {
		return fmt.Errorf("Nil %s", subject)
	}

	if len(css.Type) == 0 {
		return fmt.Errorf("Empty configuration storage type")
	}
	if len(css.SettingsKey) == 0 {
		return fmt.Errorf("%s. Empty settings key", subject)
	}

	switch css.Type {
	case ConfigurationStorageTypeFile:
		break
	case ConfigurationStorageTypeEtcd:
		for _, ep := range css.Endpoints {
			if _, err := url.ParseRequestURI(ep); err != nil {
				return fmt.Errorf("%s. Invalid etcd enpoint. %s", subject, err)
			}
		}
	case ConfigurationStorageTypeS3:
		if len(css.Region) == 0 {
			return fmt.Errorf("%s. Empty AWS region", subject)
		}
		if bucket := os.Getenv(identifoConfigBucketEnvName); len(bucket) != 0 {
			css.Bucket = bucket
		}
		if len(css.Bucket) == 0 {
			return fmt.Errorf("%s. Bucket for config is not set", subject)
		}
	default:
		return fmt.Errorf("%s. Unknown type", subject)
	}

	if err := css.KeyStorage.Validate(); err != nil {
		return err
	}
	return nil
}

const identifoStaticFilesBucketEnvName = "IDENTIFO_STATIC_FILES_BUCKET"

// Validate validates static files storage settings.
func (sfs *StaticFilesStorageSettings) Validate() error {
	subject := "StaticFilesStorageSettings"
	if sfs == nil {
		return fmt.Errorf("Nil %s", subject)
	}

	if len(sfs.Type) == 0 {
		return fmt.Errorf("%s. Empty static files storage type", subject)
	}
	if len(sfs.ServerConfigPath) == 0 {
		return fmt.Errorf("%s. Empty server config path", subject)
	}

	switch sfs.Type {
	case StaticFilesStorageTypeLocal:
		return nil
	case StaticFilesStorageTypeS3:
		if len(sfs.Region) == 0 {
			return fmt.Errorf("%s. Empty AWS region", subject)
		}
		if bucket := os.Getenv(identifoStaticFilesBucketEnvName); len(bucket) != 0 {
			sfs.Bucket = bucket
		}
		if len(sfs.Bucket) == 0 {
			return fmt.Errorf("%s. Bucket for static files is not set", subject)
		}
	case StaticFilesStorageTypeDynamoDB:
		if len(sfs.Region) == 0 {
			return fmt.Errorf("%s. Empty AWS region", subject)
		}
	default:
		return fmt.Errorf("%s. Unknown type", subject)
	}
	return nil
}

const identifoJWTKeysBucketEnvName = "IDENTIFO_JWT_KEYS_BUCKET"

// Validate validates key storage settings.
func (kss *KeyStorageSettings) Validate() error {
	subject := "KeyStorageSettings"
	if len(kss.Type) == 0 {
		return fmt.Errorf("%s. Empty key storage type", subject)
	}

	switch kss.Type {
	case KeyStorageTypeLocal:
		break
	case KeyStorageTypeS3:
		if len(kss.Region) == 0 {
			return fmt.Errorf("%s. Empty AWS region", subject)
		}
		if bucket := os.Getenv(identifoJWTKeysBucketEnvName); len(bucket) != 0 {
			kss.Bucket = bucket
		}
		if len(kss.Bucket) == 0 {
			return fmt.Errorf("%s. Bucket for keys is not set", subject)
		}
	default:
		return fmt.Errorf("%s. Unknown type '%s'", subject, kss.Type)
	}
	return nil
}

// Validate validates external services settings.
func (ess *ExternalServicesSettings) Validate() error {
	subject := "ExternalServicesSettings"
	if ess == nil {
		return fmt.Errorf("Nil %s", subject)
	}

	if err := ess.EmailService.Validate(); err != nil {
		return fmt.Errorf("%s. %s", subject, err)
	}
	if err := ess.SMSService.Validate(); err != nil {
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
	case SMSServiceTwilio:
		if len(sss.AccountSid)*len(sss.AuthToken)*len(sss.ServiceSid) == 0 {
			return fmt.Errorf("%s. Error creating Twilio SMS service, missing at least one of the parameters:"+
				"\n sidKey : %v\n tokenKey : %v\n ServiceSidKey : %v\n", subject, sss.AccountSid, sss.AuthToken, sss.ServiceSid)
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
			ess.Region = region
		}
		if sender := os.Getenv(awsSESSenderKey); len(sender) != 0 {
			ess.Sender = sender
		}
		if len(ess.Sender) == 0 {
			return fmt.Errorf("%s. Empty AWS sender", subject)
		}
		if len(ess.Region) == 0 {
			return fmt.Errorf("%s. Empty AWS region", subject)
		}
	case EmailServiceMailgun:
		if domain := os.Getenv(mailgunDomainKey); len(domain) != 0 {
			ess.Domain = domain
		}
		if publicKey := os.Getenv(mailgunPublicKey); len(publicKey) != 0 {
			ess.PublicKey = publicKey
		}
		if privateKey := os.Getenv(mailgunPrivateKey); len(privateKey) != 0 {
			ess.PrivateKey = privateKey
		}
		if sender := os.Getenv(mailgunSenderKey); len(sender) != 0 {
			ess.Sender = sender
		}

		if len(ess.Domain) == 0 {
			return fmt.Errorf("%s. Empty Mailgun domain", subject)
		}
		if len(ess.PublicKey)*len(ess.PrivateKey) == 0 {
			return fmt.Errorf("%s. At least one of the keys is empty", subject)
		}
		if len(ess.Sender) == 0 {
			return fmt.Errorf("%s. Empty Mailgun sender", subject)
		}
	default:
		return fmt.Errorf("%s. Unknown type", subject)
	}
	return nil
}
