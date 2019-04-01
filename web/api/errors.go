package api

// Error - http level error type.
type Error string

// Error implements std.Error interface.
func (e Error) Error() string { return string(e) }

const (
	// ErrorWrongInput means that request data is corrupted.
	ErrorWrongInput = Error("Wrong input data")
	// ErrorRequestSignature is a HMAC request signature error.
	ErrorRequestSignature = Error("Incorrect or empty request signature")
	// ErrorRequestInvalidAppID means that application ID header value is invalid.
	ErrorRequestInvalidAppID = Error("Incorrect or empty application ID")
	// ErrorRequestInactiveApp means that the reqesting app is inactive.
	ErrorRequestInactiveApp = Error("Requesting app is inactive")
	// ErrorRequestInvalidToken means that the token is invalid or empty.
	ErrorRequestInvalidToken = Error("Incorrect or empty Bearer token")
	// ErrorRegistrationForbidden means that registration is forbidden.
	ErrorRegistrationForbidden = Error("Registration in this app is forbidden.")

	// ErrorFederatedProviderIsNotSupported means that the federated ID provider is not supported.
	ErrorFederatedProviderIsNotSupported = Error("Federated provider is not supported")
	// ErrorFederatedProviderEmptyUserID means that the federated provider returns empty user ID, maybe access token does not have required permissions.
	ErrorFederatedProviderEmptyUserID = Error("Federated provider returns empty user ID")
)
