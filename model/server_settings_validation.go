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
func (ss *ServerSettings) Validate(rewriteDefaults bool) []error {
	result := []error{}
	if rewriteDefaults {
		ss.RewriteDefaults()
	}

	if err := ss.General.Validate(); len(err) > 0 {
		result = append(result, err...)
	}
	if err := ss.Storage.Validate(); len(err) > 0 {
		result = append(result, err...)
	}
	if err := ss.SessionStorage.Validate(); len(err) > 0 {
		result = append(result, err...)
	}
	if err := ss.LoginWebApp.Validate(); len(err) > 0 {
		result = append(result, err...)
	}
	if err := ss.AdminPanel.Validate(); err != nil {
		result = append(result, err)
	}
	if err := ss.EmailTemplates.Validate(); len(err) > 0 {
		result = append(result, err...)
	}

	if err := ss.Services.Validate(); len(err) > 0 {
		result = append(result, err...)
	}
	return result
}

// Validate validates general services settings.
func (gss *GeneralServerSettings) Validate() []error {
	subject := "GeneralServerSettings"
	result := []error{}

	if _, err := url.ParseRequestURI(gss.Host); err != nil {
		result = append(result, fmt.Errorf("%s. Host is invalid. %s", subject, err))
	}
	if len(gss.Issuer) == 0 {
		result = append(result, fmt.Errorf("%s. Issuer is not set", subject))
	}
	return result
}

// Validate validates storage settings.
func (ss *StorageSettings) Validate() []error {
	result := []error{}

	if err := ss.AppStorage.Validate(); err != nil {
		result = append(result, fmt.Errorf("AppStorage settings: %s", err))
	}
	if err := ss.UserStorage.Validate(); err != nil {
		result = append(result, fmt.Errorf("AppStorage settings: %s", err))
	}
	if err := ss.TokenStorage.Validate(); err != nil {
		result = append(result, fmt.Errorf("TokenStorage settings: %s", err))
	}
	if err := ss.TokenBlacklist.Validate(); err != nil {
		result = append(result, fmt.Errorf("TokenBlacklist settings: %s", err))
	}
	if err := ss.VerificationCodeStorage.Validate(); err != nil {
		result = append(result, fmt.Errorf("VerificationCodeStorage settings: %s", err))
	}
	if err := ss.InviteStorage.Validate(); err != nil {
		result = append(result, fmt.Errorf("InviteStorage settings: %s", err))
	}
	if ss.AppStorage.Type == DBTypeDefault ||
		ss.UserStorage.Type == DBTypeDefault ||
		ss.TokenStorage.Type == DBTypeDefault ||
		ss.TokenBlacklist.Type == DBTypeDefault ||
		ss.VerificationCodeStorage.Type == DBTypeDefault ||
		ss.InviteStorage.Type == DBTypeDefault {
		// if one of the storages is reference default storage, let' validate default storage
		if err := ss.DefaultStorage.Validate(); err != nil {
			result = append(result, fmt.Errorf("DefaultStorage settings: %s", err))
		}
		if ss.DefaultStorage.Type == DBTypeDefault {
			result = append(result, fmt.Errorf("DefaultStorage settings could not be of type Default"))
		}
	}
	return result
}

// Validate validates database settings.
func (dbs *DatabaseSettings) Validate() error {
	switch dbs.Type {
	case DBTypeFake:
	case DBTypeDefault:
	case DBTypeMem:
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
	case DBTypePlugin:
		if len(dbs.Plugin.Cmd) == 0 {
			return fmt.Errorf("empty CMD for grpc")
		}
	case DBTypeGRPC:

	default:
		return fmt.Errorf("unsupported database type %s", dbs.Type)
	}
	return nil
}

// Validate validates admin session storage settings.
func (sss *SessionStorageSettings) Validate() []error {
	subject := "SessionStorageSettings"
	result := []error{}

	if len(sss.Type) == 0 {
		result = append(result, fmt.Errorf("Empty session storage type"))
	}
	if sss.SessionDuration.Duration == 0 {
		result = append(result, fmt.Errorf("%s. Session duration is 0 seconds", subject))
	}

	switch sss.Type {
	case SessionStorageMem:
		break
	case SessionStorageRedis:
		if _, err := url.ParseRequestURI(sss.Redis.Address); err != nil {
			result = append(result, fmt.Errorf("%s. Invalid address. %s", subject, err))
		}
	case SessionStorageDynamoDB:
		if len(sss.Dynamo.Region) == 0 {
			result = append(result, fmt.Errorf("%s. Empty AWS region", subject))
		}
		if len(sss.Dynamo.Endpoint) == 0 {
			result = append(result, fmt.Errorf("%s. Empty AWS Dynamo endpoint", subject))
		}

	default:
		result = append(result, fmt.Errorf("%s. Unknown type", subject))
	}
	return result
}

// Validate validates login web app settings
func (sfs *FileStorageSettings) Validate() []error {
	subject := "LoginWebAppSettings"
	result := []error{}

	switch sfs.Type {
	case FileStorageTypeDefault:
		break
	case FileStorageTypeNone:
		break
	case FileStorageTypeLocal:
		if len(sfs.Local.Path) == 0 {
			result = append(result, fmt.Errorf("%s. empty folder", subject))
		}
	case FileStorageTypeS3:
		if len(sfs.S3.Region) == 0 {
			result = append(result, fmt.Errorf("%s. Empty AWS region", subject))
		}
		if bucket := os.Getenv(identifoLoginWebAppBucket); len(bucket) != 0 {
			sfs.S3.Bucket = bucket
		}
		if len(sfs.S3.Bucket) == 0 {
			result = append(result, fmt.Errorf("%s. empty s3 bucket for login web app", subject))
		}
	default:
		result = append(result, fmt.Errorf("%s. Unknown type", subject))
	}
	return result
}

