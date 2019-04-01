package html

// Error - http level error type.
type Error string

// Error implements std.Error interface.
func (e Error) Error() string { return string(e) }

const (
	// ErrorRegistrationForbidden means that registration is forbidden.
	ErrorRegistrationForbidden = Error("Registration in this app is forbidden.")
)
