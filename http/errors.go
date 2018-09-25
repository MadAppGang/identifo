package http

//Error - http level error type
type Error string

//Error - implementation of std.Error protocol
func (e Error) Error() string { return string(e) }

const (
	//ErrorWrongInput request data is corrupted
	ErrorWrongInput = Error("Wrong input data")
	//ErrorRequestSignature HMAC request signature error
	ErrorRequestSignature = Error("Incorrect or empty request signature")
	//ErrorRequestInvalidAppID application ID header value is invalid
	ErrorRequestInvalidAppID = Error("Incorrect or empty application ID")
	//ErrorRequestInactiveApp the reqesting app is inactive
	ErrorRequestInactiveApp = Error("Requesting app is inactive")

	//ErrorPasswordShouldHave7Letter strong password  validation
	ErrorPasswordShouldHave7Letter = Error("Password should have at least seven letters")
	//ErrorPasswordNoNumbers strong password validation
	ErrorPasswordNoNumbers = Error("Password should have at least one number")
	//ErrorPasswordNoUppercase strong password validation
	ErrorPasswordNoUppercase = Error("Pussword should have at least one uppercase symbol")
	//ErrorPasswordWrongSymbols strong password validation
	ErrorPasswordWrongSymbols = Error("Password contains wrong symbols")
)
