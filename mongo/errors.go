package mongo

//Error - domain level error type
type Error string

//Error - implementation of std.Error protocol
func (e Error) Error() string { return string(e) }

const (
	//ErrorNotFound general not found error
	ErrorNotFound = Error("not found")
	//ErrorWrongDataFormat wrong input data provided
	ErrorWrongDataFormat = Error("wrong data format")
	//ErrorNotImplemented reauested feature is not implemented yet
	ErrorNotImplemented = Error("Not implemneted")
)
