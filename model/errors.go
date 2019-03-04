package model

//Error - domain level error type
type Error string

//Error - implementation of std.Error protocol
func (e Error) Error() string { return string(e) }

const (
	// ErrorNotFound is a general not found error.
	ErrorNotFound = Error("not found")
	// ErrorWrongDataFormat is for corrupted request data.
	ErrorWrongDataFormat = Error("wrong data format")
	// ErrorUserExists is for unwanted user entry presense.
	ErrorUserExists = Error("User already exists")
	// ErrorNotImplemented is for features that are not implemented yet.
	ErrorNotImplemented = Error("Not implemented")

	// ErrorPasswordShouldHave7Letter is for failed password strength check.
	ErrorPasswordShouldHave7Letter = Error("Password should have at least seven letters")
	// ErrorPasswordNoNumbers is for failed password strength check.
	ErrorPasswordNoNumbers = Error("Password should have at least one number")
	// ErrorPasswordNoUppercase is for failed password strength check.
	ErrorPasswordNoUppercase = Error("Password should have at least one uppercase symbol")
	// ErrorPasswordWrongSymbols is for failed password strength check.
	ErrorPasswordWrongSymbols = Error("Password contains wrong symbols")
)
