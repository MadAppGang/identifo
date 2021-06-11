package boltdb

// Error - domain level error type
type Error string

// Error - implementation of std.Error protocol
func (e Error) Error() string { return string(e) }

const (
	ErrorWrongDataFormat   = Error("wrong data format") // ErrorWrongDataFormat is for corrupted request data.
	ErrorInactiveUser      = Error("user is inactive")  // ErrorInactiveUser means that user is inactive.
	ErrorEmptyAppID        = Error("empty appID param") // ErrorEmptyAppID means appID params is empty
	ErrorInactiveApp       = Error("app is inactive")   // ErrorInactiveApp means app is inactive
	ErrorEmptyDatabasePath = Error("unable to init boltdb storage with empty database path")
)
