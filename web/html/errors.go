package html

// Error - http level error type.
type Error string

// Error - implementation of std.Error protocol.
func (e Error) Error() string { return string(e) }

const (
	// ErrorRegistrationForbidden forbidden registration.
	ErrorRegistrationForbidden = Error("Registration in this app is forbidden.")
)
