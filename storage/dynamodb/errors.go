package dynamodb

// Error - domain level error type
type Error string

// Error - implementation of std.Error protocol
func (e Error) Error() string { return string(e) }

const (
	// ErrorInternalError internal error
	ErrorInternalError = Error("internal error")
	// ErrorInactiveUser means user is inactive
	ErrorInactiveUser = Error("user is inactive")
	// ErrorEmptyAppID means appID params is empty
	ErrorEmptyAppID = Error("empty appID param")
	// ErrorInactiveApp means app is inactive
	ErrorInactiveApp = Error("app is inactive")
	// ErrorEmptyEndpointRegion endpoint or region are settings
	ErrorEmptyEndpointRegion = Error("endpoint and region required for dynamodb user storage")
)
