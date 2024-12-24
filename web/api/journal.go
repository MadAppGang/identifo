package api

import (
	"github.com/madappgang/identifo/v2/logging"
	"github.com/madappgang/identifo/v2/model"
)

type AuditOperation string

const (
	AuditOperationLoginWithPassword AuditOperation = "login_with_password"
	AuditOperationLoginWithPhone    AuditOperation = "login_with_phone"
	AuditOperationLoginWith2FA      AuditOperation = "login_with_2fa"
	AuditOperationRefreshToken      AuditOperation = "refresh_token"
	AuditOperationOIDCLogin         AuditOperation = "oidc_login"
	AuditOperationFederatedLogin    AuditOperation = "federated_login"
	AuditOperationRegistration      AuditOperation = "registration"
	AuditOperationLogout            AuditOperation = "logout"
	AuditOperationImpersonatedAs    AuditOperation = "impersonated_as"
)

func (ar *Router) audit(
	op AuditOperation,
	userID, appID, device, accessRole string,
	scopes []string,
	accessToken, refreshToken string,
) {
	iss := ar.server.Services().Token.Issuer()

	auditSettings := ar.server.Settings().Audit

	accessToken = maskToken(accessToken, auditSettings.TokenRecording)
	refreshToken = maskToken(refreshToken, auditSettings.TokenRecording)

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
		"accessToken", accessToken,
		"refreshToken", refreshToken,
	)
}

func maskToken(token string, tokenRecording model.TokenRecording) string {
	switch tokenRecording {
	case model.TokenRecordingNone:
		return "<redacted>"
	case model.TokenRecordingObfuscated:
		if len(token) < 32 {
			return "<short>"
		}

		return token[:6] + "..." + token[len(token)-6:]
	case model.TokenRecordingFull:
		return token
	default:
		return "<redacted>"
	}
}
