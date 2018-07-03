package identifo

//Error - domain level error type
type Error string

//Error - implementation of std.Error protocol
func (e Error) Error() string { return string(e) }

const (
	//ErrorInternal represents internal server error, used to mask real internal problem
	ErrorInternal = Error("internal error")
)
