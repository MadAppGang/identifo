package mongo

//Error - domain level error type
type Error string

//Error - implementation of std.Error protocol
func (e Error) Error() string { return string(e) }

const (
	// ErrorInactiveUser means user is inactive
	ErrorInactiveUser = Error("User is inactive")
	// ErrorEmptyAppID means appID params is empty
	ErrorEmptyAppID = Error("Empty appID param")
	// ErrorInactiveApp means app is inactive
	ErrorInactiveApp = Error("App is inactive")
)
