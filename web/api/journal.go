package api

import (
	"net/http"

	"github.com/madappgang/identifo/v2/logging"
)

type JournalOperation string

const (
	JournalOperationLoginWithPassword JournalOperation = "login_with_password"
	JournalOperationLoginWithPhone    JournalOperation = "login_with_phone"
	JournalOperationRefreshToken      JournalOperation = "refresh_token"
	JournalOperationOIDCLogin         JournalOperation = "oidc_login"
	JournalOperationFederatedLogin    JournalOperation = "federated_login"
)

func (ar *Router) journal(
	op JournalOperation,
	userID, appID string,
	req *http.Request,
) {
	ar.logger.Info(string(op),
		logging.FieldComponent, "journal",
		logging.FieldUserID, userID,
		logging.FieldAppID, appID,
		"device", req.UserAgent())
}
