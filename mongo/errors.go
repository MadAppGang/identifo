package mongo

//Error - domain level error type
type Error string

//Error - implementation of std.Error protocol
func (e Error) Error() string { return string(e) }

const (
	// ErrorInactiveUser means user is inactive
	ErrorInactiveUser = Error("User is inactive")
)
