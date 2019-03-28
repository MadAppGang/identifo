package api

// Error - http level error type.
type Error string

// Error - implementation of std.Error protocol.
func (e Error) Error() string { return string(e) }

const (
	// ErrorWrongInput request data is corrupted.
	ErrorWrongInput = Error("Wrong input data")
	// ErrorRequestSignature HMAC request signature error.
	ErrorRequestSignature = Error("Incorrect or empty request signature")
	// ErrorRequestInvalidAppID application ID header value is invalid.
	ErrorRequestInvalidAppID = Error("Incorrect or empty application ID")
	// ErrorRequestInactiveApp the reqesting app is inactive.
	ErrorRequestInactiveApp = Error("Requesting app is inactive")
	// ErrorRequestInvalidToken invalid or empty token.
	ErrorRequestInvalidToken = Error("Incorrect or empty Bearer token")
	// ErrorRegistrationForbidden forbidden registration.
	ErrorRegistrationForbidden = Error("Registration in this app is forbidden.")

	// ErrorFederatedProviderIsNotSupported federated ID provider is not supported.
	ErrorFederatedProviderIsNotSupported = Error("Federated provider is not supported")
	// ErrorFederatedProviderEmptyUserID federated provider returns empty user ID, maybe access token does not have required permissions.
	ErrorFederatedProviderEmptyUserID = Error("Federated provider returns empty user ID")
)
