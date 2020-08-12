package middleware

//Error represents error for middleware
type Error string

//Description returns description for error
func (e Error) Description() string {
	return descriptions[e]
}

var descriptions = map[Error]string{
	ErrorTokenIsEmpty:   "Bearer JWT token is missing",
	ErrorTokenIsInvalid: "Token validation failed",
}

const (
	//ErrorTokenIsEmpty means that middleware could not find Bearer JWT token in headers
	ErrorTokenIsEmpty = "middleware.token_is_empty"
	//ErrorTokenIsInvalid token validation failed
	ErrorTokenIsInvalid = "middleware.token_is_invalid"
)
