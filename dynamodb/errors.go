package dynamodb

//Error - domain level error type
type Error string

//Error - implementation of std.Error protocol
func (e Error) Error() string { return string(e) }

const (
	//ErrorNotFound general not found error
	ErrorNotFound = Error("not found")
	//ErrorWrongDataFormat wrong input data provided
	ErrorWrongDataFormat = Error("wrong data format")
	//ErrorNotImplemented requested feature is not implemented yet
	ErrorNotImplemented = Error("Not implemented")
	//ErrorUserExists user already exists
	ErrorUserExists = Error("User is already exists")
	//ErrorInternalError internal error
	ErrorInternalError = Error("Internal error")
)
