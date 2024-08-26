package api

import (
	"github.com/madappgang/identifo/v2/logging"
)

type JournalOperation string

const (
	JournalOperationLoginWithPassword JournalOperation = "login_with_password"
	JournalOperationLoginWithPhone    JournalOperation = "login_with_phone"
	JournalOperationLoginWith2FA      JournalOperation = "login_with_2fa"
	JournalOperationRefreshToken      JournalOperation = "refresh_token"
	JournalOperationOIDCLogin         JournalOperation = "oidc_login"
	JournalOperationFederatedLogin    JournalOperation = "federated_login"
	JournalOperationRegistration      JournalOperation = "registration"
	JournalOperationLogout            JournalOperation = "logout"
	JournalOperationImpersonatedAs    JournalOperation = "impersonated_as"
)

func (ar *Router) journal(
	op JournalOperation,
	userID, appID, device, accessRole string,
	scopes []string,
) {
	iss := ar.server.Services().Token.Issuer()

	// TODO: Create an interface for the audit log
	// Implement it for logging to stdout, a database, or a remote service
	ar.logger.Info("audit_record",
		"operation", string(op),
		logging.FieldUserID, userID,
		logging.FieldAppID, appID,
		"device", device,
		"issuer", iss,
		"accessRole", accessRole,
		"scopes", scopes,
	)
}
