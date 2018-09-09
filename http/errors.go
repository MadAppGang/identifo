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
)