func (kss *AdminPanelSettings) Validate() error {
	// not it has just one bool field which is always valid
	return nil
}

// Validate validates external services settings.
func (ess *ServicesSettings) Validate() []error {
	result := []error{}

	if err := ess.Email.Validate(); len(err) > 0 {
		result = append(result, err...)
	}
	if err := ess.SMS.Validate(); len(err) > 0 {
		result = append(result, err...)
	}
	return result
}

// Validate validates SMS service settings.
func (sss *SMSServiceSettings) Validate() []error {
	subject := "SMSServiceSettings"
	result := []error{}
	if len(sss.Type) == 0 {
		return []error{fmt.Errorf("Empty SMS service type")}
	}

	switch sss.Type {
	case SMSServiceMock:
		break
	case SMSServiceNexmo:
		if len(sss.Nexmo.APIKey) == 0 {
			result = append(result, fmt.Errorf("%s. error creating Nexmo SMS service, API key is empty", subject))
		}
		if len(sss.Nexmo.APISecret) == 0 {
			result = append(result, fmt.Errorf("%s. error creating Nexmo SMS service, Nexmo secret is empty", subject))
		}
	case SMSServiceTwilio:
		if len(sss.Twilio.AccountSid) == 0 {
			result = append(result, fmt.Errorf("%s. error creating Twilio SMS service, missing Account SID", subject))
		}
		if len(sss.Twilio.AuthToken) == 0 {
			result = append(result, fmt.Errorf("%s. error creating Twilio SMS service, missing Auth Token", subject))
		}
		if len(sss.Twilio.ServiceSid) == 0 && len(sss.Twilio.SendFrom) == 0 {
			result = append(result, fmt.Errorf("%s. error creating Twilio SMS service, missing Service SID or  SendFrom", subject))
		}
	case SMSServiceRouteMobile:
		if len(sss.Routemobile.Username) == 0 {
			result = append(result, fmt.Errorf("%s. Error creating RouteMobile SMS service, missing username", subject))
		}
		if len(sss.Routemobile.Password) == 0 {
			result = append(result, fmt.Errorf("%s. Error creating RouteMobile SMS service, missing password", subject))
		}
		if len(sss.Routemobile.Source) == 0 {
			result = append(result, fmt.Errorf("%s. Error creating RouteMobile SMS service, missing source", subject))
		}
		if sss.Routemobile.Region != RouteMobileRegionUAE {
			result = append(result, fmt.Errorf("%s. Error creating RouteMobile SMS service, region %s is not supported", subject, sss.Routemobile.Region))
		}
	default:
		result = append(result, fmt.Errorf("%s. Unknown type", subject))
	}
	return result
}

const (
	// mailgunDomainKey is a name of env variable that contains Mailgun domain value.
	mailgunDomainKey = "MAILGUN_DOMAIN"
	// mailgunPrivateKey is a name of env variable that contains Mailgun private key value.
	mailgunPrivateKey = "MAILGUN_PRIVATE_KEY"
	// mailgunPublicKey is a name of env variable that contains Mailgun public key value.
	mailgunSenderKey = "MAILGUN_SENDER"

	// awsSESSenderKey is a name of env variable that contains AWS SWS sender value.
	awsSESSenderKey = "AWS_SES_SENDER"
	// awsSESRegionKey is a name of env variable that contains AWS SWS region value.
	awsSESRegionKey = "AWS_SES_REGION"
)

// Validate validates email service settings.
func (ess *EmailServiceSettings) Validate() []error {
	subject := "EmailServiceSettings"
	result := []error{}

	if len(ess.Type) == 0 {
		return []error{fmt.Errorf("%s. Empty email service type", subject)}
	}

	switch ess.Type {
	case EmailServiceMock:
		break
	case EmailServiceAWS:
		if region := os.Getenv(awsSESRegionKey); len(region) != 0 {
			ess.SES.Region = region
		}
		if sender := os.Getenv(awsSESSenderKey); len(sender) != 0 {
			ess.SES.Sender = sender
		}
		if len(ess.SES.Sender) == 0 {
			result = append(result, fmt.Errorf("%s. Empty AWS sender", subject))
		}
		if len(ess.SES.Region) == 0 {
			result = append(result, fmt.Errorf("%s. Empty AWS region", subject))
		}
	case EmailServiceMailgun:
		if domain := os.Getenv(mailgunDomainKey); len(domain) != 0 {
			ess.Mailgun.Domain = domain
		}
		if privateKey := os.Getenv(mailgunPrivateKey); len(privateKey) != 0 {
			ess.Mailgun.PrivateKey = privateKey
		}
		if sender := os.Getenv(mailgunSenderKey); len(sender) != 0 {
			ess.Mailgun.Sender = sender
		}

		if len(ess.Mailgun.Domain) == 0 {
			result = append(result, fmt.Errorf("%s. Empty Mailgun domain", subject))
		}
		if len(ess.Mailgun.PrivateKey) == 0 {
			result = append(result, fmt.Errorf("%s. Mailgun private key is empty", subject))
		}
		if len(ess.Mailgun.Sender) == 0 {
			result = append(result, fmt.Errorf("%s. Empty Mailgun sender", subject))
		}
	default:
		result = append(result, fmt.Errorf("%s. Unknown type", subject))
	}
	return result
}
