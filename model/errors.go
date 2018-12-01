package model

//Error - domain level error type
type Error string

//Error - implementation of std.Error protocol
func (e Error) Error() string { return string(e) }

const (
	//ErrorNotFound general not found error
	ErrorNotFound        = Error("not found")
	ErrorWrongDataFormat = Error("wrong data format")
	ErrorUserExists      = Error("User already exists")
	ErrorNotImplemented  = Error("Not implemented")

	//ErrorPasswordShouldHave7Letter strong password  validation
	ErrorPasswordShouldHave7Letter = Error("Password should have at least seven letters")
	//ErrorPasswordNoNumbers strong password validation
	ErrorPasswordNoNumbers = Error("Password should have at least one number")
	//ErrorPasswordNoUppercase strong password validation
	ErrorPasswordNoUppercase = Error("Password should have at least one uppercase symbol")
	//ErrorPasswordWrongSymbols strong password validation
	ErrorPasswordWrongSymbols = Error("Password contains wrong symbols")
)
