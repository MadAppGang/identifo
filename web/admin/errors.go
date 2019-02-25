package admin

//Error - http level error type
type Error string

//Error - implementation of std.Error protocol
func (e Error) Error() string { return string(e) }

const (
	// ErrorWrongInput is for corrupted request data.
	ErrorWrongInput = Error("Wrong input data")
	// ErrorRequestInvalidCookie is for invalid cookie.
	ErrorRequestInvalidCookie = Error("Invalid cookie")
	// ErrorInternalError is for internal errors.
	ErrorInternalError = Error("Internal error")
	// ErrorIncorrectLogin is for incorrect login and password.
	ErrorIncorrectLogin = Error("Incorrect login information")
	// ErrorNotAuthorized is for non-authorized access intents.
	ErrorNotAuthorized = Error("Not authorized")
)
