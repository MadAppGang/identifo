package model

//TokenService manage tokens abstraction layer
type TokenService interface {
	NewToken(u User, scopes []string, appID string) (Token, error)
	Parse(string) (Token, error)
	String(Token) (string, error)
}

//Token is app token to give user chan
type Token interface {
	Validate() error
}

//Validator calidate token with external requester
type Validator interface {
	Validate(Token) error
}

//TokenMapping is service to match tokens to services. etc
type TokenMapping interface {
}
