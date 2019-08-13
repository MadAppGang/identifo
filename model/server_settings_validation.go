package model

import (
	"fmt"
	"net/url"
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
	if len(gss.PublicKeyPath)*len(gss.PublicKeyPath) == 0 {
		return fmt.Errorf("%s. At least one of key paths is empty", subject)
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
	if dbs == nil {
		return fmt.Errorf("Nil DatabaseSettings")
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
	}
	return nil
}

// Validate validates external services settings.
func (ess *ExternalServicesSettings) Validate() error {
	subject := "ExternalServicesSettings"
	if ess == nil {
		return fmt.Errorf("Nil %s", subject)
	}

	if len(ess.MailService) == 0 {
		return fmt.Errorf("%s. Empty mail service settings", subject)
	}
	return ess.SMSService.Validate()
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
	}
	return nil
}
